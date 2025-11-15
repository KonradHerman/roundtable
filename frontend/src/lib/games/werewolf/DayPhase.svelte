<script lang="ts">
	import { gameStore } from '$lib/stores/game';
	import { session } from '$lib/stores/session';
	import { Card, Button, Badge } from '$lib/components/ui';
	import { Vote, CheckCircle2 } from 'lucide-svelte';

	export let roomState: any;
	export let wsStore: any;

	let selectedPlayer: string | null = null;
	let hasVoted = false;
	let votes: Record<string, string> = {};
	let votesRevealed = false;

	// Subscribe to vote events
	let unsubscribe = gameStore.subscribe(($game) => {
		$game.events.forEach(event => {
			if (event.type === 'vote_cast') {
				// Track that someone voted (but not who they voted for yet)
				hasVoted = hasVoted || event.actorId === $session?.playerId;
			} else if (event.type === 'votes_revealed') {
				votes = event.payload.votes || {};
				votesRevealed = true;
			}
		});
	});

	function handleVote() {
		if (!selectedPlayer || !wsStore || hasVoted) return;

		wsStore.sendAction({
			type: 'vote',
			payload: {
				targetId: selectedPlayer
			}
		});

		hasVoted = true;
	}

	function getVoteCount(playerId: string): number {
		return Object.values(votes).filter(v => v === playerId).length;
	}

	function didVoteFor(playerId: string): boolean {
		return votes[$session?.playerId || ''] === playerId;
	}

	$: players = roomState?.players || [];
</script>

<div class="space-y-6">
	<!-- Instructions -->
	<Card class="p-6 bg-amber-500/10 border-amber-500/20">
		<div class="flex items-start gap-3">
			<Vote class="w-6 h-6 text-amber-600 dark:text-amber-400 flex-shrink-0 mt-1" />
			<div>
				<h3 class="font-semibold text-lg mb-1">Discussion & Voting</h3>
				<p class="text-sm text-muted-foreground">
					Discuss with your group who you think is a werewolf, then vote to eliminate someone.
				</p>
			</div>
		</div>
	</Card>

	<!-- Voting interface -->
	{#if !votesRevealed}
		<Card class="p-6">
			<h3 class="font-semibold text-lg mb-4">
				{hasVoted ? 'Your vote has been cast' : 'Vote to eliminate'}
			</h3>

			{#if hasVoted}
				<div class="flex items-center gap-3 p-4 bg-green-500/10 border border-green-500/20 rounded-lg mb-4">
					<CheckCircle2 class="w-6 h-6 text-green-600" />
					<div>
						<p class="font-medium">Vote submitted</p>
						<p class="text-sm text-muted-foreground">Waiting for other players...</p>
					</div>
				</div>
			{:else}
				<div class="space-y-3 mb-4">
					{#each players as player}
						<button
							on:click={() => selectedPlayer = player.id}
							class="w-full p-4 text-left rounded-lg border-2 transition-all {selectedPlayer === player.id ? 'border-destructive bg-destructive/10' : 'border-border hover:border-destructive/50'}"
							disabled={hasVoted}
						>
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-3">
									<div class="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold">
										{player.displayName[0].toUpperCase()}
									</div>
									<span class="font-medium">{player.displayName}</span>
								</div>
								{#if selectedPlayer === player.id}
									<Badge variant="destructive">Selected</Badge>
								{/if}
							</div>
						</button>
					{/each}
				</div>

				<Button
					variant="destructive"
					class="w-full h-12"
					disabled={!selectedPlayer || hasVoted}
					on:click={handleVote}
				>
					Submit Vote
				</Button>
			{/if}

			<!-- Vote status -->
			<div class="mt-4 p-3 bg-muted rounded-lg">
				<p class="text-sm text-center text-muted-foreground">
					{#if hasVoted}
						Waiting for other players to vote...
					{:else}
						Select a player to eliminate
					{/if}
				</p>
			</div>
		</Card>
	{:else}
		<!-- Votes revealed -->
		<Card class="p-6">
			<h3 class="font-semibold text-lg mb-4">Votes Revealed</h3>

			<div class="space-y-2">
				{#each players as player}
					{@const voteCount = getVoteCount(player.id)}
					{@const votedByMe = didVoteFor(player.id)}
					<div class="p-4 rounded-lg border {voteCount > 0 ? 'border-destructive bg-destructive/5' : 'border-border'}">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-3">
								<div class="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold">
									{player.displayName[0].toUpperCase()}
								</div>
								<div>
									<span class="font-medium">{player.displayName}</span>
									{#if votedByMe}
										<span class="text-xs text-muted-foreground ml-2">(You voted for this player)</span>
									{/if}
								</div>
							</div>
							<Badge variant={voteCount > 0 ? 'destructive' : 'outline'}>
								{voteCount} {voteCount === 1 ? 'vote' : 'votes'}
							</Badge>
						</div>
					</div>
				{/each}
			</div>

			<div class="mt-6 p-4 bg-muted rounded-lg">
				<p class="text-sm text-center text-muted-foreground">
					Calculating results...
				</p>
			</div>
		</Card>
	{/if}
</div>
