<script lang="ts">
	import { gameStore } from '$lib/stores/game';
	import { session } from '$lib/stores/session';
	import { Card, Button, Badge } from '$lib/components/ui';
	import { Vote, CheckCircle2, Clock, Play, Pause, Plus } from 'lucide-svelte';
	import { onMount, onDestroy } from 'svelte';

	export let roomState: any;
	export let wsStore: any;
	export let timerActive: boolean = false;
	export let phaseEndsAt: Date | null = null;

	let selectedPlayer: string | null = null;
	let hasVoted = false;
	let votes: Record<string, string> = {};
	let votesRevealed = false;
	let timeRemaining: number = 0;
	let timerInterval: any = null;

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

	onMount(() => {
		timerInterval = setInterval(() => {
			if (timerActive && phaseEndsAt) {
				const now = new Date().getTime();
				const end = new Date(phaseEndsAt).getTime();
				timeRemaining = Math.max(0, Math.floor((end - now) / 1000));
			} else {
				timeRemaining = 0;
			}
		}, 100);
	});

	onDestroy(() => {
		if (timerInterval) clearInterval(timerInterval);
		if (unsubscribe) unsubscribe();
	});

	function handleVote() {
		if (!selectedPlayer || !wsStore) return;

		wsStore.sendAction({
			type: 'vote',
			payload: {
				targetId: selectedPlayer
			}
		});

		hasVoted = true;
	}

	function handleToggleTimer() {
		if (!wsStore) return;

		wsStore.sendAction({
			type: 'toggle_timer',
			payload: {
				enable: !timerActive,
				duration: 180 // 3 minutes default
			}
		});
	}

	function handleExtendTimer() {
		if (!wsStore) return;

		wsStore.sendAction({
			type: 'extend_timer',
			payload: {
				seconds: 60 // Add 1 minute
			}
		});
	}

	function formatTime(seconds: number): string {
		const mins = Math.floor(seconds / 60);
		const secs = seconds % 60;
		return `${mins}:${secs.toString().padStart(2, '0')}`;
	}

	function getVoteCount(playerId: string): number {
		return Object.values(votes).filter(v => v === playerId).length;
	}

	function didVoteFor(playerId: string): boolean {
		return votes[$session?.playerId || ''] === playerId;
	}

	$: players = roomState?.players || [];
	$: isHost = $session?.playerId === roomState?.hostId;
	$: votesSubmitted = Object.keys(votes).length;
</script>

<div class="space-y-6">
	<!-- Timer Display & Controls -->
	<Card class="p-6 bg-gruvbox-orange/20 border-gruvbox-orange">
		<div class="flex items-center justify-between flex-wrap gap-4">
			<div class="flex items-center gap-3">
				<Clock class="w-6 h-6 text-gruvbox-orange-light" />
				<div>
					<h3 class="font-semibold text-lg">Timer</h3>
					<p class="text-sm text-muted-foreground">
						{#if timerActive}
							<span class="text-gruvbox-orange-light font-mono text-xl">{formatTime(timeRemaining)}</span>
						{:else}
							<span class="text-muted-foreground">OFF</span>
						{/if}
					</p>
				</div>
			</div>
			
			{#if isHost}
				<div class="flex gap-2">
					<Button
						on:click={handleToggleTimer}
						variant="outline"
						class="flex items-center gap-2"
					>
						{#if timerActive}
							<Pause class="w-4 h-4" />
							Pause Timer
						{:else}
							<Play class="w-4 h-4" />
							Start Timer
						{/if}
					</Button>
					
					{#if timerActive}
						<Button
							on:click={handleExtendTimer}
							variant="outline"
							class="flex items-center gap-2"
						>
							<Plus class="w-4 h-4" />
							+1 Min
						</Button>
					{/if}
				</div>
			{/if}
		</div>
	</Card>

	<!-- Instructions -->
	<Card class="p-6 bg-gruvbox-yellow/20 border-gruvbox-yellow">
		<div class="flex items-start gap-3">
			<Vote class="w-6 h-6 text-gruvbox-yellow-light flex-shrink-0 mt-1" />
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
			<div class="flex items-center justify-between mb-4">
				<h3 class="font-semibold text-lg">
					{hasVoted ? 'Your vote has been cast' : 'Vote to eliminate'}
				</h3>
				<Badge variant="outline">
					{votesSubmitted} / {players.length} voted
				</Badge>
			</div>

			{#if hasVoted}
				<div class="flex items-center gap-3 p-4 bg-gruvbox-green/20 border border-gruvbox-green rounded-lg mb-4">
					<CheckCircle2 class="w-6 h-6 text-gruvbox-green-light" />
					<div>
						<p class="font-medium">Vote submitted</p>
						<p class="text-sm text-muted-foreground">You can change your vote before everyone submits</p>
					</div>
				</div>
			{/if}

			<div class="space-y-3 mb-4">
				{#each players as player}
					<button
						on:click={() => selectedPlayer = player.id}
						class="w-full p-4 text-left rounded-lg border-2 transition-all {selectedPlayer === player.id ? 'border-primary bg-primary/10' : 'border-border hover:border-primary/50'}"
					>
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-3">
								<div class="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold">
									{player.displayName[0].toUpperCase()}
								</div>
								<span class="font-medium">{player.displayName}</span>
							</div>
							{#if selectedPlayer === player.id}
								<Badge class="bg-primary">Selected</Badge>
							{/if}
						</div>
					</button>
				{/each}
			</div>

			<Button
				variant="destructive"
				class="w-full h-12"
				disabled={!selectedPlayer}
				on:click={handleVote}
			>
				{hasVoted ? 'Change Vote' : 'Submit Vote'}
			</Button>

			<!-- Vote status -->
			<div class="mt-4 p-3 bg-muted rounded-lg">
				<p class="text-sm text-center text-muted-foreground">
					{#if hasVoted}
						Waiting for other players to vote... ({votesSubmitted}/{players.length})
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
