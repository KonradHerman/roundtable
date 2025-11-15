<script lang="ts">
	import { Card, Button } from '$lib/components/ui';
	import { Eye, Users, Moon } from 'lucide-svelte';

	export let myRole: string | null;
	export let otherWerewolves: string[];
	export let otherMasons: string[];
	export let seerResult: any;
	export let roomState: any;
	export let wsStore: any;
	export let getPlayerName: (id: string) => string;

	let selectedPlayer: string | null = null;
	let actionTaken = false;

	function handleSeerView() {
		if (!selectedPlayer || !wsStore) return;

		wsStore.sendAction({
			type: 'seer_view',
			payload: {
				targetId: selectedPlayer
			}
		});

		actionTaken = true;
	}

	$: availablePlayers = roomState?.players?.filter((p: any) => p.id !== roomState?.players?.find((player: any) => player.displayName)?.id) || [];
</script>

<div class="space-y-6">
	<!-- Role-specific instructions -->
	{#if myRole === 'werewolf'}
		<Card class="p-6 bg-red-950/50 border-red-900">
			<div class="flex items-start gap-3">
				<Users class="w-6 h-6 text-red-400 flex-shrink-0 mt-1" />
				<div class="flex-1">
					<h3 class="font-semibold text-lg text-red-100 mb-2">Werewolf</h3>
					{#if otherWerewolves.length > 0}
						<p class="text-red-200 mb-3">Your fellow werewolves:</p>
						<div class="space-y-2">
							{#each otherWerewolves as werewolfId}
								<div class="bg-red-900/30 p-3 rounded-lg">
									<span class="font-medium text-red-100">
										{getPlayerName(werewolfId)}
									</span>
								</div>
							{/each}
						</div>
					{:else}
						<p class="text-red-200">You are the lone werewolf!</p>
					{/if}
				</div>
			</div>
		</Card>

	{:else if myRole === 'seer'}
		<Card class="p-6">
			<div class="space-y-4">
				<div class="flex items-center gap-3">
					<Eye class="w-6 h-6 text-purple-500" />
					<div>
						<h3 class="font-semibold text-lg">Seer</h3>
						<p class="text-sm text-muted-foreground">Choose one player to view their role</p>
					</div>
				</div>

				{#if seerResult}
					<div class="p-4 bg-purple-500/10 border border-purple-500/20 rounded-lg">
						<p class="text-sm text-muted-foreground mb-1">You viewed:</p>
						<p class="font-semibold text-lg">{getPlayerName(seerResult.targetId)}</p>
						<p class="text-purple-600 dark:text-purple-400 capitalize">
							Role: {seerResult.role}
						</p>
					</div>
				{:else}
					<div class="space-y-3">
						{#each roomState?.players || [] as player}
							<button
								on:click={() => selectedPlayer = player.id}
								class="w-full p-4 text-left rounded-lg border-2 transition-all {selectedPlayer === player.id ? 'border-purple-500 bg-purple-500/10' : 'border-border hover:border-purple-300'}"
								disabled={actionTaken}
							>
								<div class="flex items-center gap-3">
									<div class="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold">
										{player.displayName[0].toUpperCase()}
									</div>
									<span class="font-medium">{player.displayName}</span>
								</div>
							</button>
						{/each}

						<Button
							class="w-full h-12"
							disabled={!selectedPlayer || actionTaken}
							on:click={handleSeerView}
						>
							{actionTaken ? 'Action Taken' : 'View Role'}
						</Button>
					</div>
				{/if}
			</div>
		</Card>

	{:else if myRole === 'mason'}
		<Card class="p-6 bg-gray-950/50 border-gray-800">
			<div class="flex items-start gap-3">
				<Users class="w-6 h-6 text-gray-400 flex-shrink-0 mt-1" />
				<div class="flex-1">
					<h3 class="font-semibold text-lg text-gray-100 mb-2">Mason</h3>
					{#if otherMasons.length > 0}
						<p class="text-gray-200 mb-3">Your fellow mason:</p>
						<div class="space-y-2">
							{#each otherMasons as masonId}
								<div class="bg-gray-900/30 p-3 rounded-lg">
									<span class="font-medium text-gray-100">
										{getPlayerName(masonId)}
									</span>
								</div>
							{/each}
						</div>
					{:else}
						<p class="text-gray-200">You are the lone mason!</p>
					{/if}
				</div>
			</div>
		</Card>

	{:else if myRole === 'villager'}
		<Card class="p-6">
			<div class="flex items-center gap-3">
				<Moon class="w-6 h-6 text-muted-foreground" />
				<div>
					<h3 class="font-semibold text-lg">Villager</h3>
					<p class="text-muted-foreground">
						You have no special abilities. Wait for the day phase to discuss and vote.
					</p>
				</div>
			</div>
		</Card>

	{:else if myRole === 'robber'}
		<Card class="p-6">
			<div class="flex items-center gap-3">
				<span class="text-3xl">🎭</span>
				<div>
					<h3 class="font-semibold text-lg">Robber</h3>
					<p class="text-muted-foreground">
						Choose a player to swap roles with (simplified for MVP - feature coming soon)
					</p>
				</div>
			</div>
		</Card>

	{:else if myRole}
		<Card class="p-6">
			<div class="text-center">
				<h3 class="font-semibold text-lg capitalize mb-2">{myRole.replace('_', ' ')}</h3>
				<p class="text-muted-foreground">
					Role-specific actions coming soon!
				</p>
			</div>
		</Card>
	{/if}

	<!-- General info card -->
	<Card class="p-4 bg-muted/50">
		<p class="text-sm text-muted-foreground text-center">
			<strong>Remember:</strong> Stay quiet during the night phase! Only use your phone.
		</p>
	</Card>
</div>
