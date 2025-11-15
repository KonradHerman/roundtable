<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { api, type JoinRoomRequest } from '$lib/api/client';
	import { session } from '$lib/stores/session';
	import { Card, Button } from '$lib/components/ui';
	import { ArrowLeft } from 'lucide-svelte';

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

<div class="min-h-screen p-6 bg-gradient-to-br from-primary to-primary/80">
	<div class="max-w-md mx-auto py-8">
		<!-- Back button -->
		<Button
			variant="ghost"
			class="mb-6 text-white hover:bg-white/20"
			on:click={() => goto('/')}
		>
			<ArrowLeft class="w-4 h-4 mr-2" />
			Back
		</Button>

		<Card class="p-6 space-y-6">
			<div>
				<h1 class="text-2xl font-bold mb-2">Join Game</h1>
				<p class="text-muted-foreground">Enter the room code and your name</p>
			</div>

			<form on:submit|preventDefault={handleJoin} class="space-y-5">
				<!-- Room code -->
				<div class="space-y-2">
					<label for="code" class="block text-sm font-medium">
						Room Code
					</label>
					<input
						id="code"
						type="text"
						bind:value={roomCode}
						placeholder="ABC123"
						class="w-full px-4 py-3 text-center text-2xl font-mono font-bold tracking-wider uppercase rounded-lg border-2 border-input bg-background focus:border-primary focus:outline-none transition-colors"
						maxlength="6"
						disabled={loading}
						autocomplete="off"
						autocapitalize="characters"
						style="min-height: 56px;"
					/>
				</div>

				<!-- Your name -->
				<div class="space-y-2">
					<label for="name" class="block text-sm font-medium">
						Your Name
					</label>
					<input
						id="name"
						type="text"
						bind:value={displayName}
						placeholder="Enter your name"
						class="w-full px-4 py-3 text-base rounded-lg border-2 border-input bg-background focus:border-primary focus:outline-none transition-colors"
						maxlength="20"
						disabled={loading}
						autocomplete="off"
						style="min-height: 48px;"
					/>
				</div>

				<!-- Error message -->
				{#if error}
					<Card class="p-4 bg-destructive/10 border-destructive/20">
						<p class="text-sm text-destructive">{error}</p>
					</Card>
				{/if}

				<!-- Submit button -->
				<Button
					type="submit"
					class="w-full h-12 text-base"
					disabled={loading || !roomCode.trim() || !displayName.trim()}
				>
					{loading ? 'Joining...' : 'Join Room'}
				</Button>
			</form>
		</Card>
	</div>
</div>
