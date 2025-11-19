package avalon

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"

	"github.com/KonradHerman/roundtable/internal/core"
)

// Game implements the core.Game interface for Avalon
type Game struct {
	mu      sync.RWMutex
	players []*core.Player
	config  *Config
	phase   GamePhase

	// Role assignment
	roles     map[string]Role // playerID -> Role
	teams     map[string]Team // playerID -> Team
	knowledge map[string][]string // playerID -> known player IDs

	// Quest tracking
	questNumber    int           // 1-5
	questResults   []QuestResult // history of completed quests
	currentLeader  string        // playerID
	leaderIndex    int           // position in player list
	rejectionCount int           // consecutive rejections (max 5)

	// Team voting
	proposedTeam []string       // player IDs on proposed team
	teamVotes    map[string]Vote // playerID -> vote

	// Quest execution
	questCards map[string]QuestCard // playerID -> card (team members only)

	// Assassination
	assassinTarget string // playerID

	// Acknowledgments
	acknowledged map[string]bool // playerID -> acknowledged

	// Results
	winningTeam Team
	winReason   string
}

// NewGame creates a new Avalon game instance and returns it as core.Game
func NewGame() core.Game {
	return &Game{
		roles:        make(map[string]Role),
		teams:        make(map[string]Team),
		knowledge:    make(map[string][]string),
		teamVotes:    make(map[string]Vote),
		questCards:   make(map[string]QuestCard),
		acknowledged: make(map[string]bool),
		questResults: []QuestResult{},
	}
}

// Initialize sets up the game with configuration and players
func (g *Game) Initialize(config core.GameConfig, players []*core.Player) ([]core.GameEvent, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Parse config
	avalonConfig, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid config type for avalon")
	}

	// Validate config
	if err := avalonConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid avalon config: %w", err)
	}

	// Validate player count matches config
	if len(players) != len(avalonConfig.Roles) {
		return nil, fmt.Errorf("player count %d does not match role count %d", len(players), len(avalonConfig.Roles))
	}

	g.players = players
	g.config = avalonConfig
	g.phase = PhaseSetup

	// Assign roles
	if err := g.assignRoles(); err != nil {
		return nil, fmt.Errorf("failed to assign roles: %w", err)
	}

	// Select first leader randomly
	g.leaderIndex = rand.Intn(len(g.players))
	g.currentLeader = g.players[g.leaderIndex].ID
	g.questNumber = 1
	g.rejectionCount = 0

	// Create events
	events := []core.GameEvent{}

	// Game started event (public)
	gameStartedPayload := map[string]interface{}{
		"player_count": len(g.players),
		"quest_number": g.questNumber,
		"leader_id":    g.currentLeader,
	}
	gameStartedEvent, _ := core.NewPublicEvent("game_started", "system", gameStartedPayload)
	events = append(events, gameStartedEvent)

	// Role assignment events (private to each player)
	for _, player := range g.players {
		rolePayload := RoleAssignedPayload{
			Role: g.roles[player.ID],
			Team: g.teams[player.ID],
		}
		roleEvent, _ := core.NewPrivateEvent("role_assigned", "system", rolePayload, []string{player.ID})
		events = append(events, roleEvent)

		// Role knowledge event (private)
		knowledgePayload := RoleKnowledgePayload{
			KnownPlayers: g.knowledge[player.ID],
		}
		knowledgeEvent, _ := core.NewPrivateEvent("role_knowledge", "system", knowledgePayload, []string{player.ID})
		events = append(events, knowledgeEvent)
	}

	// Leader changed event (public)
	leaderPayload := LeaderChangedPayload{
		LeaderID: g.currentLeader,
	}
	leaderEvent, _ := core.NewPublicEvent("leader_changed", "system", leaderPayload)
	events = append(events, leaderEvent)

	// Advance to role reveal
	g.phase = PhaseRoleReveal
	phasePayload := map[string]interface{}{
		"phase": PhaseRoleReveal,
	}
	phaseEvent, _ := core.NewPublicEvent("phase_changed", "system", phasePayload)
	events = append(events, phaseEvent)

	return events, nil
}

