package games

import (
	"encoding/json"
	"testing"

	"github.com/KonradHerman/roundtable/internal/core"
)

func TestNewRegistry(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	if registry == nil {
		t.Fatal("NewRegistry returned nil")
	}

	// Verify werewolf and avalon are registered
	if !registry.IsRegistered("werewolf") {
		t.Error("werewolf should be registered")
	}

	if !registry.IsRegistered("avalon") {
		t.Error("avalon should be registered")
	}
}

func TestRegistry_Register(t *testing.T) {
	t.Parallel()

	registry := &Registry{
		factories: make(map[string]GameFactory),
		parsers:   make(map[string]ConfigParser),
	}

	// Create a mock game factory
	mockFactory := func() core.Game {
		return nil
	}

	// Create a mock config parser
	mockParser := func(data []byte) (core.GameConfig, error) {
		return nil, nil
	}

	// Register a new game type
	registry.Register("testgame", mockFactory, mockParser)

	// Verify it was registered
	if !registry.IsRegistered("testgame") {
		t.Error("testgame should be registered")
	}
}

func TestRegistry_IsRegistered(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	tests := []struct {
		gameType string
		want     bool
	}{
		{"werewolf", true},
		{"avalon", true},
		{"unknown", false},
		{"spyfall", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.gameType, func(t *testing.T) {
			t.Parallel()

			got := registry.IsRegistered(tt.gameType)
			if got != tt.want {
				t.Errorf("IsRegistered(%s) = %v, want %v", tt.gameType, got, tt.want)
			}
		})
	}
}

func TestRegistry_ListGames(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	games := registry.ListGames()

	if len(games) == 0 {
		t.Error("expected at least some games to be registered")
	}

	// Verify werewolf and avalon are in the list
	hasWerewolf := false
	hasAvalon := false

	for _, game := range games {
		if game == "werewolf" {
			hasWerewolf = true
		}
		if game == "avalon" {
			hasAvalon = true
		}
	}

	if !hasWerewolf {
		t.Error("werewolf should be in game list")
	}

	if !hasAvalon {
		t.Error("avalon should be in game list")
	}
}

func TestRegistry_CreateGame(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	t.Run("create werewolf game", func(t *testing.T) {
		t.Parallel()

		game, err := registry.CreateGame("werewolf")
		if err != nil {
			t.Fatalf("failed to create werewolf game: %v", err)
		}

		if game == nil {
			t.Fatal("created game is nil")
		}
	})

	t.Run("create avalon game", func(t *testing.T) {
		t.Parallel()

		game, err := registry.CreateGame("avalon")
		if err != nil {
			t.Fatalf("failed to create avalon game: %v", err)
		}

		if game == nil {
			t.Fatal("created game is nil")
		}
	})

	t.Run("fail on unknown game type", func(t *testing.T) {
		t.Parallel()

		_, err := registry.CreateGame("unknown")
		if err == nil {
			t.Fatal("expected error for unknown game type")
		}

		if !contains(err.Error(), "unknown game type") {
			t.Errorf("expected 'unknown game type' error, got: %v", err)
		}
	})

	t.Run("fail on empty game type", func(t *testing.T) {
		t.Parallel()

		_, err := registry.CreateGame("")
		if err == nil {
			t.Fatal("expected error for empty game type")
		}
	})
}

func TestRegistry_ParseConfig(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	t.Run("parse werewolf config", func(t *testing.T) {
		t.Parallel()

		configJSON := `{
			"roles": ["werewolf", "werewolf", "seer", "villager", "villager", "villager"],
			"nightDuration": 180000000000,
			"dayDuration": 300000000000
		}`

		config, err := registry.ParseConfig("werewolf", json.RawMessage(configJSON))
		if err != nil {
			t.Fatalf("failed to parse werewolf config: %v", err)
		}

		if config == nil {
			t.Fatal("parsed config is nil")
		}

		if config.GameType() != "werewolf" {
			t.Errorf("config game type = %s, want werewolf", config.GameType())
		}
	})

	t.Run("parse avalon config", func(t *testing.T) {
		t.Parallel()

		configJSON := `{
			"roles": ["merlin", "assassin", "loyal_servant", "loyal_servant", "minion"]
		}`

		config, err := registry.ParseConfig("avalon", json.RawMessage(configJSON))
		if err != nil {
			t.Fatalf("failed to parse avalon config: %v", err)
		}

		if config == nil {
			t.Fatal("parsed config is nil")
		}

		if config.GameType() != "avalon" {
			t.Errorf("config game type = %s, want avalon", config.GameType())
		}
	})

	t.Run("fail on unknown game type", func(t *testing.T) {
		t.Parallel()

		configJSON := `{"roles": []}`

		_, err := registry.ParseConfig("unknown", json.RawMessage(configJSON))
		if err == nil {
			t.Fatal("expected error for unknown game type")
		}

		if !contains(err.Error(), "unknown game type") {
			t.Errorf("expected 'unknown game type' error, got: %v", err)
		}
	})

	t.Run("fail on invalid JSON", func(t *testing.T) {
		t.Parallel()

		invalidJSON := `{invalid json`

		_, err := registry.ParseConfig("werewolf", json.RawMessage(invalidJSON))
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
	})
}

