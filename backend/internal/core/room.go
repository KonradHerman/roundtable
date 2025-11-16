package core

import (
	"errors"
	"sync"
	"time"
)

// RoomStatus represents the current state of a room.
type RoomStatus string

const (
	RoomStatusWaiting  RoomStatus = "waiting"  // Lobby, waiting for players
	RoomStatusPlaying  RoomStatus = "playing"  // Game in progress
	RoomStatusFinished RoomStatus = "finished" // Game concluded
)

// Room represents a game session container.
// It holds players, event log, and game-specific state.
type Room struct {
	mu sync.RWMutex // Protects concurrent access

	ID         string     `json:"id"`         // 6-character room code
	CreatedAt  time.Time  `json:"createdAt"`  // Room creation time
	Status     RoomStatus `json:"status"`     // Current status
	GameType   string     `json:"gameType"`   // "werewolf", "avalon", etc.
	MaxPlayers int        `json:"maxPlayers"` // Maximum allowed players

	HostID  string             `json:"hostId"`  // PlayerID of the host
	Players map[string]*Player `json:"players"` // PlayerID → Player

	EventLog []GameEvent `json:"eventLog"` // Append-only event history
	Game     Game        `json:"-"`        // Game-specific state machine
}

// NewRoom creates a new room with a generated code.
func NewRoom(roomCode string, gameType string, hostPlayer *Player, maxPlayers int) *Room {
	return &Room{
		ID:         roomCode,
		CreatedAt:  time.Now(),
		Status:     RoomStatusWaiting,
		GameType:   gameType,
		MaxPlayers: maxPlayers,
		HostID:     hostPlayer.ID,
		Players: map[string]*Player{
			hostPlayer.ID: hostPlayer,
		},
		EventLog: make([]GameEvent, 0),
	}
}

// AddPlayer adds a new player to the room.
func (r *Room) AddPlayer(player *Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Status != RoomStatusWaiting {
		return errors.New("cannot join: game already started")
	}

	if len(r.Players) >= r.MaxPlayers {
		return errors.New("room is full")
	}

	if _, exists := r.Players[player.ID]; exists {
		return errors.New("player already in room")
	}

	r.Players[player.ID] = player
	return nil
}

// RemovePlayer removes a player from the room.
func (r *Room) RemovePlayer(playerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.Players[playerID]; !exists {
		return errors.New("player not in room")
	}

	delete(r.Players, playerID)
	return nil
}

// GetPlayer retrieves a player by ID.
func (r *Room) GetPlayer(playerID string) (*Player, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	player, exists := r.Players[playerID]
	if !exists {
		return nil, errors.New("player not found")
	}

	return player, nil
}

// GetPlayerByToken finds a player by their session token.
func (r *Room) GetPlayerByToken(token string) (*Player, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, player := range r.Players {
		if player.SessionToken == token {
			return player, nil
		}
	}

	return nil, errors.New("invalid session token")
}

// GetPlayers returns all players as a slice.
func (r *Room) GetPlayers() []*Player {
	r.mu.RLock()
	defer r.mu.RUnlock()

	players := make([]*Player, 0, len(r.Players))
	for _, player := range r.Players {
		players = append(players, player)
	}

	return players
}

// IsHost checks if the given player is the host.
func (r *Room) IsHost(playerID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.HostID == playerID
}

// AppendEvent adds an event to the event log.
func (r *Room) AppendEvent(event GameEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.EventLog = append(r.EventLog, event)
}

// AppendEvents adds multiple events to the event log.
func (r *Room) AppendEvents(events []GameEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.EventLog = append(r.EventLog, events...)
}

// GetEventsForPlayer returns all events this player can see.
func (r *Room) GetEventsForPlayer(playerID string) []GameEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]GameEvent, 0, len(r.EventLog))
	for _, event := range r.EventLog {
		if event.CanPlayerSee(playerID) {
			filtered = append(filtered, event)
		}
	}

	return filtered
}

// GetPublicEvents returns all public events (for board view, spectators).
func (r *Room) GetPublicEvents() []GameEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]GameEvent, 0, len(r.EventLog))
	for _, event := range r.EventLog {
		if event.Visibility.Public {
			filtered = append(filtered, event)
		}
	}

	return filtered
}

// SetStatus updates the room status.
func (r *Room) SetStatus(status RoomStatus) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Status = status
}

// StartGame initializes the game and transitions to playing status.
func (r *Room) StartGame(game Game, config GameConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Status != RoomStatusWaiting {
		return errors.New("game already started")
	}

	// Initialize game
	players := make([]*Player, 0, len(r.Players))
	for _, player := range r.Players {
		players = append(players, player)
	}

	events, err := game.Initialize(config, players)
	if err != nil {
		return err
	}

	// If the game supports host tracking, set the host
	type HostSetter interface {
		SetHost(hostID string)
	}
	if hs, ok := game.(HostSetter); ok {
		hs.SetHost(r.HostID)
	}

	r.Game = game
	r.Status = RoomStatusPlaying
	r.EventLog = append(r.EventLog, events...)

	return nil
}

// ResetGame resets the room back to waiting status for a new game.
// Keeps players but clears game state and event log.
func (r *Room) ResetGame() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Status == RoomStatusWaiting {
		return errors.New("no game to reset")
	}

	// Clear game state
	r.Game = nil
	r.EventLog = make([]GameEvent, 0)
	r.Status = RoomStatusWaiting

	return nil
}

// ProcessAction validates and processes a player action.
func (r *Room) ProcessAction(playerID string, action Action) ([]GameEvent, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Status != RoomStatusPlaying {
		return nil, errors.New("no game in progress")
	}

	if r.Game == nil {
		return nil, errors.New("game not initialized")
	}

	// Validate action
	if err := r.Game.ValidateAction(playerID, action); err != nil {
		return nil, err
	}

	// Process action
	events, err := r.Game.ProcessAction(playerID, action)
	if err != nil {
		return nil, err
	}

	// Append events to log
	r.EventLog = append(r.EventLog, events...)

	// Check if game finished
	if r.Game.IsFinished() {
		r.Status = RoomStatusFinished
	}

	return events, nil
}

// RoomState is a snapshot of room state for client consumption.
type RoomState struct {
	ID         string     `json:"id"`
	Status     RoomStatus `json:"status"`
	GameType   string     `json:"gameType"`
	MaxPlayers int        `json:"maxPlayers"`
	HostID     string     `json:"hostId"`
	Players    []*Player  `json:"players"`
}

// GetState returns a snapshot of the room state.
func (r *Room) GetState() RoomState {
	r.mu.RLock()
	defer r.mu.RUnlock()

	players := make([]*Player, 0, len(r.Players))
	for _, player := range r.Players {
		// Don't send session tokens to clients
		playerCopy := *player
		playerCopy.SessionToken = ""
		players = append(players, &playerCopy)
	}

	return RoomState{
		ID:         r.ID,
		Status:     r.Status,
		GameType:   r.GameType,
		MaxPlayers: r.MaxPlayers,
		HostID:     r.HostID,
		Players:    players,
	}
}
