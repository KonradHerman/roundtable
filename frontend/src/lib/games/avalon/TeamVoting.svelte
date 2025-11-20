<script lang="ts">
	/**
	 * TeamVoting - All players vote to approve or reject the proposed team
	 * Votes are hidden until everyone submits
	 */
	import QuestBoard from './QuestBoard.svelte';

	let {
		currentQuest,
		questResults,
		playerCount,
		proposedTeam,
		players,
		hasVoted,
		votesSubmitted,
		totalVotes,
		rejectionCount,
		onVote
	}: {
		currentQuest: number;
		questResults: any[];
		playerCount: number;
		proposedTeam: string[];
		players: Array<{ id: string; name: string }>;
		hasVoted: boolean;
		votesSubmitted: number;
		totalVotes: number;
		rejectionCount: number;
		onVote: (vote: 'approve' | 'reject') => void;
	} = $props();

	let selectedVote = $state<'approve' | 'reject' | null>(null);

	function getPlayerName(id: string): string {
		return players.find((p) => p.id === id)?.name || 'Unknown';
	}

	function handleVote(vote: 'approve' | 'reject') {
		selectedVote = vote;
		onVote(vote);
	}
</script>

<div class="team-voting max-w-2xl mx-auto p-4">
	<QuestBoard {currentQuest} {questResults} {playerCount} />

	{#if rejectionCount > 0}
		<div class="rejection-warning bg-[#cc241d] text-white text-center py-2 px-4 rounded mb-4">
			⚠️ {rejectionCount} consecutive rejection{rejectionCount > 1 ? 's' : ''} ({5 - rejectionCount} left before Evil wins)
		</div>
	{/if}

	<!-- Proposed team display -->
	<div class="proposed-team bg-[#3c3836] p-6 rounded-lg mb-6">
		<h3 class="text-[#d79921] font-bold mb-4 text-center">Proposed Team</h3>
		<div class="team-members flex flex-wrap justify-center gap-3">
			{#each proposedTeam as playerId}
				<div class="member bg-[#504945] px-4 py-2 rounded-lg text-[#ebdbb2] font-medium">
					{getPlayerName(playerId)}
				</div>
			{/each}
		</div>
	</div>

	<!-- Vote submission -->
	{#if !hasVoted}
		<div class="vote-section bg-[#282828] p-6 rounded-lg">
			<h3 class="text-center text-[#d79921] font-bold mb-4 text-lg">Cast Your Vote</h3>
			<p class="text-center text-[#a89984] mb-6 text-sm">
				Do you approve or reject this team?
			</p>

			<div class="vote-buttons grid grid-cols-2 gap-4">
				<button
					class="vote-btn approve bg-[#458588] hover:bg-[#689d6a] text-white py-6 px-4 rounded-lg font-bold text-xl transition-all"
					onclick={() => handleVote('approve')}
				>
					✅ Approve
				</button>
				<button
					class="vote-btn reject bg-[#cc241d] hover:bg-[#fb4934] text-white py-6 px-4 rounded-lg font-bold text-xl transition-all"
					onclick={() => handleVote('reject')}
				>
					❌ Reject
				</button>
			</div>
		</div>
	{:else}
		<div class="voted bg-[#3c3836] p-8 rounded-lg text-center">
			<div class="text-4xl mb-4">
				{selectedVote === 'approve' ? '✅' : '❌'}
			</div>
			<p class="text-[#ebdbb2] font-bold text-lg mb-2">Vote Recorded</p>
			<p class="text-[#a89984]">
				You voted: <span class="font-bold">{selectedVote === 'approve' ? 'APPROVE' : 'REJECT'}</span>
			</p>
			<p class="text-[#a89984] mt-4 text-sm">
				Votes submitted: {votesSubmitted} / {totalVotes}
			</p>
		</div>
	{/if}

	<!-- Vote progress -->
	<div class="vote-progress mt-6 bg-[#3c3836] p-4 rounded-lg">
		<div class="flex justify-between text-sm text-[#a89984] mb-2">
			<span>Votes cast</span>
			<span>{votesSubmitted} / {totalVotes}</span>
		</div>
		<div class="progress-bar bg-[#282828] rounded-full h-3 overflow-hidden">
			<div
				class="progress-fill bg-[#d79921] h-full transition-all duration-300"
				style="width: {(votesSubmitted / totalVotes) * 100}%"
			></div>
		</div>
	</div>
</div>

<style>
	.vote-btn {
		cursor: pointer;
		border: none;
		min-height: 100px;
	}

	.vote-btn:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
	}

	.vote-btn:active {
		transform: translateY(0);
	}

	@media (max-width: 640px) {
		.vote-btn {
			min-height: 80px;
			font-size: 1.125rem;
		}
	}
</style>
