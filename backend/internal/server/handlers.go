package server

import (
	"encoding/json"
	"log"
	"net/http"

	"nhooyr.io/websocket"

	"github.com/yourusername/roundtable/internal/core"
	"github.com/yourusername/roundtable/internal/games"
	"github.com/yourusername/roundtable/internal/store"
	"github.com/yourusername/roundtable/internal/util"
)

// Server holds the HTTP server and its dependencies.
type Server struct {
	store      store.Store
	connMgr    *ConnectionManager
	gameRegistry *games.Registry
}

// NewServer creates a new server instance.
func NewServer(store store.Store) *Server {
	return &Server{
		store:      store,
		connMgr:    NewConnectionManager(store),
		gameRegistry: games.NewRegistry(),
	}
}

// CreateRoomRequest is the payload for creating a room.
type CreateRoomRequest struct {
	GameType   string `json:"gameType"`
	DisplayName string `json:"displayName"` // Host's display name
	MaxPlayers int    `json:"maxPlayers,omitempty"`
}

// CreateRoomResponse is the response for creating a room.
type CreateRoomResponse struct {
	RoomCode     string `json:"roomCode"`
	SessionToken string `json:"sessionToken"`
	PlayerID     string `json:"playerId"`
}

// HandleCreateRoom creates a new game room.
func (s *Server) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate game type
	if !s.gameRegistry.IsRegistered(req.GameType) {
		http.Error(w, "Unknown game type", http.StatusBadRequest)
		return
	}

	// Default max players
	if req.MaxPlayers == 0 {
		req.MaxPlayers = 10
	}

	// Create host player
	hostPlayer := core.NewPlayer(req.DisplayName)

	// Generate room code
	roomCode := util.GenerateRoomCode()

	// Create room
	room := core.NewRoom(roomCode, req.GameType, hostPlayer, req.MaxPlayers)

	// Store room
	if err := s.store.CreateRoom(room); err != nil {
		log.Printf("Failed to create room: %v", err)
		http.Error(w, "Failed to create room", http.StatusInternalServerError)
		return
	}

	// Create initial event
	event, _ := core.NewPublicEvent(core.EventPlayerJoined, "system", core.PlayerJoinedPayload{
		PlayerID:    hostPlayer.ID,
		DisplayName: hostPlayer.DisplayName,
	})
	room.AppendEvent(event)

	log.Printf("Created room %s for game %s (host: %s)", roomCode, req.GameType, hostPlayer.DisplayName)

	// Return response
	resp := CreateRoomResponse{
		RoomCode:     roomCode,
		SessionToken: hostPlayer.SessionToken,
		PlayerID:     hostPlayer.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// JoinRoomRequest is the payload for joining a room.
type JoinRoomRequest struct {
	DisplayName string `json:"displayName"`
}

// JoinRoomResponse is the response for joining a room.
type JoinRoomResponse struct {
	SessionToken string `json:"sessionToken"`
	PlayerID     string `json:"playerId"`
	RoomCode     string `json:"roomCode"`
}

// HandleJoinRoom adds a player to an existing room.
func (s *Server) HandleJoinRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract room code from URL path
	// Expected format: /api/rooms/{code}/join
	roomCode := r.PathValue("code")
	if roomCode == "" {
		http.Error(w, "Room code required", http.StatusBadRequest)
		return
	}

	var req JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.DisplayName == "" {
		http.Error(w, "Display name required", http.StatusBadRequest)
		return
	}

	// Get room
	room, err := s.store.GetRoom(roomCode)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	// Create player
	player := core.NewPlayer(req.DisplayName)

	// Add player to room
	if err := room.AddPlayer(player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create player joined event
	event, _ := core.NewPublicEvent(core.EventPlayerJoined, "system", core.PlayerJoinedPayload{
		PlayerID:    player.ID,
		DisplayName: player.DisplayName,
	})
	room.AppendEvent(event)

	// Broadcast event to connected players
	s.connMgr.BroadcastEvent(roomCode, event)
	s.connMgr.BroadcastRoomState(roomCode)

	log.Printf("Player %s joined room %s", player.DisplayName, roomCode)

	// Return response
	resp := JoinRoomResponse{
		SessionToken: player.SessionToken,
		PlayerID:     player.ID,
		RoomCode:     roomCode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// StartGameRequest is the payload for starting a game.
type StartGameRequest struct {
	Config json.RawMessage `json:"config"` // Game-specific config
}

// HandleStartGame initializes and starts the game.
func (s *Server) HandleStartGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract room code from URL path
	roomCode := r.PathValue("code")
	if roomCode == "" {
		http.Error(w, "Room code required", http.StatusBadRequest)
		return
	}

	// TODO: Extract player ID from auth header/token
	// For now, we'll validate via request body or session
	// This is simplified for MVP

	var req StartGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get room
	room, err := s.store.GetRoom(roomCode)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	// Create game instance
	game, err := s.gameRegistry.CreateGame(room.GameType)
	if err != nil {
		http.Error(w, "Failed to create game", http.StatusInternalServerError)
		return
	}

	// Parse game config
	config, err := s.gameRegistry.ParseConfig(room.GameType, req.Config)
	if err != nil {
		http.Error(w, "Invalid game configuration", http.StatusBadRequest)
		return
	}

	// Start game
	if err := room.StartGame(game, config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Broadcast events to all players
	events := room.EventLog[len(room.EventLog)-len(room.GetPlayers()):] // Get recent events
	for _, event := range events {
		s.connMgr.BroadcastEvent(roomCode, event)
	}

	log.Printf("Game started in room %s", roomCode)

	w.WriteHeader(http.StatusOK)
}

// HandleWebSocket upgrades HTTP connection to WebSocket.
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Extract room code from URL path
	roomCode := r.PathValue("code")
	if roomCode == "" {
		http.Error(w, "Room code required", http.StatusBadRequest)
		return
	}

	// Verify room exists
	if _, err := s.store.GetRoom(roomCode); err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	// Upgrade connection
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"}, // TODO: Restrict in production
	})
	if err != nil {
		log.Printf("Failed to upgrade WebSocket: %v", err)
		return
	}

	// Handle connection
	s.connMgr.HandleConnection(r.Context(), conn, roomCode)
}

// HandleGetRoom returns room information.
func (s *Server) HandleGetRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomCode := r.PathValue("code")
	if roomCode == "" {
		http.Error(w, "Room code required", http.StatusBadRequest)
		return
	}

	room, err := s.store.GetRoom(roomCode)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room.GetState())
}
