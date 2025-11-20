package avalon

// Role types
type Role string

const (
	RoleMerlin          Role = "merlin"
	RolePercival        Role = "percival"
	RoleLoyalServant    Role = "loyal_servant"
	RoleAssassin        Role = "assassin"
	RoleMorgana         Role = "morgana"
	RoleMordred         Role = "mordred"
	RoleOberon          Role = "oberon"
	RoleMinionOfMordred Role = "minion"
)

// Team types
type Team string

const (
	TeamGood Team = "good"
	TeamEvil Team = "evil"
)

// Game phases
type GamePhase string

const (
	PhaseSetup         GamePhase = "setup"
	PhaseRoleReveal    GamePhase = "role_reveal"
	PhaseTeamBuilding  GamePhase = "team_building"
	PhaseTeamVoting    GamePhase = "team_voting"
	PhaseQuestExec     GamePhase = "quest_execution"
	PhaseQuestResults  GamePhase = "quest_results"
	PhaseAssassination GamePhase = "assassination"
	PhaseFinished      GamePhase = "finished"
)

// Quest card types
type QuestCard string

const (
	CardSuccess QuestCard = "success"
	CardFail    QuestCard = "fail"
)

// Vote types
type Vote string

const (
	VoteApprove Vote = "approve"
	VoteReject  Vote = "reject"
)

// QuestResult represents the outcome of a single quest
type QuestResult struct {
	QuestNumber   int         `json:"quest_number"`
	TeamSize      int         `json:"team_size"`
	TeamMembers   []string    `json:"team_members"`   // player IDs
	Cards         []QuestCard `json:"cards"`          // shuffled results
	FailCount     int         `json:"fail_count"`     // number of fail cards
	Success       bool        `json:"success"`        // true if quest succeeded
	FailsRequired int         `json:"fails_required"` // 1 or 2 (quest 4 with 7+ players)
}

// PlayerState is the avalon-specific state for a single player
type PlayerState struct {
	Phase               GamePhase `json:"phase"`
	Role                Role      `json:"role"`
	Team                Team      `json:"team"`
	Knowledge           []string  `json:"knowledge"`            // player IDs this player knows
	HasAcknowledged     bool      `json:"hasAcknowledged"`      // acknowledged role
	HasVoted            bool      `json:"hasVoted"`             // voted on team
	HasPlayedQuestCard  bool      `json:"hasPlayedQuestCard"`   // played quest card
	IsOnProposedTeam    bool      `json:"isOnProposedTeam"`     // selected for current team
	IsCurrentLeader     bool      `json:"isCurrentLeader"`      // is the leader
	CanProposeTeam      bool      `json:"canProposeTeam"`       // can propose team (is leader + team building phase)
	CanVote             bool      `json:"canVote"`              // can vote on team
	CanPlayQuestCard    bool      `json:"canPlayQuestCard"`     // can play quest card
	CanAssassinate      bool      `json:"canAssassinate"`       // can assassinate (is assassin + assassination phase)
	QuestNumber         int       `json:"quest_number"`         // current quest (1-5)
	RejectionCount      int       `json:"rejection_count"`      // consecutive rejections
	GoodQuestWins       int       `json:"good_quest_wins"`      // good team quest wins
	EvilQuestWins       int       `json:"evil_quest_wins"`      // evil team quest wins
}

// PublicState is the avalon-specific public state
type PublicState struct {
	Phase                 GamePhase      `json:"phase"`
	PlayerCount           int            `json:"playerCount"`
	QuestNumber           int            `json:"quest_number"`          // 1-5
	RequiredTeamSize      int            `json:"required_team_size"`    // size for current quest
	CurrentLeaderID       string         `json:"current_leader_id"`     // player ID of leader
	ProposedTeam          []string       `json:"proposed_team"`         // player IDs on proposed team
	VotesSubmitted        int            `json:"votes_submitted"`       // count of votes submitted
	TotalVotes            int            `json:"total_votes"`           // total votes expected
	CardsSubmitted        int            `json:"cards_submitted"`       // quest cards submitted
	TotalCardsExpected    int            `json:"total_cards_expected"`  // quest cards expected
	QuestResults          []QuestResult  `json:"quest_results"`         // history of completed quests
	RejectionCount        int            `json:"rejection_count"`       // consecutive rejections (max 5)
	AcknowledgementsCount int            `json:"acknowledgements_count"`
	GoodQuestWins         int            `json:"good_quest_wins"`       // good team quest wins
	EvilQuestWins         int            `json:"evil_quest_wins"`       // evil team quest wins
}

// Event payloads

type RoleAssignedPayload struct {
	Role Role `json:"role"`
	Team Team `json:"team"`
}

type RoleKnowledgePayload struct {
	KnownPlayers []string `json:"known_players"` // player IDs you can see
}

type LeaderChangedPayload struct {
	LeaderID string `json:"leader_id"`
}

type TeamProposedPayload struct {
	LeaderID    string   `json:"leader_id"`
	TeamMembers []string `json:"team_members"`
	QuestNumber int      `json:"quest_number"`
	TeamSize    int      `json:"team_size"`
}

type TeamVoteCastPayload struct {
	VoterID string `json:"voter_id"`
}

type TeamVoteRecordedPayload struct {
	Vote Vote `json:"vote"` // private confirmation
}

type TeamVoteResultPayload struct {
	Approved       bool              `json:"approved"`
	Votes          map[string]Vote   `json:"votes"`          // voterID -> vote
	ApproveCount   int               `json:"approve_count"`
	RejectCount    int               `json:"reject_count"`
	RejectionCount int               `json:"rejection_count"` // consecutive rejections
}

type QuestCardPlayedPayload struct {
	PlayerID string `json:"player_id"` // who played (not which card)
}

type QuestCardRecordedPayload struct {
	Card QuestCard `json:"card"` // private confirmation
}

type QuestCompletedPayload struct {
	QuestNumber   int         `json:"quest_number"`
	TeamMembers   []string    `json:"team_members"`
	Cards         []QuestCard `json:"cards"` // shuffled
	FailCount     int         `json:"fail_count"`
	Success       bool        `json:"success"`
	FailsRequired int         `json:"fails_required"`
	GoodWins      int         `json:"good_wins"` // total good wins after this quest
	EvilWins      int         `json:"evil_wins"` // total evil wins after this quest
}

type AssassinTargetPayload struct {
	TargetID string `json:"target_id"` // who was targeted
}

type AssassinResultPayload struct {
	TargetID       string `json:"target_id"`
	TargetRole     Role   `json:"target_role"`
	WasMerlin      bool   `json:"was_merlin"`
	EvilWins       bool   `json:"evil_wins"`        // true if assassin found Merlin
}

type GameFinishedPayload struct {
	WinningTeam Team                   `json:"winning_team"`
	WinReason   string                 `json:"win_reason"`
	Roles       map[string]Role        `json:"roles"`        // all player roles revealed
	Teams       map[string]Team        `json:"teams"`        // all player teams revealed
	QuestHistory []QuestResult         `json:"quest_history"`
}

type RoleAcknowledgedPayload struct {
	PlayerID string `json:"player_id"`
	Count    int    `json:"count"`
	Total    int    `json:"total"`
}
