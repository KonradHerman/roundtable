<script lang="ts">
	import { gameStore } from '$lib/stores/game.svelte';
	import { session } from '$lib/stores/session.svelte';
	import { Card, Badge, Button } from '$lib/components/ui';
	import { Trophy, Skull } from 'lucide-svelte';
	import { confetti } from '@neoconfetti/svelte';
	import { api } from '$lib/api/client';
	import { getRoleInfo } from './roleConfig';

	let { roomState } = $props<{ roomState: any }>();
	
	let gameResults = $state<any>(null);
	let allRoles = $state<Record<string, string>>({});
	let eliminated = $state<string[]>([]);
	let winners = $state<string[]>([]);
	let winReason = $state<string>('');
	let showConfetti = $state(false);
	let hasGameFinished = $state(false);
	let isResetting = $state(false);

	let didIWin = $derived(winners.includes(session.value?.playerId || ''));
	let wasIEliminated = $derived(eliminated.includes(session.value?.playerId || ''));
	let isHost = $derived(session.value?.playerId === roomState?.hostId);

	// Reactive effect to process game events
	$effect(() => {
		gameStore.events.forEach(event => {
			if (event.type === 'roles_revealed') {
				// Roles revealed for display (this is the main result screen)
				allRoles = event.payload.roles || event.payload || {};
				// If we don't have roles yet and game finished event came first, use those
				if (Object.keys(allRoles).length === 0 && gameResults?.finalState?.roles) {
					allRoles = gameResults.finalState.roles;
				}
			} else if (event.type === 'game_finished') {
				// Game finished with winners determined (optional, for future voting feature)
				hasGameFinished = true;
				gameResults = event.payload.results;
				winners = gameResults.winners || [];
				winReason = gameResults.winReason || '';
				// Update roles from game finished if available
				if (gameResults.finalState?.roles) {
					allRoles = gameResults.finalState.roles;
				}
				eliminated = gameResults.finalState?.eliminated || [];

				// Show confetti if we won
				if (winners.includes(session.value?.playerId || '')) {
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
		return getRoleInfo(role).emoji;
	}
	
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

	<!-- Role reveal header -->
	<Card class="p-6 bg-primary text-primary-foreground border-0">
		<div class="text-center space-y-3">
			<div class="text-6xl">🎭</div>
			<div>
				<h2 class="text-3xl font-bold mb-2">Final Roles Revealed</h2>
				<p class="text-lg">
					Here are everyone's final roles after all night actions
				</p>
			</div>
		</div>
	</Card>


	<!-- All roles revealed -->
	{#if Object.keys(allRoles).length > 0}
		<Card class="p-6 bg-card border-primary">
			<div class="flex items-center gap-3 mb-6">
				<Trophy class="w-6 h-6 text-primary" />
				<h3 class="font-semibold text-xl">Everyone's Final Roles</h3>
			</div>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				{#each Object.entries(allRoles) as [playerId, role]}
					{@const player = roomState?.players?.find((p: any) => p.id === playerId)}
					{@const isMe = playerId === session.value?.playerId}
					<div class="p-4 bg-muted/50 rounded-xl border-2 {isMe ? 'border-primary' : 'border-transparent'}">
						<div class="flex items-center gap-4">
							<div class="text-5xl">
								{getRoleEmoji(role)}
							</div>
							<div class="flex-1">
								<p class="font-bold text-lg">
									{player?.displayName || 'Unknown'}
									{#if isMe}
										<span class="text-sm text-primary">(You)</span>
									{/if}
								</p>
								<p class="text-2xl font-bold capitalize text-primary mt-1">
									{role?.replace('_', ' ')}
								</p>
							</div>
						</div>
					</div>
				{/each}
			</div>
		</Card>
	{:else}
		<Card class="p-6">
			<p class="text-center text-muted-foreground">Loading roles...</p>
		</Card>
	{/if}

	<!-- Play again button -->
	{#if isHost}
		<Card class="p-6 border-primary">
			<div class="space-y-4">
				<p class="text-center text-muted-foreground">
					Now vote on who to eliminate (physically) and determine the winners!
				</p>
				<Button
					class="w-full h-14 text-lg bg-primary hover:bg-primary/90"
					on:click={handlePlayAgain}
					disabled={isResetting}
				>
					{isResetting ? 'Setting up...' : '🎮 Play Again'}
				</Button>
				<p class="text-xs text-center text-muted-foreground">
					Start a new game with the same players
				</p>
			</div>
		</Card>
	{:else}
		<Card class="p-6">
			<p class="text-center text-muted-foreground">
				Waiting for the host to start a new game...
			</p>
		</Card>
	{/if}
</div>
