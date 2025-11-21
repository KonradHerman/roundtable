package werewolf

import (
	"testing"
	"time"

	"github.com/KonradHerman/roundtable/internal/core"
)

func TestNewGame(t *testing.T) {
	t.Parallel()

	game := NewGame()

	if game == nil {
		t.Fatal("NewGame returned nil")
	}

	g := game.(*Game)
	if g.phase != PhaseSetup {
		t.Errorf("expected phase 'setup', got '%s'", g.phase)
	}

	if g.players == nil {
		t.Error("players map is nil")
	}

	if g.roleAssignments == nil {
		t.Error("roleAssignments map is nil")
	}
}

func TestConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		config      *Config
		wantErr     bool
		errContains string
	}{
		{
			name: "valid config with basic roles",
			config: &Config{
				Roles:         []RoleType{RoleWerewolf, RoleWerewolf, RoleSeer, RoleRobber, RoleVillager, RoleVillager},
				NightDuration: 3 * time.Minute,
				DayDuration:   5 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "valid config with default durations",
			config: &Config{
				Roles: []RoleType{RoleWerewolf, RoleSeer, RoleVillager},
			},
			wantErr: false,
		},
		{
			name: "fail when no roles provided",
			config: &Config{
				Roles: []RoleType{},
			},
			wantErr:     true,
			errContains: "at least one role required",
		},
		{
			name: "fail when no werewolf in roles",
			config: &Config{
				Roles: []RoleType{RoleSeer, RoleVillager, RoleVillager},
			},
			wantErr:     true,
			errContains: "at least one werewolf required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.config.Validate()

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
				// Verify default durations were set
				if tt.config.NightDuration == 0 {
					t.Error("NightDuration should have been set to default")
				}
				if tt.config.DayDuration == 0 {
					t.Error("DayDuration should have been set to default")
				}
			}
		})
	}
}

func TestGame_Initialize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		config      *Config
		players     []*core.Player
		wantErr     bool
		errContains string
	}{
		{
			name: "successfully initialize game with correct role count",
			config: &Config{
				Roles:         []RoleType{RoleWerewolf, RoleWerewolf, RoleSeer, RoleRobber, RoleVillager, RoleVillager},
				NightDuration: 3 * time.Minute,
				DayDuration:   5 * time.Minute,
			},
			players: []*core.Player{
				{ID: "p1", DisplayName: "Player1"},
				{ID: "p2", DisplayName: "Player2"},
				{ID: "p3", DisplayName: "Player3"},
			},
			wantErr: false,
		},
		{
			name: "fail with incorrect role count (too few)",
			config: &Config{
				Roles:         []RoleType{RoleWerewolf, RoleSeer},
				NightDuration: 3 * time.Minute,
				DayDuration:   5 * time.Minute,
			},
			players: []*core.Player{
				{ID: "p1", DisplayName: "Player1"},
				{ID: "p2", DisplayName: "Player2"},
			},
			wantErr:     true,
			errContains: "role count",
		},
		{
			name: "fail with incorrect role count (too many)",
			config: &Config{
				Roles:         []RoleType{RoleWerewolf, RoleSeer, RoleVillager, RoleVillager, RoleVillager, RoleVillager, RoleVillager},
				NightDuration: 3 * time.Minute,
				DayDuration:   5 * time.Minute,
			},
			players: []*core.Player{
				{ID: "p1", DisplayName: "Player1"},
				{ID: "p2", DisplayName: "Player2"},
			},
			wantErr:     true,
			errContains: "role count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			game := NewGame()
			events, err := game.Initialize(tt.config, tt.players)

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

				// Verify events were generated
				if len(events) == 0 {
					t.Error("expected events to be generated")
				}

				g := game.(*Game)

				// Verify all players received roles
				if len(g.roleAssignments) != len(tt.players) {
					t.Errorf("expected %d role assignments, got %d", len(tt.players), len(g.roleAssignments))
				}

				// Verify center cards (should be 3)
				if len(g.centerCards) != 3 {
					t.Errorf("expected 3 center cards, got %d", len(g.centerCards))
				}

				// Verify all roles from config were assigned (players + center)
				totalRolesAssigned := len(g.roleAssignments) + len(g.centerCards)
				if totalRolesAssigned != len(tt.config.Roles) {
					t.Errorf("expected %d total roles assigned, got %d", len(tt.config.Roles), totalRolesAssigned)
				}

				// Verify phase is role_reveal (waiting for role acknowledgements)
				if g.phase != PhaseRoleReveal {
					t.Errorf("expected phase 'role_reveal', got '%s'", g.phase)
				}
			}
		})
	}
}

