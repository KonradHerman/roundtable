<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { session } from '$lib/stores/session';
	import { gameStore } from '$lib/stores/game';
	import { createWebSocket } from '$lib/stores/websocket';
	import { api } from '$lib/api/client';
	import { Card, Button, Badge } from '$lib/components/ui';
	import { Users, Copy, Check } from 'lucide-svelte';
	import WerewolfGame from '$lib/games/werewolf/WerewolfGame.svelte';

	const roomCode = $page.params.code;

	let wsStore: ReturnType<typeof createWebSocket> | null = null;
	let connectionStatus: string = 'disconnected';
	let roomState: any = null;
	let gameState: any = null;
	let copied = false;

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
				const event = message.payload.event;
				gameStore.appendEvent(event);

				// Handle specific events
				if (event.type === 'game_started' || event.type === 'phase_changed') {
					// Refresh room state
					refreshRoomState();
				}
				break;

			case 'events':
				gameStore.appendEvents(message.payload.events);
				break;

			case 'error':
				console.error('Server error:', message.payload.message);
				break;
		}
	}

	async function refreshRoomState() {
		try {
			const state = await api.getRoomState(roomCode);
			roomState = state;
			gameStore.setRoomState(state);
		} catch (err) {
			console.error('Failed to refresh room state:', err);
		}
	}

	async function handleStartGame() {
		if (!$session) return;

		try {
			// Simple werewolf config with default roles
			const playerCount = roomState?.players?.length || 0;
			const roles = generateDefaultRoles(playerCount);

			await api.startGame(roomCode, {
				config: {
					roles,
					nightDuration: 30000000000,  // 30 seconds in nanoseconds
					dayDuration: 120000000000    // 2 minutes in nanoseconds
				}
			});
		} catch (err: any) {
			console.error('Failed to start game:', err);
			alert(err.message || 'Failed to start game');
		}
	}

	function generateDefaultRoles(playerCount: number): string[] {
		// One Night Werewolf requires playerCount + 3 roles (3 go to center)
		if (playerCount < 3) return [];
		
		const totalRoles = playerCount + 3; // Always 3 center cards
		const roles: string[] = [];

		// Start with core roles based on player count
		if (playerCount >= 3) {
			roles.push('werewolf', 'werewolf'); // Always 2 werewolves
			roles.push('seer'); // Always include seer
		}
		
		// Add special roles as player count increases
		if (playerCount >= 4) roles.push('robber');
		if (playerCount >= 5) roles.push('troublemaker');
		if (playerCount >= 6) roles.push('drunk');
		if (playerCount >= 7) roles.push('mason', 'mason'); // Masons come in pairs
		if (playerCount >= 8) roles.push('insomniac');
		if (playerCount >= 9) roles.push('minion');
		if (playerCount >= 10) roles.push('tanner'); // Chaos mode!

		// Fill remaining slots with villagers
		while (roles.length < totalRoles) {
			roles.push('villager');
		}

		return roles;
	}

	async function copyRoomCode() {
		try {
			await navigator.clipboard.writeText(roomCode);
			copied = true;
			setTimeout(() => copied = false, 2000);
		} catch (err) {
			console.error('Failed to copy:', err);
		}
	}

	$: isHost = $session?.playerId === roomState?.hostId;
	$: playerCount = roomState?.players?.length || 0;
	$: canStart = isHost && playerCount >= 3;
</script>

<svelte:head>
	<title>Room {roomCode} - Roundtable</title>
</svelte:head>

<div class="min-h-screen bg-background">
	<!-- Connection status banner -->
	{#if connectionStatus !== 'connected'}
		<div class="bg-yellow-500 text-white px-4 py-2 text-center text-sm font-medium">
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
	<div class="container mx-auto p-4 md:p-6 max-w-4xl">
		{#if !roomState}
			<!-- Loading -->
			<Card class="text-center py-12">
				<div class="text-4xl mb-4">🎲</div>
				<p class="text-muted-foreground">Loading room...</p>
			</Card>
		{:else if roomState.status === 'waiting'}
			<!-- Lobby view -->
			<div class="space-y-6">
				<!-- Room code card -->
				<Card class="p-6">
					<div class="text-center">
						<div class="flex items-center justify-center gap-2 mb-2">
							<h2 class="text-lg font-semibold text-muted-foreground">Room Code</h2>
							{#if isHost}
								<Badge variant="default">Host</Badge>
							{/if}
						</div>
						<button
							on:click={copyRoomCode}
							class="group relative inline-flex items-center gap-3 px-6 py-3 bg-primary/10 hover:bg-primary/20 rounded-lg transition-colors"
						>
							<span class="text-5xl font-mono font-bold tracking-wider text-primary">
								{roomCode}
							</span>
							{#if copied}
								<Check class="w-6 h-6 text-primary" />
							{:else}
								<Copy class="w-6 h-6 text-primary opacity-0 group-hover:opacity-100 transition-opacity" />
							{/if}
						</button>
						<p class="text-sm text-muted-foreground mt-3">
							Share this code with your friends
						</p>
					</div>
				</Card>

				<!-- Players card -->
				<Card class="p-6">
					<div class="flex items-center gap-2 mb-4">
						<Users class="w-5 h-5 text-muted-foreground" />
						<h2 class="font-semibold text-lg">
							Players ({playerCount}/{roomState.maxPlayers})
						</h2>
					</div>
					<div class="space-y-2">
						{#each roomState.players || [] as player}
							<div class="flex items-center gap-3 p-3 bg-muted/50 rounded-lg">
								<div class="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold text-lg">
									{player.displayName[0].toUpperCase()}
								</div>
								<div class="flex-1 min-w-0">
									<div class="font-medium truncate">{player.displayName}</div>
									<div class="text-sm text-muted-foreground">
										{player.connected ? '🟢 Online' : '🔴 Offline'}
									</div>
								</div>
								{#if player.id === roomState.hostId}
									<Badge variant="secondary">Host</Badge>
								{/if}
							</div>
						{/each}
					</div>
				</Card>

				<!-- Start game controls -->
				<Card class="p-6">
					{#if isHost}
						<div class="space-y-4">
							<div class="flex items-center justify-between p-4 bg-muted/50 rounded-lg">
								<div>
									<div class="font-medium">One Night Werewolf</div>
									<div class="text-sm text-muted-foreground">
										{playerCount} players • 3-10 recommended
									</div>
								</div>
							</div>
							<Button
								class="w-full h-14 text-lg"
								disabled={!canStart}
								on:click={handleStartGame}
							>
								Start Game
							</Button>
							{#if !canStart && playerCount < 3}
								<p class="text-sm text-muted-foreground text-center">
									Need at least 3 players to start
								</p>
							{/if}
						</div>
					{:else}
						<div class="text-center py-4">
							<p class="text-muted-foreground">
								Waiting for {roomState.players?.find(p => p.id === roomState.hostId)?.displayName || 'host'} to start the game...
							</p>
						</div>
					{/if}
				</Card>
			</div>
		{:else if roomState.status === 'playing'}
			<!-- Game view -->
			<WerewolfGame {roomCode} {roomState} {wsStore} />
		{:else if roomState.status === 'finished'}
			<!-- Results view -->
			<Card class="p-6">
				<h1 class="text-2xl font-bold mb-4">Game Finished</h1>
				<p class="text-muted-foreground">Results will appear here...</p>
				{#if isHost}
					<Button class="w-full mt-6" on:click={() => window.location.reload()}>
						Play Again
					</Button>
				{/if}
			</Card>
		{/if}
	</div>
</div>
