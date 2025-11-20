<script lang="ts">
	/**
	 * AvalonGame - Main container component that orchestrates all game phases
	 * Handles WebSocket communication and phase-specific component rendering
	 */
	import { gameStore } from '$lib/stores/game.svelte';
	import { session } from '$lib/stores/session.svelte';
	import type { AvalonRole, Team } from './roleConfig';

	import RoleReveal from './RoleReveal.svelte';
	import TeamBuilding from './TeamBuilding.svelte';
	import TeamVoting from './TeamVoting.svelte';
	import QuestExecution from './QuestExecution.svelte';
	import Assassination from './Assassination.svelte';
	import Results from './Results.svelte';

	let { roomCode, roomState, wsStore }: { roomCode: string; roomState: any; wsStore: any } =
		$props();

	// Player-specific state
	let myRole = $state<AvalonRole | null>(null);
	let myTeam = $state<Team | null>(null);
	let myKnowledge = $state<string[]>([]);
	let acknowledged = $state<boolean>(false);
	let hasVoted = $state<boolean>(false);
	let hasPlayedCard = $state<boolean>(false);

	// Game state
	let currentPhase = $state<string>('setup');
	let currentQuest = $state<number>(1);
	let currentLeaderId = $state<string>('');
	let proposedTeam = $state<string[]>([]);
	let questResults = $state<any[]>([]);
	let rejectionCount = $state<number>(0);

	// Counts
	let acknowledgementsCount = $state<number>(0);
	let totalPlayers = $state<number>(0);
	let votesSubmitted = $state<number>(0);
	let cardsSubmitted = $state<number>(0);
	let requiredTeamSize = $state<number>(0);

	// Results
	let winningTeam = $state<Team | null>(null);
	let winReason = $state<string>('');
	let allRoles = $state<Record<string, AvalonRole>>({});
	let allTeams = $state<Record<string, Team>>({});

	// Reactive effect to process game events
	$effect(() => {
		gameStore.events.forEach((event) => {
			switch (event.type) {
				case 'role_assigned':
					myRole = event.payload.role;
					myTeam = event.payload.team;
					break;

				case 'role_knowledge':
					myKnowledge = event.payload.known_players || [];
					break;

				case 'phase_changed':
					currentPhase = event.payload.phase;
					if (event.payload.quest_number) {
						currentQuest = event.payload.quest_number;
					}
					if (event.payload.team_size) {
						requiredTeamSize = event.payload.team_size;
					}
					// Reset phase-specific flags
					if (currentPhase === 'team_voting') {
						hasVoted = false;
					}
					if (currentPhase === 'quest_execution') {
						hasPlayedCard = false;
					}
					break;

				case 'leader_changed':
					currentLeaderId = event.payload.leader_id;
					break;

				case 'role_acknowledged':
					acknowledgementsCount = event.payload.count;
					totalPlayers = event.payload.total;
					if (event.payload.player_id === session.value?.playerId) {
						acknowledged = true;
					}
					break;

				case 'team_proposed':
					proposedTeam = event.payload.team_members;
					votesSubmitted = 0; // Reset vote count
					break;

				case 'team_vote_cast':
					if (event.payload.voter_id === session.value?.playerId) {
						hasVoted = true;
					}
					break;

				case 'team_vote_recorded':
					// Private confirmation
					votesSubmitted++;
					break;

				case 'team_vote_result':
					rejectionCount = event.payload.rejection_count;
					votesSubmitted = 0; // Reset
					hasVoted = false;
					break;

				case 'quest_card_played':
					if (event.payload.player_id === session.value?.playerId) {
						hasPlayedCard = true;
					}
					break;

				case 'quest_completed':
					questResults = [...questResults, event.payload];
					rejectionCount = 0; // Reset on new quest
					break;

				case 'game_finished':
					winningTeam = event.payload.winning_team;
					winReason = event.payload.win_reason;
					allRoles = event.payload.roles;
					allTeams = event.payload.teams;
					questResults = event.payload.quest_history || questResults;
					break;
			}
		});
	});

	// Derived state
	const players = $derived(roomState?.players || []);
	const playerCount = $derived(players.length);
	const myPlayerId = $derived(session.value?.playerId || '');
	const isLeader = $derived(currentLeaderId === myPlayerId);
	const isOnProposedTeam = $derived(proposedTeam.includes(myPlayerId));

	// Public state helpers
	const publicState = $derived(gameStore.publicState || {});
	const votesSubmittedPublic = $derived(publicState.votes_submitted || votesSubmitted);
	const cardsSubmittedPublic = $derived(publicState.cards_submitted || cardsSubmitted);

	// Action handlers
	function handleAcknowledgeRole() {
		if (!wsStore || acknowledged) return;
		wsStore.sendAction({ type: 'acknowledge_role', payload: {} });
	}

	function handleProposeTeam(teamMemberIds: string[]) {
		if (!wsStore) return;
		wsStore.sendAction({
			type: 'propose_team',
			payload: { team_members: teamMemberIds }
		});
	}

	function handleVoteTeam(vote: 'approve' | 'reject') {
		if (!wsStore || hasVoted) return;
		wsStore.sendAction({
			type: 'vote_team',
			payload: { vote }
		});
	}

	function handlePlayQuestCard(card: 'success' | 'fail') {
		if (!wsStore || hasPlayedCard) return;
		wsStore.sendAction({
			type: 'play_quest_card',
			payload: { card }
		});
	}

	function handleAssassinate(targetId: string) {
		if (!wsStore) return;
		wsStore.sendAction({
			type: 'assassinate',
			payload: { target_id: targetId }
		});
	}
