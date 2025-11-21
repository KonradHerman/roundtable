package core

import (
	"testing"
	"time"
)

func TestNewRoom(t *testing.T) {
	t.Parallel()

	hostPlayer := &Player{
		ID:           "host-123",
		DisplayName:  "Alice",
		SessionToken: "token-123",
	}

	room := NewRoom("ABC123", "werewolf", hostPlayer, 10)

	if room.ID != "ABC123" {
		t.Errorf("expected room ID 'ABC123', got '%s'", room.ID)
	}

	if room.GameType != "werewolf" {
		t.Errorf("expected game type 'werewolf', got '%s'", room.GameType)
	}

	if room.MaxPlayers != 10 {
		t.Errorf("expected max players 10, got %d", room.MaxPlayers)
	}

	if room.HostID != "host-123" {
		t.Errorf("expected host ID 'host-123', got '%s'", room.HostID)
	}

	if room.Status != RoomStatusWaiting {
		t.Errorf("expected status 'waiting', got '%s'", room.Status)
	}

	if len(room.Players) != 1 {
		t.Errorf("expected 1 player, got %d", len(room.Players))
	}

	if room.Players["host-123"] != hostPlayer {
		t.Error("host player not found in room")
	}

	if time.Since(room.CreatedAt) > time.Second {
		t.Error("room creation time is too old")
	}
}

func TestRoom_AddPlayer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupRoom   func() *Room
		playerToAdd *Player
		wantErr     bool
		errContains string
	}{
		{
			name: "successfully add player to waiting room",
			setupRoom: func() *Room {
				host := &Player{ID: "host", DisplayName: "Host", SessionToken: "token-host"}
				return NewRoom("ABC123", "werewolf", host, 10)
			},
			playerToAdd: &Player{ID: "player1", DisplayName: "Player1", SessionToken: "token-1"},
			wantErr:     false,
		},
		{
			name: "fail when adding to started game",
			setupRoom: func() *Room {
				host := &Player{ID: "host", DisplayName: "Host", SessionToken: "token-host"}
				room := NewRoom("ABC123", "werewolf", host, 10)
				room.Status = RoomStatusPlaying
				return room
			},
			playerToAdd: &Player{ID: "player1", DisplayName: "Player1", SessionToken: "token-1"},
			wantErr:     true,
			errContains: "game already started",
		},
		{
			name: "fail when room is full",
			setupRoom: func() *Room {
				host := &Player{ID: "host", DisplayName: "Host", SessionToken: "token-host"}
				room := NewRoom("ABC123", "werewolf", host, 2)
				room.AddPlayer(&Player{ID: "player1", DisplayName: "Player1", SessionToken: "token-1"})
				return room
			},
			playerToAdd: &Player{ID: "player2", DisplayName: "Player2", SessionToken: "token-2"},
			wantErr:     true,
			errContains: "room is full",
		},
		{
			name: "fail when player already exists",
			setupRoom: func() *Room {
				host := &Player{ID: "host", DisplayName: "Host", SessionToken: "token-host"}
				room := NewRoom("ABC123", "werewolf", host, 10)
				room.AddPlayer(&Player{ID: "player1", DisplayName: "Player1", SessionToken: "token-1"})
				return room
			},
			playerToAdd: &Player{ID: "player1", DisplayName: "Player1 Duplicate", SessionToken: "token-1-dup"},
			wantErr:     true,
			errContains: "player already in room",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			room := tt.setupRoom()
			err := room.AddPlayer(tt.playerToAdd)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if _, exists := room.Players[tt.playerToAdd.ID]; !exists {
					t.Error("player was not added to room")
				}
			}
		})
	}
}

