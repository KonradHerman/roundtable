<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { api, type JoinRoomRequest, APIError } from '$lib/api/client';
	import { session } from '$lib/stores/session.svelte';

	const roomCode = ($page.params.code || '').toUpperCase();

	let displayName = $state('');
	let loading = $state(false);
	let validating = $state(true);
	let error = $state('');
	let roomValid = $state(false);

	$effect(async () => {
		// Validate room exists and is joinable
		try {
			const roomState = await api.getRoomState(roomCode);

			// Check if room is in a joinable state
			if (roomState.status === 'waiting') {
				roomValid = true;
			} else if (roomState.status === 'finished') {
				// Room is finished/expired
				goto('/?error=room_expired');
				return;
			} else {
				// Room is in progress, but we can still allow joining
				roomValid = true;
			}
		} catch (err: any) {
			if (err instanceof APIError && err.status === 404) {
				// Room not found
				goto('/?error=room_not_found');
			} else {
				// Other error
				goto('/?error=connection_failed');
			}
			return;
		} finally {
			validating = false;
		}
	});

	async function handleJoin() {
		if (!displayName.trim()) {
			return;
		}

		loading = true;
		error = '';

		try {
			const request: JoinRoomRequest = {
				displayName: displayName.trim()
			};

			const response = await api.joinRoom(roomCode, request);

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
			if (err instanceof APIError) {
				if (err.status === 404) {
					error = 'Room not found or has ended';
				} else if (err.status === 400) {
					error = 'Room is full or cannot be joined';
				} else {
					error = err.message || 'Failed to join room';
				}
			} else {
				error = 'Failed to join room';
			}
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Join {roomCode} - Cardless</title>
</svelte:head>

<div class="flex flex-col items-center justify-center min-h-screen p-6 bg-background">
	<div class="w-full max-w-md space-y-8">
		<!-- Logo/Title -->
		<div class="text-center">
			<h1 class="text-6xl font-bold mb-2">🎲</h1>
			<h1 class="text-5xl font-bold text-primary mb-2">Cardless</h1>
			<p class="text-foreground/80 text-lg">You've been invited!</p>
		</div>

		{#if validating}
			<div class="card space-y-4">
				<div class="flex items-center justify-center py-8">
					<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
				</div>
				<p class="text-center text-muted-foreground">Validating room...</p>
			</div>
		{:else if roomValid}
			<!-- Join form -->
			<div class="card space-y-4">
				<div class="text-center">
					<p class="text-sm text-muted-foreground mb-2">Room Code</p>
					<div class="text-3xl font-mono font-bold tracking-wider text-primary">
						{roomCode}
					</div>
				</div>

				<div class="relative">
					<div class="absolute inset-0 flex items-center">
						<div class="w-full border-t border-border"></div>
					</div>
					<div class="relative flex justify-center text-sm">
						<span class="px-4 bg-card text-muted-foreground">enter your name</span>
					</div>
				</div>

				<form onsubmit={(e) => { e.preventDefault(); handleJoin(); }} class="space-y-3">
					<input
						type="text"
						bind:value={displayName}
						placeholder="Your name"
						class="input"
						maxlength="20"
						autocomplete="off"
						disabled={loading}
					/>

					{#if error}
						<div class="p-3 bg-destructive/10 border border-destructive/20 rounded-lg">
							<p class="text-sm text-destructive">{error}</p>
						</div>
					{/if}

					<button
						type="submit"
						class="btn btn-primary w-full"
						disabled={loading || !displayName.trim()}
					>
						{loading ? 'Joining...' : 'Join Game'}
					</button>
				</form>

				<div class="text-center">
					<a href="/" class="text-sm text-muted-foreground hover:text-foreground transition-colors">
						Back to home
					</a>
				</div>
			</div>
		{/if}

		<!-- Info -->
		<div class="text-center text-muted-foreground text-sm">
			<p>No signup required • Play with friends in person</p>
		</div>
	</div>
</div>
