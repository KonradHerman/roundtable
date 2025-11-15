package werewolf

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/yourusername/roundtable/internal/core"
)

// Phase represents the current game phase.
type Phase string

const (
	PhaseSetup   Phase = "setup"
	PhaseNight   Phase = "night"
	PhaseDay     Phase = "day"
	PhaseResults Phase = "results"
)

// Game implements the core.Game interface for One Night Werewolf.
type Game struct {
	config            *Config
	players           map[string]*core.Player // playerID → player
	roleAssignments   map[string]RoleType     // playerID → role
	originalRoles     map[string]RoleType     // playerID → original role (before night actions)
	votes             map[string]string       // voterID → targetID
	phase             Phase
	phaseStartedAt    time.Time
	phaseEndsAt       time.Time
	nightActionsComplete map[RoleType]bool    // Track which roles have acted
}

// NewGame creates a new werewolf game instance.
func NewGame() core.Game {
	return &Game{
		players:           make(map[string]*core.Player),
		roleAssignments:   make(map[string]RoleType),
		originalRoles:     make(map[string]RoleType),
		votes:             make(map[string]string),
		nightActionsComplete: make(map[RoleType]bool),
		phase:             PhaseSetup,
	}
}

// Initialize sets up the game with players and config.
func (g *Game) Initialize(config core.GameConfig, players []*core.Player) ([]core.GameEvent, error) {
	wConfig, ok := config.(*Config)
	if !ok {
		return nil, errors.New("invalid config type")
	}

	if err := wConfig.Validate(); err != nil {
		return nil, err
	}

	if len(players) != len(wConfig.Roles) {
		return nil, fmt.Errorf("player count (%d) must match role count (%d)", len(players), len(wConfig.Roles))
	}

	g.config = wConfig
	for _, player := range players {
		g.players[player.ID] = player
	}

	// Shuffle and assign roles
	shuffledRoles := make([]RoleType, len(wConfig.Roles))
	copy(shuffledRoles, wConfig.Roles)
	rand.Shuffle(len(shuffledRoles), func(i, j int) {
		shuffledRoles[i], shuffledRoles[j] = shuffledRoles[j], shuffledRoles[i]
	})

	events := make([]core.GameEvent, 0)

	// Game started event (public)
	gameStartedEvent, _ := core.NewPublicEvent(core.EventGameStarted, "system", core.GameStartedPayload{
		GameType:  "werewolf",
		Config:    config,
		PlayerIDs: getPlayerIDs(players),
	})
	events = append(events, gameStartedEvent)

	// Assign roles (private events)
	for i, player := range players {
		role := shuffledRoles[i]
		g.roleAssignments[player.ID] = role
		g.originalRoles[player.ID] = role

		roleEvent, _ := core.NewPrivateEvent("role_assigned", "system", RoleAssignedPayload{
			Role: role,
		}, []string{player.ID})
		events = append(events, roleEvent)
	}

	// Werewolf wakeup (show other werewolves to each werewolf)
	werewolfIDs := g.getPlayersByRole(RoleWerewolf)
	for _, werewolfID := range werewolfIDs {
		otherWerewolves := make([]string, 0)
		for _, wid := range werewolfIDs {
			if wid != werewolfID {
				otherWerewolves = append(otherWerewolves, wid)
			}
		}

		wakeupEvent, _ := core.NewPrivateEvent("werewolf_wakeup", "system", WerewolfWakeupPayload{
			OtherWerewolves: otherWerewolves,
		}, []string{werewolfID})
		events = append(events, wakeupEvent)
	}

	// Masons wakeup (show other masons)
	masonIDs := g.getPlayersByRole(RoleMason)
	for _, masonID := range masonIDs {
		otherMasons := make([]string, 0)
		for _, mid := range masonIDs {
			if mid != masonID {
				otherMasons = append(otherMasons, mid)
			}
		}

		masonEvent, _ := core.NewPrivateEvent("mason_wakeup", "system", MasonWakeupPayload{
			OtherMasons: otherMasons,
		}, []string{masonID})
		events = append(events, masonEvent)
	}

	// Start night phase
	g.phase = PhaseNight
	g.phaseStartedAt = time.Now()
	g.phaseEndsAt = g.phaseStartedAt.Add(g.config.NightDuration)

	phaseEvent, _ := core.NewPublicEvent(core.EventPhaseChanged, "system", core.PhaseChangedPayload{
		Phase: core.GamePhase{
			Name:    string(PhaseNight),
			EndsAt:  &g.phaseEndsAt,
			Message: "Night phase - roles act in secret",
		},
	})
	events = append(events, phaseEvent)

	return events, nil
}

