/**
 * Shared role configuration for consistent role display across Werewolf game
 */

export interface RoleInfo {
	emoji: string;
	color: string;
	team?: string;
	description?: string;
}

export const roleConfig: Record<string, RoleInfo> = {
	werewolf: {
		emoji: '🐺',
		color: 'bg-red-600',
		team: 'Werewolf Team',
		description: 'Find your fellow werewolves and survive the vote'
	},
	seer: {
		emoji: '🔮',
		color: 'bg-purple-600',
		team: 'Village Team',
		description: "Look at one player's role to help find the werewolves"
	},
	robber: {
		emoji: '🎭',
		color: 'bg-blue-600',
		team: 'Village Team',
		description: 'Swap roles with another player'
	},
	troublemaker: {
		emoji: '😈',
		color: 'bg-orange-600',
		team: 'Village Team',
		description: "Swap two other players' roles"
	},
	mason: {
		emoji: '🔨',
		color: 'bg-gray-600',
		team: 'Village Team',
		description: 'Know who the other mason is'
	},
	villager: {
		emoji: '👤',
		color: 'bg-green-600',
		team: 'Village Team',
		description: 'Use your wits to find the werewolves'
	},
	minion: {
		emoji: '😤',
		color: 'bg-red-700',
		team: 'Werewolf Team',
		description: "Know the werewolves but they don't know you"
	},
	tanner: {
		emoji: '🤪',
		color: 'bg-yellow-600',
		team: 'Solo',
		description: 'You win if YOU get eliminated'
	},
	drunk: {
		emoji: '🍺',
		color: 'bg-amber-600',
		team: 'Village Team',
		description: "You must swap your role but won't know your new role"
	},
	insomniac: {
		emoji: '😴',
		color: 'bg-purple-700',
		team: 'Village Team',
		description: 'Wake up last to see if your role changed'
	},
	hunter: {
		emoji: '🏹',
		color: 'bg-green-700',
		team: 'Village Team',
		description: 'If you die, the player you voted for also dies'
	}
};

export function getRoleInfo(role: string): RoleInfo {
	return roleConfig[role] || { emoji: '❓', color: 'bg-muted' };
}