</script>

<div class="avalon-game space-y-6">
	<!-- Phase header -->
	<div class="phase-header bg-[#282828] p-6 rounded-lg border-2 border-[#d79921]">
		<div class="flex items-center justify-between">
			<div>
				<h2 class="text-2xl font-bold text-[#d79921] capitalize mb-1">
					{currentPhase.replace('_', ' ')} Phase
				</h2>
				{#if currentPhase !== 'setup' && currentPhase !== 'finished'}
					<p class="text-[#a89984] text-sm">
						Quest {currentQuest} of 5
					</p>
				{/if}
			</div>
			{#if myRole}
				<div class="my-role text-right">
					<div class="text-sm text-[#a89984]">Your Role</div>
					<div class="text-lg font-bold text-[#ebdbb2]">
						{myRole.replace('_', ' ').toUpperCase()}
					</div>
				</div>
			{/if}
		</div>
	</div>

	<!-- Phase-specific content -->
	{#if currentPhase === 'role_reveal' && myRole && myTeam}
		<RoleReveal
			role={myRole}
			team={myTeam}
			knowledge={myKnowledge}
			{players}
			onAcknowledge={handleAcknowledgeRole}
		/>
	{:else if currentPhase === 'team_building'}
		<TeamBuilding
			{isLeader}
			{currentQuest}
			{questResults}
			{playerCount}
			{requiredTeamSize}
			{players}
			leaderId={currentLeaderId}
			{rejectionCount}
			onProposeTeam={handleProposeTeam}
		/>
	{:else if currentPhase === 'team_voting'}
		<TeamVoting
			{currentQuest}
			{questResults}
			{playerCount}
			{proposedTeam}
			{players}
			{hasVoted}
			votesSubmitted={votesSubmittedPublic}
			totalVotes={playerCount}
			{rejectionCount}
			onVote={handleVoteTeam}
		/>
	{:else if currentPhase === 'quest_execution' && myTeam}
		<QuestExecution
			{currentQuest}
			{questResults}
			{playerCount}
			teamMembers={proposedTeam}
			{players}
			isOnTeam={isOnProposedTeam}
			team={myTeam}
			{hasPlayedCard}
			cardsSubmitted={cardsSubmittedPublic}
			totalCardsExpected={proposedTeam.length}
			onPlayCard={handlePlayQuestCard}
		/>
	{:else if currentPhase === 'assassination'}
		<Assassination
			{currentQuest}
			{questResults}
			{playerCount}
			{players}
			isAssassin={myRole === 'assassin'}
			currentPlayerId={myPlayerId}
			onAssassinate={handleAssassinate}
		/>
	{:else if currentPhase === 'finished' && winningTeam}
		<Results
			{winningTeam}
			{winReason}
			roles={allRoles}
			teams={allTeams}
			questHistory={questResults}
			{players}
			{playerCount}
		/>
	{:else}
		<div class="loading bg-[#3c3836] p-8 rounded-lg text-center">
			<p class="text-[#a89984]">Preparing game...</p>
		</div>
	{/if}
</div>

<style>
	.avalon-game {
		max-width: 1200px;
		margin: 0 auto;
	}
</style>
