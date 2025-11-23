package store

import (
	"sync"
	"testing"
	"time"

	"github.com/KonradHerman/roundtable/internal/core"
)

func TestMemoryStore_CreateRoom(t *testing.T) {
	t.Parallel()

	t.Run("successfully create room", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()
		player := &core.Player{ID: "p1", DisplayName: "Player1"}
		room := core.NewRoom("ABC123", "werewolf", player, 10)

		err := store.CreateRoom(room)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify room can be retrieved
		retrieved, err := store.GetRoom("ABC123")
		if err != nil {
			t.Fatalf("failed to retrieve room: %v", err)
		}

		if retrieved.ID != room.ID {
			t.Errorf("retrieved room ID = %s, want %s", retrieved.ID, room.ID)
		}
	})

	t.Run("fail on duplicate room", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()
		player := &core.Player{ID: "p1", DisplayName: "Player1"}
		room := core.NewRoom("ABC123", "werewolf", player, 10)

		err := store.CreateRoom(room)
		if err != nil {
			t.Fatalf("first create failed: %v", err)
		}

		// Try to create again
		err = store.CreateRoom(room)
		if err != ErrRoomExists {
			t.Errorf("expected ErrRoomExists, got %v", err)
		}
	})
}

func TestMemoryStore_GetRoom(t *testing.T) {
	t.Parallel()

	t.Run("successfully get existing room", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()
		player := &core.Player{ID: "p1", DisplayName: "Player1"}
		room := core.NewRoom("ABC123", "werewolf", player, 10)

		err := store.CreateRoom(room)
		if err != nil {
			t.Fatalf("failed to create room: %v", err)
		}

		retrieved, err := store.GetRoom("ABC123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if retrieved.ID != "ABC123" {
			t.Errorf("room ID = %s, want ABC123", retrieved.ID)
		}

		if retrieved.GameType != "werewolf" {
			t.Errorf("game type = %s, want werewolf", retrieved.GameType)
		}
	})

	t.Run("fail on non-existent room", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()

		_, err := store.GetRoom("NONEXIST")
		if err != ErrRoomNotFound {
			t.Errorf("expected ErrRoomNotFound, got %v", err)
		}
	})
}

func TestMemoryStore_UpdateRoom(t *testing.T) {
	t.Parallel()

	t.Run("successfully update existing room", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()
		player := &core.Player{ID: "p1", DisplayName: "Player1"}
		room := core.NewRoom("ABC123", "werewolf", player, 10)

		err := store.CreateRoom(room)
		if err != nil {
			t.Fatalf("failed to create room: %v", err)
		}

		// Update room status
		room.SetStatus(core.RoomStatusPlaying)

		err = store.UpdateRoom(room)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify update (in memory store, rooms are pointers so this is automatic)
		retrieved, err := store.GetRoom("ABC123")
		if err != nil {
			t.Fatalf("failed to retrieve room: %v", err)
		}

		// Note: Since we use pointers, the room should already be updated
		if retrieved.Status != core.RoomStatusPlaying {
			t.Errorf("room status = %s, want playing", retrieved.Status)
		}
	})

	t.Run("fail on non-existent room", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()
		player := &core.Player{ID: "p1", DisplayName: "Player1"}
		room := core.NewRoom("NONEXIST", "werewolf", player, 10)

		err := store.UpdateRoom(room)
		if err != ErrRoomNotFound {
			t.Errorf("expected ErrRoomNotFound, got %v", err)
		}
	})
}

func TestMemoryStore_DeleteRoom(t *testing.T) {
	t.Parallel()

	t.Run("successfully delete existing room", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()
		player := &core.Player{ID: "p1", DisplayName: "Player1"}
		room := core.NewRoom("ABC123", "werewolf", player, 10)

		err := store.CreateRoom(room)
		if err != nil {
			t.Fatalf("failed to create room: %v", err)
		}

		err = store.DeleteRoom("ABC123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify room is gone
		_, err = store.GetRoom("ABC123")
		if err != ErrRoomNotFound {
			t.Errorf("expected ErrRoomNotFound after delete, got %v", err)
		}
	})

	t.Run("fail on non-existent room", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()

		err := store.DeleteRoom("NONEXIST")
		if err != ErrRoomNotFound {
			t.Errorf("expected ErrRoomNotFound, got %v", err)
		}
	})
}

