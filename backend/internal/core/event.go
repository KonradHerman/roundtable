package core

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// GameEvent represents an immutable fact about something that happened.
// Events are the source of truth in the event sourcing architecture.
// Current state = Initial state + Sequence of events.
type GameEvent struct {
	ID        string          `json:"id"`        // Unique event ID
	Timestamp time.Time       `json:"timestamp"` // When it happened
	Type      string          `json:"type"`      // Event type (e.g., "player_joined", "vote_cast")
	ActorID   string          `json:"actorId"`   // PlayerID who triggered it ("system" for server events)
	Payload   json.RawMessage `json:"payload"`   // Event-specific data
	Visibility EventVisibility `json:"-"`        // Who can see this event (not sent to clients)
}

// EventVisibility controls which clients receive an event.
type EventVisibility struct {
	Public      bool     // All players and spectators see it
	PlayerIDs   []string // Specific players who see it (for private info)
	SpectatorOK bool     // Spectators can see it
}

// NewEvent creates a new event with auto-generated ID and timestamp.
func NewEvent(eventType string, actorID string, payload interface{}, visibility EventVisibility) (GameEvent, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return GameEvent{}, err
	}

	return GameEvent{
		ID:         uuid.New().String(),
		Timestamp:  time.Now(),
		Type:       eventType,
		ActorID:    actorID,
		Payload:    payloadBytes,
		Visibility: visibility,
	}, nil
}

// NewPublicEvent creates an event visible to all players.
func NewPublicEvent(eventType string, actorID string, payload interface{}) (GameEvent, error) {
	return NewEvent(eventType, actorID, payload, EventVisibility{
		Public:      true,
		SpectatorOK: true,
	})
}

// NewPrivateEvent creates an event visible only to specific players.
func NewPrivateEvent(eventType string, actorID string, payload interface{}, visibleTo []string) (GameEvent, error) {
	return NewEvent(eventType, actorID, payload, EventVisibility{
		Public:      false,
		PlayerIDs:   visibleTo,
		SpectatorOK: false,
	})
}

// CanPlayerSee determines if a player should receive this event.
func (e *GameEvent) CanPlayerSee(playerID string) bool {
	if e.Visibility.Public {
		return true
	}

	for _, id := range e.Visibility.PlayerIDs {
		if id == playerID {
			return true
		}
	}

	return false
}

// Common event types (games can define their own in addition to these)
const (
	EventPlayerJoined   = "player_joined"
	EventPlayerLeft     = "player_left"
	EventPlayerReconnected = "player_reconnected"
	EventGameStarted    = "game_started"
	EventGameFinished   = "game_finished"
	EventPhaseChanged   = "phase_changed"
	EventConfigUpdated  = "config_updated"
	EventError          = "error"
)

// Common event payloads
type PlayerJoinedPayload struct {
	PlayerID    string `json:"playerId"`
	DisplayName string `json:"displayName"`
}

type PlayerLeftPayload struct {
	PlayerID string `json:"playerId"`
}

type PlayerReconnectedPayload struct {
	PlayerID string `json:"playerId"`
}

type GameStartedPayload struct {
	GameType  string      `json:"gameType"`
	Config    interface{} `json:"config"`
	PlayerIDs []string    `json:"playerIds"`
}

type GameFinishedPayload struct {
	Results GameResults `json:"results"`
}

type PhaseChangedPayload struct {
	Phase   GamePhase `json:"phase"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}
