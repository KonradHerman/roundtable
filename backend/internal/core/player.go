package core

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// Player represents an anonymous participant in a room.
// Players are identified by session tokens (no signup required).
type Player struct {
	mu sync.RWMutex // Protects Connected and LastSeenAt fields

	ID           string    `json:"id"`           // UUID
	SessionToken string    `json:"-"`            // For reconnection (never sent to clients)
	DisplayName  string    `json:"displayName"`  // User-chosen name
	Connected    bool      `json:"connected"`    // Current connection status (protected by mu)
	JoinedAt     time.Time `json:"joinedAt"`     // When they joined
	LastSeenAt   time.Time `json:"lastSeenAt"`   // Last activity timestamp (protected by mu)
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
	p.mu.Lock()
	defer p.mu.Unlock()
	p.LastSeenAt = time.Now()
}

// Disconnect marks the player as disconnected.
func (p *Player) Disconnect() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Connected = false
	p.LastSeenAt = time.Now()
}

// Reconnect marks the player as connected again.
func (p *Player) Reconnect() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Connected = true
	p.LastSeenAt = time.Now()
}

// IsConnected safely checks if the player is connected.
func (p *Player) IsConnected() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Connected
}

// GetLastSeenAt safely returns the last seen timestamp.
func (p *Player) GetLastSeenAt() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.LastSeenAt
}

// IsStale checks if the player hasn't been seen recently (for cleanup).
func (p *Player) IsStale(timeout time.Duration) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return time.Since(p.LastSeenAt) > timeout
}

// generateSessionToken creates a random token for player sessions.
// In production, use a cryptographically secure token generator.
func generateSessionToken() string {
	return uuid.New().String()
}