func TestMemoryStore_ListRooms(t *testing.T) {
	t.Parallel()

	t.Run("list empty store", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()

		rooms, err := store.ListRooms()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(rooms) != 0 {
			t.Errorf("expected 0 rooms, got %d", len(rooms))
		}
	})

	t.Run("list multiple rooms", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()

		// Create multiple rooms
		player1 := &core.Player{ID: "p1", DisplayName: "Player1"}
		room1 := core.NewRoom("ABC123", "werewolf", player1, 10)

		player2 := &core.Player{ID: "p2", DisplayName: "Player2"}
		room2 := core.NewRoom("DEF456", "avalon", player2, 10)

		player3 := &core.Player{ID: "p3", DisplayName: "Player3"}
		room3 := core.NewRoom("GHI789", "werewolf", player3, 10)

		if err := store.CreateRoom(room1); err != nil {
			t.Fatalf("failed to create room1: %v", err)
		}
		if err := store.CreateRoom(room2); err != nil {
			t.Fatalf("failed to create room2: %v", err)
		}
		if err := store.CreateRoom(room3); err != nil {
			t.Fatalf("failed to create room3: %v", err)
		}

		rooms, err := store.ListRooms()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(rooms) != 3 {
			t.Errorf("expected 3 rooms, got %d", len(rooms))
		}

		// Verify all rooms are present
		roomIDs := make(map[string]bool)
		for _, room := range rooms {
			roomIDs[room.ID] = true
		}

		expectedIDs := []string{"ABC123", "DEF456", "GHI789"}
		for _, id := range expectedIDs {
			if !roomIDs[id] {
				t.Errorf("room %s not found in list", id)
			}
		}
	})
}

func TestMemoryStore_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	t.Run("concurrent creates", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()
		var wg sync.WaitGroup
		numGoroutines := 10

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				roomCode := string(rune('A' + idx))
				player := &core.Player{ID: "p" + string(rune('0'+idx)), DisplayName: "Player" + string(rune('0'+idx))}
				room := core.NewRoom(roomCode, "werewolf", player, 10)

				err := store.CreateRoom(room)
				if err != nil {
					t.Errorf("goroutine %d: failed to create room: %v", idx, err)
				}
			}(i)
		}

		wg.Wait()

		// Verify all rooms were created
		rooms, err := store.ListRooms()
		if err != nil {
			t.Fatalf("failed to list rooms: %v", err)
		}

		if len(rooms) != numGoroutines {
			t.Errorf("expected %d rooms, got %d", numGoroutines, len(rooms))
		}
	})

	t.Run("concurrent reads and writes", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()
		player := &core.Player{ID: "p1", DisplayName: "Player1"}
		room := core.NewRoom("ABC123", "werewolf", player, 10)

		err := store.CreateRoom(room)
		if err != nil {
			t.Fatalf("failed to create room: %v", err)
		}

		var wg sync.WaitGroup
		numReaders := 20
		numWriters := 5

		// Start readers
		for i := 0; i < numReaders; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for j := 0; j < 100; j++ {
					_, err := store.GetRoom("ABC123")
					if err != nil {
						t.Errorf("reader: failed to get room: %v", err)
					}
				}
			}()
		}

		// Start writers
		for i := 0; i < numWriters; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for j := 0; j < 100; j++ {
					room.SetStatus(core.RoomStatusPlaying)
					err := store.UpdateRoom(room)
					if err != nil {
						t.Errorf("writer: failed to update room: %v", err)
					}
				}
			}()
		}

		wg.Wait()
	})

	t.Run("concurrent create and get", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()
		var wg sync.WaitGroup

		// Create rooms concurrently
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				roomCode := "ROOM" + string(rune('0'+idx))
				player := &core.Player{ID: "p" + string(rune('0'+idx)), DisplayName: "Player" + string(rune('0'+idx))}
				room := core.NewRoom(roomCode, "werewolf", player, 10)

				err := store.CreateRoom(room)
				if err != nil {
					t.Errorf("failed to create room %s: %v", roomCode, err)
				}
			}(i)
		}

		// Read rooms concurrently
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				roomCode := "ROOM" + string(rune('0'+idx))

				// Keep trying until room is created
				for attempts := 0; attempts < 100; attempts++ {
					_, err := store.GetRoom(roomCode)
					if err == nil {
						return
					}
					time.Sleep(1 * time.Millisecond)
				}
			}(i)
		}

		wg.Wait()
	})
}

