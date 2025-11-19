package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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
	// nhooyr.io/websocket uses context for timeouts, not SetReadLimit/SetReadDeadline
	// The message size limit is handled by the library's default (32KB)
	// For larger messages, the library will automatically stream them

	// Set 10-second timeout for auth message
	authCtx, authCancel := context.WithTimeout(ctx, 10*time.Second)
	defer authCancel()

	// First, client must authenticate with session token
	var authMsg ClientMessage
	if err := wsjson.Read(authCtx, conn, &authMsg); err != nil {
		slog.Error("failed to read auth message", "error", err)
		conn.Close(websocket.StatusPolicyViolation, "authentication required")
		return
	}

	if authMsg.Type != ClientMsgAuthenticate {
		slog.Warn("expected authenticate message", "got", authMsg.Type)
		conn.Close(websocket.StatusPolicyViolation, "authentication required")
		return
	}

	var authPayload AuthenticatePayload
	if err := json.Unmarshal(authMsg.Payload, &authPayload); err != nil {
		slog.Error("failed to parse auth payload", "error", err)
		conn.Close(websocket.StatusPolicyViolation, "invalid authentication")
		return
	}

	// Validate session token and get player
	room, err := cm.store.GetRoom(roomCode)
	if err != nil {
		slog.Warn("room not found", "roomCode", roomCode)
		conn.Close(websocket.StatusPolicyViolation, "room not found")
		return
	}

	player, err := room.GetPlayerByToken(authPayload.SessionToken)
	if err != nil {
		slog.Warn("invalid session token", "roomCode", roomCode)
		conn.Close(websocket.StatusPolicyViolation, "invalid session token")
		return
	}

	// Mark player as connected
	player.Reconnect()

	// Create connection context detached from HTTP request context
	// This prevents Railway's HTTP timeout from affecting the WebSocket connection
	connCtx, cancel := context.WithCancel(context.Background())
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

	slog.Info("player connected",
		"playerName", player.DisplayName,
		"playerID", player.ID,
		"roomCode", roomCode,
	)

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
		// Use the connection context without additional timeout
		// The writePump's ping will keep the connection alive
		// Railway and the websocket library handle timeouts at a lower level
		err := wsjson.Read(c.ctx, c.Conn, &msg)

		if err != nil {
			if c.ctx.Err() != nil {
				// Context cancelled, clean shutdown
				slog.Info("connection context cancelled", "playerID", c.PlayerID)
				return
			}
			// Log the specific error for debugging
			slog.Info("read error, closing connection", "playerID", c.PlayerID, "error", err)
			return
		}

		// Handle message
		cm.handleClientMessage(c, room, msg)
	}
}

// writePump sends messages to the WebSocket connection.
func (c *Connection) writePump() {
	// Use 10-second heartbeat interval for Railway compatibility
	// Send application-level heartbeat messages which proxies recognize as activity
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-c.Send:
			ctx, cancel := context.WithTimeout(c.ctx, 10*time.Second)
			err := wsjson.Write(ctx, c.Conn, msg)
			cancel()

			if err != nil {
				slog.Info("write error, closing connection", "playerID", c.PlayerID, "error", err)
				return
			}

		case <-ticker.C:
			// Send heartbeat message instead of WebSocket ping
			// Railway's proxy recognizes data frames as activity
			heartbeat, _ := NewHeartbeatMessage()
			ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
			err := wsjson.Write(ctx, c.Conn, heartbeat)
			cancel()

			if err != nil {
				slog.Info("heartbeat failed, closing connection", "playerID", c.PlayerID, "error", err)
				return
			}

		case <-c.ctx.Done():
			slog.Info("write pump context done", "playerID", c.PlayerID)
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
	// Gracefully close the WebSocket connection
	conn.Conn.Close(websocket.StatusNormalClosure, "")

	cm.mu.Lock()
	delete(cm.connections, conn.PlayerID)
	cm.mu.Unlock()

	// Mark player as disconnected
	player, _ := room.GetPlayer(conn.PlayerID)
	if player != nil {
		player.Disconnect()
		slog.Info("player disconnected",
			"playerName", player.DisplayName,
			"playerID", player.ID,
			"roomCode", room.ID,
		)
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
		slog.Error("failed to get room for broadcast", "roomCode", roomCode, "error", err)
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
			slog.Warn("failed to send event", "playerID", player.ID, "reason", "channel full")
		}
	}
}

// BroadcastRoomState sends updated room state to all connected players.
func (cm *ConnectionManager) BroadcastRoomState(roomCode string) {
	room, err := cm.store.GetRoom(roomCode)
	if err != nil {
		slog.Error("failed to get room for state broadcast", "roomCode", roomCode, "error", err)
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
			slog.Warn("failed to send room state", "playerID", player.ID, "reason", "channel full")
		}
	}
}