// assignRoles assigns roles to players with secure shuffling
func (g *Game) assignRoles() error {
	// Create a copy of roles and shuffle
	rolesCopy := make([]Role, len(g.config.Roles))
	copy(rolesCopy, g.config.Roles)

	if err := secureShuffleRoles(rolesCopy); err != nil {
		return err
	}

	// Assign shuffled roles to players
	for i, player := range g.players {
		role := rolesCopy[i]
		g.roles[player.ID] = role
		g.teams[player.ID] = getRoleTeam(role)
	}

	// Generate knowledge for each player
	for _, player := range g.players {
		g.knowledge[player.ID] = getRoleKnowledge(g.roles[player.ID], player.ID, g.roles)
	}

	return nil
}

// ValidateAction checks if a player can perform the given action
func (g *Game) ValidateAction(playerID string, action core.Action) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	switch action.Type {
	case "acknowledge_role":
		if g.phase != PhaseRoleReveal {
			return fmt.Errorf("can only acknowledge role during role reveal phase")
		}
		if g.acknowledged[playerID] {
			return fmt.Errorf("already acknowledged role")
		}

	case "propose_team":
		if g.phase != PhaseTeamBuilding {
			return fmt.Errorf("can only propose team during team building phase")
		}
		if playerID != g.currentLeader {
			return fmt.Errorf("only the leader can propose a team")
		}

	case "vote_team":
		if g.phase != PhaseTeamVoting {
			return fmt.Errorf("can only vote during team voting phase")
		}
		if _, hasVoted := g.teamVotes[playerID]; hasVoted {
			return fmt.Errorf("already voted")
		}

	case "play_quest_card":
		if g.phase != PhaseQuestExec {
			return fmt.Errorf("can only play quest cards during quest execution phase")
		}
		if !g.isOnProposedTeam(playerID) {
			return fmt.Errorf("only team members can play quest cards")
		}
		if _, hasPlayed := g.questCards[playerID]; hasPlayed {
			return fmt.Errorf("already played quest card")
		}

	case "assassinate":
		if g.phase != PhaseAssassination {
			return fmt.Errorf("can only assassinate during assassination phase")
		}
		if g.roles[playerID] != RoleAssassin {
			return fmt.Errorf("only the assassin can assassinate")
		}
		if g.assassinTarget != "" {
			return fmt.Errorf("already selected assassination target")
		}

	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}

	return nil
}

// ProcessAction executes a valid action and returns resulting events
func (g *Game) ProcessAction(playerID string, action core.Action) ([]core.GameEvent, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	switch action.Type {
	case "acknowledge_role":
		return g.processAcknowledgeRole(playerID)
	case "propose_team":
		return g.processProposeTeam(playerID, action.Payload)
	case "vote_team":
		return g.processVoteTeam(playerID, action.Payload)
	case "play_quest_card":
		return g.processPlayQuestCard(playerID, action.Payload)
	case "assassinate":
		return g.processAssassinate(playerID, action.Payload)
	default:
		return nil, fmt.Errorf("unknown action type: %s", action.Type)
	}
}

// processAcknowledgeRole handles role acknowledgment
func (g *Game) processAcknowledgeRole(playerID string) ([]core.GameEvent, error) {
	g.acknowledged[playerID] = true

	events := []core.GameEvent{}

	// Count acknowledgments
	count := 0
	for _, ack := range g.acknowledged {
		if ack {
			count++
		}
	}

	// Public event showing progress
	ackPayload := RoleAcknowledgedPayload{
		PlayerID: playerID,
		Count:    count,
		Total:    len(g.players),
	}
	ackEvent, _ := core.NewPublicEvent("role_acknowledged", playerID, ackPayload)
	events = append(events, ackEvent)

	// If all acknowledged, advance to team building
	if count == len(g.players) {
		g.phase = PhaseTeamBuilding
		phasePayload := map[string]interface{}{
			"phase":        PhaseTeamBuilding,
			"quest_number": g.questNumber,
			"team_size":    getRequiredTeamSize(len(g.players), g.questNumber),
			"leader_id":    g.currentLeader,
		}
		phaseEvent, _ := core.NewPublicEvent("phase_changed", "system", phasePayload)
		events = append(events, phaseEvent)
	}

	return events, nil
}

