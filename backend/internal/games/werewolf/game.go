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
	PhaseSetup      Phase = "setup"
	PhaseRoleReveal Phase = "role_reveal"
	PhaseNight      Phase = "night"
	PhaseDay        Phase = "day"
	PhaseResults    Phase = "results"
)

// Game implements the core.Game interface for One Night Werewolf.
type Game struct {
	config               *Config
	hostID               string                  // PlayerID of the host
	players              map[string]*core.Player // playerID → player
	roleAssignments      map[string]RoleType     // playerID → role
	originalRoles        map[string]RoleType     // playerID → original role (before night actions)
	centerCards          []RoleType              // 3 cards in the center
	roleAcknowledgements map[string]bool         // playerID → acknowledged
	votes                map[string]string       // voterID → targetID
	phase                Phase
	phaseStartedAt       time.Time
	phaseEndsAt          time.Time
	timerActive          bool              // Whether day phase timer is active
	nightActionsComplete map[RoleType]bool // Track which roles have acted
}

// NewGame creates a new werewolf game instance.
func NewGame() core.Game {
	return &Game{
		players:              make(map[string]*core.Player),
		roleAssignments:      make(map[string]RoleType),
		originalRoles:        make(map[string]RoleType),
		roleAcknowledgements: make(map[string]bool),
		votes:                make(map[string]string),
		nightActionsComplete: make(map[RoleType]bool),
		centerCards:          make([]RoleType, 0, 3),
		phase:                PhaseSetup,
		timerActive:          false,
	}
}

