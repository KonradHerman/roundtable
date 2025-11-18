package store

import (
	"container/heap"
	"sync"
	"time"

	"github.com/yourusername/roundtable/internal/core"
)

// MemoryStore is an in-memory implementation of Store.
// Suitable for MVP and single-instance deployments.
// For production multi-instance, use RedisStore.
type MemoryStore struct {
	mu        sync.RWMutex
	rooms     map[string]*core.Room // roomCode → Room
	roomItems map[string]*roomItem  // roomCode → heap item
	pq        roomHeap              // priority queue
}

// NewMemoryStore creates a new in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		rooms:     make(map[string]*core.Room),
		roomItems: make(map[string]*roomItem),
		pq:        make(roomHeap, 0),
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

	// Remove from heap if present
	if item, exists := s.roomItems[roomCode]; exists {
		heap.Remove(&s.pq, item.index)
		delete(s.roomItems, roomCode)
	}

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
		if item, exists := s.roomItems[roomCode]; exists {
			heap.Remove(&s.pq, item.index)
			delete(s.roomItems, roomCode)
		}
	}

	return nil
}

// UpdateRoomTimer updates the next phase time for a room.
func (s *MemoryStore) UpdateRoomTimer(roomID string, t time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, exists := s.rooms[roomID]
	if !exists {
		return ErrRoomNotFound
	}

	room.NextPhaseTime = t

	if t.IsZero() {
		// Remove from heap if present
		if item, exists := s.roomItems[roomID]; exists {
			heap.Remove(&s.pq, item.index)
			delete(s.roomItems, roomID)
		}
	} else {
		// Update or Add
		if item, exists := s.roomItems[roomID]; exists {
			heap.Fix(&s.pq, item.index)
		} else {
			item := &roomItem{room: room}
			heap.Push(&s.pq, item)
			s.roomItems[roomID] = item
		}
	}

	return nil
}

// PopExpiredRooms retrieves and updates rooms that have reached their phase timeout.
func (s *MemoryStore) PopExpiredRooms(until time.Time) ([]*core.Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expired []*core.Room

	for s.pq.Len() > 0 {
		item := s.pq[0]
		if item.room.NextPhaseTime.After(until) {
			break
		}

		// Remove from heap
		heap.Pop(&s.pq)
		delete(s.roomItems, item.room.ID)

		// Clear the timer so it doesn't get picked up again immediately
		// The caller is expected to set a new timer if needed
		item.room.NextPhaseTime = time.Time{}

		expired = append(expired, item.room)
	}

	return expired, nil
}

// Heap implementation

type roomItem struct {
	room  *core.Room
	index int
}

type roomHeap []*roomItem

func (h roomHeap) Len() int { return len(h) }
func (h roomHeap) Less(i, j int) bool {
	return h[i].room.NextPhaseTime.Before(h[j].room.NextPhaseTime)
}
func (h roomHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *roomHeap) Push(x any) {
	n := len(*h)
	item := x.(*roomItem)
	item.index = n
	*h = append(*h, item)
}

func (h *roomHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*h = old[0 : n-1]
	return item
}
