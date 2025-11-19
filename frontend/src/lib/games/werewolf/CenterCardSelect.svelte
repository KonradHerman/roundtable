<script lang="ts">
	import { CardBack } from '$lib/components/ui';
	import RoleCard from './RoleCard.svelte';
	
	export let cards: number[] = [0, 1, 2];
	export let selectedCards: number[] = [];
	export let flippedCards: Record<number, string> = {}; // index -> role
	export let maxSelection: number = 1;
	export let onSelect: (index: number) => void;
	export let mode: 'select' | 'reveal' = 'select'; // select for clickable, reveal for showing result

	function handleCardClick(index: number) {
		if (mode === 'reveal') return;
		onSelect(index);
	}
</script>

<div class="flex gap-4 justify-center flex-wrap">
	{#each cards as index}
		{@const isFlipped = index in flippedCards}
		{@const isSelected = selectedCards.includes(index)}
		{@const role = flippedCards[index]}
		{@const info = role ? getRoleInfo(role) : null}
		
		<button
			on:click={() => handleCardClick(index)}
			class="perspective-1000"
			disabled={mode === 'reveal'}
		>
			<div class="center-card-container {isFlipped ? 'flipped' : ''} {isSelected ? 'selected' : ''}">
				<!-- Card Back -->
				<div class="card-face card-back">
					<CardBack width={120} height={168} variant="simplified" />
				</div>

				<!-- Card Front -->
				{#if isFlipped && role}
					<div class="card-face card-front">
						<RoleCard {role} size="small" />
					</div>
				{:else}
					<div class="card-face card-front bg-card">
						<div class="center-card-front p-4 rounded-2xl border-2 border-primary flex items-center justify-center">
							<p class="text-xs text-muted-foreground">Card {index + 1}</p>
						</div>
					</div>
				{/if}
			</div>
		</button>
	{/each}
</div>

<style>
	.perspective-1000 {
		perspective: 600px;
	}

	.center-card-container {
		position: relative;
		width: 120px;
		height: 168px;
		transition: transform 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
		transform-style: preserve-3d;
		cursor: pointer;
	}

	.center-card-container:hover {
		transform: translateY(-4px);
	}

	.center-card-container.selected {
		transform: translateY(-8px) scale(1.05);
		box-shadow: 0 8px 16px rgba(215, 153, 33, 0.3);
	}

	.center-card-container.flipped {
		transform: rotateY(180deg);
	}

	.center-card-container.flipped:hover {
		transform: rotateY(180deg) translateY(-4px);
	}
	
	.center-card-container.flipped.selected {
		transform: rotateY(180deg) scale(1.05);
	}

	.card-face {
		position: absolute;
		width: 100%;
		height: 100%;
		backface-visibility: hidden;
		-webkit-backface-visibility: hidden;
	}

	.card-back {
		z-index: 2;
		transform: rotateY(0deg);
	}

	.card-front {
		transform: rotateY(180deg);
	}

	.center-card-front {
		width: 120px;
		height: 168px;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
	}

	button:disabled .center-card-container {
		cursor: default;
	}
</style>

