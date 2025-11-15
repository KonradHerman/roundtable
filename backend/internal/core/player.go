package core

import (
	"time"

	"github.com/google/uuid"
)

// Player represents an anonymous participant in a room.
// Players are identified by session tokens (no signup required).
type Player struct {
	ID           string    `json:"id"`           // UUID
	SessionToken string    `json:"-"`            // For reconnection (never sent to clients)
	DisplayName  string    `json:"displayName"`  // User-chosen name
	Connected    bool      `json:"connected"`    // Current connection status
	JoinedAt     time.Time `json:"joinedAt"`     // When they joined
	LastSeenAt   time.Time `json:"lastSeenAt"`   // Last activity timestamp
}

// NewPlayer creates a new player with generated ID and session token.
func NewPlayer(displayName string) *Player {
	return &Player{
		ID:           uuid.New().String(),
		SessionToken: generateSessionToken(),
		DisplayName:  displayName,
		Connected:    true,
		JoinedAt:     time.Now(),
		LastSeenAt:   time.Now(),
	}
}

// UpdateLastSeen marks the player as active now.
func (p *Player) UpdateLastSeen() {
	p.LastSeenAt = time.Now()
}

// Disconnect marks the player as disconnected.
func (p *Player) Disconnect() {
	p.Connected = false
	p.UpdateLastSeen()
}

// Reconnect marks the player as connected again.
func (p *Player) Reconnect() {
	p.Connected = true
	p.UpdateLastSeen()
}

// IsStale checks if the player hasn't been seen recently (for cleanup).
func (p *Player) IsStale(timeout time.Duration) bool {
	return time.Since(p.LastSeenAt) > timeout
}

// generateSessionToken creates a random token for player sessions.
// In production, use a cryptographically secure token generator.
func generateSessionToken() string {
	return uuid.New().String()
}
