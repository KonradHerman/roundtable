<script lang="ts">
	export let variant: 'ornamental' | 'simplified' = 'ornamental';
	export let width: number = 280;
	export let height: number = 392;
</script>

<svg
	{width}
	{height}
	viewBox="0 0 280 392"
	xmlns="http://www.w3.org/2000/svg"
	class="card-back-svg"
>
	<!-- Card background -->
	<rect width="280" height="392" rx="12" fill="#3c3836" />
	
	<!-- Border -->
	<rect x="4" y="4" width="272" height="384" rx="10" fill="none" stroke="#d79921" stroke-width="3" />
	
	{#if variant === 'ornamental'}
		<!-- Ornamental pattern inspired by bicycle rider back -->
		<g transform="translate(140, 196)">
			<!-- Center circle -->
			<circle cx="0" cy="0" r="8" fill="#d79921" />
			
			<!-- Four main spiral arms (rotated 90 degrees each) -->
			{#each [0, 90, 180, 270] as rotation}
				<g transform="rotate({rotation})">
					<!-- Large spiral -->
					<path
						d="M 0,0 Q 15,-20 25,-25 Q 35,-30 40,-25 Q 45,-20 42,-10 Q 39,0 30,5 Q 20,10 10,8"
						fill="none"
						stroke="#d79921"
						stroke-width="6"
						stroke-linecap="round"
					/>
					
					<!-- Inner decorative circle -->
					<circle cx="28" cy="-15" r="10" fill="#d79921" />
					<circle cx="28" cy="-15" r="6" fill="#3c3836" />
					
					<!-- Outer decorative elements -->
					<path
						d="M 45,-30 Q 50,-35 55,-35 Q 60,-35 62,-30 Q 60,-25 55,-25 Q 50,-25 45,-28"
						fill="#d79921"
					/>
					
					<!-- Club-like ornament at end -->
					<g transform="translate(65, -40)">
						<circle cx="0" cy="0" r="8" fill="#d79921" />
						<circle cx="-7" cy="-5" r="6" fill="#d79921" />
						<circle cx="7" cy="-5" r="6" fill="#d79921" />
						<path d="M -5,5 L -2,12 L 2,12 L 5,5 Z" fill="#d79921" />
					</g>
				</g>
			{/each}
			
			<!-- Diagonal decorative elements -->
			{#each [45, 135, 225, 315] as rotation}
				<g transform="rotate({rotation})">
					<!-- Small spiral accents -->
					<path
						d="M 50,0 Q 55,-8 60,-8 Q 65,-8 67,-3 Q 65,2 60,2 Q 55,2 52,-2"
						fill="#d79921"
					/>
					
					<!-- Additional club ornaments -->
					<g transform="translate(75, -10)">
						<circle cx="0" cy="0" r="5" fill="#d79921" />
						<circle cx="-4" cy="-3" r="4" fill="#d79921" />
						<circle cx="4" cy="-3" r="4" fill="#d79921" />
					</g>
				</g>
			{/each}
			
			<!-- Corner decorations -->
			{#each [0, 90, 180, 270] as rotation}
				<g transform="rotate({rotation})">
					<path
						d="M 85,-85 Q 88,-90 92,-90 Q 96,-90 98,-86 Q 96,-82 92,-82 Q 88,-82 85,-85"
						fill="#d79921"
					/>
				</g>
			{/each}
		</g>
	{:else}
		<!-- Simplified geometric pattern -->
		<g transform="translate(140, 196)">
			<!-- Center diamond -->
			<path d="M 0,-15 L 15,0 L 0,15 L -15,0 Z" fill="#d79921" />
			
			<!-- Four corner ornaments -->
			{#each [0, 90, 180, 270] as rotation}
				<g transform="rotate({rotation})">
					<!-- Decorative line -->
					<line x1="0" y1="25" x2="0" y2="60" stroke="#d79921" stroke-width="4" />
					
					<!-- End circle -->
					<circle cx="0" cy="70" r="8" fill="#d79921" />
					
					<!-- Side decorations -->
					<circle cx="-15" cy="45" r="5" fill="#d79921" />
					<circle cx="15" cy="45" r="5" fill="#d79921" />
				</g>
			{/each}
			
			<!-- Diagonal lines -->
			{#each [45, 135, 225, 315] as rotation}
				<g transform="rotate({rotation})">
					<line x1="20" y1="0" x2="50" y2="0" stroke="#d79921" stroke-width="3" />
					<circle cx="55" cy="0" r="6" fill="#d79921" />
				</g>
			{/each}
		</g>
	{/if}
	
	<!-- Inner border decoration -->
	<rect x="12" y="12" width="256" height="368" rx="8" fill="none" stroke="#d79921" stroke-width="1" opacity="0.5" />
</svg>

<style>
	.card-back-svg {
		display: block;
	}
</style>

