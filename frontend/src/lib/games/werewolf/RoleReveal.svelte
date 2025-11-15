<script lang="ts">
	import { fade, scale } from 'svelte/transition';
	import { Card } from '$lib/components/ui';

	export let role: string;

	const roleInfo: Record<string, { emoji: string; team: string; description: string; color: string }> = {
		werewolf: {
			emoji: '🐺',
			team: 'Werewolf Team',
			description: 'Find your fellow werewolves and survive the vote',
			color: 'from-red-600 to-red-800'
		},
		seer: {
			emoji: '🔮',
			team: 'Village Team',
			description: 'Look at one player\'s role to help find the werewolves',
			color: 'from-purple-600 to-purple-800'
		},
		robber: {
			emoji: '🎭',
			team: 'Village Team',
			description: 'Swap roles with another player',
			color: 'from-blue-600 to-blue-800'
		},
		troublemaker: {
			emoji: '😈',
			team: 'Village Team',
			description: 'Swap two other players\' roles',
			color: 'from-orange-600 to-orange-800'
		},
		mason: {
			emoji: '🔨',
			team: 'Village Team',
			description: 'Know who the other mason is',
			color: 'from-gray-600 to-gray-800'
		},
		villager: {
			emoji: '👤',
			team: 'Village Team',
			description: 'Use your wits to find the werewolves',
			color: 'from-green-600 to-green-800'
		},
		minion: {
			emoji: '😤',
			team: 'Werewolf Team',
			description: 'Know the werewolves but they don\'t know you',
			color: 'from-red-600 to-red-900'
		},
		tanner: {
			emoji: '🤪',
			team: 'Solo',
			description: 'You win if YOU get eliminated',
			color: 'from-yellow-600 to-yellow-800'
		},
		drunk: {
			emoji: '🍺',
			team: 'Village Team',
			description: 'You must swap your role but won\'t know your new role',
			color: 'from-amber-600 to-amber-800'
		},
		insomniac: {
			emoji: '😴',
			team: 'Village Team',
			description: 'Wake up last to see if your role changed',
			color: 'from-indigo-600 to-indigo-800'
		}
	};

	$: info = roleInfo[role] || {
		emoji: '❓',
		team: 'Unknown',
		description: 'Unknown role',
		color: 'from-gray-600 to-gray-800'
	};
</script>

<div
	class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80"
	in:fade={{ duration: 300 }}
	out:fade={{ duration: 300 }}
>
	<div in:scale={{ duration: 500, delay: 100, start: 0.5 }}>
		<Card
			class="max-w-md w-full p-8 bg-gradient-to-br {info.color} text-white border-0"
		>
			<div class="text-center space-y-6">
				<!-- Role emoji -->
				<div class="text-8xl animate-bounce">
					{info.emoji}
				</div>

				<!-- Role name -->
				<div>
					<h2 class="text-4xl font-bold mb-2 capitalize">
						{role.replace('_', ' ')}
					</h2>
					<p class="text-xl text-white/90">
						{info.team}
					</p>
				</div>

				<!-- Description -->
				<p class="text-lg text-white/90">
					{info.description}
				</p>

				<!-- Instruction -->
				<div class="pt-4 border-t border-white/20">
					<p class="text-sm text-white/75">
						Remember your role - this message will disappear soon
					</p>
				</div>
			</div>
		</Card>
	</div>
</div>

<style>
	@keyframes bounce {
		0%, 100% {
			transform: translateY(0);
		}
		50% {
			transform: translateY(-20px);
		}
	}

	.animate-bounce {
		animation: bounce 2s ease-in-out infinite;
	}
</style>