// processProposeTeam handles team proposal
func (g *Game) processProposeTeam(playerID string, payload json.RawMessage) ([]core.GameEvent, error) {
	var data struct {
		TeamMembers []string `json:"team_members"`
	}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("invalid propose_team payload: %w", err)
	}

	// Validate team size
	requiredSize := getRequiredTeamSize(len(g.players), g.questNumber)
	if len(data.TeamMembers) != requiredSize {
		return nil, fmt.Errorf("team size must be %d, got %d", requiredSize, len(data.TeamMembers))
	}

	// Validate all team members exist
	playerMap := make(map[string]bool)
	for _, p := range g.players {
		playerMap[p.ID] = true
	}
	for _, memberID := range data.TeamMembers {
		if !playerMap[memberID] {
			return nil, fmt.Errorf("invalid player ID: %s", memberID)
		}
	}

	g.proposedTeam = data.TeamMembers
	g.teamVotes = make(map[string]Vote) // Reset votes

	events := []core.GameEvent{}

	// Team proposed event (public)
	proposalPayload := TeamProposedPayload{
		LeaderID:    playerID,
		TeamMembers: data.TeamMembers,
		QuestNumber: g.questNumber,
		TeamSize:    requiredSize,
	}
	proposalEvent, _ := core.NewPublicEvent("team_proposed", playerID, proposalPayload)
	events = append(events, proposalEvent)

	// Advance to team voting
	g.phase = PhaseTeamVoting
	phasePayload := map[string]interface{}{
		"phase":        PhaseTeamVoting,
		"team_members": g.proposedTeam,
	}
	phaseEvent, _ := core.NewPublicEvent("phase_changed", "system", phasePayload)
	events = append(events, phaseEvent)

	return events, nil
}

// processVoteTeam handles team voting
func (g *Game) processVoteTeam(playerID string, payload json.RawMessage) ([]core.GameEvent, error) {
	var data struct {
		Vote string `json:"vote"` // "approve" or "reject"
	}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("invalid vote_team payload: %w", err)
	}

	vote := Vote(data.Vote)
	if vote != VoteApprove && vote != VoteReject {
		return nil, fmt.Errorf("vote must be 'approve' or 'reject', got: %s", data.Vote)
	}

	g.teamVotes[playerID] = vote

	events := []core.GameEvent{}

	// Public event (player voted, not their vote)
	voteCastPayload := TeamVoteCastPayload{
		VoterID: playerID,
	}
	voteCastEvent, _ := core.NewPublicEvent("team_vote_cast", playerID, voteCastPayload)
	events = append(events, voteCastEvent)

	// Private confirmation
	voteRecordedPayload := TeamVoteRecordedPayload{
		Vote: vote,
	}
	voteRecordedEvent, _ := core.NewPrivateEvent("team_vote_recorded", "system", voteRecordedPayload, []string{playerID})
	events = append(events, voteRecordedEvent)

	// Check if all votes are in
	if len(g.teamVotes) == len(g.players) {
		// Tally votes
		approveCount := 0
		rejectCount := 0
		for _, v := range g.teamVotes {
			if v == VoteApprove {
				approveCount++
			} else {
				rejectCount++
			}
		}

		approved := approveCount > rejectCount

		// Vote result event (public)
		voteResultPayload := TeamVoteResultPayload{
			Approved:       approved,
			Votes:          g.teamVotes,
			ApproveCount:   approveCount,
			RejectCount:    rejectCount,
			RejectionCount: g.rejectionCount,
		}
		voteResultEvent, _ := core.NewPublicEvent("team_vote_result", "system", voteResultPayload)
		events = append(events, voteResultEvent)

		if approved {
			// Team approved - advance to quest execution
			g.rejectionCount = 0 // Reset rejection count
			g.phase = PhaseQuestExec
			g.questCards = make(map[string]QuestCard) // Reset quest cards

			phasePayload := map[string]interface{}{
				"phase":        PhaseQuestExec,
				"team_members": g.proposedTeam,
			}
			phaseEvent, _ := core.NewPublicEvent("phase_changed", "system", phasePayload)
			events = append(events, phaseEvent)

		} else {
			// Team rejected
			g.rejectionCount++

			// Check for 5 consecutive rejections
			if g.rejectionCount >= 5 {
				// Evil wins automatically
				g.phase = PhaseFinished
				g.winningTeam = TeamEvil
				g.winReason = "five_consecutive_rejections"

				finishEvents := g.createGameFinishedEvents()
				events = append(events, finishEvents...)

			} else {
				// Rotate leader and try again
				g.rotateLeader()

				leaderPayload := LeaderChangedPayload{
					LeaderID: g.currentLeader,
				}
				leaderEvent, _ := core.NewPublicEvent("leader_changed", "system", leaderPayload)
				events = append(events, leaderEvent)

				// Back to team building
				g.phase = PhaseTeamBuilding
				phasePayload := map[string]interface{}{
					"phase":           PhaseTeamBuilding,
					"quest_number":    g.questNumber,
					"team_size":       getRequiredTeamSize(len(g.players), g.questNumber),
					"leader_id":       g.currentLeader,
					"rejection_count": g.rejectionCount,
				}
				phaseEvent, _ := core.NewPublicEvent("phase_changed", "system", phasePayload)
				events = append(events, phaseEvent)
			}
		}
	}

	return events, nil
}

