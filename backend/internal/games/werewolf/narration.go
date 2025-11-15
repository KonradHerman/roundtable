package werewolf

import "sort"

// GenerateNightScript creates the narration script for the night phase
// based on which roles are in play.
func GenerateNightScript(rolesInPlay []RoleType) []NightScriptStep {
	script := make([]NightScriptStep, 0)
	roleSet := make(map[RoleType]bool)
	
	// Build set of unique roles
	for _, role := range rolesInPlay {
		roleSet[role] = true
	}

	// Define the order and instructions for each role
	roleOrder := []struct {
		role        RoleType
		order       int
		instruction string
	}{
		{
			role:  RoleWerewolf,
			order: 1,
			instruction: "Werewolves, wake up and look for other Werewolves. If you are the only Werewolf, you may view one center card.",
		},
		{
			role:  RoleMinion,
			order: 2,
			instruction: "Minion, wake up. Werewolves, raise your hand so the Minion can see you. Werewolves, put your hands down. Minion, close your eyes.",
		},
		{
			role:  RoleMason,
			order: 3,
			instruction: "Masons, wake up and look for other Masons.",
		},
		{
			role:  RoleSeer,
			order: 4,
			instruction: "Seer, wake up. You may look at another player's card or two of the center cards.",
		},
		{
			role:  RoleRobber,
			order: 5,
			instruction: "Robber, wake up. You may exchange your card with another player's card, and then view your new card.",
		},
		{
			role:  RoleTroublemaker,
			order: 6,
			instruction: "Troublemaker, wake up. You may exchange cards between two other players without looking at those cards.",
		},
		{
			role:  RoleDrunk,
			order: 7,
			instruction: "Drunk, wake up and exchange your card with a card from the center without looking at your new card.",
		},
		{
			role:  RoleInsomniac,
			order: 8,
			instruction: "Insomniac, wake up and look at your card to see if it has changed.",
		},
	}

	// Add roles that are in play to the script
	for _, roleInfo := range roleOrder {
		if roleSet[roleInfo.role] {
			script = append(script, NightScriptStep{
				Role:        roleInfo.role,
				Order:       roleInfo.order,
				Instruction: roleInfo.instruction,
			})
		}
	}

	// Sort by order (should already be sorted, but to be safe)
	sort.Slice(script, func(i, j int) bool {
		return script[i].Order < script[j].Order
	})

	return script
}

// GetRoleInstructions returns detailed instructions for a specific role.
func GetRoleInstructions(role RoleType) string {
	instructions := map[RoleType]string{
		RoleWerewolf: "Wake up and look for other Werewolves. If you are the only Werewolf, you may view one center card.",
		RoleMinion: "Wake up. Werewolves, raise your hand so the Minion can see you.",
		RoleMason: "Wake up and look for other Masons.",
		RoleSeer: "Wake up. You may look at another player's card or two of the center cards.",
		RoleRobber: "Wake up. You may exchange your card with another player's card, and then view your new card.",
		RoleTroublemaker: "Wake up. You may exchange cards between two other players without looking at those cards.",
		RoleDrunk: "Wake up and exchange your card with a card from the center without looking at your new card.",
		RoleInsomniac: "Wake up and look at your card to see if it has changed.",
		RoleVillager: "You have no night action. Stay asleep.",
		RoleTanner: "You have no night action. Stay asleep.",
		RoleHunter: "You have no night action. Stay asleep.",
	}

	if instr, ok := instructions[role]; ok {
		return instr
	}
	return "No special instructions for this role."
}


