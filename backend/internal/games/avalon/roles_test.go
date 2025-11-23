package avalon

import (
	"testing"
)

func TestGetRoleTeam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		role Role
		want Team
	}{
		{RoleMerlin, TeamGood},
		{RolePercival, TeamGood},
		{RoleLoyalServant, TeamGood},
		{RoleAssassin, TeamEvil},
		{RoleMorgana, TeamEvil},
		{RoleMordred, TeamEvil},
		{RoleOberon, TeamEvil},
		{RoleMinionOfMordred, TeamEvil},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			t.Parallel()

			if got := getRoleTeam(tt.role); got != tt.want {
				t.Errorf("getRoleTeam(%s) = %v, want %v", tt.role, got, tt.want)
			}
		})
	}
}

func TestIsGoodRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		role Role
		want bool
	}{
		{RoleMerlin, true},
		{RolePercival, true},
		{RoleLoyalServant, true},
		{RoleAssassin, false},
		{RoleMorgana, false},
		{RoleMordred, false},
		{RoleOberon, false},
		{RoleMinionOfMordred, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			t.Parallel()

			if got := isGoodRole(tt.role); got != tt.want {
				t.Errorf("isGoodRole(%s) = %v, want %v", tt.role, got, tt.want)
			}
		})
	}
}

func TestIsEvilRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		role Role
		want bool
	}{
		{RoleAssassin, true},
		{RoleMorgana, true},
		{RoleMordred, true},
		{RoleOberon, true},
		{RoleMinionOfMordred, true},
		{RoleMerlin, false},
		{RolePercival, false},
		{RoleLoyalServant, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			t.Parallel()

			if got := isEvilRole(tt.role); got != tt.want {
				t.Errorf("isEvilRole(%s) = %v, want %v", tt.role, got, tt.want)
			}
		})
	}
}

