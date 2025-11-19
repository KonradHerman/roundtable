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
				type: 'dots',
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
</script>

<div class="space-y-4">
	<!-- QR Code Display -->
	<div
		class="bg-card rounded-2xl shadow-xl p-6 flex items-center justify-center border-2 border-border/50"
	>
		<div bind:this={qrContainer} class="rounded-xl max-w-full"></div>
	</div>

	<!-- Info text with click to copy -->
	<div class="text-center space-y-2">
		<p class="text-sm text-muted-foreground">Scan to join the game</p>
		<button
			type="button"
			on:click={copyInviteLink}
			class="group relative text-xs font-mono break-all px-4 py-2 rounded-lg hover:bg-muted/30 transition-colors cursor-pointer"
		>
			<span class={showCopied ? 'text-primary' : 'text-muted-foreground/60'}>
				{inviteUrl}
			</span>
			{#if showCopied}
				<span class="ml-2 text-primary">✓ Copied</span>
			{/if}
		</button>
	</div>
</div>

<style>
	/* Ensure QR code canvas inherits border radius and fits properly */
	:global(.qr-code-styling canvas) {
		border-radius: 0.75rem;
		max-width: 100%;
		height: auto;
		display: block;
	}
</style>
