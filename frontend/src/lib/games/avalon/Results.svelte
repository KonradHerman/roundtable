<script lang="ts">
	/**
	 * Results - Final game results showing winner, all roles, and quest history
	 */
	import { roleConfig, type AvalonRole, type Team } from './roleConfig';
	import QuestBoard from './QuestBoard.svelte';
	import RoleCard from './RoleCard.svelte';

	let {
		winningTeam,
		winReason,
		roles,
		teams,
		questHistory,
		players,
		playerCount
	}: {
		winningTeam: Team;
		winReason: string;
		roles: Record<string, AvalonRole>;
		teams: Record<string, Team>;
		questHistory: any[];
		players: Array<{ id: string; name: string }>;
		playerCount: number;
	} = $props();

	function getPlayerName(id: string): string {
		return players.find((p) => p.id === id)?.name || 'Unknown';
	}

	function getWinReasonText(reason: string): string {
		const reasons: Record<string, string> = {
			good_won_three_quests: 'Good won 3 quests!',
			evil_sabotaged_three_quests: 'Evil sabotaged 3 quests!',
			five_consecutive_rejections: 'Evil won after 5 consecutive rejections!',
			assassin_found_merlin: 'The Assassin correctly identified Merlin!',
			assassin_failed: 'The Assassin failed to identify Merlin!'
		};
		return reasons[reason] || reason;
	}

	// Group players by team
	const goodPlayers = $derived(
		players.filter((p) => teams[p.id] === 'good').map((p) => p.id)
	);
	const evilPlayers = $derived(
		players.filter((p) => teams[p.id] === 'evil').map((p) => p.id)
	);
</script>

<div class="results max-w-4xl mx-auto p-4">
	<!-- Winner banner -->
	<div
		class="winner-banner {winningTeam} text-white text-center py-6 px-6 rounded-xl mb-6 shadow-lg"
	>
		<div class="text-6xl mb-3">
			{winningTeam === 'good' ? '⚔️' : '💀'}
		</div>
		<h2 class="font-bold text-3xl mb-2">
			{winningTeam === 'good' ? 'Good Team Wins!' : 'Evil Team Wins!'}
		</h2>
		<p class="text-lg opacity-90">
			{getWinReasonText(winReason)}
		</p>
	</div>

	<!-- Quest history -->
	<div class="mb-6">
		<QuestBoard currentQuest={6} questResults={questHistory} {playerCount} />
	</div>

	<!-- All roles revealed -->
	<div class="roles-section bg-[#282828] p-6 rounded-lg">
		<h3 class="text-[#d79921] font-bold text-xl mb-6 text-center">All Roles Revealed</h3>

		<!-- Good team -->
		<div class="team-section mb-6">
			<h4 class="text-[#458588] font-bold text-lg mb-4 flex items-center gap-2">
				<span>⚔️</span>
				<span>Good Team</span>
			</h4>
			<div class="players-grid grid grid-cols-1 sm:grid-cols-2 gap-4">
				{#each goodPlayers as playerId}
					{@const role = roles[playerId]}
					{@const config = roleConfig[role]}
					<div class="player-role bg-[#3c3836] p-4 rounded-lg flex items-center gap-4">
						<div class="role-emoji text-4xl">{config.emoji}</div>
						<div class="flex-1">
							<div class="player-name text-[#ebdbb2] font-bold mb-1">
								{getPlayerName(playerId)}
							</div>
							<div class="role-name text-[#a89984] text-sm">{config.name}</div>
						</div>
					</div>
				{/each}
			</div>
		</div>

		<!-- Evil team -->
		<div class="team-section">
			<h4 class="text-[#cc241d] font-bold text-lg mb-4 flex items-center gap-2">
				<span>💀</span>
				<span>Evil Team</span>
			</h4>
			<div class="players-grid grid grid-cols-1 sm:grid-cols-2 gap-4">
				{#each evilPlayers as playerId}
					{@const role = roles[playerId]}
					{@const config = roleConfig[role]}
					<div class="player-role bg-[#3c3836] p-4 rounded-lg flex items-center gap-4">
						<div class="role-emoji text-4xl">{config.emoji}</div>
						<div class="flex-1">
							<div class="player-name text-[#ebdbb2] font-bold mb-1">
								{getPlayerName(playerId)}
							</div>
							<div class="role-name text-[#a89984] text-sm">{config.name}</div>
						</div>
					</div>
				{/each}
			</div>
		</div>
	</div>

	<!-- Quest details (optional expansion) -->
	{#if questHistory.length > 0}
		<details class="quest-details bg-[#3c3836] p-4 rounded-lg mt-6">
			<summary class="text-[#d79921] font-bold cursor-pointer mb-4">
				📊 Quest Details
			</summary>
			<div class="space-y-3">
				{#each questHistory as quest}
					<div class="quest-detail bg-[#282828] p-3 rounded">
						<div class="flex justify-between items-center mb-2">
							<span class="font-bold">Quest {quest.quest_number}</span>
							<span class="text-2xl">{quest.success ? '✅' : '❌'}</span>
						</div>
						<div class="text-sm text-[#a89984]">
							Team: {quest.team_members.map((id) => getPlayerName(id)).join(', ')}
						</div>
						<div class="text-sm text-[#a89984] mt-1">
							Results: {quest.cards.map((c) => (c === 'success' ? '✅' : '❌')).join(' ')}
							({quest.fail_count} fail{quest.fail_count !== 1 ? 's' : ''})
						</div>
					</div>
				{/each}
			</div>
		</details>
	{/if}
</div>

<style>
	.winner-banner.good {
		background: linear-gradient(135deg, #458588 0%, #689d6a 100%);
	}

	.winner-banner.evil {
		background: linear-gradient(135deg, #cc241d 0%, #fb4934 100%);
	}

	details summary {
		list-style: none;
	}

	details summary::-webkit-details-marker {
		display: none;
	}

	details[open] summary {
		margin-bottom: 1rem;
	}
</style>