// processPlayQuestCard handles quest card submission
func (g *Game) processPlayQuestCard(playerID string, payload json.RawMessage) ([]core.GameEvent, error) {
	var data struct {
		Card string `json:"card"` // "success" or "fail"
	}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("invalid play_quest_card payload: %w", err)
	}

	card := QuestCard(data.Card)
	if card != CardSuccess && card != CardFail {
		return nil, fmt.Errorf("card must be 'success' or 'fail', got: %s", data.Card)
	}

	// Good team can ONLY play success
	if g.teams[playerID] == TeamGood && card != CardSuccess {
		return nil, fmt.Errorf("good team players can only play success cards")
	}

	g.questCards[playerID] = card

	events := []core.GameEvent{}

	// Public event (player played, not which card)
	cardPlayedPayload := QuestCardPlayedPayload{
		PlayerID: playerID,
	}
	cardPlayedEvent, _ := core.NewPublicEvent("quest_card_played", playerID, cardPlayedPayload)
	events = append(events, cardPlayedEvent)

	// Private confirmation
	cardRecordedPayload := QuestCardRecordedPayload{
		Card: card,
	}
	cardRecordedEvent, _ := core.NewPrivateEvent("quest_card_recorded", "system", cardRecordedPayload, []string{playerID})
	events = append(events, cardRecordedEvent)

	// Check if all cards are in
	if len(g.questCards) == len(g.proposedTeam) {
		// Shuffle cards to prevent order tells
		cards := []QuestCard{}
		for _, card := range g.questCards {
			cards = append(cards, card)
		}
		g.shuffleQuestCards(cards)

		// Count fails
		failCount := 0
		for _, card := range cards {
			if card == CardFail {
				failCount++
			}
		}

		// Determine if quest succeeded
		failsRequired := getFailsRequired(len(g.players), g.questNumber)
		success := failCount < failsRequired

		// Record quest result
		result := QuestResult{
			QuestNumber:   g.questNumber,
			TeamSize:      len(g.proposedTeam),
			TeamMembers:   g.proposedTeam,
			Cards:         cards,
			FailCount:     failCount,
			Success:       success,
			FailsRequired: failsRequired,
		}
		g.questResults = append(g.questResults, result)

		// Count total wins
		goodWins, evilWins := countTeamQuests(g.questResults)

		// Quest completed event (public)
		questCompletedPayload := QuestCompletedPayload{
			QuestNumber:   g.questNumber,
			TeamMembers:   g.proposedTeam,
			Cards:         cards,
			FailCount:     failCount,
			Success:       success,
			FailsRequired: failsRequired,
			GoodWins:      goodWins,
			EvilWins:      evilWins,
		}
		questCompletedEvent, _ := core.NewPublicEvent("quest_completed", "system", questCompletedPayload)
		events = append(events, questCompletedEvent)

		// Check for game end
		if goodWins >= 3 {
			// Good won 3 quests - check for assassination
			if hasRole(g.roles, RoleMerlin) {
				// Assassination phase
				g.phase = PhaseAssassination
				phasePayload := map[string]interface{}{
					"phase": PhaseAssassination,
				}
				phaseEvent, _ := core.NewPublicEvent("phase_changed", "system", phasePayload)
				events = append(events, phaseEvent)
			} else {
				// No Merlin - Good wins
				g.phase = PhaseFinished
				g.winningTeam = TeamGood
				g.winReason = "good_won_three_quests"

				finishEvents := g.createGameFinishedEvents()
				events = append(events, finishEvents...)
			}

		} else if evilWins >= 3 {
			// Evil won 3 quests - game over
			g.phase = PhaseFinished
			g.winningTeam = TeamEvil
			g.winReason = "evil_sabotaged_three_quests"

			finishEvents := g.createGameFinishedEvents()
			events = append(events, finishEvents...)

		} else {
			// Game continues - advance to next quest
			g.questNumber++
			g.rotateLeader()
			g.rejectionCount = 0 // Reset rejection count for new quest

			leaderPayload := LeaderChangedPayload{
				LeaderID: g.currentLeader,
			}
			leaderEvent, _ := core.NewPublicEvent("leader_changed", "system", leaderPayload)
			events = append(events, leaderEvent)

			g.phase = PhaseTeamBuilding
			phasePayload := map[string]interface{}{
				"phase":        PhaseTeamBuilding,
				"quest_number": g.questNumber,
				"team_size":    getRequiredTeamSize(len(g.players), g.questNumber),
				"leader_id":    g.currentLeader,
			}
			phaseEvent, _ := core.NewPublicEvent("phase_changed", "system", phasePayload)
			events = append(events, phaseEvent)
		}
	}

	return events, nil
}

