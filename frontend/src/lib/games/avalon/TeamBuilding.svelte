<script lang="ts">
	/**
	 * TeamBuilding - Leader selects team members for the quest
	 * Only the leader can interact, others see waiting state
	 */
	import QuestBoard from './QuestBoard.svelte';

	let {
		isLeader,
		currentQuest,
		questResults,
		playerCount,
		requiredTeamSize,
		players,
		leaderId,
		rejectionCount,
		onProposeTeam
	}: {
		isLeader: boolean;
		currentQuest: number;
		questResults: any[];
		playerCount: number;
		requiredTeamSize: number;
		players: Array<{ id: string; name: string }>;
		leaderId: string;
		rejectionCount: number;
		onProposeTeam: (selectedIds: string[]) => void;
	} = $props();

	let selectedPlayers = $state<string[]>([]);

	const leaderName = $derived(players.find((p) => p.id === leaderId)?.name || 'Unknown');

	function togglePlayer(playerId: string) {
		if (selectedPlayers.includes(playerId)) {
			selectedPlayers = selectedPlayers.filter((id) => id !== playerId);
		} else if (selectedPlayers.length < requiredTeamSize) {
			selectedPlayers = [...selectedPlayers, playerId];
		}
	}

	function handlePropose() {
		if (selectedPlayers.length === requiredTeamSize) {
			onProposeTeam(selectedPlayers);
			selectedPlayers = []; // Reset for next round
		}
	}

	const canPropose = $derived(isLeader && selectedPlayers.length === requiredTeamSize);
</script>

<div class="team-building max-w-2xl mx-auto p-4">
	<QuestBoard {currentQuest} {questResults} {playerCount} />

	<div class="leader-banner bg-[#d79921] text-[#282828] text-center py-3 px-4 rounded-lg mb-6 font-bold">
		👑 {leaderName} is the Leader
	</div>

	{#if rejectionCount > 0}
		<div class="rejection-warning bg-[#cc241d] text-white text-center py-2 px-4 rounded mb-4">
			⚠️ {rejectionCount} consecutive rejection{rejectionCount > 1 ? 's' : ''} ({5 - rejectionCount} left before Evil wins)
		</div>
	{/if}

	{#if isLeader}
		<div class="instructions bg-[#3c3836] p-4 rounded-lg mb-4">
			<h3 class="text-[#d79921] mb-2 font-bold">Select {requiredTeamSize} players for the quest</h3>
			<p class="text-[#a89984] text-sm">
				Selected: {selectedPlayers.length} / {requiredTeamSize}
			</p>
		</div>

		<!-- Player selection grid -->
		<div class="player-grid grid grid-cols-2 sm:grid-cols-3 gap-3 mb-4">
			{#each players as player}
				{@const isSelected = selectedPlayers.includes(player.id)}
				<button
					class="player-btn bg-[#3c3836] text-[#ebdbb2] py-4 px-3 rounded-lg border-2 transition-all
						{isSelected ? 'border-[#d79921] bg-[#504945]' : 'border-transparent'}
						hover:bg-[#504945]"
					onclick={() => togglePlayer(player.id)}
				>
					<div class="text-center">
						<div class="text-lg mb-1">
							{isSelected ? '✓' : ''}
						</div>
						<div class="font-medium truncate">{player.name}</div>
					</div>
				</button>
			{/each}
		</div>

		<button
			class="btn-primary w-full py-4 text-lg font-bold rounded-lg"
			disabled={!canPropose}
			onclick={handlePropose}
		>
			Propose Team ({selectedPlayers.length}/{requiredTeamSize})
		</button>
	{:else}
		<div class="waiting bg-[#3c3836] p-8 rounded-lg text-center">
			<p class="text-[#a89984] text-lg">Waiting for {leaderName} to propose a team...</p>
		</div>
	{/if}
</div>

<style>
	.btn-primary {
		background: #d79921;
		color: #282828;
		border: none;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-primary:hover:not(:disabled) {
		background: #fabd2f;
		transform: translateY(-2px);
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.player-btn {
		min-height: 80px;
		cursor: pointer;
	}

	.player-btn:active {
		transform: scale(0.98);
	}
</style>
