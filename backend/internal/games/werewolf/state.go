package werewolf

import "time"

// PlayerState is the werewolf-specific state for a single player.
type PlayerState struct {
	Phase           string    `json:"phase"`
	PhaseEndsAt     time.Time `json:"phaseEndsAt"`
	YourRole        RoleType  `json:"yourRole"`
	HasVoted        bool      `json:"hasVoted"`
	HasAcknowledged bool      `json:"hasAcknowledged"`
	TimerActive     bool      `json:"timerActive"`
}

// PublicState is the werewolf-specific public state.
type PublicState struct {
	Phase                 string    `json:"phase"`
	PhaseEndsAt           time.Time `json:"phaseEndsAt"`
	PlayerCount           int       `json:"playerCount"`
	VotesSubmitted        int       `json:"votesSubmitted"`
	AcknowledgementsCount int       `json:"acknowledgementsCount"`
	TimerActive           bool      `json:"timerActive"`
}

// Event payloads

type RoleAssignedPayload struct {
	Role RoleType `json:"role"`
}

type WerewolfWakeupPayload struct {
	OtherWerewolves []string `json:"otherWerewolves"`
}

type MasonWakeupPayload struct {
	OtherMasons []string `json:"otherMasons"`
}

type VotePayload struct {
	TargetID string `json:"targetId"`
}

type VoteCastPayload struct {
	VoterID string `json:"voterId"`
}

type VotesRevealedPayload struct {
	Votes map[string]string `json:"votes"` // voterID → targetID
}

type SeerViewPayload struct {
	TargetID string `json:"targetId"`
}

type SeerResultPayload struct {
	TargetID string   `json:"targetId"`
	Role     RoleType `json:"role"`
}

type RobberSwapPayload struct {
	TargetID string `json:"targetId"`
}

type RobberResultPayload struct {
	TargetID string   `json:"targetId"`
	NewRole  RoleType `json:"newRole"`
}

// New payloads for ONUW implementation
type RoleAcknowledgedPayload struct {
	PlayerID string `json:"playerId"`
	Count    int    `json:"count"`
	Total    int    `json:"total"`
}

type NightScriptPayload struct {
	Script []NightScriptStep `json:"script"`
}

type NightScriptStep struct {
	Role        RoleType `json:"role"`
	Instruction string   `json:"instruction"`
	Order       int      `json:"order"`
}

type TimerToggledPayload struct {
	Active      bool       `json:"active"`
	PhaseEndsAt *time.Time `json:"phaseEndsAt,omitempty"`
}

type TimerExtendedPayload struct {
	PhaseEndsAt time.Time `json:"phaseEndsAt"`
	ExtendedBy  int       `json:"extendedBy"` // seconds
}

// Night action payloads

type WerewolfViewCenterPayload struct {
	CenterIndex int `json:"centerIndex"` // 0, 1, or 2
}

type WerewolfViewCenterResultPayload struct {
	CenterIndex int      `json:"centerIndex"`
	Role        RoleType `json:"role"`
}

type SeerViewCenterPayload struct {
	CenterIndices []int `json:"centerIndices"` // Must be exactly 2
}

type SeerViewCenterResultPayload struct {
	Cards []struct {
		Index int      `json:"index"`
		Role  RoleType `json:"role"`
	} `json:"cards"`
}

type TroublemakerSwapPayload struct {
	Player1ID string `json:"player1Id"`
	Player2ID string `json:"player2Id"`
}

type DrunkSwapPayload struct {
	CenterIndex int `json:"centerIndex"` // 0, 1, or 2
}

type InsomniacResultPayload struct {
	FinalRole RoleType `json:"finalRole"`
}

type RolesRevealedPayload struct {
	Roles map[string]RoleType `json:"roles"`
}