func TestRoleType_IsWerewolfTeam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		role RoleType
		want bool
	}{
		{RoleWerewolf, true},
		{RoleMinion, true},
		{RoleSeer, false},
		{RoleVillager, false},
		{RoleTanner, false},
		{RoleRobber, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			t.Parallel()

			if got := tt.role.IsWerewolfTeam(); got != tt.want {
				t.Errorf("IsWerewolfTeam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoleType_IsVillageTeam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		role RoleType
		want bool
	}{
		{RoleSeer, true},
		{RoleVillager, true},
		{RoleRobber, true},
		{RoleWerewolf, false},
		{RoleMinion, false},
		{RoleTanner, false}, // Tanner is neutral
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			t.Parallel()

			if got := tt.role.IsVillageTeam(); got != tt.want {
				t.Errorf("IsVillageTeam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoleType_HasNightAction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		role RoleType
		want bool
	}{
		{RoleSeer, true},
		{RoleRobber, true},
		{RoleTroublemaker, true},
		{RoleDrunk, true},
		{RoleInsomniac, true},
		{RoleWerewolf, false},
		{RoleVillager, false},
		{RoleMason, false},
		{RoleMinion, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			t.Parallel()

			if got := tt.role.HasNightAction(); got != tt.want {
				t.Errorf("HasNightAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGame_GetPlayerState(t *testing.T) {
	t.Parallel()

	game := NewGame()
	config := &Config{
		Roles:         []RoleType{RoleWerewolf, RoleSeer, RoleVillager, RoleVillager, RoleVillager, RoleVillager},
		NightDuration: 3 * time.Minute,
		DayDuration:   5 * time.Minute,
	}
	players := []*core.Player{
		{ID: "p1", DisplayName: "Player1"},
		{ID: "p2", DisplayName: "Player2"},
		{ID: "p3", DisplayName: "Player3"},
	}

	_, err := game.Initialize(config, players)
	if err != nil {
		t.Fatalf("failed to initialize game: %v", err)
	}

	// Get state for player 1
	state := game.GetPlayerState("p1")
	if state == nil {
		t.Fatal("GetPlayerState returned nil")
	}

	// Verify state structure
	playerState, ok := state.(PlayerState)
	if !ok {
		t.Fatalf("state is not PlayerState, got %T", state)
	}

	// Check fields
	if playerState.Phase == "" {
		t.Error("phase is empty")
	}
	if playerState.YourRole == "" {
		t.Error("yourRole is empty")
	}
}

func TestGame_GetPublicState(t *testing.T) {
	t.Parallel()

	game := NewGame()
	config := &Config{
		Roles:         []RoleType{RoleWerewolf, RoleSeer, RoleVillager, RoleVillager, RoleVillager, RoleVillager},
		NightDuration: 3 * time.Minute,
		DayDuration:   5 * time.Minute,
	}
	players := []*core.Player{
		{ID: "p1", DisplayName: "Player1"},
		{ID: "p2", DisplayName: "Player2"},
		{ID: "p3", DisplayName: "Player3"},
	}

	_, err := game.Initialize(config, players)
	if err != nil {
		t.Fatalf("failed to initialize game: %v", err)
	}

	// Get public state
	state := game.GetPublicState()
	if state == nil {
		t.Fatal("GetPublicState returned nil")
	}

	// Verify state structure
	publicState, ok := state.(PublicState)
	if !ok {
		t.Fatalf("state is not PublicState, got %T", state)
	}

	// Check fields
	if publicState.Phase == "" {
		t.Error("phase is empty")
	}

	// Verify player count
	if publicState.PlayerCount != 3 {
		t.Errorf("expected playerCount 3, got %d", publicState.PlayerCount)
	}
}

func TestGame_GetPhase(t *testing.T) {
	t.Parallel()

	game := NewGame()
	g := game.(*Game)

	tests := []struct {
		name       string
		setupPhase Phase
		want       string
	}{
		{"setup phase", PhaseSetup, "setup"},
		{"role reveal phase", PhaseRoleReveal, "role_reveal"},
		{"night phase", PhaseNight, "night"},
		{"day phase", PhaseDay, "day"},
		{"results phase", PhaseResults, "results"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.phase = tt.setupPhase
			phase := game.GetPhase()

			if phase.Name != tt.want {
				t.Errorf("expected phase '%s', got '%s'", tt.want, phase.Name)
			}
		})
	}
}

func TestGame_IsFinished(t *testing.T) {
	t.Parallel()

	game := NewGame()
	g := game.(*Game)

	tests := []struct {
		name  string
		phase Phase
		want  bool
	}{
		{"not finished in setup", PhaseSetup, false},
		{"not finished in role reveal", PhaseRoleReveal, false},
		{"not finished in night", PhaseNight, false},
		{"not finished in day", PhaseDay, false},
		{"not finished in results", PhaseResults, false}, // Game continues until explicitly reset
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.phase = tt.phase
			if got := game.IsFinished(); got != tt.want {
				t.Errorf("IsFinished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGame_RoleAssignment_Uniqueness(t *testing.T) {
	t.Parallel()

	// Run multiple times to test randomness
	for run := 0; run < 10; run++ {
		game := NewGame()
		config := &Config{
			Roles: []RoleType{
				RoleWerewolf, RoleWerewolf, RoleSeer, RoleRobber,
				RoleTroublemaker, RoleVillager, RoleVillager, RoleVillager,
			},
			NightDuration: 3 * time.Minute,
			DayDuration:   5 * time.Minute,
		}
		players := []*core.Player{
			{ID: "p1", DisplayName: "Player1"},
			{ID: "p2", DisplayName: "Player2"},
			{ID: "p3", DisplayName: "Player3"},
			{ID: "p4", DisplayName: "Player4"},
			{ID: "p5", DisplayName: "Player5"},
		}

		_, err := game.Initialize(config, players)
		if err != nil {
			t.Fatalf("run %d: failed to initialize game: %v", run, err)
		}

		g := game.(*Game)

		// Track all assigned roles
		assignedRoles := make(map[RoleType]int)
		for _, role := range g.roleAssignments {
			assignedRoles[role]++
		}
		for _, role := range g.centerCards {
			assignedRoles[role]++
		}

		// Verify each role from config appears exactly once
		expectedRoles := make(map[RoleType]int)
		for _, role := range config.Roles {
			expectedRoles[role]++
		}

		for role, expectedCount := range expectedRoles {
			if actualCount, exists := assignedRoles[role]; !exists || actualCount != expectedCount {
				t.Errorf("run %d: role %s expected %d times, got %d", run, role, expectedCount, actualCount)
			}
		}
	}
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