func TestMemoryStore_CleanupStaleRooms(t *testing.T) {
	t.Parallel()

	t.Run("cleanup finished rooms older than 1 hour", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()

		// Create a finished room with old timestamp
		player := &core.Player{ID: "p1", DisplayName: "Player1"}
		room := core.NewRoom("OLD123", "werewolf", player, 10)
		room.SetStatus(core.RoomStatusFinished)

		// Manually set old creation time
		oldTime := time.Now().Add(-2 * time.Hour)
		room.CreatedAt = oldTime

		err := store.CreateRoom(room)
		if err != nil {
			t.Fatalf("failed to create room: %v", err)
		}

		// Run cleanup
		err = store.CleanupStaleRooms()
		if err != nil {
			t.Fatalf("cleanup failed: %v", err)
		}

		// Room should be deleted
		_, err = store.GetRoom("OLD123")
		if err != ErrRoomNotFound {
			t.Errorf("expected room to be cleaned up, got error: %v", err)
		}
	})

	t.Run("preserve finished rooms newer than 1 hour", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()

		// Create a finished room with recent timestamp
		player := &core.Player{ID: "p1", DisplayName: "Player1"}
		room := core.NewRoom("NEW123", "werewolf", player, 10)
		room.SetStatus(core.RoomStatusFinished)

		err := store.CreateRoom(room)
		if err != nil {
			t.Fatalf("failed to create room: %v", err)
		}

		// Run cleanup
		err = store.CleanupStaleRooms()
		if err != nil {
			t.Fatalf("cleanup failed: %v", err)
		}

		// Room should still exist
		_, err = store.GetRoom("NEW123")
		if err != nil {
			t.Errorf("expected room to be preserved, got error: %v", err)
		}
	})

	t.Run("cleanup rooms with no connected players after 24 hours", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()

		// Create a room with old timestamp and no connected players
		player := &core.Player{ID: "p1", DisplayName: "Player1"}
		room := core.NewRoom("STALE1", "werewolf", player, 10)

		// Manually set old creation time
		oldTime := time.Now().Add(-25 * time.Hour)
		room.CreatedAt = oldTime

		err := store.CreateRoom(room)
		if err != nil {
			t.Fatalf("failed to create room: %v", err)
		}

		// Run cleanup
		err = store.CleanupStaleRooms()
		if err != nil {
			t.Fatalf("cleanup failed: %v", err)
		}

		// Room should be deleted
		_, err = store.GetRoom("STALE1")
		if err != ErrRoomNotFound {
			t.Errorf("expected stale room to be cleaned up, got error: %v", err)
		}
	})

	t.Run("preserve active rooms", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()

		// Create an active room
		player := &core.Player{ID: "p1", DisplayName: "Player1"}
		room := core.NewRoom("ACTIVE", "werewolf", player, 10)
		room.SetStatus(core.RoomStatusPlaying)

		err := store.CreateRoom(room)
		if err != nil {
			t.Fatalf("failed to create room: %v", err)
		}

		// Run cleanup
		err = store.CleanupStaleRooms()
		if err != nil {
			t.Fatalf("cleanup failed: %v", err)
		}

		// Active room should still exist
		_, err = store.GetRoom("ACTIVE")
		if err != nil {
			t.Errorf("expected active room to be preserved, got error: %v", err)
		}
	})

	t.Run("cleanup multiple rooms selectively", func(t *testing.T) {
		t.Parallel()

		store := NewMemoryStore()

		// Create mix of rooms
		// 1. Old finished room (should be deleted)
		player1 := &core.Player{ID: "p1", DisplayName: "Player1"}
		room1 := core.NewRoom("OLD_FIN", "werewolf", player1, 10)
		room1.SetStatus(core.RoomStatusFinished)
		room1.CreatedAt = time.Now().Add(-2 * time.Hour)
		store.CreateRoom(room1)

		// 2. Recent finished room (should be preserved)
		player2 := &core.Player{ID: "p2", DisplayName: "Player2"}
		room2 := core.NewRoom("NEW_FIN", "werewolf", player2, 10)
		room2.SetStatus(core.RoomStatusFinished)
		store.CreateRoom(room2)

		// 3. Active room (should be preserved)
		player3 := &core.Player{ID: "p3", DisplayName: "Player3"}
		room3 := core.NewRoom("ACTIVE", "werewolf", player3, 10)
		room3.SetStatus(core.RoomStatusPlaying)
		store.CreateRoom(room3)

		// 4. Old stale room (should be deleted)
		player4 := &core.Player{ID: "p4", DisplayName: "Player4"}
		room4 := core.NewRoom("STALE", "werewolf", player4, 10)
		room4.CreatedAt = time.Now().Add(-25 * time.Hour)
		store.CreateRoom(room4)

		// Run cleanup
		err := store.CleanupStaleRooms()
		if err != nil {
			t.Fatalf("cleanup failed: %v", err)
		}

		// Check results
		_, err = store.GetRoom("OLD_FIN")
		if err != ErrRoomNotFound {
			t.Error("expected OLD_FIN to be deleted")
		}

		_, err = store.GetRoom("NEW_FIN")
		if err != nil {
			t.Error("expected NEW_FIN to be preserved")
		}

		_, err = store.GetRoom("ACTIVE")
		if err != nil {
			t.Error("expected ACTIVE to be preserved")
		}

		_, err = store.GetRoom("STALE")
		if err != ErrRoomNotFound {
			t.Error("expected STALE to be deleted")
		}
	})
}

func TestMemoryStore_ErrorCases(t *testing.T) {
	t.Parallel()

	t.Run("all error types are distinct", func(t *testing.T) {
		t.Parallel()

		if ErrRoomNotFound == ErrRoomExists {
			t.Error("ErrRoomNotFound and ErrRoomExists should be distinct")
		}

		if ErrRoomNotFound == ErrPlayerNotFound {
			t.Error("ErrRoomNotFound and ErrPlayerNotFound should be distinct")
		}

		if ErrRoomExists == ErrPlayerNotFound {
			t.Error("ErrRoomExists and ErrPlayerNotFound should be distinct")
		}
	})
}

