<script lang="ts">
	import { goto } from '$app/navigation';
	import { api, type CreateRoomRequest } from '$lib/api/client';
	import { session } from '$lib/stores/session';

	let displayName = '';
	let selectedGame = 'werewolf';
	let maxPlayers = 10;
	let loading = false;
	let error = '';

	const games = [
		{ id: 'werewolf', name: 'One Night Werewolf', players: '3-10' },
		// Future: { id: 'avalon', name: 'Avalon', players: '5-10' },
		// Future: { id: 'bohnanza', name: 'Bohnanza', players: '3-7' }
	];

	async function handleCreate() {
		if (!displayName.trim()) {
			error = 'Please enter your name';
			return;
		}

		loading = true;
		error = '';

		try {
			const request: CreateRoomRequest = {
				gameType: selectedGame,
				displayName: displayName.trim(),
				maxPlayers
			};

			const response = await api.createRoom(request);

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
			error = err.message || 'Failed to create room';
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Create Room - Roundtable</title>
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
				<h1 class="text-2xl font-bold mb-2">Host a Game</h1>
				<p class="text-gray-600">Choose a game and enter your name</p>
			</div>

			<form on:submit|preventDefault={handleCreate} class="space-y-4">
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

				<!-- Game selection -->
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-2">
						Choose Game
					</label>
					<div class="space-y-2">
						{#each games as game}
							<label class="flex items-center p-4 border-2 rounded-xl cursor-pointer transition-colors
								{selectedGame === game.id ? 'border-primary-500 bg-primary-50' : 'border-gray-300 hover:border-gray-400'}">
								<input
									type="radio"
									name="game"
									value={game.id}
									bind:group={selectedGame}
									class="mr-3"
									disabled={loading}
								/>
								<div class="flex-1">
									<div class="font-semibold">{game.name}</div>
									<div class="text-sm text-gray-500">{game.players} players</div>
								</div>
							</label>
						{/each}
					</div>
				</div>

				<!-- Max players -->
				<div>
					<label for="maxPlayers" class="block text-sm font-medium text-gray-700 mb-2">
						Max Players: {maxPlayers}
					</label>
					<input
						id="maxPlayers"
						type="range"
						bind:value={maxPlayers}
						min="3"
						max="15"
						class="w-full"
						disabled={loading}
					/>
				</div>

				<!-- Error message -->
				{#if error}
					<div class="p-4 bg-red-50 border border-red-200 rounded-xl text-red-700 text-sm">
						{error}
					</div>
				{/if}

				<!-- Submit button -->
				<button type="submit" class="btn btn-primary w-full" disabled={loading}>
					{loading ? 'Creating...' : 'Create Room'}
				</button>
			</form>
		</div>
	</div>
</div>
