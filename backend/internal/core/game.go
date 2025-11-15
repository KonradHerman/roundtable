package core

import (
	"encoding/json"
	"time"
)

// Game represents the interface that all game implementations must satisfy.
// This abstraction allows the platform to support multiple game types
// without coupling to specific game logic.
type Game interface {
	// Initialize sets up the game with configuration and players.
	// Returns initial events (GameStarted, RoleAssigned, etc.)
	Initialize(config GameConfig, players []*Player) ([]GameEvent, error)

	// ValidateAction checks if a player can perform the given action
	// in the current game state (phase, timing, permissions).
	ValidateAction(playerID string, action Action) error

	// ProcessAction executes a valid action and returns resulting events.
	// Events are appended to the room's event log and broadcast to clients.
	ProcessAction(playerID string, action Action) ([]GameEvent, error)

	// GetPlayerState returns the game state visible to this specific player.
	// Server filters information based on what the player should see
	// (e.g., only their role, not other players' hidden info).
	GetPlayerState(playerID string) PlayerState

	// GetPublicState returns the game state visible to all players and spectators.
	// This is used for the board view and shows only public information
	// (e.g., phase, timer, public votes).
	GetPublicState() PublicState

	// GetPhase returns the current game phase for UI rendering.
	GetPhase() GamePhase

	// IsFinished returns true if the game has concluded.
	IsFinished() bool

	// GetResults returns the final game results (only valid if IsFinished).
	GetResults() GameResults

	// CheckPhaseTimeout checks if the current phase has expired and should advance.
	// Returns events if phase advances, nil if no change needed.
	CheckPhaseTimeout() ([]GameEvent, error)
}

// GameConfig is a marker interface for game-specific configuration.
// Each game implementation provides its own config type.
type GameConfig interface {
	GameType() string
	Validate() error
}

// Action represents a player's intent to do something in the game.
type Action struct {
	Type    string          `json:"type"`    // "vote", "select_target", etc.
	Payload json.RawMessage `json:"payload"` // Action-specific data
}

// PlayerState is game-specific state for a single player (filtered view).
type PlayerState interface{}

// PublicState is game-specific public state (board view, spectators).
type PublicState interface{}

// GamePhase represents the current phase of the game.
type GamePhase struct {
	Name    string     `json:"name"`    // "night", "day", "voting", etc.
	EndsAt  *time.Time `json:"endsAt"`  // Optional phase timer
	Message string     `json:"message"` // Display message for players
}

// GameResults contains the final outcome of the game.
type GameResults struct {
	Winners    []string               `json:"winners"`    // PlayerIDs of winners
	WinReason  string                 `json:"winReason"`  // "Villagers eliminated werewolf"
	FinalState map[string]interface{} `json:"finalState"` // Game-specific final state
}
