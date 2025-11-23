package avalon

import (
	"encoding/json"
	"testing"

	"github.com/KonradHerman/roundtable/internal/core"
)

func TestNewGame(t *testing.T) {
	t.Parallel()

	game := NewGame()

	if game == nil {
		t.Fatal("NewGame returned nil")
	}

	g := game.(*Game)
	if g.phase != "" && g.phase != PhaseSetup {
		t.Errorf("expected phase '' or 'setup', got '%s'", g.phase)
	}

	if g.roles == nil {
		t.Error("roles map is nil")
	}

	if g.teams == nil {
		t.Error("teams map is nil")
	}

	if g.knowledge == nil {
		t.Error("knowledge map is nil")
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
			name: "valid 5-player config with basic roles",
			config: &Config{
				Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant, RoleLoyalServant, RoleMinionOfMordred},
			},
			wantErr: false,
		},
		{
			name: "valid 7-player config",
			config: &Config{
				Roles: []Role{RoleMerlin, RolePercival, RoleLoyalServant, RoleLoyalServant, RoleAssassin, RoleMorgana, RoleMordred},
			},
			wantErr: false,
		},
		{
			name: "valid 10-player config",
			config: &Config{
				Roles: []Role{
					RoleMerlin, RolePercival, RoleLoyalServant, RoleLoyalServant, RoleLoyalServant, RoleLoyalServant,
					RoleAssassin, RoleMorgana, RoleMordred, RoleMinionOfMordred,
				},
			},
			wantErr: false,
		},
		{
			name: "fail when too few players",
			config: &Config{
				Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant},
			},
			wantErr:     true,
			errContains: "requires 5-10 players",
		},
		{
			name: "fail when too many players",
			config: &Config{
				Roles: []Role{
					RoleMerlin, RolePercival, RoleLoyalServant, RoleLoyalServant, RoleLoyalServant,
					RoleLoyalServant, RoleAssassin, RoleMorgana, RoleMordred, RoleMinionOfMordred, RoleMinionOfMordred,
				},
			},
			wantErr:     true,
			errContains: "requires 5-10 players",
		},
		{
			name: "fail when Merlin present without Assassin",
			config: &Config{
				Roles: []Role{RoleMerlin, RoleLoyalServant, RoleLoyalServant, RoleMinionOfMordred, RoleMinionOfMordred},
			},
			wantErr:     true,
			errContains: "Assassin is required when Merlin is present",
		},
		{
			name: "fail when Percival present without Merlin",
			config: &Config{
				Roles: []Role{RolePercival, RoleLoyalServant, RoleLoyalServant, RoleAssassin, RoleMinionOfMordred},
			},
			wantErr:     true,
			errContains: "Merlin is required when Percival is present",
		},
		{
			name: "fail with invalid team sizes for 5 players",
			config: &Config{
				Roles: []Role{RoleMerlin, RolePercival, RoleAssassin, RoleMorgana, RoleMordred},
			},
			wantErr:     true,
			errContains: "invalid team sizes",
		},
		{
			name: "fail with unknown role",
			config: &Config{
				Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant, Role("unknown"), RoleMinionOfMordred},
			},
			wantErr:     true,
			errContains: "unknown role",
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
			name: "successfully initialize 5-player game",
			config: &Config{
				Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant, RoleLoyalServant, RoleMinionOfMordred},
			},
			players: []*core.Player{
				{ID: "p1", DisplayName: "Player1"},
				{ID: "p2", DisplayName: "Player2"},
				{ID: "p3", DisplayName: "Player3"},
				{ID: "p4", DisplayName: "Player4"},
				{ID: "p5", DisplayName: "Player5"},
			},
			wantErr: false,
		},
		{
			name: "successfully initialize 7-player game",
			config: &Config{
				Roles: []Role{RoleMerlin, RolePercival, RoleLoyalServant, RoleLoyalServant, RoleAssassin, RoleMorgana, RoleMordred},
			},
			players: []*core.Player{
				{ID: "p1", DisplayName: "Player1"},
				{ID: "p2", DisplayName: "Player2"},
				{ID: "p3", DisplayName: "Player3"},
				{ID: "p4", DisplayName: "Player4"},
				{ID: "p5", DisplayName: "Player5"},
				{ID: "p6", DisplayName: "Player6"},
				{ID: "p7", DisplayName: "Player7"},
			},
			wantErr: false,
		},
		{
			name: "fail with player count mismatch (too few)",
			config: &Config{
				Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant, RoleLoyalServant, RoleMinionOfMordred},
			},
			players: []*core.Player{
				{ID: "p1", DisplayName: "Player1"},
				{ID: "p2", DisplayName: "Player2"},
			},
			wantErr:     true,
			errContains: "does not match role count",
		},
		{
			name: "fail with player count mismatch (too many)",
			config: &Config{
				Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant},
			},
			players: []*core.Player{
				{ID: "p1", DisplayName: "Player1"},
				{ID: "p2", DisplayName: "Player2"},
				{ID: "p3", DisplayName: "Player3"},
				{ID: "p4", DisplayName: "Player4"},
				{ID: "p5", DisplayName: "Player5"},
			},
			wantErr:     true,
			errContains: "requires 5-10 players",
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
				if len(g.roles) != len(tt.players) {
					t.Errorf("expected %d role assignments, got %d", len(tt.players), len(g.roles))
				}

				// Verify all players have teams
				if len(g.teams) != len(tt.players) {
					t.Errorf("expected %d team assignments, got %d", len(tt.players), len(g.teams))
				}

				// Verify phase is role_reveal
				if g.phase != PhaseRoleReveal {
					t.Errorf("expected phase 'role_reveal', got '%s'", g.phase)
				}

				// Verify leader was selected
				if g.currentLeader == "" {
					t.Error("no leader selected")
				}

				// Verify quest number is 1
				if g.questNumber != 1 {
					t.Errorf("expected quest number 1, got %d", g.questNumber)
				}
			}
		})
	}
}