func TestRegistry_ValidateConfig(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	t.Run("valid werewolf config", func(t *testing.T) {
		t.Parallel()

		configJSON := `{
			"roles": ["werewolf", "seer", "villager", "villager", "villager", "villager"],
			"nightDuration": 180000000000,
			"dayDuration": 300000000000
		}`

		err := registry.ValidateConfig("werewolf", json.RawMessage(configJSON))
		if err != nil {
			t.Errorf("valid config should pass validation: %v", err)
		}
	})

	t.Run("valid avalon config", func(t *testing.T) {
		t.Parallel()

		configJSON := `{
			"roles": ["merlin", "assassin", "loyal_servant", "loyal_servant", "minion"]
		}`

		err := registry.ValidateConfig("avalon", json.RawMessage(configJSON))
		if err != nil {
			t.Errorf("valid config should pass validation: %v", err)
		}
	})

	t.Run("invalid werewolf config - no werewolf", func(t *testing.T) {
		t.Parallel()

		configJSON := `{
			"roles": ["seer", "villager", "villager"]
		}`

		err := registry.ValidateConfig("werewolf", json.RawMessage(configJSON))
		if err == nil {
			t.Fatal("expected validation error for config without werewolf")
		}

		if !contains(err.Error(), "invalid configuration") {
			t.Errorf("expected validation error message, got: %v", err)
		}
	})

	t.Run("invalid avalon config - no assassin with Merlin", func(t *testing.T) {
		t.Parallel()

		configJSON := `{
			"roles": ["merlin", "loyal_servant", "loyal_servant", "minion", "minion"]
		}`

		err := registry.ValidateConfig("avalon", json.RawMessage(configJSON))
		if err == nil {
			t.Fatal("expected validation error for Merlin without Assassin")
		}
	})

	t.Run("fail on unknown game type", func(t *testing.T) {
		t.Parallel()

		configJSON := `{"roles": []}`

		err := registry.ValidateConfig("unknown", json.RawMessage(configJSON))
		if err == nil {
			t.Fatal("expected error for unknown game type")
		}
	})

	t.Run("fail on invalid JSON", func(t *testing.T) {
		t.Parallel()

		invalidJSON := `{invalid`

		err := registry.ValidateConfig("werewolf", json.RawMessage(invalidJSON))
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
	})
}

func TestRegistry_MultipleGames(t *testing.T) {
	t.Parallel()

	t.Run("create multiple game instances independently", func(t *testing.T) {
		t.Parallel()

		registry := NewRegistry()

		// Create two werewolf games
		game1, err := registry.CreateGame("werewolf")
		if err != nil {
			t.Fatalf("failed to create game1: %v", err)
		}

		game2, err := registry.CreateGame("werewolf")
		if err != nil {
			t.Fatalf("failed to create game2: %v", err)
		}

		// Verify they are different instances
		if game1 == game2 {
			t.Error("expected different game instances")
		}
	})

	t.Run("create different game types", func(t *testing.T) {
		t.Parallel()

		registry := NewRegistry()

		werewolfGame, err := registry.CreateGame("werewolf")
		if err != nil {
			t.Fatalf("failed to create werewolf game: %v", err)
		}

		avalonGame, err := registry.CreateGame("avalon")
		if err != nil {
			t.Fatalf("failed to create avalon game: %v", err)
		}

		// Verify they are different types
		if werewolfGame == avalonGame {
			t.Error("expected different game types to be different instances")
		}
	})
}

func TestRegistry_ConfigParsing(t *testing.T) {
	t.Parallel()

	t.Run("parse and validate werewolf config end-to-end", func(t *testing.T) {
		t.Parallel()

		registry := NewRegistry()

		configJSON := `{
			"roles": ["werewolf", "seer", "robber", "villager", "villager", "villager"]
		}`

		// Parse config
		config, err := registry.ParseConfig("werewolf", json.RawMessage(configJSON))
		if err != nil {
			t.Fatalf("failed to parse config: %v", err)
		}

		// Validate config
		err = config.Validate()
		if err != nil {
			t.Errorf("config validation failed: %v", err)
		}

		// Create game
		game, err := registry.CreateGame("werewolf")
		if err != nil {
			t.Fatalf("failed to create game: %v", err)
		}

		if game == nil {
			t.Fatal("game is nil")
		}
	})

	t.Run("parse and validate avalon config end-to-end", func(t *testing.T) {
		t.Parallel()

		registry := NewRegistry()

		configJSON := `{
			"roles": ["merlin", "percival", "loyal_servant", "loyal_servant", "assassin", "morgana", "mordred"]
		}`

		// Parse config
		config, err := registry.ParseConfig("avalon", json.RawMessage(configJSON))
		if err != nil {
			t.Fatalf("failed to parse config: %v", err)
		}

		// Validate config
		err = config.Validate()
		if err != nil {
			t.Errorf("config validation failed: %v", err)
		}

		// Create game
		game, err := registry.CreateGame("avalon")
		if err != nil {
			t.Fatalf("failed to create game: %v", err)
		}

		if game == nil {
			t.Fatal("game is nil")
		}
	})
}

// Helper function
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

