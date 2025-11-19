export type AvalonRole =
	| 'merlin'
	| 'percival'
	| 'loyal_servant'
	| 'assassin'
	| 'morgana'
	| 'mordred'
	| 'oberon'
	| 'minion';

export type Team = 'good' | 'evil';

export interface RoleConfig {
	name: string;
	emoji: string;
	team: Team;
	color: string;
	description: string;
	knowledge?: string;
}

export const roleConfig: Record<AvalonRole, RoleConfig> = {
	merlin: {
		name: 'Merlin',
		emoji: '🔮',
		team: 'good',
		color: 'bg-blue-600',
		description: 'Knows the forces of Evil (except Mordred)',
		knowledge:
			'You see all Evil players (except Mordred if present). Help Good win without revealing yourself!'
	},
	percival: {
		name: 'Percival',
		emoji: '👁️',
		team: 'good',
		color: 'bg-cyan-600',
		description: 'Sees Merlin and Morgana (cannot distinguish)',
		knowledge:
			'You see two powerful wizards. One is Merlin, one is Morgana. Protect Merlin!'
	},
	loyal_servant: {
		name: 'Loyal Servant of Arthur',
		emoji: '⚔️',
		team: 'good',
		color: 'bg-slate-600',
		description: 'No special knowledge, must rely on deduction',
		knowledge:
			'You have no special information. Trust your instincts and your allies!'
	},
	assassin: {
		name: 'Assassin',
		emoji: '🗡️',
		team: 'evil',
		color: 'bg-red-700',
		description: 'Can assassinate Merlin if Good wins',
		knowledge:
			'You know your Evil allies. If Good wins 3 quests, you can steal victory by identifying Merlin!'
	},
	morgana: {
		name: 'Morgana',
		emoji: '🌙',
		team: 'evil',
		color: 'bg-purple-700',
		description: 'Appears as Merlin to Percival',
		knowledge: 'You appear as Merlin to Percival. Confuse the Good team!'
	},
	mordred: {
		name: 'Mordred',
		emoji: '😈',
		team: 'evil',
		color: 'bg-orange-700',
		description: 'Hidden from Merlin',
		knowledge: 'Merlin cannot see you. Use this advantage wisely!'
	},
	oberon: {
		name: 'Oberon',
		emoji: '👻',
		team: 'evil',
		color: 'bg-gray-700',
		description: 'Unknown to other Evil players',
		knowledge:
			'You are alone. You do not know other Evil players, and they do not know you.'
	},
	minion: {
		name: 'Minion of Mordred',
		emoji: '💀',
		team: 'evil',
		color: 'bg-red-900',
		description: 'Knows other Evil players',
		knowledge: 'You know your Evil allies. Work together to sabotage the quests!'
	}
};

// Quest team sizes by player count
export const questSizes: Record<number, number[]> = {
	5: [2, 3, 2, 3, 3],
	6: [2, 3, 4, 3, 4],
	7: [2, 3, 3, 4, 4],
	8: [3, 4, 4, 5, 5],
	9: [3, 4, 4, 5, 5],
	10: [3, 4, 4, 5, 5]
};

// Returns true if quest 4 with this player count requires 2 fails
export function requiresTwoFails(playerCount: number, questNumber: number): boolean {
	return playerCount >= 7 && questNumber === 4;
}

// Returns the required team size for a quest
export function getRequiredTeamSize(playerCount: number, questNumber: number): number {
	const sizes = questSizes[playerCount] || questSizes[5];
	return sizes[questNumber - 1] || 0;
}