func TestGame_ProcessAcknowledgeRole(t *testing.T) {
	t.Parallel()

	game := NewGame()
	config := &Config{
		Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant, RoleLoyalServant, RoleMinionOfMordred},
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
		t.Fatalf("failed to initialize game: %v", err)
	}

	g := game.(*Game)

	// Verify we're in role reveal phase
	if g.phase != PhaseRoleReveal {
		t.Fatalf("expected role_reveal phase, got %s", g.phase)
	}

	// Process acknowledgments for first 4 players
	for i := 0; i < 4; i++ {
		action := core.Action{Type: "acknowledge_role"}
		events, err := game.ProcessAction(players[i].ID, action)
		if err != nil {
			t.Errorf("player %d: unexpected error: %v", i, err)
		}
		if len(events) == 0 {
			t.Errorf("player %d: expected events", i)
		}

		// Should still be in role reveal
		if g.phase != PhaseRoleReveal {
			t.Errorf("after player %d: expected role_reveal phase, got %s", i, g.phase)
		}
	}

	// Process final acknowledgment
	action := core.Action{Type: "acknowledge_role"}
	events, err := game.ProcessAction(players[4].ID, action)
	if err != nil {
		t.Fatalf("player 5: unexpected error: %v", err)
	}
	if len(events) == 0 {
		t.Error("player 5: expected events")
	}

	// Should now be in team building
	if g.phase != PhaseTeamBuilding {
		t.Errorf("expected team_building phase, got %s", g.phase)
	}
}

func TestGame_ProcessProposeTeam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		teamMembers []string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid team proposal",
			teamMembers: []string{"p1", "p2"},
			wantErr:     false,
		},
		{
			name:        "team too small",
			teamMembers: []string{"p1"},
			wantErr:     true,
			errContains: "team size must be",
		},
		{
			name:        "team too large",
			teamMembers: []string{"p1", "p2", "p3"},
			wantErr:     true,
			errContains: "team size must be",
		},
		{
			name:        "invalid player ID",
			teamMembers: []string{"p1", "invalid"},
			wantErr:     true,
			errContains: "invalid player ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			game := NewGame()
			config := &Config{
				Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant, RoleLoyalServant, RoleMinionOfMordred},
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
				t.Fatalf("failed to initialize: %v", err)
			}

			// Acknowledge all roles to advance to team building
			g := game.(*Game)
			for _, p := range players {
				g.acknowledged[p.ID] = true
			}
			g.phase = PhaseTeamBuilding

			// Create payload
			payload := map[string]interface{}{
				"team_members": tt.teamMembers,
			}
			payloadBytes, _ := json.Marshal(payload)

			action := core.Action{
				Type:    "propose_team",
				Payload: payloadBytes,
			}

			events, err := game.ProcessAction(g.currentLeader, action)

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
				if len(events) == 0 {
					t.Error("expected events")
				}

				// Should advance to team voting
				if g.phase != PhaseTeamVoting {
					t.Errorf("expected team_voting phase, got %s", g.phase)
				}

				// Verify proposed team
				if len(g.proposedTeam) != len(tt.teamMembers) {
					t.Errorf("expected %d team members, got %d", len(tt.teamMembers), len(g.proposedTeam))
				}
			}
		})
	}
}