// processAssassinate handles assassination
func (g *Game) processAssassinate(playerID string, payload json.RawMessage) ([]core.GameEvent, error) {
	var data struct {
		TargetID string `json:"target_id"`
	}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("invalid assassinate payload: %w", err)
	}

	// Validate target exists
	targetExists := false
	for _, p := range g.players {
		if p.ID == data.TargetID {
			targetExists = true
			break
		}
	}
	if !targetExists {
		return nil, fmt.Errorf("invalid target player ID: %s", data.TargetID)
	}

	g.assassinTarget = data.TargetID

	events := []core.GameEvent{}

	// Assassin target event (public)
	targetPayload := AssassinTargetPayload{
		TargetID: data.TargetID,
	}
	targetEvent, _ := core.NewPublicEvent("assassin_target", playerID, targetPayload)
	events = append(events, targetEvent)

	// Check if target is Merlin
	targetRole := g.roles[data.TargetID]
	wasMerlin := targetRole == RoleMerlin

	// Assassin result event (public)
	resultPayload := AssassinResultPayload{
		TargetID:   data.TargetID,
		TargetRole: targetRole,
		WasMerlin:  wasMerlin,
		EvilWins:   wasMerlin,
	}
	resultEvent, _ := core.NewPublicEvent("assassin_result", "system", resultPayload)
	events = append(events, resultEvent)

	// Determine winner
	if wasMerlin {
		g.winningTeam = TeamEvil
		g.winReason = "assassin_found_merlin"
	} else {
		g.winningTeam = TeamGood
		g.winReason = "assassin_failed"
	}

	g.phase = PhaseFinished

	finishEvents := g.createGameFinishedEvents()
	events = append(events, finishEvents...)

	return events, nil
}

// Helper methods

func (g *Game) rotateLeader() {
	g.leaderIndex = (g.leaderIndex + 1) % len(g.players)
	g.currentLeader = g.players[g.leaderIndex].ID
}

func (g *Game) isOnProposedTeam(playerID string) bool {
	for _, id := range g.proposedTeam {
		if id == playerID {
			return true
		}
	}
	return false
}

func (g *Game) shuffleQuestCards(cards []QuestCard) {
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}

func (g *Game) createGameFinishedEvents() []core.GameEvent {
	events := []core.GameEvent{}

	finishPayload := GameFinishedPayload{
		WinningTeam:  g.winningTeam,
		WinReason:    g.winReason,
		Roles:        g.roles,
		Teams:        g.teams,
		QuestHistory: g.questResults,
	}
	finishEvent, _ := core.NewPublicEvent("game_finished", "system", finishPayload)
	events = append(events, finishEvent)

	phasePayload := map[string]interface{}{
		"phase": PhaseFinished,
	}
	phaseEvent, _ := core.NewPublicEvent("phase_changed", "system", phasePayload)
	events = append(events, phaseEvent)

	return events
}