// SetHost sets the host player ID (called by room after initialization).
func (g *Game) SetHost(hostID string) {
	g.hostID = hostID
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

	// One Night Werewolf rule: must have exactly 3 more roles than players (center cards)
	expectedRoles := len(players) + 3
	if len(wConfig.Roles) != expectedRoles {
		return nil, fmt.Errorf("role count (%d) must be player count + 3 (%d)", len(wConfig.Roles), expectedRoles)
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

	// Assign roles to players (first N roles go to players)
	for i, player := range players {
		role := shuffledRoles[i]
		g.roleAssignments[player.ID] = role
		g.originalRoles[player.ID] = role

		roleEvent, _ := core.NewPrivateEvent("role_assigned", "system", RoleAssignedPayload{
			Role: role,
		}, []string{player.ID})
		events = append(events, roleEvent)
	}

	// The remaining 3 roles are "center cards" (not assigned to players)
	g.centerCards = shuffledRoles[len(players):]

	// Start role reveal phase (players need to acknowledge their roles)
	g.phase = PhaseRoleReveal
	g.phaseStartedAt = time.Now()

	phaseEvent, _ := core.NewPublicEvent(core.EventPhaseChanged, "system", core.PhaseChangedPayload{
		Phase: core.GamePhase{
			Name:    string(PhaseRoleReveal),
			Message: "Look at your role card and acknowledge",
		},
	})
	events = append(events, phaseEvent)

	return events, nil
}

// ValidateAction checks if a player can perform an action.
func (g *Game) ValidateAction(playerID string, action core.Action) error {
	_, exists := g.roleAssignments[playerID]
	if !exists {
		return errors.New("player not in game")
	}

	switch action.Type {
	case "acknowledge_role":
		if g.phase != PhaseRoleReveal {
			return errors.New("can only acknowledge role during role reveal phase")
		}
		if g.roleAcknowledgements[playerID] {
			return errors.New("already acknowledged")
		}
		return nil

	case "advance_phase":
		// Only host can advance phases - we'll need to check this at room level
		// For now, allow any player (room handler will check host status)
		return nil

	case "advance_to_results":
		if g.phase != PhaseDay {
			return errors.New("can only advance to results from day phase")
		}
		return nil

	case "toggle_timer":
		if g.phase != PhaseDay {
			return errors.New("can only toggle timer during day phase")
		}
		return nil

	case "extend_timer":
		if g.phase != PhaseDay {
			return errors.New("can only extend timer during day phase")
		}
		if !g.timerActive {
			return errors.New("timer is not active")
		}
		return nil

	case "vote":
		if g.phase != PhaseDay {
			return errors.New("can only vote during day phase")
		}
		// Allow vote changes - don't check if already voted
		return nil

	// Night actions
	case "werewolf_view_center":
		if g.phase != PhaseNight {
			return errors.New("can only perform night actions during night phase")
		}
		role := g.roleAssignments[playerID]
		if role != RoleWerewolf {
			return errors.New("only werewolves can view center cards")
		}
		// Check if they're the only werewolf
		werewolfCount := 0
		for _, r := range g.roleAssignments {
			if r == RoleWerewolf {
				werewolfCount++
			}
		}
		if werewolfCount != 1 {
			return errors.New("can only view center card if you are the only werewolf")
		}
		if g.nightActionsComplete[RoleWerewolf] {
			return errors.New("werewolf has already acted")
		}
		return nil

	case "seer_view_player":
		if g.phase != PhaseNight {
			return errors.New("can only perform night actions during night phase")
		}
		if g.roleAssignments[playerID] != RoleSeer {
			return errors.New("only seer can view player roles")
		}
		if g.nightActionsComplete[RoleSeer] {
			return errors.New("seer has already acted")
		}
		return nil

	case "seer_view_center":
		if g.phase != PhaseNight {
			return errors.New("can only perform night actions during night phase")
		}
		if g.roleAssignments[playerID] != RoleSeer {
			return errors.New("only seer can view center cards")
		}
		if g.nightActionsComplete[RoleSeer] {
			return errors.New("seer has already acted")
		}
		return nil

	case "robber_swap":
		if g.phase != PhaseNight {
			return errors.New("can only perform night actions during night phase")
		}
		if g.roleAssignments[playerID] != RoleRobber {
			return errors.New("only robber can swap roles")
		}
		if g.nightActionsComplete[RoleRobber] {
			return errors.New("robber has already acted")
		}
		return nil

	case "troublemaker_swap":
		if g.phase != PhaseNight {
			return errors.New("can only perform night actions during night phase")
		}
		if g.roleAssignments[playerID] != RoleTroublemaker {
			return errors.New("only troublemaker can swap players")
		}
		if g.nightActionsComplete[RoleTroublemaker] {
			return errors.New("troublemaker has already acted")
		}
		return nil

	case "drunk_swap":
		if g.phase != PhaseNight {
			return errors.New("can only perform night actions during night phase")
		}
		if g.roleAssignments[playerID] != RoleDrunk {
			return errors.New("only drunk can swap with center")
		}
		if g.nightActionsComplete[RoleDrunk] {
			return errors.New("drunk has already acted")
		}
		return nil

	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

// ProcessAction executes a validated action.
func (g *Game) ProcessAction(playerID string, action core.Action) ([]core.GameEvent, error) {
	events := make([]core.GameEvent, 0)

	switch action.Type {
	case "acknowledge_role":
		g.roleAcknowledgements[playerID] = true

		// Broadcast acknowledgement count
		ackEvent, _ := core.NewPublicEvent("role_acknowledged", "system", RoleAcknowledgedPayload{
			PlayerID: playerID,
			Count:    len(g.roleAcknowledgements),
			Total:    len(g.players),
		})
		events = append(events, ackEvent)

		// Auto-advance to night when all players acknowledged
		if len(g.roleAcknowledgements) == len(g.players) {
			nightEvents, err := g.AdvanceToNight()
			if err != nil {
				return events, err
			}
			events = append(events, nightEvents...)
		}

	case "advance_phase":
		// Advance from night to day (host only - checked at handler level)
		if g.phase == PhaseNight {
			dayEvents, err := g.AdvanceToDay()
			if err != nil {
				return nil, err
			}
			events = append(events, dayEvents...)
		}

	case "advance_to_results":
		// Advance to results phase to show all final roles
		g.phase = PhaseResults

		phaseEvent, _ := core.NewPublicEvent(core.EventPhaseChanged, "system", core.PhaseChangedPayload{
			Phase: core.GamePhase{
				Name:    string(PhaseResults),
				Message: "Role reveal - see everyone's final roles!",
			},
		})
		events = append(events, phaseEvent)

	case "toggle_timer":
		var timerPayload struct {
			Enable   bool `json:"enable"`
			Duration int  `json:"duration"` // seconds
		}
		if err := json.Unmarshal(action.Payload, &timerPayload); err != nil {
			return nil, err
		}

		duration := time.Duration(timerPayload.Duration) * time.Second
		if duration == 0 {
			duration = 3 * time.Minute // Default 3 minutes
		}

		timerEvents, err := g.ToggleTimer(timerPayload.Enable, duration)
		if err != nil {
			return nil, err
		}
		events = append(events, timerEvents...)

	case "extend_timer":
		var extendPayload struct {
			Seconds int `json:"seconds"`
		}
		if err := json.Unmarshal(action.Payload, &extendPayload); err != nil {
			return nil, err
		}

		if extendPayload.Seconds == 0 {
			extendPayload.Seconds = 60 // Default 1 minute
		}

		extendEvents, err := g.ExtendTimer(extendPayload.Seconds)
		if err != nil {
			return nil, err
		}
		events = append(events, extendEvents...)

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

	// Night actions
	case "werewolf_view_center":
		var payload WerewolfViewCenterPayload
		if err := json.Unmarshal(action.Payload, &payload); err != nil {
			return nil, err
		}

		if payload.CenterIndex < 0 || payload.CenterIndex >= len(g.centerCards) {
			return nil, errors.New("invalid center card index")
		}

		g.nightActionsComplete[RoleWerewolf] = true

		// Send private result to werewolf
		resultEvent, _ := core.NewPrivateEvent("werewolf_view_center_result", "system", WerewolfViewCenterResultPayload{
			CenterIndex: payload.CenterIndex,
			Role:        g.centerCards[payload.CenterIndex],
		}, []string{playerID})
		events = append(events, resultEvent)

	case "seer_view_player":
		var payload SeerViewPayload
		if err := json.Unmarshal(action.Payload, &payload); err != nil {
			return nil, err
		}

		// Validate target exists and is not the seer
		targetRole, exists := g.roleAssignments[payload.TargetID]
		if !exists {
			return nil, errors.New("target player not found")
		}
		if payload.TargetID == playerID {
			return nil, errors.New("seer cannot view their own card")
		}

		g.nightActionsComplete[RoleSeer] = true

		// Send private result to seer
		resultEvent, _ := core.NewPrivateEvent("seer_result", "system", SeerResultPayload{
			TargetID: payload.TargetID,
			Role:     targetRole,
		}, []string{playerID})
		events = append(events, resultEvent)

	case "seer_view_center":
		var payload SeerViewCenterPayload
		if err := json.Unmarshal(action.Payload, &payload); err != nil {
			return nil, err
		}

		if len(payload.CenterIndices) != 2 {
			return nil, errors.New("seer must view exactly 2 center cards")
		}

		for _, idx := range payload.CenterIndices {
			if idx < 0 || idx >= len(g.centerCards) {
				return nil, errors.New("invalid center card index")
			}
		}

		g.nightActionsComplete[RoleSeer] = true

		// Build result
		resultPayload := SeerViewCenterResultPayload{
			Cards: make([]struct {
				Index int      `json:"index"`
				Role  RoleType `json:"role"`
			}, 2),
		}
		for i, idx := range payload.CenterIndices {
			resultPayload.Cards[i].Index = idx
			resultPayload.Cards[i].Role = g.centerCards[idx]
		}

		// Send private result to seer
		resultEvent, _ := core.NewPrivateEvent("seer_center_result", "system", resultPayload, []string{playerID})
		events = append(events, resultEvent)

	case "robber_swap":
		var payload RobberSwapPayload
		if err := json.Unmarshal(action.Payload, &payload); err != nil {
			return nil, err
		}

		// Validate target exists and is not the robber
		if payload.TargetID == playerID {
			return nil, errors.New("robber cannot swap with themselves")
		}
		if _, exists := g.roleAssignments[payload.TargetID]; !exists {
			return nil, errors.New("target player not found")
		}

		// Perform the swap
		robberRole := g.roleAssignments[playerID]
		targetRole := g.roleAssignments[payload.TargetID]
		g.roleAssignments[playerID] = targetRole
		g.roleAssignments[payload.TargetID] = robberRole

		g.nightActionsComplete[RoleRobber] = true

		// Send private result to robber (their new role)
		resultEvent, _ := core.NewPrivateEvent("robber_result", "system", RobberResultPayload{
			TargetID: payload.TargetID,
			NewRole:  targetRole,
		}, []string{playerID})
		events = append(events, resultEvent)

	case "troublemaker_swap":
		var payload TroublemakerSwapPayload
		if err := json.Unmarshal(action.Payload, &payload); err != nil {
			return nil, err
		}

		// Validate targets
		if payload.Player1ID == playerID || payload.Player2ID == playerID {
			return nil, errors.New("troublemaker cannot swap themselves")
		}
		if payload.Player1ID == payload.Player2ID {
			return nil, errors.New("must swap two different players")
		}
		if _, exists := g.roleAssignments[payload.Player1ID]; !exists {
			return nil, errors.New("player 1 not found")
		}
		if _, exists := g.roleAssignments[payload.Player2ID]; !exists {
			return nil, errors.New("player 2 not found")
		}

		// Perform the swap
		player1Role := g.roleAssignments[payload.Player1ID]
		player2Role := g.roleAssignments[payload.Player2ID]
		g.roleAssignments[payload.Player1ID] = player2Role
		g.roleAssignments[payload.Player2ID] = player1Role

		g.nightActionsComplete[RoleTroublemaker] = true

		// Troublemaker doesn't see what they swapped, just confirmation
		confirmEvent, _ := core.NewPrivateEvent("troublemaker_confirmed", "system", map[string]interface{}{
			"player1Id": payload.Player1ID,
			"player2Id": payload.Player2ID,
		}, []string{playerID})
		events = append(events, confirmEvent)

	case "drunk_swap":
		var payload DrunkSwapPayload
		if err := json.Unmarshal(action.Payload, &payload); err != nil {
			return nil, err
		}

		if payload.CenterIndex < 0 || payload.CenterIndex >= len(g.centerCards) {
			return nil, errors.New("invalid center card index")
		}

		// Perform the swap
		drunkRole := g.roleAssignments[playerID]
		centerRole := g.centerCards[payload.CenterIndex]
		g.roleAssignments[playerID] = centerRole
		g.centerCards[payload.CenterIndex] = drunkRole

		g.nightActionsComplete[RoleDrunk] = true

		// Drunk doesn't see their new role, just confirmation
		confirmEvent, _ := core.NewPrivateEvent("drunk_confirmed", "system", map[string]interface{}{
			"centerIndex": payload.CenterIndex,
		}, []string{playerID})
		events = append(events, confirmEvent)

	default:
		return nil, fmt.Errorf("unknown action type: %s", action.Type)
	}

	return events, nil
}

// GetPlayerState returns the state visible to a specific player.
func (g *Game) GetPlayerState(playerID string) core.PlayerState {
	role := g.roleAssignments[playerID]

	state := PlayerState{
		Phase:           string(g.phase),
		PhaseEndsAt:     g.phaseEndsAt,
		YourRole:        role,
		HasVoted:        g.votes[playerID] != "",
		HasAcknowledged: g.roleAcknowledgements[playerID],
		TimerActive:     g.timerActive,
	}

	return state
}

// GetPublicState returns the state visible to all players and spectators.
func (g *Game) GetPublicState() core.PublicState {
	return PublicState{
		Phase:                 string(g.phase),
		PhaseEndsAt:           g.phaseEndsAt,
		PlayerCount:           len(g.players),
		VotesSubmitted:        len(g.votes),
		AcknowledgementsCount: len(g.roleAcknowledgements),
		TimerActive:           g.timerActive,
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
