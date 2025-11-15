<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { api, type JoinRoomRequest } from '$lib/api/client';
	import { session } from '$lib/stores/session';

	let roomCode = '';
	let displayName = '';
	let loading = false;
	let error = '';

	onMount(() => {
		// Get room code from URL query params
		const code = $page.url.searchParams.get('code');
		if (code) {
			roomCode = code.toUpperCase();
		}
	});

	async function handleJoin() {
		if (!roomCode.trim()) {
			error = 'Please enter a room code';
			return;
		}

		if (!displayName.trim()) {
			error = 'Please enter your name';
			return;
		}

		loading = true;
		error = '';

		try {
			const request: JoinRoomRequest = {
				displayName: displayName.trim()
			};

			const response = await api.joinRoom(roomCode.toUpperCase(), request);

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
	<title>Join Room - Roundtable</title>
</svelte:head>

<div class="min-h-screen p-6 bg-gradient-to-br from-primary-500 to-primary-700">
	<div class="max-w-md mx-auto py-8">
		<!-- Back button -->
		<div class="mb-6">
			<a href="/" class="text-white hover:text-primary-100 flex items-center gap-2">
				<span>←</span> Back
			</a>
		</div>

		<div class="card space-y-6">
			<div>
				<h1 class="text-2xl font-bold mb-2">Join Game</h1>
				<p class="text-gray-600">Enter the room code and your name</p>
			</div>

			<form on:submit|preventDefault={handleJoin} class="space-y-4">
				<!-- Room code -->
				<div>
					<label for="code" class="block text-sm font-medium text-gray-700 mb-2">
						Room Code
					</label>
					<input
						id="code"
						type="text"
						bind:value={roomCode}
						placeholder="ABC123"
						class="input text-center text-2xl font-mono tracking-wider uppercase"
						maxlength="6"
						disabled={loading}
						autocomplete="off"
						autocapitalize="characters"
					/>
				</div>

				<!-- Your name -->
				<div>
					<label for="name" class="block text-sm font-medium text-gray-700 mb-2">
						Your Name
					</label>
					<input
						id="name"
						type="text"
						bind:value={displayName}
						placeholder="Enter your name"
						class="input"
						maxlength="20"
						disabled={loading}
						autocomplete="off"
					/>
				</div>

				<!-- Error message -->
				{#if error}
					<div class="p-4 bg-red-50 border border-red-200 rounded-xl text-red-700 text-sm">
						{error}
					</div>
				{/if}

				<!-- Submit button -->
				<button type="submit" class="btn btn-primary w-full" disabled={loading || !roomCode.trim() || !displayName.trim()}>
					{loading ? 'Joining...' : 'Join Room'}
				</button>
			</form>
		</div>
	</div>
</div>
