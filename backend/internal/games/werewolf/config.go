package werewolf

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/KonradHerman/roundtable/internal/core"
)

// Config holds the configuration for a werewolf game.
type Config struct {
	Roles         []RoleType    `json:"roles"`         // List of roles to assign
	NightDuration time.Duration `json:"nightDuration"` // How long night phase lasts
	DayDuration   time.Duration `json:"dayDuration"`   // How long day phase lasts
}

// GameType returns the game type identifier.
func (c *Config) GameType() string {
	return "werewolf"
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if len(c.Roles) == 0 {
		return errors.New("at least one role required")
	}

	// Check for at least one werewolf
	hasWerewolf := false
	for _, role := range c.Roles {
		if role == RoleWerewolf {
			hasWerewolf = true
			break
		}
	}

	if !hasWerewolf {
		return errors.New("at least one werewolf required")
	}

	// Validate durations
	if c.NightDuration <= 0 {
		c.NightDuration = 3 * time.Minute // Default
	}

	if c.DayDuration <= 0 {
		c.DayDuration = 5 * time.Minute // Default
	}

	return nil
}

// ParseConfig parses raw JSON into a werewolf config.
func ParseConfig(data []byte) (core.GameConfig, error) {
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// RoleType represents a player role in werewolf.
type RoleType string

const (
	RoleWerewolf     RoleType = "werewolf"
	RoleSeer         RoleType = "seer"
	RoleRobber       RoleType = "robber"
	RoleTroublemaker RoleType = "troublemaker"
	RoleDrunk        RoleType = "drunk"
	RoleInsomniac    RoleType = "insomniac"
	RoleMason        RoleType = "mason"
	RoleMinion       RoleType = "minion"
	RoleHunter       RoleType = "hunter"
	RoleTanner       RoleType = "tanner"
	RoleVillager     RoleType = "villager"
)

// IsWerewolfTeam returns true if this role is on the werewolf team.
func (r RoleType) IsWerewolfTeam() bool {
	return r == RoleWerewolf || r == RoleMinion
}

// IsVillageTeam returns true if this role is on the village team.
func (r RoleType) IsVillageTeam() bool {
	return !r.IsWerewolfTeam() && r != RoleTanner
}

// HasNightAction returns true if this role acts during night.
func (r RoleType) HasNightAction() bool {
	switch r {
	case RoleSeer, RoleRobber, RoleTroublemaker, RoleDrunk, RoleInsomniac:
		return true
	default:
		return false
	}
}
