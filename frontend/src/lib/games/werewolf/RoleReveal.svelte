<script lang="ts">
	import { Button, CardBack } from '$lib/components/ui';
	import { getRoleInfo } from './roleConfig';

	let { role, acknowledged = false, acknowledgementsCount = 0, totalPlayers = 0, onAcknowledge } = $props<{
		role: string;
		acknowledged?: boolean;
		acknowledgementsCount?: number;
		totalPlayers?: number;
		onAcknowledge: () => void;
	}>();

	let cardFlipped = $state(false);
	let hasClickedReady = $state(false);

	let info = $derived(getRoleInfo(role));

	function handleShowRole() {
		if (!cardFlipped) {
			cardFlipped = true;
		}
	}

	function handleReady() {
		if (!acknowledged) {
			if (!hasClickedReady) {
				// First click of Ready - just mark as ready, show Look Again button
				hasClickedReady = true;
			} else {
				// Second click of Ready - confirm and acknowledge
				onAcknowledge();
			}
		}
	}

	function handleUnready() {
		if (acknowledged) {
			// Can't unready after acknowledging
			return;
		}
		// Reset back to showing card back
		cardFlipped = false;
		hasClickedReady = false;
	}
</script>

<div class="w-full space-y-6">
		<!-- Card Container -->
		<div class="perspective-1000">
			<div class="card-container {cardFlipped ? 'flipped' : ''}">
				<!-- Card Back -->
				<div class="card-face card-back">
					<div class="game-card border-2 border-primary">
						<CardBack width={280} height={392} variant="simplified" />
					</div>
				</div>

			<!-- Card Front -->
			<div class="card-face card-front {info.color} rounded-2xl">
				<div class="game-card p-8 flex flex-col items-center justify-between text-white">
						<div class="flex-1 flex flex-col items-center justify-center space-y-6 w-full">
							<!-- Role emoji -->
							<div class="text-8xl">
								{info.emoji}
							</div>

							<!-- Role name -->
							<div class="text-center">
								<h2 class="text-4xl font-bold mb-2 capitalize">
									{role.replace('_', ' ')}
								</h2>
								<p class="text-xl text-white/90">
									{info.team}
								</p>
							</div>

							<!-- Description -->
							<p class="text-lg text-white/90 text-center">
								{info.description}
							</p>
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- Action Buttons -->
		{#if !acknowledged}
			<div class="space-y-3">
				{#if !cardFlipped}
					<Button
						onclick={handleShowRole}
						class="w-full h-14 bg-primary hover:bg-primary/90 text-primary-foreground font-bold text-lg"
					>
						👁️ Show Role
					</Button>
					<p class="text-sm text-white/75 text-center">
						Tap to peek at your role card
					</p>
				{:else if !hasClickedReady}
					<!-- First view - only show Ready button -->
					<Button
						onclick={handleReady}
						class="w-full h-14 bg-green-600 hover:bg-green-700 text-white font-bold text-lg"
					>
						✓ Ready
					</Button>
					<p class="text-sm text-white/75 text-center">
						Ready when you've memorized your role
					</p>
				{:else}
					<!-- After clicking Ready once - show both buttons -->
					<div class="flex gap-3">
						<Button
							onclick={handleUnready}
							variant="outline"
							class="flex-1 h-14 bg-card hover:bg-muted text-foreground font-bold text-lg border-2"
						>
							↺ Look Again
						</Button>
						<Button
							onclick={handleReady}
							class="flex-1 h-14 bg-green-600 hover:bg-green-700 text-white font-bold text-lg"
						>
							✓ Confirm
						</Button>
					</div>
					<p class="text-sm text-white/75 text-center">
						Click Confirm to continue, or Look Again to review
					</p>
				{/if}
			</div>
		{:else}
			<!-- Waiting for others -->
			<div class="bg-white/10 backdrop-blur rounded-lg p-6 text-center text-white">
				<p class="font-semibold mb-3 text-lg">Waiting for other players...</p>
				<p class="text-4xl font-bold mb-2">
					{acknowledgementsCount} / {totalPlayers}
				</p>
				<p class="text-sm text-white/75">players ready</p>
			</div>
		{/if}
</div>

<style>
	.perspective-1000 {
		perspective: 1000px;
		display: flex;
		justify-content: center;
	}

	.card-container {
		position: relative;
		width: 280px;
		height: 392px;
		transition: transform 0.8s cubic-bezier(0.4, 0.0, 0.2, 1);
		transform-style: preserve-3d;
	}

	.card-container.flipped {
		transform: rotateY(180deg);
	}

	.card-face {
		position: absolute;
		width: 100%;
		height: 100%;
		backface-visibility: hidden;
		-webkit-backface-visibility: hidden;
	}

	.card-back {
		z-index: 2;
		transform: rotateY(0deg);
	}

	.card-front {
		transform: rotateY(180deg);
	}

	.game-card {
		width: 280px;
		height: 392px;
		border-radius: 12px;
		box-shadow: 0 8px 16px rgba(0, 0, 0, 0.3);
		overflow: hidden; /* Ensure rounded corners apply to content */
	}
</style>
