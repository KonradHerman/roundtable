<script lang="ts">
	import { onMount } from 'svelte';
	import { gameStore } from '$lib/stores/game';
	import { session } from '$lib/stores/session';
	import { Card, Button, Badge } from '$lib/components/ui';
	import { Moon, Sun, Clock, Eye } from 'lucide-svelte';
	import RoleReveal from './RoleReveal.svelte';
	import NightPhase from './NightPhase.svelte';
	import DayPhase from './DayPhase.svelte';
	import Results from './Results.svelte';

	// eslint-disable-next-line no-unused-vars
	export let roomCode: string;
	export let roomState: any;
	export let wsStore: any;

	let myRole: string | null = null;
	let currentPhase: string = 'setup';
	let phaseEndsAt: Date | null = null;
	let timeRemaining: number = 0;
	let otherWerewolves: string[] = [];
	let otherMasons: string[] = [];
	let seerResult: any = null;
	let showRoleReveal = true;

	// Subscribe to game events
	let unsubscribe = gameStore.subscribe(($game) => {
		// Process events to update local state
		$game.events.forEach(event => {
			if (event.type === 'role_assigned') {
				myRole = event.payload.role;
			} else if (event.type === 'werewolf_wakeup') {
				otherWerewolves = event.payload.otherWerewolves || [];
			} else if (event.type === 'mason_wakeup') {
				otherMasons = event.payload.otherMasons || [];
			} else if (event.type === 'phase_changed') {
				currentPhase = event.payload.phase.name;
				if (event.payload.phase.endsAt) {
					phaseEndsAt = new Date(event.payload.phase.endsAt);
				}
				// Hide role reveal after a few seconds
				if (currentPhase === 'night') {
					setTimeout(() => showRoleReveal = false, 5000);
				}
			} else if (event.type === 'seer_result') {
				seerResult = event.payload;
			}
		});
	});

	onMount(() => {
		// Timer countdown
		const interval = setInterval(() => {
			if (phaseEndsAt) {
				const now = new Date().getTime();
				const end = phaseEndsAt.getTime();
				timeRemaining = Math.max(0, Math.floor((end - now) / 1000));
			}
		}, 1000);

		return () => {
			clearInterval(interval);
			if (unsubscribe) unsubscribe();
		};
	});

	function formatTime(seconds: number): string {
		const mins = Math.floor(seconds / 60);
		const secs = seconds % 60;
		return `${mins}:${secs.toString().padStart(2, '0')}`;
	}

	function getPlayerName(playerId: string): string {
		const player = roomState?.players?.find((p: any) => p.id === playerId);
		return player?.displayName || 'Unknown';
	}

	$: phaseIcon = currentPhase === 'night' ? Moon : Sun;
	$: phaseColor = currentPhase === 'night' ? 'bg-indigo-600' : 'bg-amber-500';
</script>

<div class="space-y-6">
	<!-- Phase header -->
	<Card class="p-6 {phaseColor} text-white">
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-3">
				<svelte:component this={phaseIcon} class="w-8 h-8" />
				<div>
					<h2 class="text-2xl font-bold capitalize">{currentPhase} Phase</h2>
					{#if currentPhase === 'night'}
						<p class="text-white/90 text-sm">Stay quiet and use your phone...</p>
					{:else if currentPhase === 'day'}
						<p class="text-white/90 text-sm">Discuss and vote!</p>
					{/if}
				</div>
			</div>
			{#if phaseEndsAt && timeRemaining > 0}
				<div class="flex items-center gap-2 bg-black/20 px-4 py-2 rounded-lg">
					<Clock class="w-5 h-5" />
					<span class="text-2xl font-mono font-bold">{formatTime(timeRemaining)}</span>
				</div>
			{/if}
		</div>
	</Card>

	<!-- Role reveal (shows briefly when game starts) -->
	{#if showRoleReveal && myRole}
		<RoleReveal role={myRole} />
	{/if}

	<!-- Phase-specific content -->
	{#if currentPhase === 'night'}
		<NightPhase
			{myRole}
			{otherWerewolves}
			{otherMasons}
			{seerResult}
			{roomState}
			{wsStore}
			{getPlayerName}
		/>
	{:else if currentPhase === 'day'}
		<DayPhase
			{roomState}
			{wsStore}
		/>
	{:else if currentPhase === 'results'}
		<Results />
	{:else}
		<Card class="p-6">
			<p class="text-center text-muted-foreground">Preparing game...</p>
		</Card>
	{/if}
</div>