func TestGame_ProcessVoteTeam(t *testing.T) {
	t.Parallel()

	t.Run("team approved", func(t *testing.T) {
		t.Parallel()

		game := NewGame()
		config := &Config{
			Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant, RoleLoyalServant, RoleMinionOfMordred},
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
			t.Fatalf("failed to initialize: %v", err)
		}

		g := game.(*Game)
		// Setup for team voting
		g.phase = PhaseTeamVoting
		g.proposedTeam = []string{"p1", "p2"}

		// All players vote approve
		for _, p := range players {
			payload := map[string]interface{}{"vote": "approve"}
			payloadBytes, _ := json.Marshal(payload)
			action := core.Action{Type: "vote_team", Payload: payloadBytes}

			events, err := game.ProcessAction(p.ID, action)
			if err != nil {
				t.Errorf("player %s: unexpected error: %v", p.ID, err)
			}
			if len(events) == 0 {
				t.Errorf("player %s: expected events", p.ID)
			}
		}

		// Should advance to quest execution
		if g.phase != PhaseQuestExec {
			t.Errorf("expected quest_execution phase, got %s", g.phase)
		}
	})

	t.Run("team rejected", func(t *testing.T) {
		t.Parallel()

		game := NewGame()
		config := &Config{
			Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant, RoleLoyalServant, RoleMinionOfMordred},
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
			t.Fatalf("failed to initialize: %v", err)
		}

		g := game.(*Game)
		originalLeader := g.currentLeader
		g.phase = PhaseTeamVoting
		g.proposedTeam = []string{"p1", "p2"}

		// All players vote reject
		for _, p := range players {
			payload := map[string]interface{}{"vote": "reject"}
			payloadBytes, _ := json.Marshal(payload)
			action := core.Action{Type: "vote_team", Payload: payloadBytes}

			_, err := game.ProcessAction(p.ID, action)
			if err != nil {
				t.Errorf("player %s: unexpected error: %v", p.ID, err)
			}
		}

		// Should return to team building with new leader
		if g.phase != PhaseTeamBuilding {
			t.Errorf("expected team_building phase, got %s", g.phase)
		}

		// Leader should have rotated
		if g.currentLeader == originalLeader {
			t.Error("leader should have rotated")
		}

		// Rejection count should increment
		if g.rejectionCount != 1 {
			t.Errorf("expected rejection count 1, got %d", g.rejectionCount)
		}
	})

	t.Run("five consecutive rejections - evil wins", func(t *testing.T) {
		t.Parallel()

		game := NewGame()
		config := &Config{
			Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant, RoleLoyalServant, RoleMinionOfMordred},
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
			t.Fatalf("failed to initialize: %v", err)
		}

		g := game.(*Game)
		g.phase = PhaseTeamVoting
		g.proposedTeam = []string{"p1", "p2"}
		g.rejectionCount = 4 // Already at 4 rejections

		// All players vote reject (5th rejection)
		for _, p := range players {
			payload := map[string]interface{}{"vote": "reject"}
			payloadBytes, _ := json.Marshal(payload)
			action := core.Action{Type: "vote_team", Payload: payloadBytes}

			_, err := game.ProcessAction(p.ID, action)
			if err != nil {
				t.Errorf("player %s: unexpected error: %v", p.ID, err)
			}
		}

		// Game should be finished with evil victory
		if g.phase != PhaseFinished {
			t.Errorf("expected finished phase, got %s", g.phase)
		}

		if g.winningTeam != TeamEvil {
			t.Errorf("expected evil to win, got %s", g.winningTeam)
		}

		if g.winReason != "five_consecutive_rejections" {
			t.Errorf("expected win reason 'five_consecutive_rejections', got '%s'", g.winReason)
		}
	})
}

