<script lang="ts">
	import { Card, Button } from '$lib/components/ui';
	import { Moon, ChevronDown, ChevronUp } from 'lucide-svelte';
	import { session } from '$lib/stores/session';

	export let roomState: any;
	export let wsStore: any;
	export let nightScript: any[] = [];

	let scriptExpanded = true;
	let checkedSteps: Record<number, boolean> = {};

	function handleAdvanceToDay() {
		if (!wsStore) return;

		wsStore.sendAction({
			type: 'advance_phase',
			payload: {}
		});
	}

	$: isHost = $session?.playerId === roomState?.hostId;
</script>

<div class="space-y-6">
	<!-- Main instruction card -->
	<Card class="p-6 bg-gruvbox-purple/20 border-gruvbox-purple">
		<div class="flex items-center gap-3">
			<Moon class="w-8 h-8 text-gruvbox-purple-light" />
			<div>
				<h3 class="font-semibold text-xl text-foreground">Night Phase</h3>
				<p class="text-muted-foreground">Everyone close your eyes 😴</p>
			</div>
		</div>
	</Card>

	{#if isHost}
		<!-- Host: Narration script -->
		<Card class="p-6 border-primary">
			<button
				on:click={() => scriptExpanded = !scriptExpanded}
				class="w-full flex items-center justify-between mb-4"
			>
				<div class="flex items-center gap-2">
					<span class="text-2xl">📜</span>
					<h3 class="font-semibold text-lg">Narration Script (Host Only)</h3>
				</div>
				{#if scriptExpanded}
					<ChevronUp class="w-5 h-5" />
				{:else}
					<ChevronDown class="w-5 h-5" />
				{/if}
			</button>

			{#if scriptExpanded}
				<div class="space-y-3">
					<p class="text-sm text-muted-foreground mb-4">
						Read these instructions aloud in order. Check off each role as they complete their action.
					</p>

					{#if nightScript && nightScript.length > 0}
						{#each nightScript as step}
							<label class="flex items-start gap-3 p-4 bg-muted/30 rounded-lg hover:bg-muted/50 cursor-pointer transition-colors">
								<input
									type="checkbox"
									bind:checked={checkedSteps[step.order]}
									class="mt-1 w-5 h-5 rounded border-border text-primary focus:ring-primary"
								/>
								<div class="flex-1">
									<div class="font-semibold capitalize text-foreground mb-1">
										{step.order}. {step.role}
									</div>
									<p class="text-sm text-muted-foreground">
										{step.instruction}
									</p>
								</div>
							</label>
						{/each}
					{:else}
						<p class="text-muted-foreground text-center py-4">
							Script will appear here when available
						</p>
					{/if}

					<div class="pt-4 border-t border-border">
						<Button
							on:click={handleAdvanceToDay}
							class="w-full h-12 bg-primary hover:bg-primary/90 text-primary-foreground font-bold"
						>
							Advance to Day Phase →
						</Button>
						<p class="text-xs text-muted-foreground text-center mt-2">
							Click when all night actions are complete
						</p>
					</div>
				</div>
			{/if}
		</Card>
	{:else}
		<!-- Non-host: Simple wait message -->
		<Card class="p-12">
			<div class="text-center space-y-4">
				<div class="text-7xl">😴</div>
				<h3 class="text-2xl font-bold text-foreground">Keep your eyes closed</h3>
				<p class="text-lg text-muted-foreground">
					The host will narrate the night phase. Listen for your role to be called.
				</p>
				<div class="pt-6">
					<div class="inline-block px-6 py-3 bg-muted rounded-lg">
						<p class="text-sm text-muted-foreground">
							Waiting for night phase to complete...
						</p>
					</div>
				</div>
			</div>
		</Card>
	{/if}
</div>
