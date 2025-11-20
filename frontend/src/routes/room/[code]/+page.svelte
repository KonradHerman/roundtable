<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { session } from '$lib/stores/session.svelte';
	import { gameStore } from '$lib/stores/game.svelte';
	import { createWebSocket } from '$lib/stores/websocket.svelte';
	import { api } from '$lib/api/client';
	import { Card, Button, Badge } from '$lib/components/ui';
	import { Users, Copy, Check, QrCode, Share2 } from 'lucide-svelte';
	import WerewolfGame from '$lib/games/werewolf/WerewolfGame.svelte';
	import AvalonGame from '$lib/games/avalon/AvalonGame.svelte';
	import InviteQRCode from '$lib/components/InviteQRCode.svelte';
	import { browser } from '$app/environment';

	const roomCode = $page.params.code;

	let wsStore = $state<ReturnType<typeof createWebSocket> | null>(null);
	let connectionStatus = $state<string>('disconnected');
	let roomState = $state<any>(null);
	let previousStatus = $state<string | null>(null);
	let copied = $state(false);
	let selectedGame = $state<'werewolf' | 'avalon'>('werewolf');
	let showQRCode = $state(false);

	// Derived reactive values
	let isHost = $derived(session.value?.playerId === roomState?.hostId);
	let playerCount = $derived(roomState?.players?.length || 0);
	let canStart = $derived(
		isHost &&
		((selectedGame === 'werewolf' && playerCount >= 3) ||
		 (selectedGame === 'avalon' && playerCount >= 5 && playerCount <= 10))
	);
	let gameType = $derived(roomState?.gameType || 'werewolf');

	onMount(() => {
		const currentSession = session.value;

		// Verify session exists for this room
		if (!currentSession || currentSession.roomCode !== roomCode) {
			goto('/');
			return;
		}

		// Create WebSocket connection
		wsStore = createWebSocket(roomCode, currentSession.sessionToken);

		// Set up reactive effect to handle WebSocket status and messages
		let lastMessageCount = 0;
		
		$effect(() => {
			if (wsStore) {
				connectionStatus = wsStore.status;
				
				// Process new messages
				if (wsStore.messages.length > lastMessageCount) {
					const latestMessage = wsStore.messages[wsStore.messages.length - 1];
					handleMessage(latestMessage);
					lastMessageCount = wsStore.messages.length;
				}
			}
		});

		return () => {
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
				previousStatus = roomState.status;
				break;

			case 'room_state':
				const newRoomState = message.payload.roomState;
				
				// If room was reset (went from playing/finished back to waiting), clear game store
				if (previousStatus && (previousStatus === 'playing' || previousStatus === 'finished') && newRoomState.status === 'waiting') {
					console.log('Room reset detected, clearing game store');
					gameStore.reset();
				}
				
				roomState = newRoomState;
				gameStore.setRoomState(roomState);
				previousStatus = roomState.status;
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
		if (!roomCode) return;
		try {
			const state = await api.getRoomState(roomCode);
			roomState = state;
			gameStore.setRoomState(state);
		} catch (err) {
			console.error('Failed to refresh room state:', err);
		}
	}

	async function handleStartGame() {
		if (!session.value || !roomCode) return;

		try {
			const playerCount = roomState?.players?.length || 0;

			if (selectedGame === 'werewolf') {
				const roles = generateDefaultWerewolfRoles(playerCount);
				await api.startGame(roomCode, {
					gameType: 'werewolf',
					config: {
						roles,
						nightDuration: 30000000000,  // 30 seconds in nanoseconds
						dayDuration: 120000000000    // 2 minutes in nanoseconds
					}
				});
			} else if (selectedGame === 'avalon') {
				const roles = generateDefaultAvalonRoles(playerCount);
				await api.startGame(roomCode, {
					gameType: 'avalon',
					config: {
						roles
					}
				});
			}
		} catch (err: any) {
			console.error('Failed to start game:', err);
			alert(err.message || 'Failed to start game');
		}
	}

	function generateDefaultWerewolfRoles(playerCount: number): string[] {
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

	function generateDefaultAvalonRoles(playerCount: number): string[] {
		// Avalon default: Merlin + Assassin + fill with Loyal Servants and Minions
		const roles: string[] = [];

		// Team sizes based on player count
		const teamSizes: Record<number, [number, number]> = {
			5: [3, 2],   // 3 good, 2 evil
			6: [4, 2],
			7: [4, 3],
			8: [5, 3],
			9: [6, 3],
			10: [6, 4]
		};

		const [goodCount, evilCount] = teamSizes[playerCount] || [3, 2];

		// Always include Merlin and Assassin for interesting gameplay
		roles.push('merlin');
		roles.push('assassin');

		// Add Percival and Morgana for 7+ players (more interesting)
		if (playerCount >= 7) {
			roles.push('percival');
			roles.push('morgana');

			// Fill remaining
			for (let i = 0; i < goodCount - 2; i++) roles.push('loyal_servant');
			for (let i = 0; i < evilCount - 2; i++) roles.push('minion');
		} else {
			// Fill remaining with Loyal Servants and Minions
			for (let i = 0; i < goodCount - 1; i++) roles.push('loyal_servant');
			for (let i = 0; i < evilCount - 1; i++) roles.push('minion');
		}

		return roles;
	}

	async function copyRoomCode() {
		if (!roomCode) return;
		try {
			await navigator.clipboard.writeText(roomCode);
			copied = true;
			setTimeout(() => copied = false, 2000);
		} catch (err) {
			console.error('Failed to copy:', err);
		}
	}

	async function shareInviteLink() {
		if (!browser || !roomCode) return;

		const inviteUrl = `${window.location.origin}/join/${roomCode}`;

		// Check if Web Share API is available (mainly mobile)
		if (navigator.share) {
			try {
				await navigator.share({
					title: 'Join my Cardless game!',
					text: `Join room ${roomCode}`,
					url: inviteUrl
				});
			} catch (err) {
				// User cancelled or share failed
				if (err instanceof Error && err.name !== 'AbortError') {
					console.error('Share failed:', err);
					// Fallback to copy
					await navigator.clipboard.writeText(inviteUrl);
				}
			}
		} else {
			// Fallback to copy on desktop
			try {
				await navigator.clipboard.writeText(inviteUrl);
				copied = true;
				setTimeout(() => copied = false, 2000);
			} catch (err) {
				console.error('Failed to copy:', err);
			}
		}
	}
</script>

<svelte:head>
	<title>Room {roomCode} - Cardless</title>
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
						onclick={copyRoomCode}
							class="group relative inline-flex items-center gap-3 px-6 py-3 bg-primary/10 hover:bg-primary/20 rounded-lg transition-colors"
						>
							<span class="text-5xl font-mono font-bold tracking-wider text-primary">
								{roomCode}
							</span>
							{#if copied}
								<Check class="w-6 h-6 text-primary" />
							{:else}
								<Copy class="w-6 h-6 text-primary md:opacity-0 md:group-hover:opacity-100 transition-opacity" />
							{/if}
						</button>
						<p class="text-sm text-muted-foreground mt-3">
							Share this code with your friends
						</p>

						<!-- Action buttons -->
						<div class="flex gap-2 mt-4">
							<button
								onclick={shareInviteLink}
								class="flex-1 inline-flex items-center justify-center gap-2 px-4 py-2 text-sm bg-primary text-primary-foreground hover:bg-primary/90 rounded-lg transition-colors font-medium"
							>
								<Share2 class="w-4 h-4" />
								<span>Share Link</span>
							</button>

							<button
								onclick={() => showQRCode = !showQRCode}
								class="flex-1 inline-flex items-center justify-center gap-2 px-4 py-2 text-sm bg-secondary/10 hover:bg-secondary/20 text-secondary rounded-lg transition-colors"
							>
								<QrCode class="w-4 h-4" />
								<span>{showQRCode ? 'Hide' : 'Show'} QR Code</span>
							</button>
						</div>
					</div>
				</Card>

				<!-- QR Code card (collapsible) -->
				{#if showQRCode && roomCode}
					<Card class="p-6">
						<InviteQRCode roomCode={roomCode} />
					</Card>
				{/if}

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
							<!-- Game selection -->
							<div>
								<label class="text-sm font-medium mb-2 block">Select Game</label>
								<div class="grid grid-cols-2 gap-3">
									<button
										class="p-4 rounded-lg border-2 transition-all {selectedGame === 'werewolf' ? 'border-primary bg-primary/10' : 'border-muted-foreground/20 hover:border-primary/50'}"
										onclick={() => selectedGame = 'werewolf'}
									>
										<div class="text-3xl mb-2">🐺</div>
										<div class="font-medium">One Night Werewolf</div>
										<div class="text-xs text-muted-foreground mt-1">3-10 players</div>
									</button>
									<button
										class="p-4 rounded-lg border-2 transition-all {selectedGame === 'avalon' ? 'border-primary bg-primary/10' : 'border-muted-foreground/20 hover:border-primary/50'}"
										onclick={() => selectedGame = 'avalon'}
									>
										<div class="text-3xl mb-2">🗡️</div>
										<div class="font-medium">Avalon</div>
										<div class="text-xs text-muted-foreground mt-1">5-10 players</div>
									</button>
								</div>
							</div>

							<!-- Selected game info -->
							<div class="p-4 bg-muted/50 rounded-lg">
								<div class="font-medium">
									{selectedGame === 'werewolf' ? 'One Night Werewolf' : 'Avalon (The Resistance)'}
								</div>
								<div class="text-sm text-muted-foreground">
									{playerCount} players • {selectedGame === 'werewolf' ? '3-10' : '5-10'} recommended
								</div>
							</div>

							<Button
								class="w-full h-14 text-lg"
								disabled={!canStart}
								onclick={handleStartGame}
							>
								Start Game
							</Button>
							{#if !canStart}
								<p class="text-sm text-muted-foreground text-center">
									{#if selectedGame === 'werewolf' && playerCount < 3}
										Need at least 3 players to start Werewolf
									{:else if selectedGame === 'avalon' && playerCount < 5}
										Need at least 5 players to start Avalon
									{:else if selectedGame === 'avalon' && playerCount > 10}
										Avalon supports maximum 10 players
									{/if}
								</p>
							{/if}
						</div>
					{:else}
						<div class="text-center py-4">
						<p class="text-muted-foreground">
							Waiting for {roomState.players?.find((p: any) => p.id === roomState.hostId)?.displayName || 'host'} to start the game...
						</p>
						</div>
					{/if}
				</Card>
			</div>
		{:else if roomState.status === 'playing' && roomCode && wsStore}
			<!-- Game view -->
			{#if gameType === 'werewolf'}
				<WerewolfGame {roomCode} {roomState} {wsStore} />
			{:else if gameType === 'avalon'}
				<AvalonGame {roomCode} {roomState} {wsStore} />
			{:else}
				<Card class="p-6">
					<p class="text-center text-muted-foreground">Unknown game type: {gameType}</p>
				</Card>
			{/if}
		{:else if roomState.status === 'finished'}
			<!-- Results view -->
			<Card class="p-6">
				<h1 class="text-2xl font-bold mb-4">Game Finished</h1>
				<p class="text-muted-foreground">Results will appear here...</p>
				{#if isHost}
					<Button class="w-full mt-6" onclick={() => window.location.reload()}>
						Play Again
					</Button>
				{/if}
			</Card>
		{/if}
	</div>
</div>
