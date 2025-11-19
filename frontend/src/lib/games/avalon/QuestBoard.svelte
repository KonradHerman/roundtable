<script lang="ts">
	/**
	 * QuestBoard - Visual display of all 5 quests with their status
	 * Shows current quest, completed quests, team sizes, and 2-fail indicator
	 */
	import { questSizes, requiresTwoFails } from './roleConfig';

	let {
		currentQuest,
		questResults,
		playerCount
	}: {
		currentQuest: number;
		questResults: Array<{
			quest_number: number;
			success: boolean;
			team_size: number;
			cards: ('success' | 'fail')[];
			fail_count: number;
		}>;
		playerCount: number;
	} = $props();

	const teamSizes = $derived(questSizes[playerCount] || questSizes[5]);

	function getQuestStatus(questNum: number): 'pending' | 'success' | 'fail' {
		const result = questResults.find((r) => r.quest_number === questNum);
		if (!result) return 'pending';
		return result.success ? 'success' : 'fail';
	}

	// Count wins for each team
	const goodWins = $derived(questResults.filter((r) => r.success).length);
	const evilWins = $derived(questResults.filter((r) => !r.success).length);
</script>

<div class="quest-board bg-[#282828] p-6 rounded-lg mb-6">
	<h3 class="board-title text-center text-[#d79921] font-bold text-xl mb-4">Quest Progress</h3>

	<!-- Score display -->
	<div class="score-display flex justify-center gap-8 mb-4 text-lg font-bold">
		<div class="text-[#458588]">Good: {goodWins}</div>
		<div class="text-[#a89984]">/</div>
		<div class="text-[#cc241d]">Evil: {evilWins}</div>
	</div>

	<!-- Quest track -->
	<div class="quests flex justify-between gap-2">
		{#each [1, 2, 3, 4, 5] as questNum}
			{@const status = getQuestStatus(questNum)}
			{@const isActive = questNum === currentQuest}
			{@const needsTwoFails = requiresTwoFails(playerCount, questNum)}

			<div
				class="quest flex-1 bg-[#3c3836] p-3 rounded text-center border-2 transition-all
					{isActive ? 'border-[#d79921] bg-[#504945]' : 'border-transparent'}
					{status === 'success' ? 'bg-[#458588]' : ''}
					{status === 'fail' ? 'bg-[#cc241d]' : ''}"
			>
				<div class="quest-number font-bold text-sm mb-1">Quest {questNum}</div>
				<div class="team-size text-xs text-[#a89984] mb-2">
					Team: {teamSizes[questNum - 1]}
					{#if needsTwoFails}
						<span class="special text-[#fe8019] font-bold">**</span>
					{/if}
				</div>

				{#if status !== 'pending'}
					<div class="result-icon text-2xl">
						{status === 'success' ? '✅' : '❌'}
					</div>
				{:else if isActive}
					<div class="current-indicator text-[#d79921] text-xl">▶</div>
				{/if}
			</div>
		{/each}
	</div>

	{#if requiresTwoFails(playerCount, 4)}
		<p class="footnote text-center mt-3 text-xs text-[#a89984]">
			** Quest 4 requires 2 FAIL cards to fail
		</p>
	{/if}
</div>

<style>
	.quest {
		min-width: 60px;
	}

	@media (max-width: 640px) {
		.quests {
			gap: 0.25rem;
		}

		.quest {
			padding: 0.5rem 0.25rem;
			min-width: 50px;
		}

		.quest-number {
			font-size: 0.75rem;
		}

		.team-size {
			font-size: 0.625rem;
		}
	}
</style>
