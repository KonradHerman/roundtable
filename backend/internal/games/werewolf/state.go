package werewolf

import "time"

// PlayerState is the werewolf-specific state for a single player.
type PlayerState struct {
	Phase            string    `json:"phase"`
	PhaseEndsAt      time.Time `json:"phaseEndsAt"`
	YourRole         RoleType  `json:"yourRole"`
	HasVoted         bool      `json:"hasVoted"`
	OtherWerewolves  []string  `json:"otherWerewolves,omitempty"`
	OtherMasons      []string  `json:"otherMasons,omitempty"`
}

// PublicState is the werewolf-specific public state.
type PublicState struct {
	Phase          string    `json:"phase"`
	PhaseEndsAt    time.Time `json:"phaseEndsAt"`
	PlayerCount    int       `json:"playerCount"`
	VotesSubmitted int       `json:"votesSubmitted"`
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
