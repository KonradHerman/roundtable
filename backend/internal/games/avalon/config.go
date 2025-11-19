package avalon

import (
	"encoding/json"
	"fmt"

	"github.com/KonradHerman/roundtable/internal/core"
)

// Config represents the configuration for an Avalon game
type Config struct {
	Roles []Role `json:"roles"`
}

// GameType returns "avalon"
func (c *Config) GameType() string {
	return "avalon"
}

// Validate ensures the configuration is valid for the given player count
func (c *Config) Validate() error {
	if len(c.Roles) < 5 || len(c.Roles) > 10 {
		return fmt.Errorf("Avalon requires 5-10 players, got %d", len(c.Roles))
	}

	// Count good and evil roles
	goodCount := 0
	evilCount := 0
	hasMerlin := false
	hasAssassin := false
	hasPercival := false

	for _, role := range c.Roles {
		if isGoodRole(role) {
			goodCount++
		} else if isEvilRole(role) {
			evilCount++
		} else {
			return fmt.Errorf("unknown role: %s", role)
		}

		// Track special roles
		if role == RoleMerlin {
			hasMerlin = true
		}
		if role == RoleAssassin {
			hasAssassin = true
		}
		if role == RolePercival {
			hasPercival = true
		}
	}

	// Validate team sizes based on player count
	expectedGood, expectedEvil := getExpectedTeamSizes(len(c.Roles))
	if goodCount != expectedGood || evilCount != expectedEvil {
		return fmt.Errorf(
			"invalid team sizes for %d players: expected %d good, %d evil; got %d good, %d evil",
			len(c.Roles), expectedGood, expectedEvil, goodCount, evilCount,
		)
	}

	// If Merlin is present, Assassin must be present
	if hasMerlin && !hasAssassin {
		return fmt.Errorf("Assassin is required when Merlin is present")
	}

	// If Percival is present, Merlin must be present
	if hasPercival && !hasMerlin {
		return fmt.Errorf("Merlin is required when Percival is present")
	}

	// Warning: If Morgana is present without Percival, she has no purpose
	// (But we allow it - just less optimal)

	return nil
}

// isGoodRole returns true if the role is on the Good team
func isGoodRole(role Role) bool {
	return role == RoleMerlin || role == RolePercival || role == RoleLoyalServant
}

// isEvilRole returns true if the role is on the Evil team
func isEvilRole(role Role) bool {
	return role == RoleAssassin || role == RoleMorgana || role == RoleMordred ||
		role == RoleOberon || role == RoleMinionOfMordred
}

// getExpectedTeamSizes returns the correct team distribution for player count
func getExpectedTeamSizes(playerCount int) (good int, evil int) {
	switch playerCount {
	case 5, 6:
		return playerCount - 2, 2
	case 7, 8, 9:
		return playerCount - 3, 3
	case 10:
		return 6, 4
	default:
		return 0, 0
	}
}

// ParseConfig parses a JSON config for Avalon and returns core.GameConfig interface
// This matches the ConfigParser signature in the game registry
func ParseConfig(data []byte) (core.GameConfig, error) {
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse avalon config: %w", err)
	}
	return &config, nil
}

// DefaultConfig returns a default configuration for the given player count
// Uses standard roles: Merlin + Assassin + rest as Loyal Servants and Minions
func DefaultConfig(playerCount int) *Config {
	goodCount, evilCount := getExpectedTeamSizes(playerCount)

	roles := make([]Role, 0, playerCount)

	// Always include Merlin and Assassin for standard game
	roles = append(roles, RoleMerlin)
	goodCount--

	roles = append(roles, RoleAssassin)
	evilCount--

	// Fill remaining good roles with Loyal Servants
	for i := 0; i < goodCount; i++ {
		roles = append(roles, RoleLoyalServant)
	}

	// Fill remaining evil roles with Minions
	for i := 0; i < evilCount; i++ {
		roles = append(roles, RoleMinionOfMordred)
	}

	return &Config{Roles: roles}
}

// DefaultConfigWithPercival returns a config with Percival and Morgana
func DefaultConfigWithPercival(playerCount int) *Config {
	goodCount, evilCount := getExpectedTeamSizes(playerCount)

	roles := make([]Role, 0, playerCount)

	// Add special roles
	roles = append(roles, RoleMerlin, RolePercival)
	goodCount -= 2

	roles = append(roles, RoleAssassin, RoleMorgana)
	evilCount -= 2

	// Fill remaining roles
	for i := 0; i < goodCount; i++ {
		roles = append(roles, RoleLoyalServant)
	}
	for i := 0; i < evilCount; i++ {
		roles = append(roles, RoleMinionOfMordred)
	}

	return &Config{Roles: roles}
}