func TestRoom_RemovePlayer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupRoom   func() *Room
		playerID    string
		wantErr     bool
		errContains string
	}{
		{
			name: "successfully remove player",
			setupRoom: func() *Room {
				host := &Player{ID: "host", DisplayName: "Host", SessionToken: "token-host"}
				room := NewRoom("ABC123", "werewolf", host, 10)
				room.AddPlayer(&Player{ID: "player1", DisplayName: "Player1", SessionToken: "token-1"})
				return room
			},
			playerID: "player1",
			wantErr:  false,
		},
		{
			name: "fail when player not in room",
			setupRoom: func() *Room {
				host := &Player{ID: "host", DisplayName: "Host", SessionToken: "token-host"}
				return NewRoom("ABC123", "werewolf", host, 10)
			},
			playerID:    "nonexistent",
			wantErr:     true,
			errContains: "player not in room",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			room := tt.setupRoom()
			initialCount := len(room.Players)

			err := room.RemovePlayer(tt.playerID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if _, exists := room.Players[tt.playerID]; exists {
					t.Error("player was not removed from room")
				}
				if len(room.Players) != initialCount-1 {
					t.Errorf("expected %d players, got %d", initialCount-1, len(room.Players))
				}
			}
		})
	}
}

func TestRoom_GetPlayer(t *testing.T) {
	t.Parallel()

	host := &Player{ID: "host", DisplayName: "Host", SessionToken: "token-host"}
	player1 := &Player{ID: "player1", DisplayName: "Player1", SessionToken: "token-1"}
	room := NewRoom("ABC123", "werewolf", host, 10)
	room.AddPlayer(player1)

	tests := []struct {
		name        string
		playerID    string
		wantPlayer  *Player
		wantErr     bool
		errContains string
	}{
		{
			name:       "get existing host player",
			playerID:   "host",
			wantPlayer: host,
			wantErr:    false,
		},
		{
			name:       "get existing non-host player",
			playerID:   "player1",
			wantPlayer: player1,
			wantErr:    false,
		},
		{
			name:        "fail when player not found",
			playerID:    "nonexistent",
			wantPlayer:  nil,
			wantErr:     true,
			errContains: "player not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			player, err := room.GetPlayer(tt.playerID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if player.ID != tt.wantPlayer.ID {
					t.Errorf("expected player ID '%s', got '%s'", tt.wantPlayer.ID, player.ID)
				}
			}
		})
	}
}

func TestRoom_GetPlayerByToken(t *testing.T) {
	t.Parallel()

	host := &Player{ID: "host", DisplayName: "Host", SessionToken: "token-host"}
	player1 := &Player{ID: "player1", DisplayName: "Player1", SessionToken: "token-1"}
	room := NewRoom("ABC123", "werewolf", host, 10)
	room.AddPlayer(player1)

	tests := []struct {
		name        string
		token       string
		wantPlayer  *Player
		wantErr     bool
		errContains string
	}{
		{
			name:       "get player by valid token",
			token:      "token-1",
			wantPlayer: player1,
			wantErr:    false,
		},
		{
			name:       "get host by valid token",
			token:      "token-host",
			wantPlayer: host,
			wantErr:    false,
		},
		{
			name:        "fail when token not found",
			token:       "invalid-token",
			wantPlayer:  nil,
			wantErr:     true,
			errContains: "invalid session token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			player, err := room.GetPlayerByToken(tt.token)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if player.ID != tt.wantPlayer.ID {
					t.Errorf("expected player ID '%s', got '%s'", tt.wantPlayer.ID, player.ID)
				}
			}
		})
	}
}

func TestRoom_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	host := &Player{ID: "host", DisplayName: "Host", SessionToken: "token-host"}
	room := NewRoom("ABC123", "werewolf", host, 100)

	// Test concurrent player additions
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			player := &Player{
				ID:           string(rune('A' + id)),
				DisplayName:  "Player",
				SessionToken: string(rune('A' + id)),
			}
			room.AddPlayer(player)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all players were added (host + 10 players = 11)
	if len(room.Players) != 11 {
		t.Errorf("expected 11 players, got %d", len(room.Players))
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && hasSubstring(s, substr)))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
