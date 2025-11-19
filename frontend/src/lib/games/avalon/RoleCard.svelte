<script lang="ts">
	/**
	 * Reusable role card component for Avalon game
	 * Ensures consistent styling for role reveals and results
	 */
	import { roleConfig, type AvalonRole } from './roleConfig';

	let {
		role,
		size = 'medium',
		showName = true
	}: {
		role: AvalonRole;
		size?: 'small' | 'medium' | 'large';
		showName?: boolean;
	} = $props();

	const info = $derived(roleConfig[role]);
	const dimensions = $derived(
		{
			small: { width: '120px', height: '168px', emoji: 'text-5xl', text: 'text-sm' },
			medium: { width: '200px', height: '280px', emoji: 'text-8xl', text: 'text-2xl' },
			large: { width: '280px', height: '392px', emoji: 'text-9xl', text: 'text-3xl' }
		}[size]
	);
</script>

<div
	class="role-card {info.color} rounded-2xl shadow-lg text-white"
	style="width: {dimensions.width}; height: {dimensions.height}"
>
	<div class="h-full p-4 flex flex-col items-center justify-center gap-4">
		<div class={dimensions.emoji}>
			{info.emoji}
		</div>
		{#if showName}
			<p class="{dimensions.text} font-bold text-center">
				{info.name}
			</p>
		{/if}
	</div>
</div>

<style>
	.role-card {
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
	}
</style>