func TestGame_ProcessPlayQuestCard(t *testing.T) {
	t.Parallel()

	t.Run("good player can only play success", func(t *testing.T) {
		t.Parallel()

		game := NewGame()
		config := &Config{
			Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant, RoleLoyalServant, RoleMinionOfMordred},
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
			t.Fatalf("failed to initialize: %v", err)
		}

		g := game.(*Game)

		// Find a good player
		var goodPlayerID string
		for pid, team := range g.teams {
			if team == TeamGood {
				goodPlayerID = pid
				break
			}
		}

		if goodPlayerID == "" {
			t.Fatal("no good player found")
		}

		// Setup for quest execution with good player on team
		g.phase = PhaseQuestExec
		g.proposedTeam = []string{goodPlayerID}

		// Try to play fail card as good player
		payload := map[string]interface{}{"card": "fail"}
		payloadBytes, _ := json.Marshal(payload)
		action := core.Action{Type: "play_quest_card", Payload: payloadBytes}

		_, err = game.ProcessAction(goodPlayerID, action)
		if err == nil {
			t.Error("expected error when good player plays fail card")
		}
		if !contains(err.Error(), "can only play success") {
			t.Errorf("expected error about success-only, got: %v", err)
		}
	})

	t.Run("quest succeeds with all success cards", func(t *testing.T) {
		t.Parallel()

		game := NewGame()
		config := &Config{
			Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant, RoleLoyalServant, RoleMinionOfMordred},
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
			t.Fatalf("failed to initialize: %v", err)
		}

		g := game.(*Game)
		g.phase = PhaseQuestExec
		g.proposedTeam = []string{"p1", "p2"}

		// Both players play success
		for _, pid := range g.proposedTeam {
			payload := map[string]interface{}{"card": "success"}
			payloadBytes, _ := json.Marshal(payload)
			action := core.Action{Type: "play_quest_card", Payload: payloadBytes}

			_, err := game.ProcessAction(pid, action)
			if err != nil {
				t.Errorf("player %s: unexpected error: %v", pid, err)
			}
		}

		// Quest should succeed
		if len(g.questResults) != 1 {
			t.Fatalf("expected 1 quest result, got %d", len(g.questResults))
		}

		if !g.questResults[0].Success {
			t.Error("expected quest to succeed")
		}

		// Should advance to next quest
		if g.questNumber != 2 {
			t.Errorf("expected quest number 2, got %d", g.questNumber)
		}
	})
}

func TestGame_ProcessAssassinate(t *testing.T) {
	t.Parallel()

	t.Run("assassin finds Merlin - evil wins", func(t *testing.T) {
		t.Parallel()

		game := NewGame()
		config := &Config{
			Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant, RoleLoyalServant, RoleMinionOfMordred},
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
			t.Fatalf("failed to initialize: %v", err)
		}

		g := game.(*Game)

		// Find assassin and Merlin
		var assassinID, merlinID string
		for pid, role := range g.roles {
			if role == RoleAssassin {
				assassinID = pid
			}
			if role == RoleMerlin {
				merlinID = pid
			}
		}

		if assassinID == "" || merlinID == "" {
			t.Fatal("couldn't find assassin or Merlin")
		}

		// Setup assassination phase
		g.phase = PhaseAssassination

		// Assassin targets Merlin
		payload := map[string]interface{}{"target_id": merlinID}
		payloadBytes, _ := json.Marshal(payload)
		action := core.Action{Type: "assassinate", Payload: payloadBytes}

		events, err := game.ProcessAction(assassinID, action)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(events) == 0 {
			t.Error("expected events")
		}

		// Evil should win
		if g.winningTeam != TeamEvil {
			t.Errorf("expected evil to win, got %s", g.winningTeam)
		}

		if g.winReason != "assassin_found_merlin" {
			t.Errorf("expected win reason 'assassin_found_merlin', got '%s'", g.winReason)
		}

		if g.phase != PhaseFinished {
			t.Errorf("expected finished phase, got %s", g.phase)
		}
	})

	t.Run("assassin misses Merlin - good wins", func(t *testing.T) {
		t.Parallel()

		game := NewGame()
		config := &Config{
			Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant, RoleLoyalServant, RoleMinionOfMordred},
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
			t.Fatalf("failed to initialize: %v", err)
		}

		g := game.(*Game)

		// Find assassin and a non-Merlin player
		var assassinID, targetID string
		for pid, role := range g.roles {
			if role == RoleAssassin {
				assassinID = pid
			}
			if role != RoleAssassin && role != RoleMerlin && targetID == "" {
				targetID = pid
			}
		}

		if assassinID == "" || targetID == "" {
			t.Fatal("couldn't find assassin or target")
		}

		// Setup assassination phase
		g.phase = PhaseAssassination

		// Assassin targets wrong player
		payload := map[string]interface{}{"target_id": targetID}
		payloadBytes, _ := json.Marshal(payload)
		action := core.Action{Type: "assassinate", Payload: payloadBytes}

		_, err = game.ProcessAction(assassinID, action)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Good should win
		if g.winningTeam != TeamGood {
			t.Errorf("expected good to win, got %s", g.winningTeam)
		}

		if g.winReason != "assassin_failed" {
			t.Errorf("expected win reason 'assassin_failed', got '%s'", g.winReason)
		}
	})
}

