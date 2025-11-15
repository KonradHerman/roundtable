<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { session } from '$lib/stores/session';
	import { gameStore } from '$lib/stores/game';
	import { createWebSocket } from '$lib/stores/websocket';

	const roomCode = $page.params.code;

	let wsStore: ReturnType<typeof createWebSocket> | null = null;
	let connectionStatus: string = 'disconnected';
	let roomState: any = null;

	onMount(() => {
		const currentSession = $session;

		// Verify session exists for this room
		if (!currentSession || currentSession.roomCode !== roomCode) {
			goto('/');
			return;
		}

		// Create WebSocket connection
		wsStore = createWebSocket(roomCode, currentSession.sessionToken);

		// Subscribe to WebSocket messages
		const unsubscribe = wsStore.subscribe(($ws) => {
			connectionStatus = $ws.status;

			// Process messages
			if ($ws.messages.length > 0) {
				const latestMessage = $ws.messages[$ws.messages.length - 1];
				handleMessage(latestMessage);
			}
		});

		return () => {
			unsubscribe();
			if (wsStore) {
				wsStore.disconnect();
			}
		};
	});

	function handleMessage(message: any) {
		console.log('Received message:', message);

		switch (message.type) {
			case 'authenticated':
				roomState = message.payload.roomState;
				gameStore.setRoomState(roomState);
				break;

			case 'room_state':
				roomState = message.payload.roomState;
				gameStore.setRoomState(roomState);
				break;

			case 'event':
				gameStore.appendEvent(message.payload.event);
				break;

			case 'events':
				gameStore.appendEvents(message.payload.events);
				break;

			case 'error':
				console.error('Server error:', message.payload.message);
				break;
		}
	}
</script>

<svelte:head>
	<title>Room {roomCode} - Roundtable</title>
</svelte:head>

<div class="min-h-screen bg-gray-50">
	<!-- Connection status banner -->
	{#if connectionStatus !== 'connected'}
		<div class="bg-yellow-500 text-white px-4 py-2 text-center text-sm">
			{#if connectionStatus === 'connecting'}
				Connecting...
			{:else if connectionStatus === 'reconnecting'}
				Reconnecting...
			{:else}
				Disconnected
			{/if}
		</div>
	{/if}

	<!-- Main content -->
	<div class="container mx-auto p-6 max-w-4xl">
		{#if !roomState}
			<!-- Loading -->
			<div class="card text-center py-12">
				<div class="text-4xl mb-4">🎲</div>
				<p class="text-gray-600">Loading room...</p>
			</div>
		{:else if roomState.status === 'waiting'}
			<!-- Lobby view -->
			<div class="card">
				<div class="text-center mb-6">
					<h1 class="text-3xl font-bold mb-2">Room Code</h1>
					<div class="text-5xl font-mono font-bold tracking-wider text-primary-600 mb-4">
						{roomCode}
					</div>
					<p class="text-gray-600">Waiting for players to join...</p>
				</div>

				<!-- Players list -->
				<div class="space-y-2 mb-6">
					<h2 class="font-semibold text-lg mb-2">
						Players ({roomState.players?.length || 0}/{roomState.maxPlayers})
					</h2>
					{#each roomState.players || [] as player}
						<div class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
							<div class="w-10 h-10 rounded-full bg-primary-500 text-white flex items-center justify-center font-bold">
								{player.displayName[0].toUpperCase()}
							</div>
							<div class="flex-1">
								<div class="font-medium">{player.displayName}</div>
								<div class="text-sm text-gray-500">
									{player.connected ? '🟢 Online' : '🔴 Offline'}
								</div>
							</div>
							{#if player.id === roomState.hostId}
								<span class="text-xs bg-primary-100 text-primary-700 px-2 py-1 rounded-full">Host</span>
							{/if}
						</div>
					{/each}
				</div>

				<!-- Host controls -->
				{#if $session?.playerId === roomState.hostId}
					<div class="border-t pt-4">
						<button
							class="btn btn-primary w-full"
							disabled={!roomState.players || roomState.players.length < 3}
						>
							Start Game
						</button>
						{#if roomState.players && roomState.players.length < 3}
							<p class="text-sm text-gray-500 text-center mt-2">
								Need at least 3 players to start
							</p>
						{/if}
					</div>
				{:else}
					<div class="text-center text-gray-500 text-sm">
						Waiting for host to start the game...
					</div>
				{/if}
			</div>
		{:else if roomState.status === 'playing'}
			<!-- Game view -->
			<div class="card">
				<h1 class="text-2xl font-bold mb-4">Game in Progress</h1>
				<p class="text-gray-600">Game UI will go here...</p>
			</div>
		{:else if roomState.status === 'finished'}
			<!-- Results view -->
			<div class="card">
				<h1 class="text-2xl font-bold mb-4">Game Finished</h1>
				<p class="text-gray-600">Results will go here...</p>
			</div>
		{/if}
	</div>
</div>
