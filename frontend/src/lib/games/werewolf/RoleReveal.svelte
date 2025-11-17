<script lang="ts">
	import { Button, CardBack } from '$lib/components/ui';

	export let role: string;
	export let acknowledged: boolean = false;
	export let acknowledgementsCount: number = 0;
	export let totalPlayers: number = 0;
	export let onAcknowledge: () => void;

	let cardFlipped = false;
	let readyToAcknowledge = false;

	const roleInfo: Record<string, { emoji: string; team: string; description: string; color: string }> = {
		werewolf: {
			emoji: '🐺',
			team: 'Werewolf Team',
			description: 'Find your fellow werewolves and survive the vote',
			color: 'bg-gruvbox-red'
		},
		seer: {
			emoji: '🔮',
			team: 'Village Team',
			description: 'Look at one player\'s role to help find the werewolves',
			color: 'bg-gruvbox-purple'
		},
		robber: {
			emoji: '🎭',
			team: 'Village Team',
			description: 'Swap roles with another player',
			color: 'bg-gruvbox-blue'
		},
		troublemaker: {
			emoji: '😈',
			team: 'Village Team',
			description: 'Swap two other players\' roles',
			color: 'bg-gruvbox-orange'
		},
		mason: {
			emoji: '🔨',
			team: 'Village Team',
			description: 'Know who the other mason is',
			color: 'bg-muted'
		},
		villager: {
			emoji: '👤',
			team: 'Village Team',
			description: 'Use your wits to find the werewolves',
			color: 'bg-gruvbox-green'
		},
		minion: {
			emoji: '😤',
			team: 'Werewolf Team',
			description: 'Know the werewolves but they don\'t know you',
			color: 'bg-gruvbox-red'
		},
		tanner: {
			emoji: '🤪',
			team: 'Solo',
			description: 'You win if YOU get eliminated',
			color: 'bg-gruvbox-yellow'
		},
		drunk: {
			emoji: '🍺',
			team: 'Village Team',
			description: 'You must swap your role but won\'t know your new role',
			color: 'bg-gruvbox-orange'
		},
		insomniac: {
			emoji: '😴',
			team: 'Village Team',
			description: 'Wake up last to see if your role changed',
			color: 'bg-gruvbox-purple-light'
		}
	};

	$: info = roleInfo[role] || {
		emoji: '❓',
		team: 'Unknown',
		description: 'Unknown role',
		color: 'bg-muted'
	};

	function handleShowRole() {
		if (!cardFlipped) {
			cardFlipped = true;
			readyToAcknowledge = true;
		}
	}

	function handleReady() {
		if (readyToAcknowledge && !acknowledged) {
			onAcknowledge();
		}
	}

	function handleUnready() {
		if (acknowledged) {
			// Can't unready after acknowledging
			return;
		}
		cardFlipped = !cardFlipped;
	}
</script>

<div class="w-full space-y-6">
		<!-- Card Container -->
		<div class="perspective-1000">
			<div class="card-container {cardFlipped ? 'flipped' : ''}">
				<!-- Card Back -->
				<div class="card-face card-back">
					<div class="game-card border-2 border-primary">
						<CardBack width={280} height={392} variant="simplified" />
					</div>
				</div>

				<!-- Card Front -->
				<div class="card-face card-front {info.color}">
					<div class="game-card p-8 rounded-2xl flex flex-col items-center justify-between text-white">
						<div class="flex-1 flex flex-col items-center justify-center space-y-6 w-full">
							<!-- Role emoji -->
							<div class="text-8xl">
								{info.emoji}
							</div>

							<!-- Role name -->
							<div class="text-center">
								<h2 class="text-4xl font-bold mb-2 capitalize">
									{role.replace('_', ' ')}
								</h2>
								<p class="text-xl text-white/90">
									{info.team}
								</p>
							</div>

							<!-- Description -->
							<p class="text-lg text-white/90 text-center">
								{info.description}
							</p>
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- Action Buttons -->
		{#if !acknowledged}
			<div class="space-y-3">
				{#if !cardFlipped}
					<Button
						on:click={handleShowRole}
						class="w-full h-14 bg-primary hover:bg-primary/90 text-primary-foreground font-bold text-lg"
					>
						👁️ Show Role
					</Button>
					<p class="text-sm text-white/75 text-center">
						Tap to peek at your role card
					</p>
				{:else}
					<div class="flex gap-3">
						<Button
							on:click={handleUnready}
							variant="outline"
							class="flex-1 h-14 bg-card hover:bg-muted text-foreground font-bold text-lg border-2"
						>
							↺ Look Again
						</Button>
						<Button
							on:click={handleReady}
							class="flex-1 h-14 bg-green-600 hover:bg-green-700 text-white font-bold text-lg"
						>
							✓ Ready
						</Button>
					</div>
					<p class="text-sm text-white/75 text-center">
						Ready when you've memorized your role
					</p>
				{/if}
			</div>
		{:else}
			<!-- Waiting for others -->
			<div class="bg-white/10 backdrop-blur rounded-lg p-6 text-center text-white">
				<p class="font-semibold mb-3 text-lg">Waiting for other players...</p>
				<p class="text-4xl font-bold mb-2">
					{acknowledgementsCount} / {totalPlayers}
				</p>
				<p class="text-sm text-white/75">players ready</p>
			</div>
		{/if}
</div>

<style>
	.perspective-1000 {
		perspective: 1000px;
		display: flex;
		justify-content: center;
	}

	.card-container {
		position: relative;
		width: 280px;
		height: 392px;
		transition: transform 0.8s cubic-bezier(0.4, 0.0, 0.2, 1);
		transform-style: preserve-3d;
	}

	.card-container.flipped {
		transform: rotateY(180deg);
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

	.game-card {
		width: 280px;
		height: 392px;
		border-radius: 12px;
		box-shadow: 0 8px 16px rgba(0, 0, 0, 0.3);
		overflow: hidden; /* Ensure rounded corners apply to content */
	}
</style>
