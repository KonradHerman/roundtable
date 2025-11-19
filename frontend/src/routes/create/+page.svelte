<script lang="ts">
	import { goto } from '$app/navigation';
	import { api, type CreateRoomRequest } from '$lib/api/client';
	import { session } from '$lib/stores/session.svelte';
	import { Card, Button } from '$lib/components/ui';
	import { ArrowLeft } from 'lucide-svelte';

	let displayName = '';
	let selectedGame = 'werewolf';
	let loading = false;
	let error = '';

	const games = [
		{ id: 'werewolf', name: 'One Night Werewolf', players: '3-10', emoji: '🐺' },
		// Future: { id: 'avalon', name: 'Avalon', players: '5-10', emoji: '⚔️' },
		// Future: { id: 'bohnanza', name: 'Bohnanza', players: '3-7', emoji: '🌱' }
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
				maxPlayers: 15
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
	<title>Create Room - Cardless</title>
</svelte:head>

<div class="min-h-screen p-6 bg-background">
	<div class="max-w-md mx-auto py-8">
		<!-- Back button -->
		<Button
			variant="ghost"
			class="mb-6"
			on:click={() => goto('/')}
		>
			<ArrowLeft class="w-4 h-4 mr-2" />
			Back
		</Button>

		<Card class="p-6 space-y-6">
			<div>
				<h1 class="text-2xl font-bold mb-2">Host a Game</h1>
				<p class="text-muted-foreground">Choose a game and enter your name</p>
			</div>

			<form on:submit|preventDefault={handleCreate} class="space-y-5">
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

				<!-- Game selection -->
				<div class="space-y-2">
					<div class="block text-sm font-medium mb-2">
						Choose Game
					</div>
					<div class="space-y-2">
						{#each games as game}
							<label class="flex items-center gap-3 p-4 border-2 rounded-lg cursor-pointer transition-all hover:border-primary/50 {selectedGame === game.id ? 'border-primary bg-primary/5' : 'border-input'}">
								<input
									type="radio"
									name="game"
									value={game.id}
									bind:group={selectedGame}
									class="sr-only"
									disabled={loading}
								/>
								<span class="text-3xl">{game.emoji}</span>
								<div class="flex-1">
									<div class="font-semibold">{game.name}</div>
									<div class="text-sm text-muted-foreground">{game.players} players</div>
								</div>
								{#if selectedGame === game.id}
									<div class="w-5 h-5 rounded-full bg-primary flex items-center justify-center">
										<div class="w-2 h-2 rounded-full bg-white"></div>
									</div>
								{/if}
							</label>
						{/each}
					</div>
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
					disabled={loading || !displayName.trim()}
				>
					{loading ? 'Creating...' : 'Create Room'}
				</Button>
			</form>
		</Card>
	</div>
</div>
