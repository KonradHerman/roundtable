<script lang="ts">
	/**
	 * Reusable role card component used throughout Werewolf game
	 * Ensures consistent styling for role reveals, center cards, and action results
	 */
	import { getRoleInfo } from './roleConfig';
	
	export let role: string;
	export let size: 'small' | 'medium' | 'large' = 'medium';
	export let showName: boolean = true;

	$: info = getRoleInfo(role);
	$: dimensions = {
		small: { width: '120px', height: '168px', emoji: 'text-5xl', text: 'text-sm' },
		medium: { width: '200px', height: '280px', emoji: 'text-8xl', text: 'text-2xl' },
		large: { width: '280px', height: '392px', emoji: 'text-9xl', text: 'text-3xl' }
	}[size];
</script>

<div 
	class="role-card {info.color} rounded-2xl shadow-lg text-white"
	style="width: {dimensions.width}; height: {dimensions.height}"
>
	<div class="h-full p-4 flex flex-col items-center justify-center gap-4">
		<div class="{dimensions.emoji}">
			{info.emoji}
		</div>
		{#if showName}
			<p class="{dimensions.text} font-bold capitalize text-center">
				{role.replace('_', ' ')}
			</p>
		{/if}
	</div>
</div>

<style>
	.role-card {
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
	}
</style>

