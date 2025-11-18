<script lang="ts">
	import { CardBack, Button } from '$lib/components/ui';
	
	interface Props {
		cards?: number[];
		selectedCards?: number[];
		flippedCards?: Record<number, string>;
		maxSelection?: number;
		onSelect: (index: number) => void;
		mode?: 'select' | 'reveal';
	}

	let { 
		cards = [0, 1, 2],
		selectedCards = [],
		flippedCards = {},
		maxSelection = 1,
		onSelect,
		mode = 'select'
	} = $props<Props>();

	function handleCardClick(index: number) {
		if (mode === 'reveal') return;
		onSelect(index);
	}

	const roleInfo: Record<string, { emoji: string; color: string }> = {
		werewolf: { emoji: '🐺', color: 'bg-red-600' },
		seer: { emoji: '🔮', color: 'bg-purple-600' },
		robber: { emoji: '🎭', color: 'bg-blue-600' },
		troublemaker: { emoji: '😈', color: 'bg-orange-600' },
		mason: { emoji: '🔨', color: 'bg-gray-600' },
		villager: { emoji: '👤', color: 'bg-green-600' },
		minion: { emoji: '😤', color: 'bg-red-700' },
		tanner: { emoji: '🤪', color: 'bg-yellow-600' },
		drunk: { emoji: '🍺', color: 'bg-amber-600' },
		insomniac: { emoji: '😴', color: 'bg-purple-700' }
	};

	function getRoleInfo(role: string) {
		return roleInfo[role] || { emoji: '❓', color: 'bg-muted' };
	}
</script>

<div class="flex gap-4 justify-center flex-wrap">
	{#each cards as index}
		{@const isFlipped = index in flippedCards}
		{@const isSelected = selectedCards.includes(index)}
		{@const role = flippedCards[index]}
		{@const info = role ? getRoleInfo(role) : null}
		
		<button
			onclick={() => handleCardClick(index)}
			class="perspective-1000"
			disabled={mode === 'reveal'}
		>
			<div class="center-card-container {isFlipped ? 'flipped' : ''} {isSelected ? 'selected' : ''}">
				<!-- Card Back -->
				<div class="card-face card-back">
					<CardBack width={120} height={168} variant="simplified" />
				</div>

				<!-- Card Front -->
				{#if isFlipped && info}
					<div class="card-face card-front {info.color}">
						<div class="center-card-front p-4 rounded-lg flex flex-col items-center justify-center text-white">
							<div class="text-5xl mb-2">
								{info.emoji}
							</div>
							<p class="text-sm font-bold capitalize text-center">
								{role.replace('_', ' ')}
							</p>
						</div>
					</div>
				{:else}
					<div class="card-face card-front bg-card">
						<div class="center-card-front p-4 rounded-lg border-2 border-primary flex items-center justify-center">
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

