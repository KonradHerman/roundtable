<script lang="ts">
	import { Card, Button, Badge } from '$lib/components/ui';
	import { Moon, ChevronDown, ChevronUp, Eye, Users, ArrowRightLeft } from 'lucide-svelte';
	import { session } from '$lib/stores/session';
	import { gameStore } from '$lib/stores/game';
	import { onMount } from 'svelte';

	export let roomState: any;
	export let wsStore: any;
	export let nightScript: any[] = [];

	let scriptExpanded = true;
	let checkedSteps: Record<number, boolean> = {};
	
	// Role-specific state
	let myRole: string | null = null;
	let otherWerewolves: string[] = [];
	let otherMasons: string[] = [];
	let selectedPlayer: string | null = null;
	let selectedCenterCard: number | null = null;
	let selectedCenterCards: number[] = [];
	let selectedPlayer1: string | null = null;
	let selectedPlayer2: string | null = null;
	let actionResult: any = null;
	let hasActed = false;

	// Subscribe to game events to get role-specific info
	let unsubscribe = gameStore.subscribe(($game) => {
		$game.events.forEach((event: any) => {
			if (event.type === 'role_assigned') {
				myRole = event.payload.role;
			} else if (event.type === 'werewolf_wakeup') {
				otherWerewolves = event.payload.otherWerewolves || [];
			} else if (event.type === 'mason_wakeup') {
				otherMasons = event.payload.otherMasons || [];
			} else if (event.type === 'werewolf_view_center_result') {
				actionResult = event.payload;
				hasActed = true;
			} else if (event.type === 'seer_result') {
				actionResult = event.payload;
				hasActed = true;
			} else if (event.type === 'seer_center_result') {
				actionResult = event.payload;
				hasActed = true;
			} else if (event.type === 'robber_result') {
				actionResult = event.payload;
				hasActed = true;
			} else if (event.type === 'troublemaker_confirmed') {
				actionResult = event.payload;
				hasActed = true;
			} else if (event.type === 'drunk_confirmed') {
				actionResult = event.payload;
				hasActed = true;
			}
		});
	});

	onMount(() => {
		return () => {
			if (unsubscribe) unsubscribe();
		};
	});

	function handleAdvanceToDay() {
		if (!wsStore) return;
		wsStore.sendAction({
			type: 'advance_phase',
			payload: {}
		});
	}

	// Role-specific action handlers
	function handleWerewolfViewCenter(centerIndex: number) {
		wsStore.sendAction({
			type: 'werewolf_view_center',
			payload: { centerIndex }
		});
	}

	function handleSeerViewPlayer(playerId: string) {
		wsStore.sendAction({
			type: 'seer_view_player',
			payload: { targetId: playerId }
		});
	}

	function handleSeerViewCenter(indices: number[]) {
		wsStore.sendAction({
			type: 'seer_view_center',
			payload: { centerIndices: indices }
		});
	}

	function handleRobberSwap(playerId: string) {
		wsStore.sendAction({
			type: 'robber_swap',
			payload: { targetId: playerId }
		});
	}

	function handleTroublemakerSwap(player1Id: string, player2Id: string) {
		wsStore.sendAction({
			type: 'troublemaker_swap',
			payload: { player1Id, player2Id }
		});
	}

	function handleDrunkSwap(centerIndex: number) {
		wsStore.sendAction({
			type: 'drunk_swap',
			payload: { centerIndex }
		});
	}

	function getPlayerName(playerId: string): string {
		const player = roomState?.players?.find((p: any) => p.id === playerId);
		return player?.displayName || 'Unknown';
	}

	$: isHost = $session?.playerId === roomState?.hostId;
	$: otherPlayers = roomState?.players?.filter((p: any) => p.id !== $session?.playerId) || [];
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

	<!-- Host: Narration script -->
	{#if isHost}
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
	{/if}

	<!-- Player role-specific UIs -->
	{#if !isHost || myRole}
		{#if myRole === 'werewolf'}
			<Card class="p-6 bg-red-500/10 border-red-500/30">
				<div class="flex items-center gap-3 mb-4">
					<Users class="w-6 h-6 text-red-400" />
					<div>
						<h3 class="font-semibold text-lg">Werewolf</h3>
						<p class="text-sm text-muted-foreground">Find your pack</p>
					</div>
				</div>

				{#if otherWerewolves.length > 0}
					<div class="space-y-2">
						<p class="text-sm font-medium">Other werewolves:</p>
						{#each otherWerewolves as wwId}
							<div class="p-3 bg-red-500/5 rounded-lg border border-red-500/20">
								<p class="font-medium">{getPlayerName(wwId)}</p>
							</div>
						{/each}
					</div>
				{:else}
					<div class="space-y-4">
						<p class="text-sm">You are the only werewolf! You may view one center card.</p>
						{#if !hasActed}
							<div class="grid grid-cols-3 gap-2">
								{#each [0, 1, 2] as index}
									<button
										on:click={() => handleWerewolfViewCenter(index)}
										class="p-4 bg-muted hover:bg-primary/20 rounded-lg border-2 border-border hover:border-primary transition-all"
									>
										<p class="text-sm font-medium">Card {index + 1}</p>
									</button>
								{/each}
							</div>
						{:else if actionResult}
							<div class="p-4 bg-green-500/10 border border-green-500/30 rounded-lg">
								<p class="text-sm font-medium">Center card {actionResult.centerIndex + 1}: <span class="capitalize">{actionResult.role}</span></p>
							</div>
						{/if}
					</div>
				{/if}
			</Card>

		{:else if myRole === 'seer'}
			<Card class="p-6 bg-blue-500/10 border-blue-500/30">
				<div class="flex items-center gap-3 mb-4">
					<Eye class="w-6 h-6 text-blue-400" />
					<div>
						<h3 class="font-semibold text-lg">Seer</h3>
						<p class="text-sm text-muted-foreground">View one player or two center cards</p>
					</div>
				</div>

				{#if !hasActed}
					<div class="space-y-4">
						<div>
							<p class="text-sm font-medium mb-2">View a player's card:</p>
							<div class="space-y-2">
								{#each otherPlayers as player}
									<button
										on:click={() => handleSeerViewPlayer(player.id)}
										class="w-full p-3 bg-muted hover:bg-primary/20 rounded-lg border-2 border-border hover:border-primary transition-all text-left"
									>
										{player.displayName}
									</button>
								{/each}
							</div>
						</div>

						<div class="relative">
							<div class="absolute inset-0 flex items-center">
								<div class="w-full border-t border-border"></div>
							</div>
							<div class="relative flex justify-center text-sm">
								<span class="px-2 bg-card text-muted-foreground">OR</span>
							</div>
						</div>

						<div>
							<p class="text-sm font-medium mb-2">View two center cards:</p>
							<div class="grid grid-cols-3 gap-2">
								{#each [0, 1, 2] as index}
									<button
										on:click={() => {
											if (selectedCenterCards.includes(index)) {
												selectedCenterCards = selectedCenterCards.filter(i => i !== index);
											} else if (selectedCenterCards.length < 2) {
												selectedCenterCards = [...selectedCenterCards, index];
											}
										}}
										class="p-4 rounded-lg border-2 transition-all {selectedCenterCards.includes(index) ? 'bg-primary/20 border-primary' : 'bg-muted border-border hover:border-primary/50'}"
									>
										<p class="text-sm font-medium">Card {index + 1}</p>
									</button>
								{/each}
							</div>
							{#if selectedCenterCards.length === 2}
								<Button
									on:click={() => handleSeerViewCenter(selectedCenterCards)}
									class="w-full mt-3"
								>
									View Selected Cards
								</Button>
							{/if}
						</div>
					</div>
				{:else if actionResult}
					<div class="p-4 bg-green-500/10 border border-green-500/30 rounded-lg">
						{#if actionResult.targetId}
							<p class="font-medium">{getPlayerName(actionResult.targetId)} is: <span class="capitalize text-lg">{actionResult.role}</span></p>
						{:else if actionResult.cards}
							<p class="font-medium mb-2">Center cards:</p>
							{#each actionResult.cards as card}
								<p>Card {card.index + 1}: <span class="capitalize">{card.role}</span></p>
							{/each}
						{/if}
					</div>
				{/if}
			</Card>

		{:else if myRole === 'robber'}
			<Card class="p-6 bg-orange-500/10 border-orange-500/30">
				<div class="flex items-center gap-3 mb-4">
					<ArrowRightLeft class="w-6 h-6 text-orange-400" />
					<div>
						<h3 class="font-semibold text-lg">Robber</h3>
						<p class="text-sm text-muted-foreground">Swap with another player and see your new role</p>
					</div>
				</div>

				{#if !hasActed}
					<div class="space-y-2">
						{#each otherPlayers as player}
							<button
								on:click={() => handleRobberSwap(player.id)}
								class="w-full p-3 bg-muted hover:bg-primary/20 rounded-lg border-2 border-border hover:border-primary transition-all text-left"
							>
								Swap with {player.displayName}
							</button>
						{/each}
					</div>
				{:else if actionResult}
					<div class="p-4 bg-green-500/10 border border-green-500/30 rounded-lg">
						<p class="font-medium">You swapped with {getPlayerName(actionResult.targetId)}</p>
						<p class="text-lg mt-2">Your new role: <span class="capitalize font-bold">{actionResult.newRole}</span></p>
					</div>
				{/if}
			</Card>

		{:else if myRole === 'troublemaker'}
			<Card class="p-6 bg-purple-500/10 border-purple-500/30">
				<div class="flex items-center gap-3 mb-4">
					<ArrowRightLeft class="w-6 h-6 text-purple-400" />
					<div>
						<h3 class="font-semibold text-lg">Troublemaker</h3>
						<p class="text-sm text-muted-foreground">Swap two other players (you don't see what)</p>
					</div>
				</div>

				{#if !hasActed}
					<div class="space-y-4">
						<div>
							<p class="text-sm font-medium mb-2">Select first player:</p>
							<div class="space-y-2">
								{#each otherPlayers as player}
									<button
										on:click={() => selectedPlayer1 = player.id}
										class="w-full p-3 rounded-lg border-2 transition-all text-left {selectedPlayer1 === player.id ? 'bg-primary/20 border-primary' : 'bg-muted border-border hover:border-primary/50'}"
									>
										{player.displayName}
									</button>
								{/each}
							</div>
						</div>

						{#if selectedPlayer1}
							<div>
								<p class="text-sm font-medium mb-2">Select second player:</p>
								<div class="space-y-2">
									{#each otherPlayers.filter(p => p.id !== selectedPlayer1) as player}
										<button
											on:click={() => selectedPlayer2 = player.id}
											class="w-full p-3 rounded-lg border-2 transition-all text-left {selectedPlayer2 === player.id ? 'bg-primary/20 border-primary' : 'bg-muted border-border hover:border-primary/50'}"
										>
											{player.displayName}
										</button>
									{/each}
								</div>
							</div>
						{/if}

						{#if selectedPlayer1 && selectedPlayer2}
							<Button
								on:click={() => handleTroublemakerSwap(selectedPlayer1, selectedPlayer2)}
								class="w-full"
							>
								Swap These Players
							</Button>
						{/if}
					</div>
				{:else if actionResult}
					<div class="p-4 bg-green-500/10 border border-green-500/30 rounded-lg">
						<p class="font-medium">✓ You swapped {getPlayerName(actionResult.player1Id)} and {getPlayerName(actionResult.player2Id)}</p>
						<p class="text-sm text-muted-foreground mt-2">You don't know what roles they had</p>
					</div>
				{/if}
			</Card>

		{:else if myRole === 'drunk'}
			<Card class="p-6 bg-amber-500/10 border-amber-500/30">
				<div class="flex items-center gap-3 mb-4">
					<span class="text-2xl">🍺</span>
					<div>
						<h3 class="font-semibold text-lg">Drunk</h3>
						<p class="text-sm text-muted-foreground">Swap with a center card (you don't see your new role)</p>
					</div>
				</div>

				{#if !hasActed}
					<div class="grid grid-cols-3 gap-2">
						{#each [0, 1, 2] as index}
							<button
								on:click={() => handleDrunkSwap(index)}
								class="p-4 bg-muted hover:bg-primary/20 rounded-lg border-2 border-border hover:border-primary transition-all"
							>
								<p class="text-sm font-medium">Card {index + 1}</p>
							</button>
						{/each}
					</div>
				{:else}
					<div class="p-4 bg-green-500/10 border border-green-500/30 rounded-lg">
						<p class="font-medium">✓ You swapped with center card {actionResult.centerIndex + 1}</p>
						<p class="text-sm text-muted-foreground mt-2">You don't know your new role!</p>
					</div>
				{/if}
			</Card>

		{:else if myRole === 'mason'}
			<Card class="p-6 bg-gray-500/10 border-gray-500/30">
				<div class="flex items-center gap-3 mb-4">
					<Users class="w-6 h-6 text-gray-400" />
					<div>
						<h3 class="font-semibold text-lg">Mason</h3>
						<p class="text-sm text-muted-foreground">Find the other mason(s)</p>
					</div>
				</div>

				{#if otherMasons.length > 0}
					<div class="space-y-2">
						<p class="text-sm font-medium">Other mason(s):</p>
						{#each otherMasons as masonId}
							<div class="p-3 bg-gray-500/5 rounded-lg border border-gray-500/20">
								<p class="font-medium">{getPlayerName(masonId)}</p>
							</div>
						{/each}
					</div>
				{:else}
					<p class="text-muted-foreground">You are the only mason in the game.</p>
				{/if}
			</Card>

		{:else if myRole === 'villager' || myRole === 'tanner' || myRole === 'hunter' || myRole === 'insomniac'}
			<Card class="p-6">
				<div class="text-center space-y-4">
					<div class="text-6xl">😴</div>
					<h3 class="text-xl font-bold capitalize">{myRole}</h3>
					<p class="text-muted-foreground">
						{#if myRole === 'insomniac'}
							You'll see your role at the end of the night phase
						{:else if myRole === 'tanner'}
							You have no night action. You want to get eliminated!
						{:else if myRole === 'hunter'}
							You have no night action. If you die, the player you voted for also dies.
						{:else}
							You have no night action. Keep your eyes closed.
						{/if}
					</p>
				</div>
			</Card>
		{/if}
	{/if}
</div>