// Interface implementation

func (g *Game) GetPlayerState(playerID string) core.PlayerState {
	g.mu.RLock()
	defer g.mu.RUnlock()

	goodWins, evilWins := countTeamQuests(g.questResults)

	return PlayerState{
		Phase:              g.phase,
		Role:               g.roles[playerID],
		Team:               g.teams[playerID],
		Knowledge:          g.knowledge[playerID],
		HasAcknowledged:    g.acknowledged[playerID],
		HasVoted:           g.teamVotes[playerID] != "",
		HasPlayedQuestCard: g.questCards[playerID] != "",
		IsOnProposedTeam:   g.isOnProposedTeam(playerID),
		IsCurrentLeader:    g.currentLeader == playerID,
		CanProposeTeam:     g.currentLeader == playerID && g.phase == PhaseTeamBuilding,
		CanVote:            g.phase == PhaseTeamVoting && g.teamVotes[playerID] == "",
		CanPlayQuestCard:   g.phase == PhaseQuestExec && g.isOnProposedTeam(playerID) && g.questCards[playerID] == "",
		CanAssassinate:     g.phase == PhaseAssassination && g.roles[playerID] == RoleAssassin,
		QuestNumber:        g.questNumber,
		RejectionCount:     g.rejectionCount,
		GoodQuestWins:      goodWins,
		EvilQuestWins:      evilWins,
	}
}

func (g *Game) GetPublicState() core.PublicState {
	g.mu.RLock()
	defer g.mu.RUnlock()

	goodWins, evilWins := countTeamQuests(g.questResults)

	return PublicState{
		Phase:                 g.phase,
		PlayerCount:           len(g.players),
		QuestNumber:           g.questNumber,
		RequiredTeamSize:      getRequiredTeamSize(len(g.players), g.questNumber),
		CurrentLeaderID:       g.currentLeader,
		ProposedTeam:          g.proposedTeam,
		VotesSubmitted:        len(g.teamVotes),
		TotalVotes:            len(g.players),
		CardsSubmitted:        len(g.questCards),
		TotalCardsExpected:    len(g.proposedTeam),
		QuestResults:          g.questResults,
		RejectionCount:        g.rejectionCount,
		AcknowledgementsCount: g.countAcknowledgments(),
		GoodQuestWins:         goodWins,
		EvilQuestWins:         evilWins,
	}
}

func (g *Game) countAcknowledgments() int {
	count := 0
	for _, ack := range g.acknowledged {
		if ack {
			count++
		}
	}
	return count
}

func (g *Game) GetPhase() core.GamePhase {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return core.GamePhase{
		Name:    string(g.phase),
		EndsAt:  nil, // Avalon has no phase timers
		Message: g.getPhaseMessage(),
	}
}

func (g *Game) getPhaseMessage() string {
	switch g.phase {
	case PhaseSetup:
		return "Setting up game..."
	case PhaseRoleReveal:
		return "Review your role"
	case PhaseTeamBuilding:
		return fmt.Sprintf("Quest %d: Leader selects team", g.questNumber)
	case PhaseTeamVoting:
		return "Vote to approve or reject the team"
	case PhaseQuestExec:
		return "Team members: play your quest cards"
	case PhaseQuestResults:
		return "Quest results revealed"
	case PhaseAssassination:
		return "The Assassin is choosing their target..."
	case PhaseFinished:
		return "Game finished"
	default:
		return ""
	}
}

func (g *Game) IsFinished() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.phase == PhaseFinished
}

func (g *Game) GetResults() core.GameResults {
	g.mu.RLock()
	defer g.mu.RUnlock()

	winners := []string{}
	for _, player := range g.players {
		if g.teams[player.ID] == g.winningTeam {
			winners = append(winners, player.ID)
		}
	}

	return core.GameResults{
		Winners:   winners,
		WinReason: g.winReason,
		FinalState: map[string]interface{}{
			"winning_team":  g.winningTeam,
			"quest_history": g.questResults,
			"roles":         g.roles,
		},
	}
}

func (g *Game) CheckPhaseTimeout() ([]core.GameEvent, error) {
	// Avalon has no phase timers - all phases are player-driven
	return nil, nil
}
