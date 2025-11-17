<script lang="ts">
	import { CardBack } from '$lib/components/ui';
	
	export let players: any[] = [];
	export let selectedPlayerId: string | null = null;
	export let onSelect: (playerId: string) => void;
	export let currentPlayerEmoji: string = '🎭';

	function handleCardClick(playerId: string) {
		onSelect(playerId);
	}
</script>

<div class="flex gap-4 justify-center flex-wrap">
	{#each players as player}
		{@const isSelected = selectedPlayerId === player.id}
		
		<button
			on:click={() => handleCardClick(player.id)}
			class="perspective-card"
		>
			<div class="player-card {isSelected ? 'selected' : ''}">
				<!-- Card representation -->
				<div class="relative">
					<CardBack width={100} height={140} variant="simplified" />
					
					<!-- Player initial overlay -->
					<div class="absolute inset-0 flex items-center justify-center">
						<div class="w-12 h-12 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold text-xl shadow-lg">
							{player.displayName[0].toUpperCase()}
						</div>
					</div>
				</div>
				
				<!-- Player name below card -->
				<p class="mt-2 text-sm font-medium text-center truncate max-w-[100px]">
					{player.displayName}
				</p>
			</div>
		</button>
	{/each}
	
	<!-- Current player card -->
	<div class="flex items-center mx-2">
		<div class="text-4xl text-muted-foreground">{currentPlayerEmoji}</div>
	</div>
	
	<div class="player-card-you">
		<div class="relative">
			<CardBack width={100} height={140} variant="simplified" />
			
			<!-- You indicator -->
			<div class="absolute inset-0 flex items-center justify-center">
				<div class="text-2xl font-bold text-primary">YOU</div>
			</div>
		</div>
		
		<p class="mt-2 text-sm font-medium text-center">You</p>
	</div>
</div>

<style>
	.perspective-card {
		perspective: 600px;
	}

	.player-card {
		transition: transform 0.2s ease;
		cursor: pointer;
	}

	.player-card:hover {
		transform: translateY(-8px);
	}

	.player-card.selected {
		transform: translateY(-12px) scale(1.05);
	}

	.player-card-you {
		opacity: 0.8;
	}
</style>

