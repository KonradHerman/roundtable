package avalon

// QuestConfig defines team sizes and fail requirements per player count
type QuestConfig struct {
	PlayerCount int
	TeamSizes   [5]int // Quest 1-5 team sizes
	Quest4Fails int    // Number of fails required for quest 4 (1 or 2)
}

// Quest configuration based on player count
var questConfigs = map[int]QuestConfig{
	5:  {PlayerCount: 5, TeamSizes: [5]int{2, 3, 2, 3, 3}, Quest4Fails: 1},
	6:  {PlayerCount: 6, TeamSizes: [5]int{2, 3, 4, 3, 4}, Quest4Fails: 1},
	7:  {PlayerCount: 7, TeamSizes: [5]int{2, 3, 3, 4, 4}, Quest4Fails: 2},
	8:  {PlayerCount: 8, TeamSizes: [5]int{3, 4, 4, 5, 5}, Quest4Fails: 2},
	9:  {PlayerCount: 9, TeamSizes: [5]int{3, 4, 4, 5, 5}, Quest4Fails: 2},
	10: {PlayerCount: 10, TeamSizes: [5]int{3, 4, 4, 5, 5}, Quest4Fails: 2},
}

// getQuestConfig returns the quest configuration for the given player count
func getQuestConfig(playerCount int) QuestConfig {
	config, ok := questConfigs[playerCount]
	if !ok {
		// Default to 5-player config if invalid
		return questConfigs[5]
	}
	return config
}

// getRequiredTeamSize returns the team size required for the given quest
func getRequiredTeamSize(playerCount int, questNumber int) int {
	config := getQuestConfig(playerCount)
	if questNumber < 1 || questNumber > 5 {
		return 0
	}
	return config.TeamSizes[questNumber-1]
}

// getFailsRequired returns the number of fail cards needed to fail the quest
// Quest 4 with 7+ players requires 2 fails, all others require 1
func getFailsRequired(playerCount int, questNumber int) int {
	if questNumber == 4 && playerCount >= 7 {
		return 2
	}
	return 1
}

// requiresTwoFails returns true if the quest requires 2 fail cards
func requiresTwoFails(playerCount int, questNumber int) bool {
	return getFailsRequired(playerCount, questNumber) == 2
}