func TestGetRoleKnowledge_Merlin(t *testing.T) {
	t.Parallel()

	roles := map[string]Role{
		"p1": RoleMerlin,
		"p2": RoleAssassin,
		"p3": RoleMorgana,
		"p4": RoleMordred,
		"p5": RoleLoyalServant,
	}

	knowledge := getRoleKnowledge(RoleMerlin, "p1", roles)

	// Merlin should see Assassin and Morgana but NOT Mordred
	expectedKnown := map[string]bool{
		"p2": true, // Assassin
		"p3": true, // Morgana
	}

	expectedNotKnown := map[string]bool{
		"p1": true, // Self
		"p4": true, // Mordred (hidden from Merlin)
		"p5": true, // Loyal Servant
	}

	for _, pid := range knowledge {
		if !expectedKnown[pid] {
			t.Errorf("Merlin should not know player %s", pid)
		}
	}

	for pid := range expectedKnown {
		found := false
		for _, knownPID := range knowledge {
			if knownPID == pid {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Merlin should know player %s", pid)
		}
	}

	for pid := range expectedNotKnown {
		for _, knownPID := range knowledge {
			if knownPID == pid {
				t.Errorf("Merlin should not know player %s", pid)
			}
		}
	}
}

func TestGetRoleKnowledge_Percival(t *testing.T) {
	t.Parallel()

	roles := map[string]Role{
		"p1": RolePercival,
		"p2": RoleMerlin,
		"p3": RoleMorgana,
		"p4": RoleAssassin,
		"p5": RoleLoyalServant,
	}

	knowledge := getRoleKnowledge(RolePercival, "p1", roles)

	// Percival should see Merlin and Morgana (can't distinguish)
	if len(knowledge) != 2 {
		t.Errorf("Percival should know 2 players, got %d", len(knowledge))
	}

	expectedKnown := map[string]bool{
		"p2": true, // Merlin
		"p3": true, // Morgana
	}

	for _, pid := range knowledge {
		if !expectedKnown[pid] {
			t.Errorf("Percival saw unexpected player %s", pid)
		}
	}
}

func TestGetRoleKnowledge_EvilSeesEachOther(t *testing.T) {
	t.Parallel()

	roles := map[string]Role{
		"p1": RoleAssassin,
		"p2": RoleMorgana,
		"p3": RoleMordred,
		"p4": RoleMinionOfMordred,
		"p5": RoleOberon,
		"p6": RoleMerlin,
	}

	// Assassin should see all evil except Oberon
	knowledge := getRoleKnowledge(RoleAssassin, "p1", roles)

	expectedKnown := map[string]bool{
		"p2": true, // Morgana
		"p3": true, // Mordred
		"p4": true, // Minion
	}

	if len(knowledge) != len(expectedKnown) {
		t.Errorf("Assassin should know %d players, got %d", len(expectedKnown), len(knowledge))
	}

	for _, pid := range knowledge {
		if !expectedKnown[pid] {
			t.Errorf("Assassin saw unexpected player %s", pid)
		}
	}

	// Verify Oberon is NOT in the knowledge
	for _, pid := range knowledge {
		if pid == "p5" {
			t.Error("Assassin should not see Oberon")
		}
	}
}

func TestGetRoleKnowledge_Oberon(t *testing.T) {
	t.Parallel()

	roles := map[string]Role{
		"p1": RoleOberon,
		"p2": RoleAssassin,
		"p3": RoleMorgana,
		"p4": RoleMerlin,
		"p5": RoleLoyalServant,
	}

	knowledge := getRoleKnowledge(RoleOberon, "p1", roles)

	// Oberon should see nothing (lone wolf)
	if len(knowledge) != 0 {
		t.Errorf("Oberon should know 0 players, got %d", len(knowledge))
	}
}

func TestGetRoleKnowledge_LoyalServant(t *testing.T) {
	t.Parallel()

	roles := map[string]Role{
		"p1": RoleLoyalServant,
		"p2": RoleMerlin,
		"p3": RoleAssassin,
		"p4": RoleMorgana,
		"p5": RoleLoyalServant,
	}

	knowledge := getRoleKnowledge(RoleLoyalServant, "p1", roles)

	// Loyal Servant should see nothing
	if len(knowledge) != 0 {
		t.Errorf("Loyal Servant should know 0 players, got %d", len(knowledge))
	}
}

func TestSecureShuffleRoles(t *testing.T) {
	t.Parallel()

	t.Run("shuffle changes order", func(t *testing.T) {
		t.Parallel()

		// Create a slice of roles
		roles := []Role{
			RoleMerlin, RolePercival, RoleLoyalServant, RoleLoyalServant,
			RoleAssassin, RoleMorgana, RoleMordred,
		}

		// Keep original order
		original := make([]Role, len(roles))
		copy(original, roles)

		// Shuffle multiple times and verify it's actually random
		changedAtLeastOnce := false
		for i := 0; i < 20; i++ {
			rolesCopy := make([]Role, len(original))
			copy(rolesCopy, original)

			err := secureShuffleRoles(rolesCopy)
			if err != nil {
				t.Fatalf("secureShuffleRoles failed: %v", err)
			}

			// Check if order changed
			isDifferent := false
			for j := range rolesCopy {
				if rolesCopy[j] != original[j] {
					isDifferent = true
					break
				}
			}

			if isDifferent {
				changedAtLeastOnce = true
				break
			}
		}

		if !changedAtLeastOnce {
			t.Error("shuffle never changed the order in 20 attempts")
		}
	})

	t.Run("shuffle preserves all roles", func(t *testing.T) {
		t.Parallel()

		roles := []Role{
			RoleMerlin, RolePercival, RoleLoyalServant, RoleLoyalServant,
			RoleAssassin, RoleMorgana, RoleMordred,
		}

		err := secureShuffleRoles(roles)
		if err != nil {
			t.Fatalf("secureShuffleRoles failed: %v", err)
		}

		// Count roles
		roleCounts := make(map[Role]int)
		for _, role := range roles {
			roleCounts[role]++
		}

		// Verify expected counts
		if roleCounts[RoleMerlin] != 1 {
			t.Errorf("expected 1 Merlin, got %d", roleCounts[RoleMerlin])
		}
		if roleCounts[RoleLoyalServant] != 2 {
			t.Errorf("expected 2 Loyal Servants, got %d", roleCounts[RoleLoyalServant])
		}
		if roleCounts[RoleAssassin] != 1 {
			t.Errorf("expected 1 Assassin, got %d", roleCounts[RoleAssassin])
		}
	})

	t.Run("shuffle works with single element", func(t *testing.T) {
		t.Parallel()

		roles := []Role{RoleMerlin}

		err := secureShuffleRoles(roles)
		if err != nil {
			t.Fatalf("secureShuffleRoles failed: %v", err)
		}

		if len(roles) != 1 || roles[0] != RoleMerlin {
			t.Error("single element shuffle failed")
		}
	})

	t.Run("shuffle works with empty slice", func(t *testing.T) {
		t.Parallel()

		roles := []Role{}

		err := secureShuffleRoles(roles)
		if err != nil {
			t.Fatalf("secureShuffleRoles failed: %v", err)
		}

		if len(roles) != 0 {
			t.Error("empty slice shuffle failed")
		}
	})
}

func TestHasRole(t *testing.T) {
	t.Parallel()

	roles := map[string]Role{
		"p1": RoleMerlin,
		"p2": RoleAssassin,
		"p3": RoleLoyalServant,
	}

	tests := []struct {
		name string
		role Role
		want bool
	}{
		{"has Merlin", RoleMerlin, true},
		{"has Assassin", RoleAssassin, true},
		{"has Loyal Servant", RoleLoyalServant, true},
		{"does not have Morgana", RoleMorgana, false},
		{"does not have Percival", RolePercival, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := hasRole(roles, tt.role); got != tt.want {
				t.Errorf("hasRole(%s) = %v, want %v", tt.role, got, tt.want)
			}
		})
	}
}

