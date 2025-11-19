<script lang="ts">
	/**
	 * QuestExecution - Team members secretly play Success or Fail cards
	 * Good players can ONLY play Success, Evil players can choose
	 */
	import QuestBoard from './QuestBoard.svelte';

	let {
		currentQuest,
		questResults,
		playerCount,
		teamMembers,
		players,
		isOnTeam,
		team,
		hasPlayedCard,
		cardsSubmitted,
		totalCardsExpected,
		onPlayCard
	}: {
		currentQuest: number;
		questResults: any[];
		playerCount: number;
		teamMembers: string[];
		players: Array<{ id: string; name: string }>;
		isOnTeam: boolean;
		team: 'good' | 'evil';
		hasPlayedCard: boolean;
		cardsSubmitted: number;
		totalCardsExpected: number;
		onPlayCard: (card: 'success' | 'fail') => void;
	} = $props();

	let selectedCard = $state<'success' | 'fail' | null>(null);

	function getPlayerName(id: string): string {
		return players.find((p) => p.id === id)?.name || 'Unknown';
	}

	function handlePlayCard(card: 'success' | 'fail') {
		selectedCard = card;
		onPlayCard(card);
	}

	const canPlayFail = $derived(team === 'evil');
</script>

<div class="quest-execution max-w-2xl mx-auto p-4">
	<QuestBoard {currentQuest} {questResults} {playerCount} />

	<!-- Team members display -->
	<div class="team-display bg-[#3c3836] p-6 rounded-lg mb-6">
		<h3 class="text-[#d79921] font-bold mb-4 text-center">Quest Team</h3>
		<div class="members flex flex-wrap justify-center gap-3">
			{#each teamMembers as memberId}
				<div class="member bg-[#504945] px-4 py-2 rounded-lg text-[#ebdbb2] font-medium">
					{getPlayerName(memberId)}
				</div>
			{/each}
		</div>
	</div>

	{#if isOnTeam}
		{#if !hasPlayedCard}
			<!-- Card selection -->
			<div class="card-selection bg-[#282828] p-6 rounded-lg">
				<h3 class="text-center text-[#d79921] font-bold mb-4 text-lg">Play Your Quest Card</h3>

				{#if team === 'good'}
					<p class="text-center text-[#a89984] mb-6 text-sm">
						As a Good player, you must play <span class="text-[#458588] font-bold">Success</span>
					</p>
				{:else}
					<p class="text-center text-[#a89984] mb-6 text-sm">
						As Evil, you may choose to <span class="text-[#458588]">help</span> or <span class="text-[#cc241d]">sabotage</span>
					</p>
				{/if}

				<div class="quest-cards grid {canPlayFail ? 'grid-cols-2' : 'grid-cols-1'} gap-4">
					<!-- Success card -->
					<button
						class="quest-card success bg-[#458588] hover:bg-[#689d6a] text-white py-8 px-6 rounded-xl font-bold text-2xl transition-all border-4 border-transparent hover:border-[#8ec07c]"
						onclick={() => handlePlayCard('success')}
					>
						<div class="text-5xl mb-2">✅</div>
						<div>SUCCESS</div>
					</button>

					<!-- Fail card (Evil only) -->
					{#if canPlayFail}
						<button
							class="quest-card fail bg-[#cc241d] hover:bg-[#fb4934] text-white py-8 px-6 rounded-xl font-bold text-2xl transition-all border-4 border-transparent hover:border-[#fb4934]"
							onclick={() => handlePlayCard('fail')}
						>
							<div class="text-5xl mb-2">❌</div>
							<div>FAIL</div>
						</button>
					{/if}
				</div>

				{#if !canPlayFail}
					<p class="text-center text-[#a89984] mt-4 text-xs italic">
						Good players cannot play Fail cards
					</p>
				{/if}
			</div>
		{:else}
			<!-- Card played confirmation -->
			<div class="card-played bg-[#3c3836] p-8 rounded-lg text-center">
				<div class="text-6xl mb-4">
					{selectedCard === 'success' ? '✅' : '❌'}
				</div>
				<p class="text-[#ebdbb2] font-bold text-xl mb-2">Card Played</p>
				<p class="text-[#a89984]">
					You played: <span class="font-bold uppercase">{selectedCard}</span>
				</p>
				<p class="text-[#a89984] mt-4 text-sm">
					Waiting for other team members...
				</p>
			</div>
		{/if}
	{:else}
		<!-- Not on team - waiting -->
		<div class="waiting bg-[#3c3836] p-8 rounded-lg text-center">
			<p class="text-[#a89984] text-lg mb-4">The quest team is playing their cards...</p>
			<p class="text-[#a89984] text-sm">
				Cards played: {cardsSubmitted} / {totalCardsExpected}
			</p>
		</div>
	{/if}

	<!-- Progress indicator -->
	<div class="progress-section mt-6 bg-[#3c3836] p-4 rounded-lg">
		<div class="flex justify-between text-sm text-[#a89984] mb-2">
			<span>Quest cards played</span>
			<span>{cardsSubmitted} / {totalCardsExpected}</span>
		</div>
		<div class="progress-bar bg-[#282828] rounded-full h-3 overflow-hidden">
			<div
				class="progress-fill bg-[#d79921] h-full transition-all duration-300"
				style="width: {(cardsSubmitted / totalCardsExpected) * 100}%"
			></div>
		</div>
	</div>
</div>

<style>
	.quest-card {
		cursor: pointer;
		border: 4px solid transparent;
		min-height: 180px;
	}

	.quest-card:hover {
		transform: translateY(-4px);
		box-shadow: 0 8px 16px rgba(0, 0, 0, 0.3);
	}

	.quest-card:active {
		transform: translateY(-2px);
	}

	@media (max-width: 640px) {
		.quest-card {
			min-height: 140px;
			font-size: 1.25rem;
			padding: 1.5rem 1rem;
		}

		.quest-card .text-5xl {
			font-size: 2.5rem;
		}
	}
</style>
