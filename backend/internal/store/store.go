package store

import (
	"errors"
	"time"

	"github.com/yourusername/roundtable/internal/core"
)

var (
	ErrRoomNotFound   = errors.New("room not found")
	ErrRoomExists     = errors.New("room already exists")
	ErrPlayerNotFound = errors.New("player not found")
)

// Store defines the interface for room persistence.
// Implementations: in-memory (MVP), Redis (production).
type Store interface {
	// CreateRoom stores a new room.
	CreateRoom(room *core.Room) error

	// GetRoom retrieves a room by its code.
	GetRoom(roomCode string) (*core.Room, error)

	// UpdateRoom updates an existing room.
	UpdateRoom(room *core.Room) error

	// DeleteRoom removes a room.
	DeleteRoom(roomCode string) error

	// ListRooms returns all active rooms (for admin/debugging).
	ListRooms() ([]*core.Room, error)

	// CleanupStaleRooms removes rooms that haven't been active recently.
	CleanupStaleRooms() error

	// UpdateRoomTimer updates the next phase time for a room.
	UpdateRoomTimer(roomID string, t time.Time) error

	// PopExpiredRooms retrieves and updates rooms that have reached their phase timeout.
	PopExpiredRooms(until time.Time) ([]*core.Room, error)
}
