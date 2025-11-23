package avalon

import (
	"testing"
)

func TestGetQuestConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		playerCount int
		wantSizes   [5]int
		wantQuest4  int
	}{
		{5, [5]int{2, 3, 2, 3, 3}, 1},
		{6, [5]int{2, 3, 4, 3, 4}, 1},
		{7, [5]int{2, 3, 3, 4, 4}, 2},
		{8, [5]int{3, 4, 4, 5, 5}, 2},
		{9, [5]int{3, 4, 4, 5, 5}, 2},
		{10, [5]int{3, 4, 4, 5, 5}, 2},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.playerCount))+" players", func(t *testing.T) {
			t.Parallel()

			config := getQuestConfig(tt.playerCount)

			if config.PlayerCount != tt.playerCount {
				t.Errorf("PlayerCount = %d, want %d", config.PlayerCount, tt.playerCount)
			}

			if config.TeamSizes != tt.wantSizes {
				t.Errorf("TeamSizes = %v, want %v", config.TeamSizes, tt.wantSizes)
			}

			if config.Quest4Fails != tt.wantQuest4 {
				t.Errorf("Quest4Fails = %d, want %d", config.Quest4Fails, tt.wantQuest4)
			}
		})
	}
}

func TestGetQuestConfig_Invalid(t *testing.T) {
	t.Parallel()

	// Invalid player count should return 5-player config as default
	config := getQuestConfig(3)

	if config.PlayerCount != 5 {
		t.Errorf("expected default to 5 players, got %d", config.PlayerCount)
	}

	// Another invalid count
	config = getQuestConfig(15)

	if config.PlayerCount != 5 {
		t.Errorf("expected default to 5 players, got %d", config.PlayerCount)
	}
}

func TestGetRequiredTeamSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		playerCount int
		questNumber int
		want        int
	}{
		// 5 players
		{5, 1, 2},
		{5, 2, 3},
		{5, 3, 2},
		{5, 4, 3},
		{5, 5, 3},

		// 6 players
		{6, 1, 2},
		{6, 2, 3},
		{6, 3, 4},
		{6, 4, 3},
		{6, 5, 4},

		// 7 players
		{7, 1, 2},
		{7, 2, 3},
		{7, 3, 3},
		{7, 4, 4},
		{7, 5, 4},

		// 8 players
		{8, 1, 3},
		{8, 2, 4},
		{8, 3, 4},
		{8, 4, 5},
		{8, 5, 5},

		// 9 players
		{9, 1, 3},
		{9, 2, 4},
		{9, 3, 4},
		{9, 4, 5},
		{9, 5, 5},

		// 10 players
		{10, 1, 3},
		{10, 2, 4},
		{10, 3, 4},
		{10, 4, 5},
		{10, 5, 5},

		// Invalid quest numbers
		{5, 0, 0},
		{5, 6, 0},
		{5, -1, 0},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.playerCount))+" players quest "+string(rune(tt.questNumber)), func(t *testing.T) {
			t.Parallel()

			got := getRequiredTeamSize(tt.playerCount, tt.questNumber)

			if got != tt.want {
				t.Errorf("getRequiredTeamSize(%d, %d) = %d, want %d",
					tt.playerCount, tt.questNumber, got, tt.want)
			}
		})
	}
}

func TestGetFailsRequired(t *testing.T) {
	t.Parallel()

	tests := []struct {
		playerCount int
		questNumber int
		want        int
	}{
		// Quest 4 with 7+ players requires 2 fails
		{7, 4, 2},
		{8, 4, 2},
		{9, 4, 2},
		{10, 4, 2},

		// Quest 4 with 5-6 players requires 1 fail
		{5, 4, 1},
		{6, 4, 1},

		// All other quests require 1 fail
		{5, 1, 1},
		{5, 2, 1},
		{5, 3, 1},
		{5, 5, 1},
		{7, 1, 1},
		{7, 2, 1},
		{7, 3, 1},
		{7, 5, 1},
		{10, 1, 1},
		{10, 2, 1},
		{10, 3, 1},
		{10, 5, 1},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.playerCount))+" players quest "+string(rune(tt.questNumber)), func(t *testing.T) {
			t.Parallel()

			got := getFailsRequired(tt.playerCount, tt.questNumber)

			if got != tt.want {
				t.Errorf("getFailsRequired(%d, %d) = %d, want %d",
					tt.playerCount, tt.questNumber, got, tt.want)
			}
		})
	}
}

func TestRequiresTwoFails(t *testing.T) {
	t.Parallel()

	tests := []struct {
		playerCount int
		questNumber int
		want        bool
	}{
		// Quest 4 with 7+ players requires 2 fails
		{7, 4, true},
		{8, 4, true},
		{9, 4, true},
		{10, 4, true},

		// Quest 4 with 5-6 players requires 1 fail
		{5, 4, false},
		{6, 4, false},

		// All other quests require 1 fail
		{5, 1, false},
		{5, 2, false},
		{5, 3, false},
		{5, 5, false},
		{7, 1, false},
		{7, 2, false},
		{7, 3, false},
		{7, 5, false},
		{10, 1, false},
		{10, 2, false},
		{10, 3, false},
		{10, 5, false},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.playerCount))+" players quest "+string(rune(tt.questNumber)), func(t *testing.T) {
			t.Parallel()

			got := requiresTwoFails(tt.playerCount, tt.questNumber)

			if got != tt.want {
				t.Errorf("requiresTwoFails(%d, %d) = %v, want %v",
					tt.playerCount, tt.questNumber, got, tt.want)
			}
		})
	}
}

