<script lang="ts">
	import { session } from '$lib/stores/session.svelte';
	import { Card, Button } from '$lib/components/ui';
	import { Users, Clock, Play, Pause, Plus } from 'lucide-svelte';

	let { roomState, wsStore, timerActive = false, phaseEndsAt = null } = $props<{
		roomState: any;
		wsStore: any;
		timerActive?: boolean;
		phaseEndsAt?: Date | null;
	}>();

	let timeRemaining = $state<number>(0);
	let timerInterval: any = null;

	let isHost = $derived(session.value?.playerId === roomState?.hostId);

	// No voting events to track - voting is physical!

	$effect(() => {
		timerInterval = setInterval(() => {
			if (timerActive && phaseEndsAt) {
				const now = new Date().getTime();
				const end = new Date(phaseEndsAt).getTime();
				timeRemaining = Math.max(0, Math.floor((end - now) / 1000));
			} else {
				timeRemaining = 0;
			}
		}, 100);

		return () => {
			if (timerInterval) clearInterval(timerInterval);
		};
	});

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

	function handleRevealRoles() {
		if (!wsStore) return;
		
		// Advance to results phase to show all roles
		wsStore.sendAction({
			type: 'advance_to_results',
			payload: {}
		});
	}
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
						onclick={handleToggleTimer}
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
							onclick={handleExtendTimer}
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
			<Users class="w-6 h-6 text-gruvbox-yellow-light flex-shrink-0 mt-1" />
			<div>
				<h3 class="font-semibold text-lg mb-1">Discussion Time</h3>
				<p class="text-sm text-muted-foreground">
					Discuss with your group who you think is a werewolf. Share what you learned during the night phase (or bluff!).
				</p>
			</div>
		</div>
	</Card>

	<!-- Discussion guide -->
	<Card class="p-6">
		<h3 class="font-semibold text-lg mb-4">Discussion Guide</h3>
		
		<div class="space-y-3 text-sm">
			<div class="flex items-start gap-2">
				<span class="text-lg">1️⃣</span>
				<p><strong>Share information:</strong> What role did you start as? What did you see or do?</p>
			</div>
			
			<div class="flex items-start gap-2">
				<span class="text-lg">2️⃣</span>
				<p><strong>Look for inconsistencies:</strong> Are people's stories adding up?</p>
			</div>
			
			<div class="flex items-start gap-2">
				<span class="text-lg">3️⃣</span>
				<p><strong>Deduce final roles:</strong> Based on night actions, who has what role now?</p>
			</div>
			
			<div class="flex items-start gap-2">
				<span class="text-lg">4️⃣</span>
				<p><strong>Vote when ready:</strong> Everyone point at who to eliminate simultaneously!</p>
			</div>
		</div>

		{#if isHost}
			<div class="mt-6 pt-6 border-t border-border">
				<p class="text-sm text-muted-foreground mb-3">
					When discussion is done and everyone has voted physically:
				</p>
				<Button
					onclick={handleRevealRoles}
					class="w-full h-12 bg-primary hover:bg-primary/90"
				>
					Reveal All Roles →
				</Button>
			</div>
		{:else}
			<div class="mt-6 pt-6 border-t border-border">
				<p class="text-sm text-center text-muted-foreground">
					Waiting for the host to reveal roles...
				</p>
			</div>
		{/if}
	</Card>
</div>
