<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';
	import { globalConnectionStatus } from '$lib/stores/websocket';
	import { fade } from 'svelte/transition';
	import { WifiOff, Loader2 } from 'lucide-svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
	}

	let { children }: Props = $props();

	// Use $derived to properly subscribe to the store in Svelte 5
	let status = $derived($globalConnectionStatus);

	// Only show connection status when in a room
	let isInRoom = $derived($page.url.pathname.startsWith('/room/'));
</script>

<div class="min-h-screen relative">
	{#if isInRoom && (status === 'reconnecting' || status === 'disconnected')}
		<div
			transition:fade
			class="fixed top-0 left-0 right-0 z-50 p-2 text-center text-white font-medium flex items-center justify-center gap-2 {status === 'reconnecting' ? 'bg-yellow-500' : 'bg-red-500'}"
		>
			{#if status === 'reconnecting'}
				<Loader2 class="w-4 h-4 animate-spin" />
				<span>Reconnecting...</span>
			{:else}
				<WifiOff class="w-4 h-4" />
				<span>Connection Lost</span>
			{/if}
		</div>
	{/if}

	{@render children()}
</div>
