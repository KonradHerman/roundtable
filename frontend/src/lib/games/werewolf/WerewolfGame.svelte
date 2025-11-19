<script lang="ts">
	import { gameStore } from '$lib/stores/game.svelte';
	import { session } from '$lib/stores/session.svelte';
	import { Card, Badge } from '$lib/components/ui';
	import { Moon, Sun, Eye, Clock } from 'lucide-svelte';
	import RoleReveal from './RoleReveal.svelte';
	import NightPhase from './NightPhase.svelte';
	import DayPhase from './DayPhase.svelte';
	import Results from './Results.svelte';

	let { roomCode, roomState, wsStore } = $props<{
		roomCode: string;
		roomState: any;
		wsStore: any;
	}>();

	let myRole = $state<string | null>(null);
	let currentPhase = $state<string>('setup');
	let phaseEndsAt = $state<Date | null>(null);
	let acknowledged = $state<boolean>(false);
	let acknowledgementsCount = $state<number>(0);
	let totalPlayers = $state<number>(0);
	let nightScript = $state<any[]>([]);
	let timerActive = $state<boolean>(false);

	// Reactive effect to process game events
	$effect(() => {
		gameStore.events.forEach(event => {
			if (event.type === 'role_assigned') {
				myRole = event.payload.role;
			} else if (event.type === 'phase_changed') {
				currentPhase = event.payload.phase.name;
				if (event.payload.phase.endsAt) {
					phaseEndsAt = new Date(event.payload.phase.endsAt);
				}
			} else if (event.type === 'role_acknowledged') {
				acknowledgementsCount = event.payload.count;
				totalPlayers = event.payload.total;
				if (event.payload.playerId === session.value?.playerId) {
					acknowledged = true;
				}
			} else if (event.type === 'night_script') {
				nightScript = event.payload.script || [];
			} else if (event.type === 'timer_toggled') {
				timerActive = event.payload.active;
				if (event.payload.phaseEndsAt) {
					phaseEndsAt = new Date(event.payload.phaseEndsAt);
				}
			} else if (event.type === 'timer_extended') {
				phaseEndsAt = new Date(event.payload.phaseEndsAt);
			}
		});
	});

	function handleAcknowledgeRole() {
		if (!wsStore || acknowledged) return;

		wsStore.sendAction({
			type: 'acknowledge_role',
			payload: {}
		});
	}

	let phaseIcon = $derived(currentPhase === 'night' ? Moon : currentPhase === 'day' ? Sun : currentPhase === 'role_reveal' ? Eye : Clock);
	let phaseColor = $derived(currentPhase === 'night' ? 'bg-gruvbox-purple' : currentPhase === 'day' ? 'bg-gruvbox-yellow' : currentPhase === 'role_reveal' ? 'bg-gruvbox-blue' : 'bg-muted');
</script>

<div class="space-y-6">
	<!-- Phase header -->
	<Card class="p-6 {phaseColor} text-white border-0">
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-3">
				<phaseIcon class="w-8 h-8"></phaseIcon>
				<div>
					<h2 class="text-2xl font-bold capitalize">
						{currentPhase.replace('_', ' ')} Phase
					</h2>
					{#if currentPhase === 'role_reveal'}
						<p class="text-white/90 text-sm">Look at your role card</p>
					{:else if currentPhase === 'night'}
						<p class="text-white/90 text-sm">Everyone close your eyes</p>
					{:else if currentPhase === 'day'}
						<p class="text-white/90 text-sm">Discuss and vote!</p>
					{/if}
				</div>
			</div>
		</div>
	</Card>

	<!-- Phase-specific content -->
	{#if currentPhase === 'role_reveal' && myRole}
		<RoleReveal 
			role={myRole}
			{acknowledged}
			{acknowledgementsCount}
			{totalPlayers}
			onAcknowledge={handleAcknowledgeRole}
		/>
	{:else if currentPhase === 'night'}
		<NightPhase
			{roomState}
			{wsStore}
			{nightScript}
		/>
	{:else if currentPhase === 'day'}
		<DayPhase
			{roomState}
			{wsStore}
			{timerActive}
			{phaseEndsAt}
		/>
	{:else if currentPhase === 'results'}
		<Results {roomState} />
	{:else}
		<Card class="p-6">
			<p class="text-center text-muted-foreground">Preparing game...</p>
		</Card>
	{/if}
</div>
