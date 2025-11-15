<script lang="ts">
	import { gameStore } from '$lib/stores/game';
	import { session } from '$lib/stores/session';
	import { Card, Badge, Button } from '$lib/components/ui';
	import { Trophy, Skull } from 'lucide-svelte';
	import { confetti } from '@neoconfetti/svelte';

	let gameResults: any = null;
	let allRoles: Record<string, string> = {};
	let eliminated: string[] = [];
	let winners: string[] = [];
	let winReason: string = '';
	let showConfetti = false;

	// Subscribe to game finished event
	let unsubscribe = gameStore.subscribe(($game) => {
		$game.events.forEach(event => {
			if (event.type === 'game_finished') {
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
</script>

<div class="space-y-6">
	{#if showConfetti}
		<div class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 pointer-events-none">
			<div use:confetti={{
				particleCount: 100,
				spread: 70,
				colors: ['#0ea5e9', '#f59e0b', '#10b981']
			}}></div>
		</div>
	{/if}

	<!-- Results header -->
	<Card class="p-6 {didIWin ? 'bg-green-500' : 'bg-red-500'} text-white border-0">
		<div class="text-center space-y-3">
			<div class="text-6xl">
				{didIWin ? '🎉' : '😔'}
			</div>
			<div>
				<h2 class="text-3xl font-bold mb-2">
					{didIWin ? 'You Won!' : 'You Lost'}
				</h2>
				<p class="text-lg text-white/90">
					{winReason}
				</p>
			</div>
		</div>
	</Card>

	<!-- Eliminated players -->
	{#if eliminated.length > 0}
		<Card class="p-6">
			<div class="flex items-center gap-3 mb-4">
				<Skull class="w-6 h-6 text-destructive" />
				<h3 class="font-semibold text-lg">Eliminated</h3>
			</div>
			<div class="space-y-2">
				{#each eliminated as playerId}
					{@const player = gameResults?.finalState?.players?.find((p: any) => p.id === playerId)}
					{@const role = allRoles[playerId]}
					<div class="flex items-center justify-between p-3 bg-destructive/10 border border-destructive/20 rounded-lg">
						<div class="flex items-center gap-3">
							<span class="text-2xl">{getRoleEmoji(role)}</span>
							<div>
								<p class="font-medium">{player?.displayName || 'Unknown'}</p>
								<p class="text-sm text-muted-foreground capitalize">{role?.replace('_', ' ')}</p>
							</div>
						</div>
						<Badge variant="destructive">Eliminated</Badge>
					</div>
				{/each}
			</div>
		</Card>
	{/if}

	<!-- All roles revealed -->
	<Card class="p-6">
		<div class="flex items-center gap-3 mb-4">
			<Trophy class="w-6 h-6 text-primary" />
			<h3 class="font-semibold text-lg">All Roles</h3>
		</div>
		<div class="grid grid-cols-1 md:grid-cols-2 gap-3">
			{#each Object.entries(allRoles) as [playerId, role]}
				{@const player = gameResults?.finalState?.players?.find((p: any) => p.id === playerId)}
				{@const isWinner = winners.includes(playerId)}
				<div class="flex items-center justify-between p-3 bg-muted/50 rounded-lg {isWinner ? 'ring-2 ring-green-500' : ''}">
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
						<Badge variant={getRoleBadgeVariant(role)}>
							{role}
						</Badge>
						{#if isWinner}
							<Badge variant="outline" class="bg-green-500/10 text-green-700 border-green-500/20">
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
		<Button
			class="w-full h-12"
			on:click={() => window.location.reload()}
		>
			Back to Lobby
		</Button>
		<p class="text-sm text-center text-muted-foreground mt-3">
			The host can start a new game from the lobby
		</p>
	</Card>
</div>
