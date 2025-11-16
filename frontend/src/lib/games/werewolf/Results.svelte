<script lang="ts">
	import { gameStore } from '$lib/stores/game';
	import { session } from '$lib/stores/session';
	import { Card, Badge, Button } from '$lib/components/ui';
	import { Trophy, Skull } from 'lucide-svelte';
	import { confetti } from '@neoconfetti/svelte';

	export let roomState: any;
	
	import { api } from '$lib/api/client';
	
	let gameResults: any = null;
	let allRoles: Record<string, string> = {};
	let eliminated: string[] = [];
	let winners: string[] = [];
	let winReason: string = '';
	let showConfetti = false;
	let hasGameFinished = false;

	// Subscribe to game events
	let unsubscribe = gameStore.subscribe(($game) => {
		$game.events.forEach(event => {
			if (event.type === 'roles_revealed') {
				// Roles revealed for display (before voting/game finished)
				allRoles = event.payload.roles || {};
			} else if (event.type === 'game_finished') {
				// Game finished with winners determined
				hasGameFinished = true;
				gameResults = event.payload.results;
				winners = gameResults.winners || [];
				winReason = gameResults.winReason || '';
				allRoles = gameResults.finalState?.roles || {};
				eliminated = gameResults.finalState?.eliminated || [];

				// Show confetti if we won
				if (winners.includes($session?.playerId || '')) {
					showConfetti = true;
					setTimeout(() => showConfetti = false, 5000);
				}
			}
		});
	});

	function getRoleBadgeVariant(role: string): 'default' | 'destructive' | 'secondary' {
		if (role === 'werewolf' || role === 'minion') return 'destructive';
		if (role === 'villager') return 'secondary';
		return 'default';
	}

	function getRoleEmoji(role: string): string {
		const emojis: Record<string, string> = {
			werewolf: '🐺',
			seer: '🔮',
			robber: '🎭',
			troublemaker: '😈',
			mason: '🔨',
			villager: '👤',
			minion: '😤',
			tanner: '🤪',
			drunk: '🍺',
			insomniac: '😴'
		};
		return emojis[role] || '❓';
	}

	$: didIWin = winners.includes($session?.playerId || '');
	$: wasIEliminated = eliminated.includes($session?.playerId || '');
	$: isHost = $session?.playerId === roomState?.hostId;
	
	let isResetting = false;
	
	async function handlePlayAgain() {
		if (!roomState?.id) return;
		
		isResetting = true;
		try {
			await api.resetGame(roomState.id);
			// Room will be reset and clients will receive updated state via WebSocket
		} catch (error) {
			console.error('Failed to reset game:', error);
			alert('Failed to start a new game. Please try again.');
		} finally {
			isResetting = false;
		}
	}
</script>

<div class="space-y-6">
	{#if showConfetti}
		<div class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 pointer-events-none">
			<div use:confetti={{
				particleCount: 100,
				colors: ['#d79921', '#b8bb26', '#83a598', '#d3869b']
			}}></div>
		</div>
	{/if}

	<!-- Results header - only show if game has finished with winners -->
	{#if hasGameFinished}
		<Card class="p-6 {didIWin ? 'bg-gruvbox-green' : 'bg-gruvbox-red'} text-white border-0">
			<div class="text-center space-y-3">
				<div class="text-6xl">
					{didIWin ? '🎉' : '😔'}
				</div>
				<div>
					<h3 class="text-3xl font-bold mb-2">
						{didIWin ? 'You Won!' : 'You Lost'}
					</h3>
					<p class="text-lg text-white/90">
						{winReason}
					</p>
				</div>
			</div>
		</Card>
	{:else}
		<!-- Role reveal header (no winners yet) -->
		<Card class="p-6 bg-primary text-primary-foreground border-0">
			<div class="text-center space-y-3">
				<div class="text-6xl">🎭</div>
				<div>
					<h2 class="text-3xl font-bold mb-2">Role Reveal</h2>
					<p class="text-lg">
						Everyone's final roles are shown below. Determine the winner based on your physical votes!
					</p>
				</div>
			</div>
		</Card>
	{/if}

	<!-- Eliminated players -->
	{#if hasGameFinished && eliminated.length > 0}
		<Card class="p-6 bg-gruvbox-red/10 border-gruvbox-red">
			<div class="flex items-center gap-3 mb-4">
				<Skull class="w-6 h-6 text-gruvbox-red-light" />
				<h3 class="font-semibold text-lg">Eliminated</h3>
			</div>
			<div class="space-y-2">
				{#each eliminated as playerId}
					{@const player = gameResults?.finalState?.players?.find((p: any) => p.id === playerId)}
					{@const role = allRoles[playerId]}
					<div class="flex items-center justify-between p-3 bg-gruvbox-red/20 border border-gruvbox-red rounded-lg">
						<div class="flex items-center gap-3">
							<span class="text-2xl">{getRoleEmoji(role)}</span>
							<div>
								<p class="font-medium">{player?.displayName || 'Unknown'}</p>
								<p class="text-sm text-muted-foreground capitalize">{role?.replace('_', ' ')}</p>
							</div>
						</div>
						<Badge class="bg-gruvbox-red text-white">Eliminated</Badge>
					</div>
				{/each}
			</div>
		</Card>
	{/if}

	<!-- All roles revealed -->
	<Card class="p-6 bg-card border-primary">
		<div class="flex items-center gap-3 mb-4">
			<Trophy class="w-6 h-6 text-primary" />
			<h3 class="font-semibold text-lg">All Final Roles</h3>
		</div>
		<div class="grid grid-cols-1 md:grid-cols-2 gap-3">
			{#each Object.entries(allRoles) as [playerId, role]}
				{@const player = roomState?.players?.find((p: any) => p.id === playerId)}
				{@const isWinner = hasGameFinished && winners.includes(playerId)}
				<div class="flex items-center justify-between p-3 bg-muted/50 rounded-lg {isWinner ? 'ring-2 ring-gruvbox-green' : ''}">
					<div class="flex items-center gap-3">
						<span class="text-2xl">{getRoleEmoji(role)}</span>
						<div>
							<p class="font-medium">
								{player?.displayName || 'Unknown'}
								{#if playerId === $session?.playerId}
									<span class="text-xs text-muted-foreground">(You)</span>
								{/if}
							</p>
							<p class="text-sm text-muted-foreground capitalize">{role?.replace('_', ' ')}</p>
						</div>
					</div>
					<div class="flex flex-col items-end gap-1">
						<Badge variant={getRoleBadgeVariant(role)} class="capitalize">
							{role}
						</Badge>
						{#if isWinner}
							<Badge class="bg-gruvbox-green text-white border-gruvbox-green">
								Winner
							</Badge>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	</Card>

	<!-- Play again button -->
	<Card class="p-6">
		<div class="space-y-4">
			{#if !hasGameFinished}
				<p class="text-sm text-center text-muted-foreground">
					Discuss who won based on your physical votes and the roles revealed above!
				</p>
			{/if}
			
			{#if isHost}
				<Button
					class="w-full h-12 bg-primary hover:bg-primary/90"
					on:click={handlePlayAgain}
					disabled={isResetting}
				>
					{isResetting ? 'Resetting...' : '🎮 Play Again'}
				</Button>
				<p class="text-xs text-center text-muted-foreground">
					Start a new game with the same players
				</p>
			{:else}
				<p class="text-sm text-center text-muted-foreground">
					Waiting for the host to start a new game...
				</p>
			{/if}
		</div>
	</Card>
</div>
