# Games Roadmap

> **Play at [cardless.games](https://cardless.games)**

This document outlines detailed implementation plans for games to be added to the Cardless platform.

## Current Status

### ✅ Implemented
- **One Night Werewolf** - Social deduction with hidden roles and night actions

### 🔄 In Progress
- Completing digital night actions for all Werewolf roles
- Role reveal and play again features

### 📋 Planned
See below for detailed plans

---

## Game 2: Avalon (The Resistance)

**Priority**: HIGH - validates multi-game architecture

**Player Count**: 5-10 players

**Game Length**: 15-30 minutes

### Overview
Avalon is a social deduction game where players are divided into two teams: Good (loyal servants of Arthur) and Evil (minions of Mordred). The good team tries to complete three quests successfully, while the evil team tries to sabotage them. The twist: Merlin knows who is evil, but the Assassin can win by identifying Merlin.

### Core Mechanics
1. **Role Assignment**: Players receive roles (Merlin, Percival, Morgana, Assassin, etc.)
2. **Quest Voting**: Leader proposes a team for each quest
3. **Team Approval**: Everyone votes to approve or reject the team
4. **Quest Execution**: Selected players secretly support or sabotage
5. **Assassin Endgame**: If good wins 3 quests, Assassin can identify Merlin to steal victory

### Technical Requirements

#### Backend
- Quest tracking system (5 quests, each with team size requirements)
- Team proposal and voting logic
- Quest card playing (success/fail)
- Leader rotation mechanism
- Assassin targeting system
- Role-specific information filtering (Merlin sees evil, Percival sees Merlin + Morgana)

#### Frontend Components
- **Role Reveal**: Show player's role and what they know
- **Quest Board**: Visual display of quest history (5 quests)
- **Team Selection**: Leader UI to select team members
- **Team Voting**: All players vote to approve/reject team
- **Quest Cards**: Selected players play success/fail cards
- **Assassin Phase**: Assassin selects who they think is Merlin
- **Results**: Show final roles and explain win condition

#### Events
- `role_assigned` (private, with role-specific info)
- `leader_changed`
- `team_proposed`
- `team_vote_cast`
- `team_approved` / `team_rejected`
- `quest_card_played`
- `quest_completed` (with success/fail counts)
- `assassin_targets`
- `game_finished` (with winning team)

### Implementation Estimate
- **Backend**: 2 days (game logic, quest system, role mechanics)
- **Frontend**: 2-3 days (quest board UI, team selection, voting)
- **Testing**: 1 day (multiplayer testing)
- **Total**: 5-6 days

### Design Notes
- Reuse lobby and role reveal patterns from Werewolf
- Quest board should be prominent and visual
- Team selection should be easy on mobile (tap to toggle)
- Make leader role obvious to all players
- Provide clear guidance on team sizes per quest

---

## Game 3: Spyfall

**Priority**: HIGH - different mechanics from Werewolf/Avalon

**Player Count**: 3-8 players

**Game Length**: 6-8 minutes per round

### Overview
All players are at the same location except one spy. Everyone knows the location except the spy. Players ask each other questions to figure out who the spy is without revealing the location. The spy tries to blend in and guess the location.

### Core Mechanics
1. **Location Assignment**: One random player is the spy, others get the location
2. **Question Rounds**: Players take turns asking each other questions
3. **Timer**: 6-8 minute round timer
4. **Accusation**: Players can vote to accuse someone of being the spy
5. **Spy Guess**: If accused, spy can guess the location to win

### Technical Requirements

#### Backend
- Location database (30-50 locations with role variations)
- Random spy selection
- Turn order management
- Timer system
- Voting system for accusations
- Spy location guessing

#### Frontend Components
- **Role Reveal**: Spy sees "You are the SPY!" / Others see location + role
- **Location Display**: Persistent display of your location (not shown to spy)
- **Question Timer**: Countdown timer for the round
- **Turn Indicator**: Shows whose turn it is to ask/answer
- **Accusation UI**: Vote for who you think is the spy
- **Spy Guess**: Spy selects location from list if accused
- **Location List**: Reference sheet of all possible locations (for spy)
- **Results**: Reveal spy and location, show who won

#### Events
- `location_assigned` (private)
- `turn_changed`
- `timer_started`
- `timer_extended`
- `accusation_started`
- `accusation_vote_cast`
- `spy_accused`
- `spy_guesses_location`
- `round_finished`

### Implementation Estimate
- **Backend**: 2 days (location system, turn management)
- **Frontend**: 2 days (clean timer UI, location display, accusation flow)
- **Content**: 0.5 days (create location database)
- **Testing**: 0.5 days
- **Total**: 4-5 days

### Design Notes
- Clean, minimal UI - this is a conversation game
- Timer should be prominent but not distracting
- Location should be easy to reference (always visible for non-spies)
- Provide location list reference for spy
- Consider adding role variants per location for replayability

---

## Game 4: Skull

**Priority**: MEDIUM - pure bluffing mechanics

**Player Count**: 3-6 players

**Game Length**: 15-30 minutes

### Overview
Each player has four discs: three flowers and one skull. Players take turns placing discs face-down. At any point, a player can make a bid to flip over a number of discs, trying to reveal only flowers. If successful, they earn a point. If they reveal a skull, they lose a disc.

### Core Mechanics
1. **Setup**: Each player has 4 discs (3 flowers, 1 skull)
2. **Playing Phase**: Players take turns playing discs face-down
3. **Bidding Phase**: Players bid how many flowers they can reveal
4. **Reveal Phase**: Highest bidder reveals discs (their own first, then others)
5. **Winning**: First to win 2 points wins, OR last player with discs remaining

### Technical Requirements

#### Backend
- Disc state management per player (4 discs each)
- Playing phase logic
- Bidding system (increasing bids)
- Reveal logic (must reveal all your own first)
- Point tracking
- Disc loss mechanic (random if multiple discs, or player choice)

#### Frontend Components
- **Player Mat**: Visual representation of your 4 discs
- **Playing Area**: Central area showing all played discs (face-down)
- **Play Disc**: UI to select flower or skull to play
- **Bidding UI**: Make or pass on bids (must increase)
- **Reveal**: Animation of discs being flipped
- **Scoring**: Track points and remaining discs per player

#### Events
- `disc_played`
- `bidding_started`
- `bid_made`
- `bid_passed`
- `bidding_finished`
- `disc_revealed`
- `challenge_succeeded` / `challenge_failed`
- `point_scored`
- `disc_lost`
- `round_finished`

### Implementation Estimate
- **Backend**: 1.5 days (disc management, bidding, reveal logic)
- **Frontend**: 2 days (visual disc representation, animations)
- **Testing**: 0.5 days
- **Total**: 3-4 days

### Design Notes
- Visual design is key - make discs feel tactile
- Clear indication of whose turn it is
- Show all players' remaining disc counts
- Smooth animations for reveals
- Consider adding player choice for which disc to lose

---

## Game 5: Wavelength

**Priority**: MEDIUM - unique UI challenge

**Player Count**: 2-12 players (team-based)

**Game Length**: 30-45 minutes

### Overview
A team game where players give clues on a spectrum. One player sees a hidden target on a spectrum (e.g., "Cold ← → Hot") and gives a one-word clue. Their team guesses where on the spectrum the target is.

### Core Mechanics
1. **Teams**: Players divided into two teams
2. **Spectrum Cards**: Cards with opposing concepts (e.g., "Underrated ← → Overrated")
3. **Target**: Hidden target on the spectrum (left, center, or right zone)
4. **Clue**: Clue-giver provides one word/phrase
5. **Guess**: Team discusses and places guess on spectrum
6. **Scoring**: Points based on accuracy (4 points for bullseye, 3 for close, 2 for zone)

### Technical Requirements

#### Backend
- Spectrum card database (100+ spectrums)
- Random target generation (with fuzzy zones for scoring)
- Team management
- Turn/clue-giver rotation
- Scoring system (zone-based)

#### Frontend Components
- **Team Assignment**: Divide players into teams
- **Spectrum Display**: Visual spectrum with endpoints
- **Clue Entry**: Clue-giver enters clue
- **Dial/Guess UI**: Interactive dial for team to place guess
- **Reveal**: Animated reveal of target location
- **Scoring**: Visual feedback on accuracy
- **Score Board**: Track team scores

#### Events
- `teams_assigned`
- `spectrum_drawn` (clue-giver sees target)
- `clue_given`
- `guess_placed`
- `target_revealed`
- `points_scored`
- `turn_finished`

### Implementation Estimate
- **Backend**: 2 days (spectrum system, scoring logic)
- **Frontend**: 3 days (dial UI component, animations)
- **Content**: 1 day (create spectrum database)
- **Testing**: 1 day
- **Total**: 6-7 days

### Design Notes
- The dial/spectrum UI is the centerpiece - needs to be intuitive
- Consider using a circular dial or linear slider
- Animated reveal should build tension
- Provide spectrum card reference list
- Consider allowing teams to debate before locking in guess

---

## Game 6: Coup

**Priority**: LOW - complex state management

**Player Count**: 2-6 players

**Game Length**: 15 minutes

### Overview
Bluffing game where each player has two secret character cards. On your turn, you can claim any character's ability. Other players can challenge you. Last player with influence (cards) wins.

### Core Mechanics
1. **Characters**: Duke, Assassin, Captain, Ambassador, Contessa (each with unique abilities)
2. **Actions**: Income, Foreign Aid, Coup, or character abilities
3. **Bluffing**: Can claim any character ability without having it
4. **Challenging**: Call out bluffs (loser loses influence)
5. **Blocking**: Some actions can be blocked by certain characters
6. **Elimination**: Lose both influence cards = eliminated

### Technical Requirements

#### Backend
- Character cards (3 copies of each character)
- Coin management
- Action validation and resolution
- Challenge system
- Block system
- Influence loss (reveal cards)
- Complex state machine for action resolution

#### Frontend Components
- **Influence Display**: Your two cards (face-up when lost)
- **Coin Counter**: Track coins per player
- **Action UI**: Select action to take
- **Challenge/Block UI**: React to other players' actions
- **Character Reference**: Quick reference for abilities
- **Game Log**: Recent actions for context

#### Events
- `cards_dealt` (private)
- `action_declared`
- `challenge_declared`
- `block_declared`
- `challenge_resolved`
- `coins_changed`
- `influence_lost`
- `player_eliminated`
- `game_finished`

### Implementation Estimate
- **Backend**: 3 days (complex action resolution, challenge system)
- **Frontend**: 2 days (action UI, challenge/block flow)
- **Testing**: 1 day (many edge cases)
- **Total**: 5-6 days

### Design Notes
- Needs clear action log to track complex interactions
- Character reference should always be accessible
- Challenge/block prompts need clear timeouts
- Consider "undo last action" for mistakes
- This is the most complex game - save for later

---

## Implementation Priority

### Phase 1 (After Werewolf Complete)
1. **Avalon** (5-6 days) - Validates platform for different game types
2. **Spyfall** (4-5 days) - Different social mechanics

### Phase 2 (After Platform Validated)
3. **Skull** (3-4 days) - Pure bluffing, simpler than others
4. **Wavelength** (6-7 days) - Team game, unique UI

### Phase 3 (Long-term)
5. **Coup** (5-6 days) - Most complex, save for when platform is mature

---

## Common Patterns to Extract

As we implement more games, look for patterns to extract into shared components:

### UI Components
- Timer component (used in Werewolf, Spyfall, Wavelength)
- Voting system (used in Werewolf, Avalon, Skull)
- Team assignment (used in Wavelength, potentially others)
- Role reveal screen (used in all hidden role games)
- Turn indicator (used in Spyfall, Skull, Coup)

### Backend Patterns
- Role assignment algorithms
- Turn rotation logic
- Voting collection and resolution
- Timer management
- Phase transitions

### Testing Strategies
- Multi-player simulation framework
- Common edge case scenarios (disconnect, rejoin, timeout)
- Game-agnostic test utilities

---

## Game Submission Guidelines

Want to add a new game? Consider these questions:

### Suitability Checklist
- ✅ Works well in-person (not remote)
- ✅ Benefits from digital assistant (complex rules, hidden info, or timing)
- ✅ Supports 3-10 players
- ✅ Games last 5-45 minutes
- ✅ Minimal physical components needed
- ✅ Not trademarked/proprietary (or we have permission)

### Implementation Complexity
- **Low** (2-4 days): Simple rules, minimal state
- **Medium** (4-7 days): Moderate complexity, some unique UI
- **High** (7+ days): Complex interactions, custom UI components

### Documentation Required
- Game overview and rules
- Player count and game length
- Core mechanics breakdown
- Backend requirements
- Frontend components needed
- Event types
- Estimated implementation time

---

**Last Updated**: November 2024

For questions or game suggestions, open an issue or check [CONTRIBUTING.md](CONTRIBUTING.md).

