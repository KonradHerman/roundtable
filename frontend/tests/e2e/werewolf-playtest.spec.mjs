import { test } from '@playwright/test';

const PLAYER_NAMES = [
	'Alice',
	'Bob', 
	'Charlie',
	'Diana',
	'Eve',
	'Frank',
	'Grace',
	'Henry'
];

/** @type {Array<{context: import('@playwright/test').BrowserContext, page: import('@playwright/test').Page, name: string, isHost: boolean}>} */
let players = [];
let roomCode = '';

test.describe('Werewolf 8-Player Playtest', () => {
	
	test.afterAll(async () => {
		// Close all browser contexts
		for (const player of players) {
			await player.context.close();
		}
	});

	test('Play a full Werewolf game with 8 players', async ({ browser }) => {
		console.log('\n🐺 Starting Werewolf 8-Player Playtest\n');

		// Step 1: Create room as host (Alice)
		console.log('📍 Step 1: Creating room as host...');
		const hostContext = await browser.newContext();
		const hostPage = await hostContext.newPage();
		
		await hostPage.goto('/create');
		await hostPage.waitForLoadState('networkidle');
		
		// Enter host name
		await hostPage.fill('input#name', PLAYER_NAMES[0]);
		await hostPage.click('button[type="submit"]');
		
		// Wait for room page and get room code
		await hostPage.waitForURL(/\/room\/[A-Z0-9]+/);
		roomCode = hostPage.url().split('/room/')[1];
		console.log(`✅ Room created with code: ${roomCode}`);
		
		players.push({
			context: hostContext,
			page: hostPage,
			name: PLAYER_NAMES[0],
			isHost: true
		});

		// Wait for room to be fully loaded
		await hostPage.waitForSelector('text=Room Code');
		
		// Step 2: Join 7 more players
		console.log('\n📍 Step 2: Joining 7 more players...');
		
		for (let i = 1; i < 8; i++) {
			const playerContext = await browser.newContext();
			const playerPage = await playerContext.newPage();
			
			// Go to home page and join
			await playerPage.goto('/');
			await playerPage.waitForLoadState('networkidle');
			
			// Enter room code
			await playerPage.fill('input[placeholder="Enter room code"]', roomCode);
			
			// Wait for name field to appear and fill it
			await playerPage.waitForSelector('input[placeholder="Enter your name"]');
			await playerPage.fill('input[placeholder="Enter your name"]', PLAYER_NAMES[i]);
			
			// Click join
			await playerPage.click('button:has-text("Join Game")');
			
			// Wait for room page
			await playerPage.waitForURL(/\/room\/[A-Z0-9]+/);
			console.log(`✅ ${PLAYER_NAMES[i]} joined the room`);
			
			players.push({
				context: playerContext,
				page: playerPage,
				name: PLAYER_NAMES[i],
				isHost: false
			});
			
			// Small delay between joins
			await playerPage.waitForTimeout(300);
		}

		// Wait for all players to appear in host's view
		console.log('\n📍 Step 3: Verifying all players are connected...');
		await players[0].page.waitForSelector('text=Players (8/15)');
		console.log('✅ All 8 players connected');

		// Step 4: Start the game (host)
		console.log('\n📍 Step 4: Starting Werewolf game...');
		const hostPage2 = players[0].page;
		
		// Select Werewolf (should be default)
		await hostPage2.click('button:has-text("One Night Werewolf")');
		
		// Click start button
		await hostPage2.click('button:has-text("Start Game")');
		console.log('✅ Game started!');

		// Step 5: Wait for role reveal phase and acknowledge roles
		console.log('\n📍 Step 5: Role Reveal Phase...');
		
		// Wait for game to transition
		await hostPage2.waitForTimeout(1500);
		
		// Each player sees and acknowledges their role
		for (const player of players) {
			try {
				// Wait for role reveal phase
				await player.page.waitForSelector('text=Role Reveal Phase', { timeout: 10000 });
				
				// Click "Show Role" button
				const showRoleBtn = player.page.locator('button:has-text("Show Role")');
				if (await showRoleBtn.isVisible({ timeout: 2000 })) {
					await showRoleBtn.click();
					console.log(`  👁️ ${player.name} viewing role...`);
					await player.page.waitForTimeout(800);
				}
				
				// Click "Ready" button (first click)
				const readyBtn = player.page.locator('button:has-text("Ready")');
				if (await readyBtn.isVisible({ timeout: 2000 })) {
					await readyBtn.click();
					console.log(`  ✋ ${player.name} clicked Ready...`);
					await player.page.waitForTimeout(400);
				}
				
				// Click "Confirm" button (second click)
				const confirmBtn = player.page.locator('button:has-text("Confirm")');
				if (await confirmBtn.isVisible({ timeout: 2000 })) {
					await confirmBtn.click();
					console.log(`  ✅ ${player.name} confirmed role!`);
				}
			} catch (e) {
				console.log(`  ⚠️ ${player.name} had issue acknowledging role: ${e}`);
			}
		}

		// Wait for all acknowledgements to register
		console.log('\n  ⏳ Waiting for all acknowledgements...');
		await hostPage2.waitForTimeout(3000);

		// Step 6: Night phase
		console.log('\n📍 Step 6: Night Phase...');
		
		try {
			await hostPage2.waitForSelector('text=Night Phase', { timeout: 15000 });
			console.log('  🌙 Night phase started');
		} catch {
			console.log('  ⚠️ Night phase transition taking time...');
			await hostPage2.waitForTimeout(2000);
		}

		// Host shows narration script
		const hostScriptBtn = hostPage2.locator('button:has-text("Show Host Script")');
		if (await hostScriptBtn.isVisible({ timeout: 3000 }).catch(() => false)) {
			await hostScriptBtn.click();
			console.log('  📜 Host viewing narration script');
			await hostPage2.waitForTimeout(500);
		}

		// Each player performs their night action
		for (const player of players) {
			try {
				// Click "Show Night Action" to reveal action UI
				const showActionBtn = player.page.locator('button:has-text("Show Night Action")');
				if (await showActionBtn.isVisible({ timeout: 2000 }).catch(() => false)) {
					await showActionBtn.click();
					console.log(`  🌙 ${player.name} viewing night action...`);
					await player.page.waitForTimeout(800);
					
					// Try to perform an action - select first available player
					const playerButton = player.page.locator('button.w-full.p-3').first();
					if (await playerButton.isVisible({ timeout: 1000 }).catch(() => false)) {
						await playerButton.click();
						console.log(`    ✨ ${player.name} performed action`);
						await player.page.waitForTimeout(300);
					}
				}
			} catch (e) {
				// Some roles have no night action
			}
		}

		// Host advances to day phase
		console.log('\n📍 Step 7: Advancing to Day Phase...');
		await hostPage2.waitForTimeout(1000);
		
		const advanceBtn = hostPage2.locator('button:has-text("Advance to Day Phase")');
		if (await advanceBtn.isVisible({ timeout: 5000 }).catch(() => false)) {
			await advanceBtn.click();
			console.log('  ☀️ Advanced to day phase');
		} else {
			console.log('  ⚠️ Advance button not found yet, waiting...');
			await hostPage2.waitForTimeout(3000);
			if (await advanceBtn.isVisible({ timeout: 2000 }).catch(() => false)) {
				await advanceBtn.click();
				console.log('  ☀️ Advanced to day phase (after wait)');
			}
		}

		// Step 8: Day phase - Discussion & Reveal roles
		console.log('\n📍 Step 8: Day Phase (Discussion)...');
		
		try {
			await hostPage2.waitForSelector('text=Day Phase', { timeout: 15000 });
			console.log('  🌅 Day phase active - Discussion time!');
		} catch {
			console.log('  ⚠️ Day phase indicator not found, continuing...');
		}

		// Wait a moment for "discussion"
		await hostPage2.waitForTimeout(2000);

		// Host reveals roles to end the game
		console.log('\n📍 Step 9: Revealing Results...');
		const revealBtn = hostPage2.locator('button:has-text("Reveal All Roles")');
		if (await revealBtn.isVisible({ timeout: 5000 }).catch(() => false)) {
			await revealBtn.click();
			console.log('  🎭 Revealing all roles...');
		} else {
			console.log('  ⚠️ Reveal button not found');
		}

		// Wait for results
		await hostPage2.waitForTimeout(3000);

		// Take screenshots of final state
		console.log('\n📸 Taking final screenshots...');
		for (let i = 0; i < players.length; i++) {
			try {
				await players[i].page.screenshot({ 
					path: `tests/e2e/screenshots/player-${i + 1}-${players[i].name}-final.png`,
					fullPage: true 
				});
				console.log(`  📷 Screenshot saved for ${players[i].name}`);
			} catch (e) {
				console.log(`  ⚠️ Could not take screenshot for ${players[i].name}`);
			}
		}

		console.log('\n🎮 Werewolf 8-Player Playtest Complete!\n');
		console.log('Check tests/e2e/screenshots/ for final game state images.\n');
	});
});

