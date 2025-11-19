# Avalon (The Resistance) - Prototype Plan

> **Game Type**: Social Deduction, Team-Based, Quest Voting
> **Players**: 5-10
> **Duration**: 15-30 minutes
> **Priority**: HIGH - Phase 2 validation game

---

## Table of Contents
1. [Game Rules Reference](#game-rules-reference)
2. [Design Principles Application](#design-principles-application)
3. [Complete Event Flow](#complete-event-flow)
4. [Phase Breakdown](#phase-breakdown)
5. [Backend Implementation](#backend-implementation)
6. [Frontend Implementation](#frontend-implementation)
7. [State Management](#state-management)
8. [Implementation Checklist](#implementation-checklist)

---

## Game Rules Reference

### Overview
Avalon is a social deduction game of hidden loyalty. Players are divided into two teams: **Good** (loyal servants of Arthur) and **Evil** (minions of Mordred). The Good team tries to successfully complete three out of five quests, while the Evil team tries to sabotage them. The twist: certain Good players know who is Evil, but must remain hidden.

### Setup
- **5 Players**: 3 Good, 2 Evil (quest sizes: 2, 3, 2, 3, 3)
- **6 Players**: 4 Good, 2 Evil (quest sizes: 2, 3, 4, 3, 4)
- **7 Players**: 4 Good, 3 Evil (quest sizes: 2, 3, 3, 4*, 4)
- **8-10 Players**: See official rules (7+ quest 4 requires 2 fails)

\* Quest 4 with 7+ players requires **2 FAIL cards** to fail

### Roles

#### Good Team Roles
- **Merlin** 🔮 (Required)
  - Sees all Evil players (except Mordred if in game)
  - If Assassin identifies Merlin after Good wins 3 quests, Evil wins instead
  - Must help Good team win without revealing identity

- **Percival** 👁️ (Optional)
  - Sees Merlin and Morgana (can't distinguish which is which)
  - Helps protect Merlin by creating confusion

- **Loyal Servant of Arthur** ⚔️ (Filler role)
  - No special knowledge
  - Must use deduction and trust to identify team

#### Evil Team Roles
- **Assassin** 🗡️ (Required if Merlin is in game)
  - Sees all Evil teammates
  - If Good wins 3 quests, can attempt to assassinate Merlin
  - If assassinates correctly, Evil wins

- **Morgana** 🌙 (Optional)
  - Appears as Merlin to Percival
  - Sees all Evil teammates
  - Creates confusion for Percival

- **Mordred** 😈 (Optional)
  - Hidden from Merlin
  - Sees all Evil teammates
  - Invisible to Merlin's knowledge

- **Oberon** 👻 (Optional)
  - Does NOT see Evil teammates
  - Invisible to other Evil players
  - Lone wolf saboteur

- **Minion of Mordred** 💀 (Filler role)
  - Sees all Evil teammates (except Oberon)
  - No special abilities

### Game Flow

#### 1. Role Assignment Phase
- Players receive roles
- Players with special knowledge see relevant information:
  - **Merlin** sees all Evil (except Mordred)
  - **Evil players** see each other (except Oberon)
  - **Percival** sees Merlin and Morgana

#### 2. Quest Phase (Repeat 5 times or until one team wins 3)

**Step A: Team Building**
- Current **Leader** proposes team of required size
- Leader rotates clockwise each round (including failed votes)

**Step B: Team Approval Vote**
- All players simultaneously vote **Approve** or **Reject**
- If majority approves → proceed to Quest
- If majority rejects → leader rotates, try again
- If **5 consecutive rejections** → Evil wins automatically

**Step C: Quest Execution**
- Only team members participate
- Each selected player secretly plays **SUCCESS** or **FAIL** card
- Good players can ONLY play Success
- Evil players can choose Success or Fail
- Cards shuffled and revealed:
  - **1+ Fail cards** → Quest fails (2+ fails for Quest 4 with 7+ players)
  - **All Success cards** → Quest succeeds

**Step D: Quest Resolution**
- Update quest track
- First team to 3 wins... or do they?

#### 3. Assassination Phase (Only if Good wins 3 quests)
- **Assassin** must identify **Merlin**
- Assassin selects one player
- If correct → **Evil wins**
- If incorrect → **Good wins**

### Win Conditions
1. **Evil Wins If**:
   - 3 quests fail
   - 5 consecutive team rejections
   - Assassin correctly identifies Merlin after Good wins 3 quests

2. **Good Wins If**:
   - 3 quests succeed AND Assassin fails to identify Merlin (or Merlin not in game)

---

## Design Principles Application

### 1. Enhancing In-Person Interactions ✅
**Digital**:
- Role assignment with hidden information
- Quest card submissions (prevents counting/stacking tells)
- Vote counting and revelation
- Quest history tracking
- Assassination selection

**Physical**:
- Team discussions and deliberations
- Leader proposals (verbal announcement)
- Social reads and body language
- Accusation and defense
- Victory celebration

### 2. Tracking and Automating Card States ✅
- **Role cards**: Digital assignment ensures proper distribution and hidden knowledge
- **Quest cards**: Digital submission prevents physical tells (card orientation, timing)
- **Vote tokens**: Digital voting prevents early reveals and ensures simultaneous submission
- **Quest track**: Automated tracking of successes/failures
- **Leader token**: Automatic rotation system

### 3. Preserving Social Deduction Core
- App handles **mechanics** (quest results, vote counting)
- Players handle **social** (accusations, defenses, tells, reads)
- No AI analysis or hints about who is lying
- Results displayed neutrally without commentary

### 4. Mobile-First UX
- Large touch targets for team selection (48x48px minimum)
- Simple approve/reject buttons
- Clear visual quest board
- One-handed operation possible
- Quick reference for role abilities

---

## Complete Event Flow

### Event Types & Visibility

#### Public Events (Everyone sees)
```typescript
// Setup & Phase Changes
EVENT_GAME_STARTED          // Game initialized with player count
EVENT_PHASE_CHANGED         // Phase transition
EVENT_LEADER_CHANGED        // New leader assigned
EVENT_GAME_FINISHED         // Game concluded with winning team

// Team Building
EVENT_TEAM_PROPOSED         // Leader proposes team (player IDs visible)
EVENT_TEAM_VOTE_CAST        // Someone voted (not their vote, just that they voted)
EVENT_TEAM_APPROVED         // Team vote passed
EVENT_TEAM_REJECTED         // Team vote failed
EVENT_REJECTION_COUNT       // Track consecutive rejections (max 5)

// Quest Execution
EVENT_QUEST_CARD_PLAYED     // Someone played card (not which card)
EVENT_QUEST_COMPLETED       // Quest result with success/fail counts

// Assassination
EVENT_ASSASSIN_TARGET       // Assassin selected target (public)
EVENT_ASSASSINATION_RESULT  // Success or failure
```

#### Private Events (Specific players only)
```typescript
// Role Assignment
role_assigned               // Your role and team (private to each player)
role_knowledge              // What you see based on your role:
                           // - Merlin: sees Evil (except Mordred)
                           // - Percival: sees Merlin + Morgana
                           // - Evil: sees Evil teammates (except Oberon)
                           // - Oberon: sees nothing
                           // - Good: sees nothing

// Voting
team_vote_recorded         // Your vote was recorded (confirmation)
quest_card_recorded        // Your quest card was recorded (confirmation)
```

---

## Phase Breakdown

### Phase 1: SETUP
**Duration**: Until game starts
**Transitions to**: ROLE_REVEAL

**Actions**:
- Host configures roles (must follow player count rules)
- Validates configuration (e.g., Assassin required if Merlin present)
- Assigns first leader (random)

**Events**:
- `EVENT_GAME_STARTED`
- `role_assigned` (private to each player)
- `role_knowledge` (private, role-specific)
- `EVENT_LEADER_CHANGED`
- `EVENT_PHASE_CHANGED` → ROLE_REVEAL

**Backend State**:
```go
type AvalonState struct {
    Players        []*Player
    Roles          map[string]Role  // playerID -> Role
    CurrentLeader  string           // playerID
    QuestNumber    int              // 1-5
    QuestResults   []QuestResult    // History
    RejectionCount int              // Consecutive rejections (max 5)
    Phase          GamePhase
}
```

---

### Phase 2: ROLE_REVEAL
**Duration**: Until all players acknowledge
**Transitions to**: TEAM_BUILDING

**UI Display**:
- **Your Role**: Large card showing role name, emoji, team
- **What You Know**: Role-specific information
  - Merlin: "You see the forces of Evil: [names]"
  - Percival: "You see two powerful wizards: [names]. One is Merlin, one is Morgana."
  - Evil: "Your Evil allies are: [names]"
  - Loyal Servant: "You know nothing. Trust your instincts."
- **Acknowledge Button**: "I understand my role"

**Actions**:
- Players read and acknowledge roles
- Auto-advance when all acknowledged

**Events**:
- `player_acknowledged_role` (public count only)
- `EVENT_PHASE_CHANGED` → TEAM_BUILDING

---

### Phase 3: TEAM_BUILDING
**Duration**: Until leader proposes team
**Transitions to**: TEAM_VOTING

**UI Display**:
- **Quest Board**: Shows current quest number and required team size
- **Leader Indicator**: Highlight current leader
- **Team Selection** (leader only):
  - Grid of player avatars
  - Tap to select/deselect
  - Counter showing X/Y selected
  - "Propose Team" button (enabled when correct count)

**Actions**:
- Leader selects players for quest (exact count required)
- Leader confirms proposal
- System validates team composition

**Events**:
- `EVENT_TEAM_PROPOSED` (public: leader ID, selected player IDs)
- `EVENT_PHASE_CHANGED` → TEAM_VOTING

**Validation**:
- Team size must match quest requirements
- Only current leader can propose
- Cannot propose duplicate players

---

### Phase 4: TEAM_VOTING
**Duration**: Until all players vote
**Transitions to**: TEAM_BUILDING (if rejected) or QUEST_EXECUTION (if approved)

**UI Display**:
- **Proposed Team**: List of selected players
- **Vote Buttons**: APPROVE (✅) / REJECT (❌)
- **Vote Count**: "3/7 players have voted"
- **Rejection Counter**: "Rejection 2/5" (if applicable)

**Actions**:
- All players vote simultaneously
- Votes hidden until all submitted
- System tallies votes
- Majority required to approve

**Events**:
- `EVENT_TEAM_VOTE_CAST` (public: player voted, not their vote)
- `team_vote_recorded` (private: confirmation to voter)
- `EVENT_TEAM_APPROVED` or `EVENT_TEAM_REJECTED`
  - Payload: vote breakdown (X approve, Y reject)
- `EVENT_REJECTION_COUNT` (if rejected)

**Transitions**:
- **If Approved**: → QUEST_EXECUTION
- **If Rejected**:
  - Increment rejection count
  - Rotate leader
  - If 5 rejections → GAME_FINISHED (Evil wins)
  - Else → TEAM_BUILDING

---

### Phase 5: QUEST_EXECUTION
**Duration**: Until all team members play cards
**Transitions to**: QUEST_RESULTS

**UI Display**:
- **Team Members Only**: Quest card selection
  - "SUCCESS" card (always available)
  - "FAIL" card (only if Evil team)
- **Other Players**: Waiting indicator
- **Card Count**: "3/4 cards played"

**Actions**:
- Team members select and submit cards
- Good players can ONLY submit Success
- Evil players choose Success or Fail
- System shuffles cards to prevent order tells

**Events**:
- `EVENT_QUEST_CARD_PLAYED` (public: someone played, not which card)
- `quest_card_recorded` (private: confirmation)
- `EVENT_QUEST_COMPLETED`
  - Payload: shuffled results (e.g., [Success, Success, Fail])

**Validation**:
- Good team can only play Success
- Evil team can play either
- Exactly one card per team member

---

### Phase 6: QUEST_RESULTS
**Duration**: Brief pause for reading
**Transitions to**: TEAM_BUILDING (if game continues) or ASSASSINATION (if Good wins 3) or GAME_FINISHED

**UI Display**:
- **Quest Board**: Updated with result
- **Cards Revealed**: Visual display of shuffled results
  - "3 SUCCESS, 1 FAIL" → Quest Failed
- **Score Track**:
  - Good: ⭕⭕❌ (2 success, 1 fail)
  - Evil: ❌⭕⭕ (1 fail, 2 success)

**Actions**:
- Display quest results
- Update quest track
- Auto-advance after 3 seconds

**Events**:
- `EVENT_PHASE_CHANGED`

**Transitions**:
- If Good has 3 successes AND Merlin in game → ASSASSINATION
- If Good has 3 successes AND no Merlin → GAME_FINISHED (Good wins)
- If Evil has 3 fails → GAME_FINISHED (Evil wins)
- Else → rotate leader, increment quest number, → TEAM_BUILDING

---

### Phase 7: ASSASSINATION (Only if Good wins 3 quests with Merlin present)
**Duration**: Until Assassin selects target
**Transitions to**: GAME_FINISHED

**UI Display**:
- **Assassin Only**: "Select who you believe is Merlin"
  - Grid of all players (except self)
  - Confirm selection button
- **Other Players**: "The Assassin is choosing their target..."

**Actions**:
- Assassin selects one player
- Confirms selection
- System reveals if correct

**Events**:
- `EVENT_ASSASSIN_TARGET` (public: selected player)
- `EVENT_ASSASSINATION_RESULT` (success/failure)
- `EVENT_GAME_FINISHED`

**Win Condition Check**:
- If target is Merlin → Evil wins
- If target is not Merlin → Good wins

---

### Phase 8: GAME_FINISHED
**Duration**: Indefinite
**Transitions to**: New game (if host chooses)

**UI Display**:
- **Final Roles**: All players revealed with roles
- **Winning Team**: Banner with result
- **Win Condition**: Explanation
  - "Good won 3 quests, and the Assassin failed to identify Merlin!"
  - "Evil won! The Assassin correctly identified Merlin."
  - "Evil sabotaged 3 quests!"
  - "Evil won after 5 consecutive team rejections."
- **Quest History**: Full quest board with results
- **Play Again Button** (host only)

**Events**:
- None (terminal phase)

---

## Backend Implementation

### File Structure
```
backend/internal/games/avalon/
├── game.go           # Implements core.Game interface
├── state.go          # Avalon-specific state types
├── config.go         # Role configuration and validation
├── phases.go         # Phase transition logic
├── roles.go          # Role definitions and knowledge
├── quests.go         # Quest logic (team sizes, fail requirements)
└── voting.go         # Team vote and quest card logic
```

### State Types

```go
// state.go
package avalon

type Role string

const (
    RoleMerlin          Role = "merlin"
    RolePercival        Role = "percival"
    RoleLoyalServant    Role = "loyal_servant"
    RoleAssassin        Role = "assassin"
    RoleMorgana         Role = "morgana"
    RoleMordred         Role = "mordred"
    RoleOberon          Role = "oberon"
    RoleMinionOfMordred Role = "minion"
)

type Team string

const (
    TeamGood Team = "good"
    TeamEvil Team = "evil"
)

type GamePhase string

const (
    PhaseSetup         GamePhase = "setup"
    PhaseRoleReveal    GamePhase = "role_reveal"
    PhaseTeamBuilding  GamePhase = "team_building"
    PhaseTeamVoting    GamePhase = "team_voting"
    PhaseQuestExec     GamePhase = "quest_execution"
    PhaseQuestResults  GamePhase = "quest_results"
    PhaseAssassination GamePhase = "assassination"
    PhaseFinished      GamePhase = "finished"
)

type QuestCard string

const (
    CardSuccess QuestCard = "success"
    CardFail    QuestCard = "fail"
)

type Vote string

const (
    VoteApprove Vote = "approve"
    VoteReject  Vote = "reject"
)

type QuestResult struct {
    QuestNumber   int         `json:"quest_number"`
    TeamSize      int         `json:"team_size"`
    TeamMembers   []string    `json:"team_members"`   // player IDs
    Cards         []QuestCard `json:"cards"`          // shuffled results
    Success       bool        `json:"success"`
    FailsRequired int         `json:"fails_required"` // 1 or 2
}

type AvalonState struct {
    // Core game state
    Players        []*core.Player
    Roles          map[string]Role        // playerID -> Role
    Teams          map[string]Team        // playerID -> Team
    Config         *AvalonConfig
    Phase          GamePhase

    // Quest tracking
    QuestNumber    int                    // 1-5
    QuestResults   []QuestResult
    CurrentLeader  string                 // playerID
    LeaderIndex    int                    // position in player list
    RejectionCount int                    // consecutive rejections (max 5)

    // Team voting
    ProposedTeam   []string               // player IDs
    TeamVotes      map[string]Vote        // playerID -> vote

    // Quest execution
    QuestCards     map[string]QuestCard   // playerID -> card (team members only)

    // Assassination
    AssassinTarget string                 // playerID

    // Acknowledgments
    Acknowledged   map[string]bool        // playerID -> acknowledged
}

type AvalonConfig struct {
    Roles           []Role  `json:"roles"`
    IncludeMerlin   bool    `json:"include_merlin"`
    IncludePercival bool    `json:"include_percival"`
    IncludeMorgana  bool    `json:"include_morgana"`
    IncludeMordred  bool    `json:"include_mordred"`
    IncludeOberon   bool    `json:"include_oberon"`
}
```

### Role Knowledge Logic

```go
// roles.go
package avalon

// GetRoleKnowledge returns what a player sees based on their role
func (g *AvalonGame) GetRoleKnowledge(playerID string) []string {
    role := g.state.Roles[playerID]

    switch role {
    case RoleMerlin:
        // Merlin sees all Evil except Mordred
        return g.getEvilPlayersExcept(RoleMordred)

    case RolePercival:
        // Percival sees Merlin and Morgana (can't distinguish)
        return g.getPlayersWithRoles([]Role{RoleMerlin, RoleMorgana})

    case RoleAssassin, RoleMorgana, RoleMordred, RoleMinionOfMordred:
        // Evil sees each other (except Oberon)
        return g.getEvilPlayersExcept(RoleOberon)

    case RoleOberon:
        // Oberon sees nothing
        return []string{}

    default:
        // Loyal servants see nothing
        return []string{}
    }
}

func (g *AvalonGame) getEvilPlayersExcept(exceptRole Role) []string {
    var result []string
    for pid, role := range g.state.Roles {
        if g.state.Teams[pid] == TeamEvil && role != exceptRole {
            result = append(result, pid)
        }
    }
    return result
}

func (g *AvalonGame) getPlayersWithRoles(roles []Role) []string {
    var result []string
    for pid, role := range g.state.Roles {
        for _, r := range roles {
            if role == r {
                result = append(result, pid)
                break
            }
        }
    }
    return result
}
```

### Quest Configuration

```go
// quests.go
package avalon

// QuestConfig defines team sizes and fail requirements per player count
type QuestConfig struct {
    PlayerCount int
    TeamSizes   [5]int  // Quest 1-5 team sizes
    Quest4Fails int     // Number of fails required for quest 4 (1 or 2)
}

var questConfigs = map[int]QuestConfig{
    5:  {PlayerCount: 5, TeamSizes: [5]int{2, 3, 2, 3, 3}, Quest4Fails: 1},
    6:  {PlayerCount: 6, TeamSizes: [5]int{2, 3, 4, 3, 4}, Quest4Fails: 1},
    7:  {PlayerCount: 7, TeamSizes: [5]int{2, 3, 3, 4, 4}, Quest4Fails: 2},
    8:  {PlayerCount: 8, TeamSizes: [5]int{3, 4, 4, 5, 5}, Quest4Fails: 2},
    9:  {PlayerCount: 9, TeamSizes: [5]int{3, 4, 4, 5, 5}, Quest4Fails: 2},
    10: {PlayerCount: 10, TeamSizes: [5]int{3, 4, 4, 5, 5}, Quest4Fails: 2},
}

func (g *AvalonGame) GetCurrentQuestConfig() QuestConfig {
    return questConfigs[len(g.state.Players)]
}

func (g *AvalonGame) GetRequiredTeamSize() int {
    config := g.GetCurrentQuestConfig()
    return config.TeamSizes[g.state.QuestNumber-1]
}

func (g *AvalonGame) GetFailsRequired() int {
    if g.state.QuestNumber == 4 {
        return g.GetCurrentQuestConfig().Quest4Fails
    }
    return 1
}
```

### Configuration Validation

```go
// config.go
package avalon

func ValidateConfig(playerCount int, config *AvalonConfig) error {
    // Validate player count
    if playerCount < 5 || playerCount > 10 {
        return fmt.Errorf("Avalon requires 5-10 players, got %d", playerCount)
    }

    // Count good and evil roles
    goodCount := 0
    evilCount := 0

    for _, role := range config.Roles {
        if isGoodRole(role) {
            goodCount++
        } else {
            evilCount++
        }
    }

    // Validate team sizes based on player count
    expectedGood, expectedEvil := getExpectedTeamSizes(playerCount)
    if goodCount != expectedGood || evilCount != expectedEvil {
        return fmt.Errorf(
            "Invalid team sizes for %d players: expected %d good, %d evil; got %d good, %d evil",
            playerCount, expectedGood, expectedEvil, goodCount, evilCount,
        )
    }

    // If Merlin is present, Assassin must be present
    hasMerlin := contains(config.Roles, RoleMerlin)
    hasAssassin := contains(config.Roles, RoleAssassin)
    if hasMerlin && !hasAssassin {
        return fmt.Errorf("Assassin is required when Merlin is present")
    }

    // If Percival is present, Merlin must be present
    hasPercival := contains(config.Roles, RolePercival)
    if hasPercival && !hasMerlin {
        return fmt.Errorf("Merlin is required when Percival is present")
    }

    // If Morgana is present, Percival should be present (warning only)
    hasMorgana := contains(config.Roles, RoleMorgana)
    if hasMorgana && !hasPercival {
        // This is allowed but recommended to include Percival
    }

    return nil
}

func getExpectedTeamSizes(playerCount int) (good int, evil int) {
    switch playerCount {
    case 5, 6:
        return playerCount - 2, 2
    case 7, 8, 9:
        return playerCount - 3, 3
    case 10:
        return 6, 4
    default:
        return 0, 0
    }
}
```

### Phase Transitions

```go
// phases.go
package avalon

func (g *AvalonGame) AdvanceToRoleReveal() ([]core.GameEvent, error) {
    g.state.Phase = PhaseRoleReveal

    events := []core.GameEvent{}

    // Phase change event
    phaseEvent, _ := core.NewPublicEvent(
        "phase_changed",
        "system",
        map[string]interface{}{
            "phase": PhaseRoleReveal,
        },
    )
    events = append(events, phaseEvent)

    return events, nil
}

func (g *AvalonGame) AdvanceToTeamBuilding() ([]core.GameEvent, error) {
    g.state.Phase = PhaseTeamBuilding
    g.state.ProposedTeam = nil
    g.state.TeamVotes = make(map[string]Vote)

    events := []core.GameEvent{}

    // Phase change
    phaseEvent, _ := core.NewPublicEvent(
        "phase_changed",
        "system",
        map[string]interface{}{
            "phase":        PhaseTeamBuilding,
            "quest_number": g.state.QuestNumber,
            "team_size":    g.GetRequiredTeamSize(),
        },
    )
    events = append(events, phaseEvent)

    return events, nil
}

func (g *AvalonGame) AdvanceToQuestExecution() ([]core.GameEvent, error) {
    g.state.Phase = PhaseQuestExec
    g.state.QuestCards = make(map[string]QuestCard)

    events := []core.GameEvent{}

    phaseEvent, _ := core.NewPublicEvent(
        "phase_changed",
        "system",
        map[string]interface{}{
            "phase":        PhaseQuestExec,
            "team_members": g.state.ProposedTeam,
        },
    )
    events = append(events, phaseEvent)

    return events, nil
}

func (g *AvalonGame) CheckGameEnd() (bool, Team, string) {
    // Count quest successes and failures
    goodWins := 0
    evilWins := 0

    for _, result := range g.state.QuestResults {
        if result.Success {
            goodWins++
        } else {
            evilWins++
        }
    }

    // Check for 3 wins
    if goodWins >= 3 {
        // Good wins, but assassination may reverse
        if g.hasRole(RoleMerlin) {
            return false, "", "assassination_required"
        }
        return true, TeamGood, "good_won_quests"
    }

    if evilWins >= 3 {
        return true, TeamEvil, "evil_sabotaged_quests"
    }

    // Check rejection limit
    if g.state.RejectionCount >= 5 {
        return true, TeamEvil, "five_rejections"
    }

    return false, "", ""
}
```

---

## Frontend Implementation

### File Structure
```
frontend/src/lib/games/avalon/
├── AvalonGame.svelte           # Main game container
├── RoleReveal.svelte           # Role assignment + knowledge display
├── TeamBuilding.svelte         # Leader team selection
├── TeamVoting.svelte           # All players vote on team
├── QuestExecution.svelte       # Team members play cards
├── QuestResults.svelte         # Display quest outcome
├── Assassination.svelte        # Assassin selection
├── Results.svelte              # Final game results
├── QuestBoard.svelte           # Visual quest track (1-5)
├── PlayerSelector.svelte       # Reusable player selection grid
├── RoleCard.svelte             # Role display component
└── roleConfig.ts               # Role definitions and styling
```

### Role Configuration

```typescript
// roleConfig.ts
export type AvalonRole =
    | 'merlin'
    | 'percival'
    | 'loyal_servant'
    | 'assassin'
    | 'morgana'
    | 'mordred'
    | 'oberon'
    | 'minion';

export type Team = 'good' | 'evil';

export interface RoleConfig {
    name: string;
    emoji: string;
    team: Team;
    color: string;
    description: string;
    knowledge?: string;
}

export const roleConfig: Record<AvalonRole, RoleConfig> = {
    merlin: {
        name: 'Merlin',
        emoji: '🔮',
        team: 'good',
        color: 'bg-blue-600',
        description: 'Knows the forces of Evil (except Mordred)',
        knowledge: 'You see all Evil players (except Mordred if present). Help Good win without revealing yourself!'
    },
    percival: {
        name: 'Percival',
        emoji: '👁️',
        team: 'good',
        color: 'bg-cyan-600',
        description: 'Sees Merlin and Morgana (cannot distinguish)',
        knowledge: 'You see two powerful wizards. One is Merlin, one is Morgana. Protect Merlin!'
    },
    loyal_servant: {
        name: 'Loyal Servant of Arthur',
        emoji: '⚔️',
        team: 'good',
        color: 'bg-slate-600',
        description: 'No special knowledge, must rely on deduction',
        knowledge: 'You have no special information. Trust your instincts and your allies!'
    },
    assassin: {
        name: 'Assassin',
        emoji: '🗡️',
        team: 'evil',
        color: 'bg-red-700',
        description: 'Can assassinate Merlin if Good wins',
        knowledge: 'You know your Evil allies. If Good wins 3 quests, you can steal victory by identifying Merlin!'
    },
    morgana: {
        name: 'Morgana',
        emoji: '🌙',
        team: 'evil',
        color: 'bg-purple-700',
        description: 'Appears as Merlin to Percival',
        knowledge: 'You appear as Merlin to Percival. Confuse the Good team!'
    },
    mordred: {
        name: 'Mordred',
        emoji: '😈',
        team: 'evil',
        color: 'bg-orange-700',
        description: 'Hidden from Merlin',
        knowledge: 'Merlin cannot see you. Use this advantage wisely!'
    },
    oberon: {
        name: 'Oberon',
        emoji: '👻',
        team: 'evil',
        color: 'bg-gray-700',
        description: 'Unknown to other Evil players',
        knowledge: 'You are alone. You do not know other Evil players, and they do not know you.'
    },
    minion: {
        name: 'Minion of Mordred',
        emoji: '💀',
        team: 'evil',
        color: 'bg-red-900',
        description: 'Knows other Evil players',
        knowledge: 'You know your Evil allies. Work together to sabotage the quests!'
    }
};

export const questSizes: Record<number, number[]> = {
    5:  [2, 3, 2, 3, 3],
    6:  [2, 3, 4, 3, 4],
    7:  [2, 3, 3, 4, 4],
    8:  [3, 4, 4, 5, 5],
    9:  [3, 4, 4, 5, 5],
    10: [3, 4, 4, 5, 5],
};

export function requiresTwoFails(playerCount: number, questNumber: number): boolean {
    return playerCount >= 7 && questNumber === 4;
}
```

### Component: RoleReveal.svelte

```svelte
<script lang="ts">
  import { roleConfig, type AvalonRole } from './roleConfig';
  import RoleCard from './RoleCard.svelte';

  let {
    role,
    team,
    knowledge = [],
    players,
    onAcknowledge
  } = $props<{
    role: AvalonRole;
    team: 'good' | 'evil';
    knowledge: string[];
    players: Array<{ id: string; name: string }>;
    onAcknowledge: () => void;
  }>();

  let acknowledged = $state(false);

  const config = $derived(roleConfig[role]);

  function getPlayerName(id: string): string {
    return players.find(p => p.id === id)?.name || 'Unknown';
  }

  function handleAcknowledge() {
    acknowledged = true;
    onAcknowledge();
  }
</script>

<div class="role-reveal">
  <h2 class="text-center mb-4">Your Role</h2>

  <RoleCard {role} large />

  <div class="team-banner {team}">
    {team === 'good' ? '⚔️ Good Team' : '💀 Evil Team'}
  </div>

  <div class="knowledge-section">
    <h3>What You Know</h3>
    <p class="description">{config.knowledge}</p>

    {#if knowledge.length > 0}
      <div class="known-players">
        {#each knowledge as playerId}
          <div class="player-badge">
            {getPlayerName(playerId)}
          </div>
        {/each}
      </div>
    {/if}
  </div>

  {#if !acknowledged}
    <button
      class="btn-primary w-full mt-6"
      onclick={handleAcknowledge}
    >
      I Understand My Role
    </button>
  {:else}
    <div class="text-center text-muted">
      Waiting for other players...
    </div>
  {/if}
</div>

<style>
  .role-reveal {
    max-width: 480px;
    margin: 0 auto;
    padding: 2rem;
  }

  .team-banner {
    text-align: center;
    padding: 1rem;
    border-radius: 0.5rem;
    font-weight: bold;
    margin: 1rem 0;
  }

  .team-banner.good {
    background: #458588;
    color: white;
  }

  .team-banner.evil {
    background: #cc241d;
    color: white;
  }

  .knowledge-section {
    background: #3c3836;
    padding: 1.5rem;
    border-radius: 0.5rem;
    margin-top: 1.5rem;
  }

  .knowledge-section h3 {
    margin-bottom: 0.75rem;
    color: #d79921;
  }

  .description {
    color: #ebdbb2;
    line-height: 1.6;
  }

  .known-players {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-top: 1rem;
  }

  .player-badge {
    background: #504945;
    padding: 0.5rem 1rem;
    border-radius: 0.25rem;
    color: #fbf1c7;
  }
</style>
```

### Component: QuestBoard.svelte

```svelte
<script lang="ts">
  let {
    currentQuest,
    questResults,
    playerCount
  } = $props<{
    currentQuest: number;
    questResults: Array<{
      quest_number: number;
      success: boolean;
      team_size: number;
      cards: ('success' | 'fail')[];
    }>;
    playerCount: number;
  }>();

  import { questSizes, requiresTwoFails } from './roleConfig';

  const teamSizes = $derived(questSizes[playerCount] || questSizes[5]);

  function getQuestStatus(questNum: number) {
    const result = questResults.find(r => r.quest_number === questNum);
    if (!result) return 'pending';
    return result.success ? 'success' : 'fail';
  }
</script>

<div class="quest-board">
  <h3 class="board-title">Quest Progress</h3>

  <div class="quests">
    {#each [1, 2, 3, 4, 5] as questNum}
      <div
        class="quest"
        class:active={questNum === currentQuest}
        class:success={getQuestStatus(questNum) === 'success'}
        class:fail={getQuestStatus(questNum) === 'fail'}
      >
        <div class="quest-number">Quest {questNum}</div>
        <div class="team-size">
          Team: {teamSizes[questNum - 1]}
          {#if requiresTwoFails(playerCount, questNum)}
            <span class="special">*</span>
          {/if}
        </div>

        {#if getQuestStatus(questNum) !== 'pending'}
          <div class="result-icon">
            {getQuestStatus(questNum) === 'success' ? '✅' : '❌'}
          </div>
        {:else if questNum === currentQuest}
          <div class="current-indicator">▶</div>
        {/if}
      </div>
    {/each}
  </div>

  {#if requiresTwoFails(playerCount, 4)}
    <p class="footnote">* Quest 4 requires 2 FAIL cards to fail</p>
  {/if}
</div>

<style>
  .quest-board {
    background: #282828;
    padding: 1.5rem;
    border-radius: 0.5rem;
    margin-bottom: 1.5rem;
  }

  .board-title {
    text-align: center;
    color: #d79921;
    margin-bottom: 1rem;
  }

  .quests {
    display: flex;
    justify-content: space-between;
    gap: 0.5rem;
  }

  .quest {
    flex: 1;
    background: #3c3836;
    padding: 1rem 0.5rem;
    border-radius: 0.25rem;
    text-align: center;
    border: 2px solid transparent;
    transition: all 0.2s;
  }

  .quest.active {
    border-color: #d79921;
    background: #504945;
  }

  .quest.success {
    background: #458588;
  }

  .quest.fail {
    background: #cc241d;
  }

  .quest-number {
    font-weight: bold;
    font-size: 0.875rem;
    margin-bottom: 0.25rem;
  }

  .team-size {
    font-size: 0.75rem;
    color: #a89984;
  }

  .special {
    color: #fe8019;
    font-weight: bold;
  }

  .result-icon {
    font-size: 1.5rem;
    margin-top: 0.5rem;
  }

  .current-indicator {
    font-size: 1.25rem;
    color: #d79921;
    margin-top: 0.5rem;
  }

  .footnote {
    margin-top: 0.75rem;
    font-size: 0.75rem;
    color: #a89984;
    text-align: center;
  }
</style>
```

### Component: TeamBuilding.svelte

```svelte
<script lang="ts">
  import PlayerSelector from './PlayerSelector.svelte';
  import QuestBoard from './QuestBoard.svelte';

  let {
    isLeader,
    currentQuest,
    questResults,
    playerCount,
    requiredTeamSize,
    players,
    leaderId,
    onProposeTeam
  } = $props<{
    isLeader: boolean;
    currentQuest: number;
    questResults: any[];
    playerCount: number;
    requiredTeamSize: number;
    players: Array<{ id: string; name: string }>;
    leaderId: string;
    onProposeTeam: (selectedIds: string[]) => void;
  }>();

  let selectedPlayers = $state<string[]>([]);

  const leaderName = $derived(
    players.find(p => p.id === leaderId)?.name || 'Unknown'
  );

  function togglePlayer(playerId: string) {
    if (selectedPlayers.includes(playerId)) {
      selectedPlayers = selectedPlayers.filter(id => id !== playerId);
    } else if (selectedPlayers.length < requiredTeamSize) {
      selectedPlayers = [...selectedPlayers, playerId];
    }
  }

  function handlePropose() {
    if (selectedPlayers.length === requiredTeamSize) {
      onProposeTeam(selectedPlayers);
    }
  }

  const canPropose = $derived(
    isLeader && selectedPlayers.length === requiredTeamSize
  );
</script>

<div class="team-building">
  <QuestBoard {currentQuest} {questResults} {playerCount} />

  <div class="leader-banner">
    👑 {leaderName} is the Leader
  </div>

  {#if isLeader}
    <div class="instructions">
      <h3>Select {requiredTeamSize} players for the quest</h3>
      <p class="counter">
        {selectedPlayers.length} / {requiredTeamSize} selected
      </p>
    </div>

    <PlayerSelector
      {players}
      {selectedPlayers}
      onToggle={togglePlayer}
    />

    <button
      class="btn-primary w-full mt-4"
      disabled={!canPropose}
      onclick={handlePropose}
    >
      Propose Team
    </button>
  {:else}
    <div class="waiting">
      <p>Waiting for {leaderName} to propose a team...</p>
    </div>
  {/if}
</div>

<style>
  .team-building {
    max-width: 600px;
    margin: 0 auto;
    padding: 1rem;
  }

  .leader-banner {
    background: #d79921;
    color: #282828;
    text-align: center;
    padding: 1rem;
    border-radius: 0.5rem;
    font-weight: bold;
    margin-bottom: 1.5rem;
  }

  .instructions h3 {
    color: #d79921;
    margin-bottom: 0.5rem;
  }

  .counter {
    color: #a89984;
    font-size: 0.875rem;
  }

  .waiting {
    text-align: center;
    padding: 3rem 1rem;
    color: #a89984;
  }
</style>
```

---

## State Management

### Game Store Integration

```typescript
// Add to existing game.svelte.ts
export interface AvalonPlayerState {
    role: string;
    team: 'good' | 'evil';
    knowledge: string[];  // player IDs you know
    hasVoted: boolean;
    hasPlayedCard: boolean;
    hasAcknowledged: boolean;
}

export interface AvalonPublicState {
    phase: string;
    currentLeader: string;
    questNumber: number;
    questResults: QuestResult[];
    proposedTeam: string[];
    votesSubmitted: number;
    cardsSubmitted: number;
    rejectionCount: number;
}

// Event handlers
function handleAvalonEvent(event: GameEvent) {
    switch (event.type) {
        case 'role_assigned':
            playerState.role = event.payload.role;
            playerState.team = event.payload.team;
            break;

        case 'role_knowledge':
            playerState.knowledge = event.payload.known_players;
            break;

        case 'team_proposed':
            publicState.proposedTeam = event.payload.team_members;
            break;

        case 'quest_completed':
            publicState.questResults.push(event.payload);
            break;

        // ... handle all events
    }
}
```

---

## Implementation Checklist

### Backend Tasks
- [ ] Create `backend/internal/games/avalon/` directory
- [ ] Implement state types (`state.go`)
- [ ] Implement role logic and knowledge (`roles.go`)
- [ ] Implement quest configuration (`quests.go`)
- [ ] Implement config validation (`config.go`)
- [ ] Implement voting logic (`voting.go`)
- [ ] Implement phase transitions (`phases.go`)
- [ ] Implement `Game` interface (`game.go`)
  - [ ] `Initialize()`
  - [ ] `ValidateAction()`
  - [ ] `ProcessAction()`
  - [ ] `GetPlayerState()`
  - [ ] `GetPublicState()`
  - [ ] `GetPhase()`
  - [ ] `IsFinished()`
  - [ ] `GetResults()`
- [ ] Register game in `games/registry.go`
- [ ] Write unit tests for role knowledge
- [ ] Write unit tests for voting logic
- [ ] Write unit tests for quest resolution

### Frontend Tasks
- [ ] Create `frontend/src/lib/games/avalon/` directory
- [ ] Implement `roleConfig.ts` with all roles
- [ ] Implement `RoleCard.svelte` component
- [ ] Implement `RoleReveal.svelte` component
- [ ] Implement `QuestBoard.svelte` component
- [ ] Implement `PlayerSelector.svelte` component
- [ ] Implement `TeamBuilding.svelte` component
- [ ] Implement `TeamVoting.svelte` component
- [ ] Implement `QuestExecution.svelte` component
- [ ] Implement `QuestResults.svelte` component
- [ ] Implement `Assassination.svelte` component
- [ ] Implement `Results.svelte` component
- [ ] Implement `AvalonGame.svelte` main container
- [ ] Add event handlers to game store
- [ ] Add Avalon to game selection UI

### Testing Tasks
- [ ] Test role assignment randomization
- [ ] Test role knowledge (Merlin, Percival, Evil)
- [ ] Test team building and validation
- [ ] Test team voting (approve/reject)
- [ ] Test quest card submissions (good can only play success)
- [ ] Test quest resolution (1 fail vs 2 fail for quest 4)
- [ ] Test 5 rejections auto-loss
- [ ] Test assassination phase (Merlin identification)
- [ ] Test all win conditions
- [ ] Multiplayer playtest with 5-10 players
- [ ] Mobile UI testing

### Polish Tasks
- [ ] Add animations for quest results
- [ ] Add sound effects (optional)
- [ ] Add "How to Play" guide
- [ ] Add role reference sheet
- [ ] Optimize for mobile (touch targets, one-handed use)
- [ ] Test reconnection handling
- [ ] Performance testing with 10 players

---

## Estimated Timeline

- **Backend**: 2 days (game logic, voting, quests, roles)
- **Frontend**: 2-3 days (UI components, quest board, team selection)
- **Testing**: 1 day (unit tests, multiplayer testing)
- **Polish**: 0.5 days (animations, mobile optimization)

**Total**: 5-6 days

---

## Notes

### Key Differences from Werewolf
- **Team-based** instead of individual roles
- **Multi-round** (5 quests) instead of single night
- **Leader rotation** mechanic
- **Voting on teams** before quest execution
- **Assassination endgame** mechanic

### Shared Patterns
- Role assignment with hidden information ✅
- Role reveal phase ✅
- Event sourcing architecture ✅
- Mobile-first UI ✅
- In-person social focus ✅

### Technical Challenges
1. **Quest board visualization** - needs to be clear and prominent
2. **Team selection UX** - must be intuitive on mobile
3. **Vote timing** - all players must vote before reveal
4. **Assassination phase** - only triggers conditionally
5. **5 rejection rule** - edge case win condition

---

**Document Version**: 1.0
**Last Updated**: 2025-11-19
**Status**: Ready for Implementation