func TestGame_GetPlayerState(t *testing.T) {
	t.Parallel()

	game := NewGame()
	config := &Config{
		Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant, RoleLoyalServant, RoleMinionOfMordred},
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
	if playerState.Role == "" {
		t.Error("role is empty")
	}
	if playerState.Team == "" {
		t.Error("team is empty")
	}
	if playerState.QuestNumber != 1 {
		t.Errorf("expected quest number 1, got %d", playerState.QuestNumber)
	}
}

func TestGame_GetPublicState(t *testing.T) {
	t.Parallel()

	game := NewGame()
	config := &Config{
		Roles: []Role{RoleMerlin, RoleAssassin, RoleLoyalServant, RoleLoyalServant, RoleMinionOfMordred},
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

	if publicState.PlayerCount != 5 {
		t.Errorf("expected playerCount 5, got %d", publicState.PlayerCount)
	}

	if publicState.QuestNumber != 1 {
		t.Errorf("expected quest number 1, got %d", publicState.QuestNumber)
	}

	if publicState.CurrentLeaderID == "" {
		t.Error("current leader ID is empty")
	}
}

func TestGame_GetPhase(t *testing.T) {
	t.Parallel()

	game := NewGame()
	g := game.(*Game)

	tests := []struct {
		name       string
		setupPhase GamePhase
		want       string
	}{
		{"setup phase", PhaseSetup, "setup"},
		{"role reveal phase", PhaseRoleReveal, "role_reveal"},
		{"team building phase", PhaseTeamBuilding, "team_building"},
		{"team voting phase", PhaseTeamVoting, "team_voting"},
		{"quest execution phase", PhaseQuestExec, "quest_execution"},
		{"assassination phase", PhaseAssassination, "assassination"},
		{"finished phase", PhaseFinished, "finished"},
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
		phase GamePhase
		want  bool
	}{
		{"not finished in setup", PhaseSetup, false},
		{"not finished in role reveal", PhaseRoleReveal, false},
		{"not finished in team building", PhaseTeamBuilding, false},
		{"not finished in team voting", PhaseTeamVoting, false},
		{"not finished in quest execution", PhaseQuestExec, false},
		{"not finished in assassination", PhaseAssassination, false},
		{"finished", PhaseFinished, true},
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
			Roles: []Role{
				RoleMerlin, RolePercival, RoleLoyalServant, RoleLoyalServant,
				RoleAssassin, RoleMorgana, RoleMordred,
			},
		}
		players := []*core.Player{
			{ID: "p1", DisplayName: "Player1"},
			{ID: "p2", DisplayName: "Player2"},
			{ID: "p3", DisplayName: "Player3"},
			{ID: "p4", DisplayName: "Player4"},
			{ID: "p5", DisplayName: "Player5"},
			{ID: "p6", DisplayName: "Player6"},
			{ID: "p7", DisplayName: "Player7"},
		}

		_, err := game.Initialize(config, players)
		if err != nil {
			t.Fatalf("run %d: failed to initialize game: %v", run, err)
		}

		g := game.(*Game)

		// Track all assigned roles
		assignedRoles := make(map[Role]int)
		for _, role := range g.roles {
			assignedRoles[role]++
		}

		// Verify each role from config appears exactly once
		expectedRoles := make(map[Role]int)
		for _, role := range config.Roles {
			expectedRoles[role]++
		}

		for role, expectedCount := range expectedRoles {
			if actualCount, exists := assignedRoles[role]; !exists || actualCount != expectedCount {
				t.Errorf("run %d: role %s expected %d times, got %d", run, role, expectedCount, actualCount)
			}
		}

		// Verify all players have different roles (since we have unique roles)
		roleSet := make(map[Role]bool)
		for _, role := range g.roles {
			if roleSet[role] && expectedRoles[role] == 1 {
				t.Errorf("run %d: duplicate role %s found", run, role)
			}
			roleSet[role] = true
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

