package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"github.com/yourusername/roundtable/internal/core"
	"github.com/yourusername/roundtable/internal/store"
)

// ConnectionManager manages WebSocket connections for all rooms.
type ConnectionManager struct {
	store       store.Store
	connections map[string]*Connection // playerID → Connection
	mu          sync.RWMutex
}

// NewConnectionManager creates a new connection manager.
func NewConnectionManager(store store.Store) *ConnectionManager {
	return &ConnectionManager{
		store:       store,
		connections: make(map[string]*Connection),
	}
}

// Connection represents a single WebSocket connection.
type Connection struct {
	PlayerID string
	RoomCode string
	Conn     *websocket.Conn
	Send     chan ServerMessage
	ctx      context.Context
	cancel   context.CancelFunc
}

// HandleConnection manages a WebSocket connection lifecycle.
func (cm *ConnectionManager) HandleConnection(ctx context.Context, conn *websocket.Conn, roomCode string) {
	// First, client must authenticate with session token
	var authMsg ClientMessage
	if err := wsjson.Read(ctx, conn, &authMsg); err != nil {
		log.Printf("Failed to read auth message: %v", err)
		conn.Close(websocket.StatusPolicyViolation, "authentication required")
		return
	}

	if authMsg.Type != ClientMsgAuthenticate {
		log.Printf("Expected authenticate message, got: %s", authMsg.Type)
		conn.Close(websocket.StatusPolicyViolation, "authentication required")
		return
	}

	var authPayload AuthenticatePayload
	if err := json.Unmarshal(authMsg.Payload, &authPayload); err != nil {
		log.Printf("Failed to parse auth payload: %v", err)
		conn.Close(websocket.StatusPolicyViolation, "invalid authentication")
		return
	}

	// Validate session token and get player
	room, err := cm.store.GetRoom(roomCode)
	if err != nil {
		log.Printf("Room not found: %s", roomCode)
		conn.Close(websocket.StatusPolicyViolation, "room not found")
		return
	}

	player, err := room.GetPlayerByToken(authPayload.SessionToken)
	if err != nil {
		log.Printf("Invalid session token for room %s", roomCode)
		conn.Close(websocket.StatusPolicyViolation, "invalid session token")
		return
	}

	// Mark player as connected
	player.Reconnect()

	// Create connection context
	connCtx, cancel := context.WithCancel(ctx)
	connection := &Connection{
		PlayerID: player.ID,
		RoomCode: roomCode,
		Conn:     conn,
		Send:     make(chan ServerMessage, 256),
		ctx:      connCtx,
		cancel:   cancel,
	}

	// Register connection
	cm.mu.Lock()
	// Close existing connection if player reconnecting
	if existingConn, exists := cm.connections[player.ID]; exists {
		existingConn.Close()
	}
	cm.connections[player.ID] = connection
	cm.mu.Unlock()

	log.Printf("Player %s (%s) connected to room %s", player.DisplayName, player.ID, roomCode)

	// Send authenticated message with current state
	authResponse, _ := NewAuthenticatedMessage(player.ID, room.GetState())
	connection.Send <- authResponse

	// Send event history for this player
	events := room.GetEventsForPlayer(player.ID)
	if len(events) > 0 {
		eventsMsg, _ := NewEventsMessage(events)
		connection.Send <- eventsMsg
	}

	// Broadcast player reconnected event if game in progress
	if room.Status == core.RoomStatusPlaying {
		event, _ := core.NewPublicEvent(core.EventPlayerReconnected, "system", core.PlayerReconnectedPayload{
			PlayerID: player.ID,
		})
		room.AppendEvent(event)
		cm.BroadcastEvent(roomCode, event)
	}

	// Start read and write pumps
	go connection.writePump()
	connection.readPump(cm, room)

	// Cleanup on disconnect
	cm.handleDisconnect(connection, room)
}

// readPump reads messages from the WebSocket connection.
func (c *Connection) readPump(cm *ConnectionManager, room *core.Room) {
	defer c.cancel()

	for {
		var msg ClientMessage
		if err := wsjson.Read(c.ctx, c.Conn, &msg); err != nil {
			if c.ctx.Err() != nil {
				// Context cancelled, clean shutdown
				return
			}
			log.Printf("Read error for player %s: %v", c.PlayerID, err)
			return
		}

		// Handle message
		cm.handleClientMessage(c, room, msg)
	}
}

