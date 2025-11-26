<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { api, type JoinRoomRequest } from '$lib/api/client';
	import { session } from '$lib/stores/session.svelte';

	let joinCode = $state('');
	let displayName = $state('');
	let loading = $state(false);
	let error = $state('');

	// Handle error from redirect (expired/invalid room)
	$effect(() => {
		const errorParam = $page.url.searchParams.get('error');
		if (errorParam === 'room_expired') {
			error = 'That room has ended';
		} else if (errorParam === 'room_not_found') {
			error = 'Room not found';
		} else if (errorParam === 'connection_failed') {
			error = 'Failed to connect to server';
		}
	});

	async function handleJoin() {
		if (!joinCode.trim() || !displayName.trim()) {
			return;
		}

		loading = true;
		error = '';

		try {
			const request: JoinRoomRequest = {
				displayName: displayName.trim()
			};

			const response = await api.joinRoom(joinCode.toUpperCase(), request);

			// Save session
			session.set({
				playerId: response.playerId,
				sessionToken: response.sessionToken,
				roomCode: response.roomCode,
				displayName: displayName.trim()
			});

			// Redirect to room
			goto(`/room/${response.roomCode}`);
		} catch (err: any) {
			if (err.status === 404) {
				error = 'Room not found';
			} else {
				error = err.message || 'Failed to join room';
			}
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Cardless - Party Games</title>
</svelte:head>

<div class="flex flex-col items-center justify-center min-h-screen p-6 bg-background">
	<div class="w-full max-w-md space-y-8">
		<!-- Logo/Title -->
		<div class="text-center">
			<h1 class="text-6xl font-bold mb-2">🎲</h1>
			<h1 class="text-5xl font-bold text-primary mb-2">Cardless</h1>
			<p class="text-foreground/80 text-lg">Party games without the cards</p>
		</div>

		<!-- Main actions -->
		<div class="card space-y-4">
			<a href="/create" class="btn btn-primary w-full block text-center">
				Host a Game
			</a>

			<div class="relative">
				<div class="absolute inset-0 flex items-center">
					<div class="w-full border-t border-border"></div>
				</div>
				<div class="relative flex justify-center text-sm">
					<span class="px-4 bg-card text-muted-foreground">or join with code</span>
				</div>
			</div>

			<form onsubmit={(e) => { e.preventDefault(); handleJoin(); }} class="space-y-3">
				{#if error}
					<div class="p-3 bg-destructive/10 border border-destructive/20 rounded-lg">
						<p class="text-sm text-destructive">{error}</p>
					</div>
				{/if}

				<input
					type="text"
					bind:value={joinCode}
					placeholder="Enter room code"
					class="input text-center text-2xl font-mono tracking-wider uppercase"
					maxlength="6"
					autocomplete="off"
					autocapitalize="characters"
					disabled={loading}
				/>

				{#if joinCode.trim()}
					<input
						type="text"
						bind:value={displayName}
						placeholder="Enter your name"
						class="input"
						maxlength="20"
						autocomplete="off"
						disabled={loading}
					/>

					<button type="submit" class="btn btn-secondary w-full" disabled={loading || !displayName.trim()}>
						{loading ? 'Joining...' : 'Join Game'}
					</button>
				{/if}
			</form>
		</div>

		<!-- Info -->
		<div class="text-center text-muted-foreground text-sm">
			<p>No signup required • Play with friends in person</p>
		</div>
	</div>
</div>
