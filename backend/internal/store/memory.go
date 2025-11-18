package store

import (
	"sync"
	"time"

	"github.com/yourusername/roundtable/internal/core"
)

// MemoryStore is an in-memory implementation of Store.
// Suitable for MVP and single-instance deployments.
// For production multi-instance, use RedisStore.
type MemoryStore struct {
	mu    sync.RWMutex
	rooms map[string]*core.Room // roomCode → Room
}

// NewMemoryStore creates a new in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		rooms: make(map[string]*core.Room),
	}
}

// CreateRoom stores a new room.
func (s *MemoryStore) CreateRoom(room *core.Room) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.rooms[room.ID]; exists {
		return ErrRoomExists
	}

	s.rooms[room.ID] = room
	return nil
}

// GetRoom retrieves a room by code.
func (s *MemoryStore) GetRoom(roomCode string) (*core.Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, exists := s.rooms[roomCode]
	if !exists {
		return nil, ErrRoomNotFound
	}

	return room, nil
}

// UpdateRoom updates an existing room.
// In the in-memory implementation, rooms are pointers so this is a no-op.
// This method exists for interface compatibility with Redis implementation.
func (s *MemoryStore) UpdateRoom(room *core.Room) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.rooms[room.ID]; !exists {
		return ErrRoomNotFound
	}

	// Room is already updated since we store pointers
	return nil
}

// DeleteRoom removes a room.
func (s *MemoryStore) DeleteRoom(roomCode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.rooms[roomCode]; !exists {
		return ErrRoomNotFound
	}

	delete(s.rooms, roomCode)
	return nil
}

// ListRooms returns all active rooms.
func (s *MemoryStore) ListRooms() ([]*core.Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rooms := make([]*core.Room, 0, len(s.rooms))
	for _, room := range s.rooms {
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// CleanupStaleRooms removes rooms that are finished or have no active players.
func (s *MemoryStore) CleanupStaleRooms() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	const staleTimeout = 24 * time.Hour

	toDelete := make([]string, 0)

	for roomCode, room := range s.rooms {
		// Get room info safely with internal locking
		status, createdAt, anyConnected := room.GetCleanupInfo()
		
		// Delete finished rooms older than 1 hour
		if status == core.RoomStatusFinished && time.Since(createdAt) > 1*time.Hour {
			toDelete = append(toDelete, roomCode)
			continue
		}

		// Delete rooms with no connected players older than staleTimeout
		if !anyConnected && time.Since(createdAt) > staleTimeout {
			toDelete = append(toDelete, roomCode)
		}
	}

	for _, roomCode := range toDelete {
		delete(s.rooms, roomCode)
	}

	return nil
}
