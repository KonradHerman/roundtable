<script lang="ts">
	import { goto } from '$app/navigation';

	let joinCode = '';

	function handleJoin() {
		if (joinCode.trim()) {
			goto(`/join?code=${joinCode.toUpperCase()}`);
		}
	}
</script>

<svelte:head>
	<title>Roundtable - Party Games</title>
</svelte:head>

<div class="flex flex-col items-center justify-center min-h-screen p-6 bg-gradient-to-br from-primary-500 to-primary-700">
	<div class="w-full max-w-md space-y-8">
		<!-- Logo/Title -->
		<div class="text-center">
			<h1 class="text-6xl font-bold text-white mb-2">🎲</h1>
			<h1 class="text-4xl font-bold text-white mb-2">Roundtable</h1>
			<p class="text-primary-100 text-lg">Party games for friends</p>
		</div>

		<!-- Main actions -->
		<div class="card space-y-4">
			<a href="/create" class="btn btn-primary w-full block text-center">
				Host a Game
			</a>

			<div class="relative">
				<div class="absolute inset-0 flex items-center">
					<div class="w-full border-t border-gray-300"></div>
				</div>
				<div class="relative flex justify-center text-sm">
					<span class="px-4 bg-white text-gray-500">or join with code</span>
				</div>
			</div>

			<form on:submit|preventDefault={handleJoin} class="space-y-3">
				<input
					type="text"
					bind:value={joinCode}
					placeholder="Enter room code"
					class="input text-center text-2xl font-mono tracking-wider uppercase"
					maxlength="6"
					autocomplete="off"
					autocapitalize="characters"
				/>
				<button type="submit" class="btn btn-secondary w-full" disabled={!joinCode.trim()}>
					Join Game
				</button>
			</form>
		</div>

		<!-- Info -->
		<div class="text-center text-primary-100 text-sm">
			<p>No signup required • Play with friends in person</p>
		</div>
	</div>
</div>