// writePump sends messages to the WebSocket connection.
func (c *Connection) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-c.Send:
			ctx, cancel := context.WithTimeout(c.ctx, 10*time.Second)
			err := wsjson.Write(ctx, c.Conn, msg)
			cancel()

			if err != nil {
				log.Printf("Write error for player %s: %v", c.PlayerID, err)
				return
			}

		case <-ticker.C:
			// Send ping to keep connection alive
			ctx, cancel := context.WithTimeout(c.ctx, 10*time.Second)
			err := c.Conn.Ping(ctx)
			cancel()

			if err != nil {
				log.Printf("Ping error for player %s: %v", c.PlayerID, err)
				return
			}

		case <-c.ctx.Done():
			return
		}
	}
}

// handleClientMessage processes incoming client messages.
func (cm *ConnectionManager) handleClientMessage(conn *Connection, room *core.Room, msg ClientMessage) {
	switch msg.Type {
	case ClientMsgPing:
		pong, _ := NewPongMessage()
		conn.Send <- pong

	case ClientMsgAction:
		var actionPayload ActionPayload
		if err := json.Unmarshal(msg.Payload, &actionPayload); err != nil {
			errMsg, _ := NewErrorMessage("Invalid action payload")
			conn.Send <- errMsg
			return
		}

		// Process action
		events, err := room.ProcessAction(conn.PlayerID, actionPayload.Action)
		if err != nil {
			errMsg, _ := NewErrorMessage(fmt.Sprintf("Action failed: %v", err))
			conn.Send <- errMsg
			return
		}

		// Broadcast events to affected players
		for _, event := range events {
			cm.BroadcastEvent(room.ID, event)
		}

	default:
		errMsg, _ := NewErrorMessage(fmt.Sprintf("Unknown message type: %s", msg.Type))
		conn.Send <- errMsg
	}
}

// handleDisconnect cleans up after a connection closes.
func (cm *ConnectionManager) handleDisconnect(conn *Connection, room *core.Room) {
	cm.mu.Lock()
	delete(cm.connections, conn.PlayerID)
	cm.mu.Unlock()

	// Mark player as disconnected
	player, _ := room.GetPlayer(conn.PlayerID)
	if player != nil {
		player.Disconnect()
		log.Printf("Player %s (%s) disconnected from room %s", player.DisplayName, player.ID, room.ID)
	}
}

// Close closes the connection.
func (c *Connection) Close() {
	c.cancel()
	c.Conn.Close(websocket.StatusNormalClosure, "connection closed")
}

// BroadcastEvent sends an event to all players who can see it.
func (cm *ConnectionManager) BroadcastEvent(roomCode string, event core.GameEvent) {
	room, err := cm.store.GetRoom(roomCode)
	if err != nil {
		log.Printf("Failed to get room %s: %v", roomCode, err)
		return
	}

	eventMsg, _ := NewEventMessage(event)

	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for _, player := range room.GetPlayers() {
		if !event.CanPlayerSee(player.ID) {
			continue
		}

		conn, exists := cm.connections[player.ID]
		if !exists || conn.RoomCode != roomCode {
			continue
		}

		select {
		case conn.Send <- eventMsg:
		default:
			log.Printf("Failed to send event to player %s: channel full", player.ID)
		}
	}
}

// BroadcastRoomState sends updated room state to all connected players.
func (cm *ConnectionManager) BroadcastRoomState(roomCode string) {
	room, err := cm.store.GetRoom(roomCode)
	if err != nil {
		log.Printf("Failed to get room %s: %v", roomCode, err)
		return
	}

	stateMsg, _ := NewRoomStateMessage(room.GetState())

	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for _, player := range room.GetPlayers() {
		conn, exists := cm.connections[player.ID]
		if !exists || conn.RoomCode != roomCode {
			continue
		}

		select {
		case conn.Send <- stateMsg:
		default:
			log.Printf("Failed to send room state to player %s: channel full", player.ID)
		}
	}
}
