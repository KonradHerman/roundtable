<script lang="ts">
	import { onMount } from 'svelte';
	import QRCodeStyling from 'qr-code-styling';
	import { browser } from '$app/environment';

	export let roomCode: string;
	export let size = 280;

	let qrContainer: HTMLDivElement;
	let qrCode: QRCodeStyling | null = null;
	let showCopied = false;

	// Get the full invite URL
	$: inviteUrl = browser ? `${window.location.origin}/join/${roomCode}` : '';

	onMount(() => {
		if (!browser) return;

		// Create QR code with soft aesthetic styling
		qrCode = new QRCodeStyling({
			width: size,
			height: size,
			data: inviteUrl,
			margin: 10,
			qrOptions: {
				typeNumber: 0,
				mode: 'Byte',
				errorCorrectionLevel: 'M'
			},
			imageOptions: {
				hideBackgroundDots: true,
				imageSize: 0.4,
				margin: 4
			},
			dotsOptions: {
				type: 'rounded',
				color: '#d79921', // primary golden yellow
				gradient: {
					type: 'linear',
					rotation: 45,
					colorStops: [
						{ offset: 0, color: '#d79921' },
						{ offset: 1, color: '#fabd2f' }
					]
				}
			},
			backgroundOptions: {
				color: '#3c3836' // card background
			},
			cornersSquareOptions: {
				type: 'extra-rounded',
				color: '#d79921',
				gradient: {
					type: 'linear',
					rotation: 45,
					colorStops: [
						{ offset: 0, color: '#d79921' },
						{ offset: 1, color: '#fabd2f' }
					]
				}
			},
			cornersDotOptions: {
				type: 'dot',
				color: '#fabd2f'
			}
		});

		// Append to container
		qrCode.append(qrContainer);
	});

	async function copyInviteLink() {
		if (!browser) return;

		try {
			await navigator.clipboard.writeText(inviteUrl);
			showCopied = true;
			setTimeout(() => {
				showCopied = false;
			}, 2000);
		} catch (err) {
			console.error('Failed to copy:', err);
		}
	}

	async function downloadQR() {
		if (!qrCode) return;
		qrCode.download({
			name: `cardless-${roomCode}`,
			extension: 'png'
		});
	}
</script>

<div class="space-y-4">
	<!-- QR Code Display -->
	<div class="relative">
		<div
			class="bg-card rounded-2xl shadow-xl p-6 flex items-center justify-center overflow-hidden border-2 border-border/50"
		>
			<div bind:this={qrContainer} class="rounded-xl overflow-hidden"></div>
		</div>

		<!-- Decorative corners for extra soft aesthetic -->
		<div class="absolute -top-1 -left-1 w-6 h-6 border-t-2 border-l-2 border-primary/20 rounded-tl-lg"></div>
		<div class="absolute -top-1 -right-1 w-6 h-6 border-t-2 border-r-2 border-primary/20 rounded-tr-lg"></div>
		<div class="absolute -bottom-1 -left-1 w-6 h-6 border-b-2 border-l-2 border-primary/20 rounded-bl-lg"></div>
		<div class="absolute -bottom-1 -right-1 w-6 h-6 border-b-2 border-r-2 border-primary/20 rounded-br-lg"></div>
	</div>

	<!-- Info text -->
	<div class="text-center space-y-2">
		<p class="text-sm text-muted-foreground">Scan to join the game</p>
		<p class="text-xs text-muted-foreground/60 font-mono break-all px-4">
			{inviteUrl}
		</p>
	</div>

	<!-- Actions -->
	<div class="flex gap-2">
		<button
			type="button"
			on:click={copyInviteLink}
			class="btn btn-secondary flex-1 text-sm py-3 relative"
		>
			{#if showCopied}
				<span class="flex items-center justify-center gap-2">
					<span>✓</span>
					<span>Copied!</span>
				</span>
			{:else}
				<span class="flex items-center justify-center gap-2">
					<span>🔗</span>
					<span>Copy Link</span>
				</span>
			{/if}
		</button>

		<button
			type="button"
			on:click={downloadQR}
			class="btn btn-secondary flex-1 text-sm py-3"
		>
			<span class="flex items-center justify-center gap-2">
				<span>⬇️</span>
				<span>Download</span>
			</span>
		</button>
	</div>
</div>

<style>
	/* Ensure QR code canvas inherits border radius */
	:global(.qr-code-styling canvas) {
		border-radius: 0.75rem;
	}
</style>
