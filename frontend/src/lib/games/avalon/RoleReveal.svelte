<script lang="ts">
	/**
	 * Role Reveal component - shows player their role, team, and knowledge
	 * Each player sees this on their own device after game starts
	 */
	import { roleConfig, type AvalonRole, type Team } from './roleConfig';
	import RoleCard from './RoleCard.svelte';

	let {
		role,
		team,
		knowledge = [],
		players,
		onAcknowledge
	}: {
		role: AvalonRole;
		team: Team;
		knowledge: string[];
		players: Array<{ id: string; name: string }>;
		onAcknowledge: () => void;
	} = $props();

	let acknowledged = $state(false);

	const config = $derived(roleConfig[role]);

	function getPlayerName(id: string): string {
		return players.find((p) => p.id === id)?.name || 'Unknown';
	}

	function handleAcknowledge() {
		acknowledged = true;
		onAcknowledge();
	}
</script>

<div class="role-reveal max-w-md mx-auto p-6">
	<h2 class="text-center mb-6 text-2xl font-bold text-[#d79921]">Your Role</h2>

	<div class="flex justify-center mb-6">
		<RoleCard {role} size="large" />
	</div>

	<div class="team-banner {team} text-center py-4 px-6 rounded-lg mb-6 font-bold text-lg">
		{team === 'good' ? '⚔️ Good Team' : '💀 Evil Team'}
	</div>

	<div class="knowledge-section bg-[#3c3836] p-6 rounded-lg">
		<h3 class="text-[#d79921] mb-3 font-bold">What You Know</h3>
		<p class="text-[#ebdbb2] leading-relaxed mb-4">
			{config.knowledge}
		</p>

		{#if knowledge.length > 0}
			<div class="known-players flex flex-wrap gap-2 mt-4">
				{#each knowledge as playerId}
					<div class="player-badge bg-[#504945] px-4 py-2 rounded text-[#fbf1c7]">
						{getPlayerName(playerId)}
					</div>
				{/each}
			</div>
		{/if}
	</div>

	{#if !acknowledged}
		<button class="btn-primary w-full mt-8 py-4 text-lg font-bold" onclick={handleAcknowledge}>
			I Understand My Role
		</button>
	{:else}
		<div class="text-center text-[#a89984] mt-8">Waiting for other players...</div>
	{/if}
</div>

<style>
	.team-banner.good {
		background: #458588;
		color: white;
	}

	.team-banner.evil {
		background: #cc241d;
		color: white;
	}

	.btn-primary {
		background: #d79921;
		color: #282828;
		border: none;
		border-radius: 0.5rem;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-primary:hover:not(:disabled) {
		background: #fabd2f;
		transform: translateY(-2px);
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
</style>
