package avalon

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// getRoleTeam returns the team for a given role
func getRoleTeam(role Role) Team {
	if isGoodRole(role) {
		return TeamGood
	}
	return TeamEvil
}

// getRoleKnowledge returns the player IDs that this role can see
// This implements the complex visibility rules of Avalon:
// - Merlin sees all Evil except Mordred
// - Percival sees Merlin and Morgana (can't distinguish)
// - Evil sees each other except Oberon
// - Oberon sees nothing
// - Good (non-Merlin/Percival) sees nothing
func getRoleKnowledge(role Role, playerID string, roles map[string]Role) []string {
	knowledge := []string{}

	switch role {
	case RoleMerlin:
		// Merlin sees all Evil players except Mordred
		for pid, r := range roles {
			if pid != playerID && isEvilRole(r) && r != RoleMordred {
				knowledge = append(knowledge, pid)
			}
		}

	case RolePercival:
		// Percival sees Merlin and Morgana (cannot distinguish which is which)
		for pid, r := range roles {
			if pid != playerID && (r == RoleMerlin || r == RoleMorgana) {
				knowledge = append(knowledge, pid)
			}
		}

	case RoleAssassin, RoleMorgana, RoleMordred, RoleMinionOfMordred:
		// Evil players see each other except Oberon
		for pid, r := range roles {
			if pid != playerID && isEvilRole(r) && r != RoleOberon {
				knowledge = append(knowledge, pid)
			}
		}

	case RoleOberon:
		// Oberon sees nothing (lone wolf)
		// knowledge remains empty

	case RoleLoyalServant:
		// Loyal servants see nothing
		// knowledge remains empty

	default:
		// Unknown role - sees nothing
		// knowledge remains empty
	}

	return knowledge
}

// secureShuffleRoles shuffles roles using crypto/rand for unpredictability
func secureShuffleRoles(roles []Role) error {
	n := len(roles)
	for i := n - 1; i > 0; i-- {
		// Generate cryptographically secure random number
		maxBig := big.NewInt(int64(i + 1))
		jBig, err := rand.Int(rand.Reader, maxBig)
		if err != nil {
			return fmt.Errorf("failed to generate secure random number: %w", err)
		}
		j := int(jBig.Int64())

		// Swap roles[i] and roles[j]
		roles[i], roles[j] = roles[j], roles[i]
	}
	return nil
}

// hasRole returns true if any player has the given role
func hasRole(roles map[string]Role, role Role) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// getPlayerWithRole returns the player ID with the given role, or empty string
func getPlayerWithRole(roles map[string]Role, role Role) string {
	for pid, r := range roles {
		if r == role {
			return pid
		}
	}
	return ""
}

// countTeamQuests counts how many quests each team has won
func countTeamQuests(results []QuestResult) (goodWins int, evilWins int) {
	for _, result := range results {
		if result.Success {
			goodWins++
		} else {
			evilWins++
		}
	}
	return
}

// getRoleDescription returns a human-readable description of the role
func getRoleDescription(role Role) string {
	descriptions := map[Role]string{
		RoleMerlin:          "Knows the forces of Evil (except Mordred). Help Good win without revealing yourself!",
		RolePercival:        "You see two powerful wizards. One is Merlin, one is Morgana. Protect Merlin!",
		RoleLoyalServant:    "You have no special information. Trust your instincts and your allies!",
		RoleAssassin:        "You know your Evil allies. If Good wins 3 quests, you can steal victory by identifying Merlin!",
		RoleMorgana:         "You appear as Merlin to Percival. Confuse the Good team!",
		RoleMordred:         "Merlin cannot see you. Use this advantage wisely!",
		RoleOberon:          "You are alone. You do not know other Evil players, and they do not know you.",
		RoleMinionOfMordred: "You know your Evil allies. Work together to sabotage the quests!",
	}
	return descriptions[role]
}
