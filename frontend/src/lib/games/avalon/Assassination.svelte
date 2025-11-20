<script lang="ts">
	/**
	 * Assassination - The Assassin attempts to identify Merlin
	 * Only triggers if Good wins 3 quests and Merlin is in the game
	 */
	import QuestBoard from './QuestBoard.svelte';

	let {
		currentQuest,
		questResults,
		playerCount,
		players,
		isAssassin,
		currentPlayerId,
		onAssassinate
	}: {
		currentQuest: number;
		questResults: any[];
		playerCount: number;
		players: Array<{ id: string; name: string }>;
		isAssassin: boolean;
		currentPlayerId: string;
		onAssassinate: (targetId: string) => void;
	} = $props();

	let selectedTarget = $state<string | null>(null);
	let confirmed = $state(false);

	function handleSelectTarget(playerId: string) {
		if (playerId !== currentPlayerId) {
			selectedTarget = playerId;
		}
	}

	function handleConfirm() {
		if (selectedTarget) {
			confirmed = true;
			onAssassinate(selectedTarget);
		}
	}

	const targetName = $derived(
		selectedTarget ? players.find((p) => p.id === selectedTarget)?.name : null
	);
</script>

<div class="assassination max-w-2xl mx-auto p-4">
	<QuestBoard {currentQuest} {questResults} {playerCount} />

	<div class="assassin-banner bg-[#cc241d] text-white text-center py-4 px-6 rounded-lg mb-6">
		<div class="text-3xl mb-2">🗡️</div>
		<h2 class="font-bold text-xl">Assassination Phase</h2>
		<p class="text-sm mt-2 opacity-90">Good won 3 quests, but Evil has one final chance...</p>
	</div>

	{#if isAssassin}
		{#if !confirmed}
			<!-- Assassin selection UI -->
			<div class="assassin-selection bg-[#282828] p-6 rounded-lg">
				<h3 class="text-[#d79921] font-bold mb-4 text-center text-lg">
					Select Who You Believe is Merlin
				</h3>
				<p class="text-[#a89984] text-center mb-6 text-sm">
					If you identify Merlin correctly, Evil wins. Choose wisely...
				</p>

				<!-- Player grid (excluding self) -->
				<div class="player-grid grid grid-cols-2 sm:grid-cols-3 gap-3 mb-6">
					{#each players as player}
						{#if player.id !== currentPlayerId}
							{@const isSelected = selectedTarget === player.id}
							<button
								class="player-btn bg-[#3c3836] text-[#ebdbb2] py-4 px-3 rounded-lg border-2 transition-all
									{isSelected ? 'border-[#cc241d] bg-[#cc241d] text-white' : 'border-transparent'}
									hover:bg-[#504945]"
								onclick={() => handleSelectTarget(player.id)}
							>
								<div class="text-center">
									<div class="text-xl mb-1">
										{isSelected ? '🗡️' : ''}
									</div>
									<div class="font-medium truncate">{player.name}</div>
								</div>
							</button>
						{/if}
					{/each}
				</div>

				{#if selectedTarget}
					<div class="confirmation bg-[#3c3836] p-4 rounded-lg mb-4">
						<p class="text-center text-[#ebdbb2] mb-2">
							Assassinate <span class="font-bold text-[#cc241d]">{targetName}</span>?
						</p>
						<p class="text-center text-[#a89984] text-xs">
							This decision cannot be undone
						</p>
					</div>
				{/if}

				<button
					class="btn-assassinate w-full py-4 text-lg font-bold rounded-lg"
					disabled={!selectedTarget}
					onclick={handleConfirm}
				>
					{selectedTarget ? `🗡️ Assassinate ${targetName}` : 'Select a Target'}
				</button>
			</div>
		{:else}
			<!-- Confirmation shown -->
			<div class="confirmed bg-[#3c3836] p-8 rounded-lg text-center">
				<div class="text-6xl mb-4">🗡️</div>
				<p class="text-[#ebdbb2] font-bold text-xl mb-2">Target Selected</p>
				<p class="text-[#a89984]">
					You have chosen to assassinate <span class="font-bold">{targetName}</span>
				</p>
				<p class="text-[#a89984] mt-4 text-sm">Revealing results...</p>
			</div>
		{/if}
	{:else}
		<!-- Non-assassin waiting -->
		<div class="waiting bg-[#3c3836] p-8 rounded-lg text-center">
			<div class="text-6xl mb-4">🗡️</div>
			<p class="text-[#a89984] text-lg mb-2">The Assassin is choosing their target...</p>
			<p class="text-[#a89984] text-sm">
				If the Assassin correctly identifies Merlin, Evil wins
			</p>
		</div>
	{/if}
</div>

<style>
	.btn-assassinate {
		background: #cc241d;
		color: white;
		border: none;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-assassinate:hover:not(:disabled) {
		background: #fb4934;
		transform: translateY(-2px);
	}

	.btn-assassinate:disabled {
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