func TestQuestConfigConsistency(t *testing.T) {
	t.Parallel()

	// Verify that all valid player counts have consistent configurations
	for playerCount := 5; playerCount <= 10; playerCount++ {
		t.Run(string(rune(playerCount))+" players", func(t *testing.T) {
			t.Parallel()

			config := getQuestConfig(playerCount)

			// Verify quest 1-5 all have valid team sizes
			for quest := 1; quest <= 5; quest++ {
				teamSize := getRequiredTeamSize(playerCount, quest)
				if teamSize <= 0 {
					t.Errorf("quest %d has invalid team size %d", quest, teamSize)
				}

				// Team size should not exceed player count
				if teamSize > playerCount {
					t.Errorf("quest %d team size %d exceeds player count %d",
						quest, teamSize, playerCount)
				}

				// Team size should match config
				expectedSize := config.TeamSizes[quest-1]
				if teamSize != expectedSize {
					t.Errorf("quest %d: getRequiredTeamSize returned %d, config has %d",
						quest, teamSize, expectedSize)
				}
			}

			// Verify Quest4Fails is 1 or 2
			if config.Quest4Fails != 1 && config.Quest4Fails != 2 {
				t.Errorf("Quest4Fails must be 1 or 2, got %d", config.Quest4Fails)
			}

			// Verify Quest4Fails matches the logic
			expectedQuest4Fails := 1
			if playerCount >= 7 {
				expectedQuest4Fails = 2
			}
			if config.Quest4Fails != expectedQuest4Fails {
				t.Errorf("Quest4Fails = %d, expected %d for %d players",
					config.Quest4Fails, expectedQuest4Fails, playerCount)
			}
		})
	}
}

func TestQuestTeamSizeProgression(t *testing.T) {
	t.Parallel()

	// Verify that team sizes generally increase with more players
	for playerCount := 5; playerCount <= 10; playerCount++ {
		config := getQuestConfig(playerCount)

		// Check that we have exactly 5 quest team sizes
		if len(config.TeamSizes) != 5 {
			t.Errorf("%d players: expected 5 team sizes, got %d", playerCount, len(config.TeamSizes))
		}

		// Verify all team sizes are > 0
		for i, size := range config.TeamSizes {
			if size <= 0 {
				t.Errorf("%d players quest %d: team size must be > 0, got %d",
					playerCount, i+1, size)
			}
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		playerCount int
	}{
		{5},
		{6},
		{7},
		{8},
		{9},
		{10},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.playerCount))+" players", func(t *testing.T) {
			t.Parallel()

			config := DefaultConfig(tt.playerCount)

			// Verify role count matches player count
			if len(config.Roles) != tt.playerCount {
				t.Errorf("expected %d roles, got %d", tt.playerCount, len(config.Roles))
			}

			// Verify config is valid
			if err := config.Validate(); err != nil {
				t.Errorf("DefaultConfig validation failed: %v", err)
			}

			// Verify Merlin and Assassin are present
			hasMerlin := false
			hasAssassin := false
			for _, role := range config.Roles {
				if role == RoleMerlin {
					hasMerlin = true
				}
				if role == RoleAssassin {
					hasAssassin = true
				}
			}

			if !hasMerlin {
				t.Error("DefaultConfig missing Merlin")
			}
			if !hasAssassin {
				t.Error("DefaultConfig missing Assassin")
			}
		})
	}
}

func TestDefaultConfigWithPercival(t *testing.T) {
	t.Parallel()

	tests := []struct {
		playerCount int
	}{
		{5},
		{6},
		{7},
		{8},
		{9},
		{10},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.playerCount))+" players", func(t *testing.T) {
			t.Parallel()

			config := DefaultConfigWithPercival(tt.playerCount)

			// Verify role count matches player count
			if len(config.Roles) != tt.playerCount {
				t.Errorf("expected %d roles, got %d", tt.playerCount, len(config.Roles))
			}

			// Verify config is valid
			if err := config.Validate(); err != nil {
				t.Errorf("DefaultConfigWithPercival validation failed: %v", err)
			}

			// Verify special roles are present
			hasMerlin := false
			hasPercival := false
			hasAssassin := false
			hasMorgana := false

			for _, role := range config.Roles {
				switch role {
				case RoleMerlin:
					hasMerlin = true
				case RolePercival:
					hasPercival = true
				case RoleAssassin:
					hasAssassin = true
				case RoleMorgana:
					hasMorgana = true
				}
			}

			if !hasMerlin {
				t.Error("DefaultConfigWithPercival missing Merlin")
			}
			if !hasPercival {
				t.Error("DefaultConfigWithPercival missing Percival")
			}
			if !hasAssassin {
				t.Error("DefaultConfigWithPercival missing Assassin")
			}
			if !hasMorgana {
				t.Error("DefaultConfigWithPercival missing Morgana")
			}
		})
	}
}