func TestGetPlayerWithRole(t *testing.T) {
	t.Parallel()

	roles := map[string]Role{
		"p1": RoleMerlin,
		"p2": RoleAssassin,
		"p3": RoleLoyalServant,
	}

	tests := []struct {
		name string
		role Role
		want string
	}{
		{"find Merlin", RoleMerlin, "p1"},
		{"find Assassin", RoleAssassin, "p2"},
		{"find Loyal Servant", RoleLoyalServant, "p3"},
		{"role not found", RoleMorgana, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := getPlayerWithRole(roles, tt.role); got != tt.want {
				t.Errorf("getPlayerWithRole(%s) = %s, want %s", tt.role, got, tt.want)
			}
		})
	}
}

func TestCountTeamQuests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		results       []QuestResult
		wantGoodWins  int
		wantEvilWins  int
	}{
		{
			name:         "no quests completed",
			results:      []QuestResult{},
			wantGoodWins: 0,
			wantEvilWins: 0,
		},
		{
			name: "all good wins",
			results: []QuestResult{
				{Success: true},
				{Success: true},
				{Success: true},
			},
			wantGoodWins: 3,
			wantEvilWins: 0,
		},
		{
			name: "all evil wins",
			results: []QuestResult{
				{Success: false},
				{Success: false},
				{Success: false},
			},
			wantGoodWins: 0,
			wantEvilWins: 3,
		},
		{
			name: "mixed results",
			results: []QuestResult{
				{Success: true},
				{Success: false},
				{Success: true},
				{Success: false},
			},
			wantGoodWins: 2,
			wantEvilWins: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			goodWins, evilWins := countTeamQuests(tt.results)

			if goodWins != tt.wantGoodWins {
				t.Errorf("goodWins = %d, want %d", goodWins, tt.wantGoodWins)
			}
			if evilWins != tt.wantEvilWins {
				t.Errorf("evilWins = %d, want %d", evilWins, tt.wantEvilWins)
			}
		})
	}
}

func TestGetRoleDescription(t *testing.T) {
	t.Parallel()

	tests := []struct {
		role Role
	}{
		{RoleMerlin},
		{RolePercival},
		{RoleLoyalServant},
		{RoleAssassin},
		{RoleMorgana},
		{RoleMordred},
		{RoleOberon},
		{RoleMinionOfMordred},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			t.Parallel()

			desc := getRoleDescription(tt.role)

			// Just verify it returns something non-empty
			if desc == "" {
				t.Errorf("getRoleDescription(%s) returned empty string", tt.role)
			}
		})
	}
}

func TestGetExpectedTeamSizes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		playerCount int
		wantGood    int
		wantEvil    int
	}{
		{5, 3, 2},
		{6, 4, 2},
		{7, 4, 3},
		{8, 5, 3},
		{9, 6, 3},
		{10, 6, 4},
		{4, 0, 0},  // Invalid
		{11, 0, 0}, // Invalid
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.playerCount))+" players", func(t *testing.T) {
			t.Parallel()

			good, evil := getExpectedTeamSizes(tt.playerCount)

			if good != tt.wantGood {
				t.Errorf("good = %d, want %d", good, tt.wantGood)
			}
			if evil != tt.wantEvil {
				t.Errorf("evil = %d, want %d", evil, tt.wantEvil)
			}
		})
	}
}