// ValidateAction checks if a player can perform an action.
func (g *Game) ValidateAction(playerID string, action core.Action) error {
	role, exists := g.roleAssignments[playerID]
	if !exists {
		return errors.New("player not in game")
	}

	switch action.Type {
	case "vote":
		if g.phase != PhaseDay {
			return errors.New("can only vote during day phase")
		}
		if _, hasVoted := g.votes[playerID]; hasVoted {
			return errors.New("already voted")
		}
		return nil

	case "seer_view":
		if role != RoleSeer {
			return errors.New("only seer can view roles")
		}
		if g.phase != PhaseNight {
			return errors.New("can only view during night")
		}
		if g.nightActionsComplete[RoleSeer] {
			return errors.New("seer has already acted")
		}
		return nil

	// Add more action validations for other roles...

	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

// ProcessAction executes a validated action.
func (g *Game) ProcessAction(playerID string, action core.Action) ([]core.GameEvent, error) {
	events := make([]core.GameEvent, 0)

	switch action.Type {
	case "vote":
		var votePayload VotePayload
		if err := json.Unmarshal(action.Payload, &votePayload); err != nil {
			return nil, err
		}

		g.votes[playerID] = votePayload.TargetID

		// Create vote cast event (public, but target hidden until all vote)
		voteEvent, _ := core.NewPublicEvent("vote_cast", playerID, VoteCastPayload{
			VoterID: playerID,
		})
		events = append(events, voteEvent)

		// Check if everyone has voted
		if len(g.votes) == len(g.players) {
			// Reveal votes
			voteRevealEvent, _ := core.NewPublicEvent("votes_revealed", "system", VotesRevealedPayload{
				Votes: g.votes,
			})
			events = append(events, voteRevealEvent)

			// Calculate results
			results := g.calculateResults()
			resultsEvent, _ := core.NewPublicEvent(core.EventGameFinished, "system", core.GameFinishedPayload{
				Results: results,
			})
			events = append(events, resultsEvent)

			g.phase = PhaseResults
		}

	case "seer_view":
		var seerPayload SeerViewPayload
		if err := json.Unmarshal(action.Payload, &seerPayload); err != nil {
			return nil, err
		}

		targetRole := g.roleAssignments[seerPayload.TargetID]

		// Send role to seer (private)
		seerResultEvent, _ := core.NewPrivateEvent("seer_result", "system", SeerResultPayload{
			TargetID: seerPayload.TargetID,
			Role:     targetRole,
		}, []string{playerID})
		events = append(events, seerResultEvent)

		g.nightActionsComplete[RoleSeer] = true

	// Add more action implementations for other roles...

	default:
		return nil, fmt.Errorf("unknown action type: %s", action.Type)
	}

	return events, nil
}

// GetPlayerState returns the state visible to a specific player.
func (g *Game) GetPlayerState(playerID string) core.PlayerState {
	role := g.roleAssignments[playerID]

	state := PlayerState{
		Phase:       string(g.phase),
		PhaseEndsAt: g.phaseEndsAt,
		YourRole:    role,
		HasVoted:    g.votes[playerID] != "",
	}

	// Add role-specific information
	if g.phase == PhaseNight {
		switch role {
		case RoleWerewolf:
			state.OtherWerewolves = g.getPlayersByRole(RoleWerewolf)
			// Remove self
			for i, id := range state.OtherWerewolves {
				if id == playerID {
					state.OtherWerewolves = append(state.OtherWerewolves[:i], state.OtherWerewolves[i+1:]...)
					break
				}
			}
		case RoleMason:
			state.OtherMasons = g.getPlayersByRole(RoleMason)
			for i, id := range state.OtherMasons {
				if id == playerID {
					state.OtherMasons = append(state.OtherMasons[:i], state.OtherMasons[i+1:]...)
					break
				}
			}
		}
	}

	return state
}

// GetPublicState returns the state visible to all players and spectators.
func (g *Game) GetPublicState() core.PublicState {
	return PublicState{
		Phase:          string(g.phase),
		PhaseEndsAt:    g.phaseEndsAt,
		PlayerCount:    len(g.players),
		VotesSubmitted: len(g.votes),
	}
}

// GetPhase returns the current game phase.
func (g *Game) GetPhase() core.GamePhase {
	return core.GamePhase{
		Name:   string(g.phase),
		EndsAt: &g.phaseEndsAt,
	}
}

// IsFinished returns true if the game has concluded.
func (g *Game) IsFinished() bool {
	return g.phase == PhaseResults
}

// GetResults returns the final game results.
func (g *Game) GetResults() core.GameResults {
	return g.calculateResults()
}

// calculateResults determines the winner based on votes.
func (g *Game) calculateResults() core.GameResults {
	// Count votes
	voteCounts := make(map[string]int)
	for _, targetID := range g.votes {
		voteCounts[targetID]++
	}

	// Find player(s) with most votes
	maxVotes := 0
	eliminated := make([]string, 0)
	for playerID, count := range voteCounts {
		if count > maxVotes {
			maxVotes = count
			eliminated = []string{playerID}
		} else if count == maxVotes {
			eliminated = append(eliminated, playerID)
		}
	}

	// Determine winners based on roles and eliminations
	werewolfDied := false
	for _, playerID := range eliminated {
		role := g.roleAssignments[playerID]
		if role.IsWerewolfTeam() {
			werewolfDied = true
		}
	}

	winners := make([]string, 0)
	var winReason string

	if werewolfDied {
		// Village team wins
		for playerID, role := range g.roleAssignments {
			if role.IsVillageTeam() {
				winners = append(winners, playerID)
			}
		}
		winReason = "Village team eliminated a werewolf!"
	} else {
		// Werewolf team wins
		for playerID, role := range g.roleAssignments {
			if role.IsWerewolfTeam() {
				winners = append(winners, playerID)
			}
		}
		winReason = "Werewolf team survived!"
	}

	// Check for tanner win (tanner wins alone if they die)
	for _, playerID := range eliminated {
		if g.roleAssignments[playerID] == RoleTanner {
			winners = []string{playerID}
			winReason = "Tanner wins by getting eliminated!"
			break
		}
	}

	return core.GameResults{
		Winners:   winners,
		WinReason: winReason,
		FinalState: map[string]interface{}{
			"votes":      g.votes,
			"eliminated": eliminated,
			"roles":      g.roleAssignments,
		},
	}
}

// Helper functions

func (g *Game) getPlayersByRole(role RoleType) []string {
	players := make([]string, 0)
	for playerID, playerRole := range g.roleAssignments {
		if playerRole == role {
			players = append(players, playerID)
		}
	}
	return players
}

func getPlayerIDs(players []*core.Player) []string {
	ids := make([]string, len(players))
	for i, player := range players {
		ids[i] = player.ID
	}
	return ids
}
