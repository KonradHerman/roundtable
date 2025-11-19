package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"nhooyr.io/websocket"

	"github.com/yourusername/roundtable/internal/core"
	"github.com/yourusername/roundtable/internal/games"
	"github.com/yourusername/roundtable/internal/store"
	"github.com/yourusername/roundtable/internal/util"
)

// Server holds the HTTP server and its dependencies.
type Server struct {
	store        store.Store
	connMgr      *ConnectionManager
	gameRegistry *games.Registry
}

// NewServer creates a new server instance.
func NewServer(store store.Store) *Server {
	return &Server{
		store:        store,
		connMgr:      NewConnectionManager(store),
		gameRegistry: games.NewRegistry(),
	}
}

// ConnectionManager returns the connection manager (for phase checks).
func (s *Server) ConnectionManager() *ConnectionManager {
	return s.connMgr
}

// CreateRoomRequest is the payload for creating a room.
type CreateRoomRequest struct {
	GameType    string `json:"gameType"`
	DisplayName string `json:"displayName"` // Host's display name
	MaxPlayers  int    `json:"maxPlayers,omitempty"`
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

	// Limit request body to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1*1024*1024)

	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Request too large or malformed", http.StatusBadRequest)
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
		slog.Error("failed to create room", "error", err)
		http.Error(w, "Failed to create room", http.StatusInternalServerError)
		return
	}

	// Create initial event
	event, _ := core.NewPublicEvent(core.EventPlayerJoined, "system", core.PlayerJoinedPayload{
		PlayerID:    hostPlayer.ID,
		DisplayName: hostPlayer.DisplayName,
	})
	room.AppendEvent(event)

	slog.Info("created room",
		"roomCode", roomCode,
		"gameType", req.GameType,
		"hostName", hostPlayer.DisplayName,
		"hostID", hostPlayer.ID,
	)

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

	// Limit request body to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1*1024*1024)

	// Extract room code from URL path
	// Expected format: /api/rooms/{code}/join
	roomCode := r.PathValue("code")
	if roomCode == "" {
		http.Error(w, "Room code required", http.StatusBadRequest)
		return
	}

	var req JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Request too large or malformed", http.StatusBadRequest)
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

	slog.Info("player joined room",
		"playerName", player.DisplayName,
		"playerID", player.ID,
		"roomCode", roomCode,
	)

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

	// Limit request body to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1*1024*1024)

	// Extract room code from URL path
	roomCode := r.PathValue("code")
	if roomCode == "" {
		http.Error(w, "Room code required", http.StatusBadRequest)
		return
	}

	var req StartGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Request too large or malformed", http.StatusBadRequest)
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

	// Get the current event log length before starting
	eventLogLengthBefore := room.GetEventLogLength()

	// Start game
	if err := room.StartGame(game, config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Broadcast all new events that were created during game start
	newEvents := room.GetEventsSince(eventLogLengthBefore)
	for _, event := range newEvents {
		s.connMgr.BroadcastEvent(roomCode, event)
	}

	slog.Info("game started", "roomCode", roomCode, "gameType", room.GameType)

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}

// HandleResetGame resets the room back to lobby for a new game.
func (s *Server) HandleResetGame(w http.ResponseWriter, r *http.Request) {
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

	// Get room
	room, err := s.store.GetRoom(roomCode)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	// Reset the game
	if err := room.ResetGame(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Broadcast updated room state to all players
	s.connMgr.BroadcastRoomState(roomCode)

	slog.Info("room reset for new game", "roomCode", roomCode)

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "reset"})
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

	// Get allowed origins from environment
	allowedOrigins := getWebSocketOrigins()

	// Upgrade connection with origin restrictions
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: allowedOrigins,
	})
	if err != nil {
		slog.Error("failed to upgrade WebSocket", "error", err, "remoteAddr", r.RemoteAddr)
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

// getWebSocketOrigins returns allowed WebSocket origin patterns from environment.
func getWebSocketOrigins() []string {
	// Check if we're running on Railway
	inRailway := os.Getenv("RAILWAY_ENVIRONMENT") != "" || os.Getenv("RAILWAY_PROJECT_ID") != ""

	// Start with configured origins (use ALLOWED_ORIGIN for consistency with CORS)
	originsEnv := os.Getenv("ALLOWED_ORIGIN")
	if originsEnv == "" {
		originsEnv = os.Getenv("ALLOWED_ORIGINS") // Fallback to plural for backwards compatibility
	}

	patterns := []string{}

	// Parse configured origins
	if originsEnv != "" {
		origins := strings.Split(originsEnv, ",")
		for _, origin := range origins {
			origin = strings.TrimSpace(origin)
			if origin == "" {
				continue
			}

			// If it's a wildcard, use it as-is
			if origin == "*" {
				return []string{"*"}
			}

			// Convert full URL to origin pattern
			// Remove protocol (https:// or http://)
			pattern := strings.TrimPrefix(strings.TrimPrefix(origin, "https://"), "http://")
			patterns = append(patterns, pattern)
		}
	}

	// In Railway environment, also allow all *.up.railway.app origins
	if inRailway {
		patterns = append(patterns, "*.up.railway.app")
		slog.Info("WebSocket origins configured for Railway", "patterns", patterns)
	} else {
		// Dev default - allow localhost on any port
		if len(patterns) == 0 {
			patterns = []string{"localhost:*", "127.0.0.1:*"}
		} else {
			// Add localhost to custom patterns for development
			patterns = append(patterns, "localhost:*", "127.0.0.1:*")
		}
		slog.Info("WebSocket origins configured", "patterns", patterns)
	}

	return patterns
}
