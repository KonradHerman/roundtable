<script lang="ts">
	import { fade } from 'svelte/transition';
	import { Button } from '$lib/components/ui';
	import { onMount } from 'svelte';

	export let role: string;
	export let acknowledged: boolean = false;
	export let acknowledgementsCount: number = 0;
	export let totalPlayers: number = 0;
	export let onAcknowledge: () => void;

	let flipped = false;

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

	onMount(() => {
		// Flip the card after a short delay
		setTimeout(() => {
			flipped = true;
		}, 300);
	});
</script>

<div
	class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/90"
	in:fade={{ duration: 300 }}
	out:fade={{ duration: 300 }}
>
	<div class="perspective-1000 w-full max-w-md">
		<div class="card-container {flipped ? 'flipped' : ''}">
			<!-- Card Back -->
			<div class="card-face card-back">
				<div class="p-8 bg-card border-2 border-primary rounded-2xl h-[500px] flex items-center justify-center">
					<div class="text-center">
						<div class="text-6xl mb-4">🎴</div>
						<p class="text-xl font-bold text-foreground">Your Role</p>
					</div>
				</div>
			</div>

			<!-- Card Front -->
			<div class="card-face card-front {info.color}">
				<div class="p-8 rounded-2xl h-[500px] flex flex-col items-center justify-between text-white">
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

					{#if !acknowledged}
						<!-- Acknowledgement button -->
						<div class="w-full space-y-3">
							<Button
								on:click={onAcknowledge}
								class="w-full h-12 bg-white hover:bg-white/90 text-gray-900 font-bold"
							>
								I've Seen My Role ✓
							</Button>
							<p class="text-sm text-white/75 text-center">
								Click when you've memorized your role
							</p>
						</div>
					{:else}
						<!-- Waiting for others -->
						<div class="w-full">
							<div class="bg-white/20 backdrop-blur rounded-lg p-4 text-center">
								<p class="font-semibold mb-2">Waiting for other players...</p>
								<p class="text-2xl font-bold">
									{acknowledgementsCount} / {totalPlayers}
								</p>
								<p class="text-sm text-white/75 mt-1">players ready</p>
							</div>
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	.perspective-1000 {
		perspective: 1000px;
	}

	.card-container {
		position: relative;
		width: 100%;
		height: 500px;
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
</style>
