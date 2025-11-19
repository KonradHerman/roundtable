# Spyfall - Prototype Plan

> **Game Type**: Social Deduction, Question & Answer, Location Guessing
> **Players**: 3-8
> **Duration**: 6-8 minutes per round
> **Priority**: HIGH - Phase 2 validation game

---

## Table of Contents
1. [Game Rules Reference](#game-rules-reference)
2. [Design Principles Application](#design-principles-application)
3. [Complete Event Flow](#complete-event-flow)
4. [Phase Breakdown](#phase-breakdown)
5. [Backend Implementation](#backend-implementation)
6. [Frontend Implementation](#frontend-implementation)
7. [Location Database](#location-database)
8. [State Management](#state-management)
9. [Implementation Checklist](#implementation-checklist)

---

## Game Rules Reference

### Overview
Spyfall is a conversation-based social deduction game. All players are at the same location except one **spy**. Everyone knows the location and their role at that location except the spy. Players ask each other questions to figure out who the spy is without revealing the location. The spy tries to blend in and guess the location.

### Setup
1. **Random location** selected from database (30-50 locations)
2. **One random player** assigned as the spy
3. **All other players** receive the location and a specific role at that location
4. **Timer** starts (configurable: 6-8 minutes)

### Game Flow

#### Phase 1: Role Reveal (Private)
- **Spy** sees: "You are the SPY! Figure out the location."
- **Non-spies** see: Location name + their specific role
  - Example: "Beach - Lifeguard"

#### Phase 2: Question Round
- **First player** (random or clockwise from dealer) starts
- **Active player**:
  - Asks any other player one question
  - Questions should help identify the spy without revealing location
  - Example: "Do you enjoy working here?" "Is it usually crowded?"
- **Answering player**:
  - Answers the question (trying to prove they're not the spy)
  - Becomes the next active player (asks the next question)
- **Timer** runs continuously (6-8 minutes)
- **Turn order** is dynamic (each answerer becomes next questioner)

#### Phase 3: Accusation (Optional)
- **Any player** can call for an accusation vote at any time
- **All players** vote on who they think is the spy
- If **majority** agrees on one player:
  - **If accused is the spy** → Non-spies win
  - **If accused is NOT the spy** → Spy gets one guess at location
    - If spy guesses correctly → Spy wins
    - If spy guesses incorrectly → Non-spies win
- If **no majority** or **tie** → game continues

#### Phase 4: Time Expires
- If timer runs out before accusation:
  - **Spy** can guess the location
  - If correct → Spy wins
  - If incorrect → Non-spies win

### Scoring (Optional)
- **Spy Wins**: Spy gets points equal to number of players
- **Non-Spies Win**: Each non-spy gets 1 point, +1 if they voted for spy

### Win Conditions
1. **Spy Wins If**:
   - Timer expires and spy correctly guesses location
   - Accused but correctly guesses location

2. **Non-Spies Win If**:
   - Correctly identify and vote out the spy
   - Spy guesses location incorrectly

---

## Design Principles Application

### 1. Enhancing In-Person Interactions ✅
**Digital**:
- Role and location assignment (prevents shuffling/peeking)
- Timer management (start/pause/extend)
- Location database (no physical cards needed)
- Voting collection (simultaneous votes)
- Spy location guess (hidden from others)

**Physical**:
- Question asking and answering (voice only)
- Reading body language and tells
- Social deduction through conversation
- Bluffing and misdirection
- Dramatic accusations

### 2. Tracking and Automating Card States ✅
- **Location cards**: Digital assignment ensures no peeking
- **Role cards**: Each player at location has unique role
- **Spy card**: Secret assignment with no tells
- **Timer**: Automated countdown with controls
- **Voting**: Simultaneous submission prevents bandwagoning

### 3. Preserving Social Core
- App handles **assignment** and **timing**
- Players handle **interrogation** and **deduction**
- No AI hints or assistance
- Focus on conversation, not UI

### 4. Mobile-First UX
- Location always visible (for non-spies)
- Quick reference location list (for spy)
- Large vote buttons
- Prominent timer display
- One-handed friendly interface

---

## Complete Event Flow

### Event Types & Visibility

#### Public Events (Everyone sees)
```typescript
// Setup & Phase Changes
EVENT_GAME_STARTED          // Game initialized
EVENT_PHASE_CHANGED         // Phase transition
EVENT_GAME_FINISHED         // Game concluded

// Timer Management
EVENT_TIMER_STARTED         // Timer began
EVENT_TIMER_PAUSED          // Timer paused
EVENT_TIMER_EXTENDED        // Time added (e.g., +30 seconds)
EVENT_TIMER_EXPIRED         // Time ran out

// Turn Management
EVENT_TURN_CHANGED          // New player's turn to ask question
EVENT_QUESTION_ASKED        // Question asked (text + asker + answerer)
EVENT_ANSWER_GIVEN          // Answer given (text)

// Accusation
EVENT_ACCUSATION_STARTED    // Someone called for accusation vote
EVENT_ACCUSATION_VOTE_CAST  // Someone voted (not their vote)
EVENT_ACCUSATION_RESOLVED   // Vote results revealed
EVENT_SPY_REVEALED          // Spy identity revealed

// Location Guess
EVENT_SPY_GUESSES_LOCATION  // Spy made location guess
EVENT_GUESS_RESULT          // Guess was correct/incorrect
```

#### Private Events (Specific players only)
```typescript
// Role Assignment
location_assigned           // Your location + role (non-spies only)
spy_assigned                // You are the spy (spy only)

// Voting
accusation_vote_recorded    // Your vote was recorded (confirmation)

// Location Guess
location_list               // All possible locations (spy only, at start)
```

---

## Phase Breakdown

### Phase 1: SETUP
**Duration**: Until game starts
**Transitions to**: ROLE_REVEAL

**Actions**:
- Host selects timer duration (6, 7, or 8 minutes)
- Host starts game
- System randomly selects location
- System randomly selects one player as spy
- System assigns roles to non-spies

**Events**:
- `EVENT_GAME_STARTED`
- `spy_assigned` (private to spy)
- `location_assigned` (private to each non-spy)
- `location_list` (private to spy - all possible locations)
- `EVENT_PHASE_CHANGED` → ROLE_REVEAL

**Backend State**:
```go
type SpyfallState struct {
    Players         []*core.Player
    Location        *Location
    SpyID           string
    Roles           map[string]string  // playerID -> role
    TimerDuration   int                // seconds
    TimerStarted    time.Time
    TimerPaused     bool
    TimeRemaining   int                // seconds
    CurrentTurn     string             // playerID
    Phase           GamePhase
}

type Location struct {
    Name  string
    Roles []string  // e.g., ["Lifeguard", "Surfer", "Ice Cream Vendor"]
}
```

---

### Phase 2: ROLE_REVEAL
**Duration**: Until all players acknowledge
**Transitions to**: QUESTION_ROUND

**UI Display**:

**For Spy**:
```
┌─────────────────────────┐
│    You are the SPY! 🕵️   │
├─────────────────────────┤
│ Figure out the location │
│ by asking questions.    │
│                         │
│ Tap here to see all     │
│ possible locations ↓    │
└─────────────────────────┘
```

**For Non-Spies**:
```
┌─────────────────────────┐
│      🏖️  Beach           │
├─────────────────────────┤
│   Your Role: Lifeguard  │
│                         │
│ Ask questions to find   │
│ the spy without         │
│ revealing the location! │
└─────────────────────────┘
```

**Actions**:
- Players view their role
- Players acknowledge role
- Auto-advance when all acknowledged

**Events**:
- `player_acknowledged_role` (public count only)
- `EVENT_PHASE_CHANGED` → QUESTION_ROUND

---

### Phase 3: QUESTION_ROUND
**Duration**: Until accusation or timer expires
**Transitions to**: ACCUSATION or LOCATION_GUESS

**UI Display**:

**Header**:
```
┌─────────────────────────────┐
│  ⏱️  5:32 remaining          │
│  📍 Beach (Your location)   │  ← (not shown to spy)
│  🎭 Lifeguard (Your role)   │  ← (not shown to spy)
└─────────────────────────────┘
```

**Turn Indicator**:
```
┌─────────────────────────────┐
│  Current Turn: Alice         │
│  Alice is asking Bob         │
└─────────────────────────────┘
```

**Question Log** (Recent 3-5 questions):
```
Alice → Bob: "Do you enjoy working here?"
Bob: "It depends on the weather."

Bob → Charlie: "Is it crowded today?"
Charlie: "Usually, especially on weekends."
```

**Action Buttons**:
```
┌──────────────────┬──────────────────┐
│  📍 Location     │  ⏸️  Pause Timer  │  ← (host only)
│     List         │                  │
├──────────────────┴──────────────────┤
│       🔍 Call for Accusation        │
└─────────────────────────────────────┘
```

**Actions**:
- **Current turn player**: Selects another player, types question, submits
- **Answering player**: Types answer, submits (becomes next turn)
- **Any player**: Can call for accusation vote
- **Host**: Can pause/resume/extend timer

**Events**:
- `EVENT_TURN_CHANGED`
- `EVENT_QUESTION_ASKED`
- `EVENT_ANSWER_GIVEN`
- `EVENT_TIMER_PAUSED` / `EVENT_TIMER_RESUMED`
- `EVENT_TIMER_EXTENDED`
- `EVENT_ACCUSATION_STARTED` (if player calls for accusation)

**Timer Behavior**:
- Counts down continuously
- Host can pause (e.g., for bathroom break)
- Host can extend (+30 sec, +1 min, +2 min)
- Auto-transitions to LOCATION_GUESS when timer expires

**Transition Triggers**:
- **Timer expires** → LOCATION_GUESS (spy guesses)
- **Player calls accusation** → ACCUSATION

---

### Phase 4: ACCUSATION
**Duration**: Until all players vote
**Transitions to**: LOCATION_GUESS (if non-spy accused) or GAME_FINISHED

**UI Display**:
```
┌─────────────────────────────────┐
│   🔍 Accusation Vote             │
├─────────────────────────────────┤
│ Who do you think is the spy?    │
│                                 │
│  [Alice]  [Bob]  [Charlie]      │
│  [David]  [Eve]  [Frank]        │
│                                 │
│ Votes submitted: 4/6            │
└─────────────────────────────────┘
```

**Actions**:
- All players vote simultaneously
- Votes hidden until all submitted
- System tallies votes
- Majority required to accuse

**Events**:
- `EVENT_ACCUSATION_VOTE_CAST` (public: player voted)
- `accusation_vote_recorded` (private: confirmation)
- `EVENT_ACCUSATION_RESOLVED`
  - Payload: vote breakdown (player → vote count)
- `EVENT_SPY_REVEALED` (if spy identified)

**Vote Resolution**:
- **Majority votes for spy** → EVENT_GAME_FINISHED (non-spies win)
- **Majority votes for non-spy** → LOCATION_GUESS (spy gets guess)
- **No majority / tie** → QUESTION_ROUND (game continues)

---

### Phase 5: LOCATION_GUESS
**Duration**: Until spy guesses or timer expires
**Transitions to**: GAME_FINISHED

**Trigger Conditions**:
1. Timer expired
2. Accusation vote identified non-spy

**UI Display** (Spy only):
```
┌─────────────────────────────────┐
│   Guess the Location! 🕵️         │
├─────────────────────────────────┤
│ Select the location:            │
│                                 │
│  🏖️  Beach                       │
│  🏥 Hospital                     │
│  ✈️  Airplane                    │
│  🎭 Theater                      │
│  🏫 School                       │
│  ... (scroll for more)          │
│                                 │
│  [Confirm Guess]                │
└─────────────────────────────────┘
```

**UI Display** (Others):
```
┌─────────────────────────────────┐
│   The spy is guessing the       │
│   location...                   │
└─────────────────────────────────┘
```

**Actions**:
- Spy selects location from full list
- Spy confirms guess
- System checks if correct

**Events**:
- `EVENT_SPY_GUESSES_LOCATION` (public: location guessed)
- `EVENT_GUESS_RESULT` (correct/incorrect)
- `EVENT_GAME_FINISHED`

**Win Determination**:
- **Correct guess** → Spy wins
- **Incorrect guess** → Non-spies win

---

### Phase 6: GAME_FINISHED
**Duration**: Indefinite
**Transitions to**: New round (if host chooses)

**UI Display**:
```
┌─────────────────────────────────┐
│   🎉 Non-Spies Win! 🎉           │
├─────────────────────────────────┤
│ The location was: Beach 🏖️       │
│ The spy was: Bob 🕵️              │
│                                 │
│ Bob guessed: Hospital ❌        │
├─────────────────────────────────┤
│ Roles:                          │
│  • Alice - Lifeguard            │
│  • Bob - SPY                    │
│  • Charlie - Surfer             │
│  • David - Ice Cream Vendor     │
└─────────────────────────────────┘

[Play Another Round]  [Change Game]
```

**Events**:
- None (terminal phase)

**Actions**:
- Display winning team
- Reveal all roles and location
- Explain win condition
- Host can start new round (new location, new spy)

---

## Backend Implementation

### File Structure
```
backend/internal/games/spyfall/
├── game.go           # Implements core.Game interface
├── state.go          # Spyfall-specific state types
├── config.go         # Game configuration (timer duration)
├── phases.go         # Phase transition logic
├── locations.go      # Location database and selection
├── timer.go          # Timer management logic
└── voting.go         # Accusation voting logic
```

### State Types

```go
// state.go
package spyfall

type GamePhase string

const (
    PhaseSetup         GamePhase = "setup"
    PhaseRoleReveal    GamePhase = "role_reveal"
    PhaseQuestionRound GamePhase = "question_round"
    PhaseAccusation    GamePhase = "accusation"
    PhaseLocationGuess GamePhase = "location_guess"
    PhaseFinished      GamePhase = "finished"
)

type Location struct {
    Name  string   `json:"name"`
    Icon  string   `json:"icon"`  // emoji
    Roles []string `json:"roles"` // possible roles at this location
}

type QuestionAnswer struct {
    AskerID   string `json:"asker_id"`
    AnswererID string `json:"answerer_id"`
    Question  string `json:"question"`
    Answer    string `json:"answer"`
    Timestamp int64  `json:"timestamp"`
}

type SpyfallState struct {
    // Core game state
    Players   []*core.Player
    Phase     GamePhase
    Config    *SpyfallConfig

    // Role assignment
    Location  *Location
    SpyID     string
    Roles     map[string]string  // playerID -> role at location

    // Timer
    TimerDuration  int              // total seconds
    TimerStartTime time.Time
    TimerPaused    bool
    PausedAt       time.Time
    PausedDuration time.Duration
    TimeExtended   int              // extra seconds added

    // Question round
    CurrentTurn     string           // playerID whose turn to ask
    QuestionHistory []QuestionAnswer

    // Accusation
    AccusationVotes map[string]string  // voterID -> suspectID
    AccusedPlayerID string

    // Location guess
    SpyGuess        string           // location name

    // Acknowledgments
    Acknowledged    map[string]bool  // playerID -> acknowledged
}

type SpyfallConfig struct {
    TimerMinutes    int    `json:"timer_minutes"`    // 6, 7, or 8
    AllowExtensions bool   `json:"allow_extensions"` // can host extend time?
    LocationSet     string `json:"location_set"`     // "classic", "expanded", "custom"
}
```

### Location Database

```go
// locations.go
package spyfall

var ClassicLocations = []Location{
    {
        Name: "Airplane",
        Icon: "✈️",
        Roles: []string{
            "First Class Passenger",
            "Economy Passenger",
            "Flight Attendant",
            "Pilot",
            "Co-Pilot",
            "Air Marshal",
            "Mechanic",
        },
    },
    {
        Name: "Bank",
        Icon: "🏦",
        Roles: []string{
            "Manager",
            "Teller",
            "Security Guard",
            "Customer",
            "Robber",
            "Consultant",
            "Armored Truck Driver",
        },
    },
    {
        Name: "Beach",
        Icon: "🏖️",
        Roles: []string{
            "Lifeguard",
            "Surfer",
            "Ice Cream Vendor",
            "Sunbather",
            "Beach Volleyball Player",
            "Photographer",
            "Kite Surfer",
        },
    },
    {
        Name: "Casino",
        Icon: "🎰",
        Roles: []string{
            "Dealer",
            "Gambler",
            "Bartender",
            "Security",
            "Pit Boss",
            "Waitress",
            "High Roller",
        },
    },
    {
        Name: "Circus",
        Icon: "🎪",
        Roles: []string{
            "Acrobat",
            "Clown",
            "Magician",
            "Animal Trainer",
            "Juggler",
            "Ringmaster",
            "Ticket Seller",
        },
    },
    {
        Name: "Hospital",
        Icon: "🏥",
        Roles: []string{
            "Doctor",
            "Nurse",
            "Patient",
            "Surgeon",
            "Anesthesiologist",
            "Receptionist",
            "Paramedic",
        },
    },
    {
        Name: "Hotel",
        Icon: "🏨",
        Roles: []string{
            "Guest",
            "Receptionist",
            "Bellhop",
            "Housekeeper",
            "Concierge",
            "Security",
            "Manager",
        },
    },
    {
        Name: "Military Base",
        Icon: "🪖",
        Roles: []string{
            "Soldier",
            "Commander",
            "Medic",
            "Engineer",
            "Sniper",
            "Cook",
            "Drill Sergeant",
        },
    },
    {
        Name: "Movie Studio",
        Icon: "🎬",
        Roles: []string{
            "Director",
            "Actor",
            "Cameraman",
            "Producer",
            "Stunt Double",
            "Sound Engineer",
            "Costume Designer",
        },
    },
    {
        Name: "Ocean Liner",
        Icon: "🚢",
        Roles: []string{
            "Captain",
            "Passenger",
            "Bartender",
            "Musician",
            "Waiter",
            "Cook",
            "Engineer",
        },
    },
    {
        Name: "Passenger Train",
        Icon: "🚂",
        Roles: []string{
            "Conductor",
            "Passenger",
            "Engineer",
            "Ticket Inspector",
            "Waiter",
            "Mechanic",
            "Stoker",
        },
    },
    {
        Name: "Pirate Ship",
        Icon: "🏴‍☠️",
        Roles: []string{
            "Captain",
            "First Mate",
            "Cook",
            "Sailor",
            "Prisoner",
            "Lookout",
            "Cabin Boy",
        },
    },
    {
        Name: "Polar Station",
        Icon: "🧊",
        Roles: []string{
            "Scientist",
            "Researcher",
            "Medic",
            "Geologist",
            "Meteorologist",
            "Cook",
            "Radio Operator",
        },
    },
    {
        Name: "Restaurant",
        Icon: "🍽️",
        Roles: []string{
            "Chef",
            "Waiter",
            "Customer",
            "Sommelier",
            "Host",
            "Dishwasher",
            "Food Critic",
        },
    },
    {
        Name: "School",
        Icon: "🏫",
        Roles: []string{
            "Teacher",
            "Student",
            "Principal",
            "Janitor",
            "Librarian",
            "Gym Teacher",
            "Lunch Lady",
        },
    },
    {
        Name: "Space Station",
        Icon: "🛰️",
        Roles: []string{
            "Commander",
            "Astronaut",
            "Scientist",
            "Engineer",
            "Doctor",
            "Alien",
            "Space Tourist",
        },
    },
    {
        Name: "Submarine",
        Icon: "🔱",
        Roles: []string{
            "Captain",
            "Sonar Technician",
            "Cook",
            "Sailor",
            "Radioman",
            "Navigator",
            "Torpedo Operator",
        },
    },
    {
        Name: "Supermarket",
        Icon: "🛒",
        Roles: []string{
            "Cashier",
            "Customer",
            "Manager",
            "Stocker",
            "Baker",
            "Butcher",
            "Security",
        },
    },
    {
        Name: "Theater",
        Icon: "🎭",
        Roles: []string{
            "Actor",
            "Director",
            "Prompter",
            "Audience Member",
            "Ticket Seller",
            "Coat Check",
            "Spotlight Operator",
        },
    },
    {
        Name: "University",
        Icon: "🎓",
        Roles: []string{
            "Professor",
            "Student",
            "Dean",
            "Librarian",
            "Janitor",
            "Graduate Student",
            "Researcher",
        },
    },
};

// Random location selection with crypto/rand
func SelectRandomLocation(locations []Location) *Location {
    if len(locations) == 0 {
        return nil
    }

    // Use crypto/rand for unpredictable selection
    n, err := rand.Int(rand.Reader, big.NewInt(int64(len(locations))))
    if err != nil {
        // Fallback to math/rand if crypto fails
        return &locations[mathrand.Intn(len(locations))]
    }

    return &locations[n.Int64()]
}

// Assign random role to each player (except spy)
func (g *SpyfallGame) AssignRoles() {
    roles := make([]string, len(g.state.Location.Roles))
    copy(roles, g.state.Location.Roles)

    // Shuffle roles
    secureShuffleRoles(roles)

    // Assign to non-spy players
    roleIndex := 0
    for _, player := range g.state.Players {
        if player.ID != g.state.SpyID {
            g.state.Roles[player.ID] = roles[roleIndex]
            roleIndex++
        }
    }
}
```

### Timer Management

```go
// timer.go
package spyfall

func (g *SpyfallGame) StartTimer() {
    g.state.TimerStartTime = time.Now()
    g.state.TimerPaused = false
}

func (g *SpyfallGame) PauseTimer() {
    if !g.state.TimerPaused {
        g.state.TimerPaused = true
        g.state.PausedAt = time.Now()
    }
}

func (g *SpyfallGame) ResumeTimer() {
    if g.state.TimerPaused {
        pauseDuration := time.Since(g.state.PausedAt)
        g.state.PausedDuration += pauseDuration
        g.state.TimerPaused = false
    }
}

func (g *SpyfallGame) ExtendTimer(seconds int) {
    g.state.TimeExtended += seconds
}

func (g *SpyfallGame) GetTimeRemaining() int {
    totalDuration := time.Duration(g.state.TimerDuration + g.state.TimeExtended) * time.Second

    elapsed := time.Since(g.state.TimerStartTime) - g.state.PausedDuration

    if g.state.TimerPaused {
        elapsed = g.state.PausedAt.Sub(g.state.TimerStartTime) - g.state.PausedDuration
    }

    remaining := totalDuration - elapsed

    if remaining < 0 {
        return 0
    }

    return int(remaining.Seconds())
}

func (g *SpyfallGame) IsTimerExpired() bool {
    return g.GetTimeRemaining() <= 0
}

// Called periodically to check for timer expiration
func (g *SpyfallGame) CheckPhaseTimeout() ([]core.GameEvent, error) {
    if g.state.Phase == PhaseQuestionRound && g.IsTimerExpired() {
        // Timer expired - move to location guess
        return g.AdvanceToLocationGuess()
    }
    return nil, nil
}
```

### Action Processing

```go
// game.go
package spyfall

func (g *SpyfallGame) ProcessAction(playerID string, action core.Action) ([]core.GameEvent, error) {
    switch action.Type {
    case "acknowledge_role":
        return g.processAcknowledgeRole(playerID)

    case "ask_question":
        return g.processAskQuestion(playerID, action.Payload)

    case "answer_question":
        return g.processAnswerQuestion(playerID, action.Payload)

    case "call_accusation":
        return g.processCallAccusation(playerID)

    case "vote_accusation":
        return g.processAccusationVote(playerID, action.Payload)

    case "guess_location":
        return g.processLocationGuess(playerID, action.Payload)

    case "pause_timer":
        return g.processPauseTimer(playerID)

    case "resume_timer":
        return g.processResumeTimer(playerID)

    case "extend_timer":
        return g.processExtendTimer(playerID, action.Payload)

    default:
        return nil, fmt.Errorf("unknown action type: %s", action.Type)
    }
}

func (g *SpyfallGame) processAskQuestion(playerID string, payload map[string]interface{}) ([]core.GameEvent, error) {
    // Validate it's player's turn
    if g.state.CurrentTurn != playerID {
        return nil, fmt.Errorf("not your turn to ask")
    }

    question := payload["question"].(string)
    answererID := payload["answerer_id"].(string)

    // Validate answerer exists
    if !g.playerExists(answererID) {
        return nil, fmt.Errorf("invalid answerer")
    }

    // Record question
    qa := QuestionAnswer{
        AskerID:    playerID,
        AnswererID: answererID,
        Question:   question,
        Timestamp:  time.Now().Unix(),
    }
    g.state.QuestionHistory = append(g.state.QuestionHistory, qa)

    // Create event
    event, _ := core.NewPublicEvent(
        "question_asked",
        playerID,
        map[string]interface{}{
            "asker_id":    playerID,
            "answerer_id": answererID,
            "question":    question,
        },
    )

    return []core.GameEvent{event}, nil
}
```

---

## Frontend Implementation

### File Structure
```
frontend/src/lib/games/spyfall/
├── SpyfallGame.svelte           # Main game container
├── RoleReveal.svelte            # Show location/role or spy status
├── QuestionRound.svelte         # Main game phase with timer
├── Accusation.svelte            # Voting phase
├── LocationGuess.svelte         # Spy guesses location
├── Results.svelte               # Final results
├── TimerDisplay.svelte          # Countdown timer
├── QuestionLog.svelte           # Recent Q&A history
├── LocationList.svelte          # Reference sheet (for spy)
├── PlayerGrid.svelte            # Reusable player selection
└── locationConfig.ts            # Location database (frontend copy)
```

### Location Config (Frontend)

```typescript
// locationConfig.ts
export interface Location {
    name: string;
    icon: string;
    roles: string[];
}

export const locations: Location[] = [
    {
        name: "Airplane",
        icon: "✈️",
        roles: [
            "First Class Passenger",
            "Economy Passenger",
            "Flight Attendant",
            "Pilot",
            "Co-Pilot",
            "Air Marshal",
            "Mechanic"
        ]
    },
    {
        name: "Bank",
        icon: "🏦",
        roles: [
            "Manager",
            "Teller",
            "Security Guard",
            "Customer",
            "Robber",
            "Consultant",
            "Armored Truck Driver"
        ]
    },
    {
        name: "Beach",
        icon: "🏖️",
        roles: [
            "Lifeguard",
            "Surfer",
            "Ice Cream Vendor",
            "Sunbather",
            "Beach Volleyball Player",
            "Photographer",
            "Kite Surfer"
        ]
    },
    // ... (all 20+ locations)
];

export function getLocationByName(name: string): Location | undefined {
    return locations.find(loc => loc.name === name);
}
```

### Component: RoleReveal.svelte

```svelte
<script lang="ts">
  import { getLocationByName, type Location } from './locationConfig';

  let {
    isSpy,
    location,
    role,
    allLocations,
    onAcknowledge
  } = $props<{
    isSpy: boolean;
    location?: string;
    role?: string;
    allLocations: string[];
    onAcknowledge: () => void;
  }>();

  let showingLocationList = $state(false);
  let acknowledged = $state(false);

  const locationData = $derived(
    location ? getLocationByName(location) : null
  );

  function handleAcknowledge() {
    acknowledged = true;
    onAcknowledge();
  }
</script>

<div class="role-reveal">
  {#if isSpy}
    <!-- SPY VIEW -->
    <div class="spy-card">
      <div class="spy-icon">🕵️</div>
      <h1>You are the SPY!</h1>
      <p>Figure out the location by asking questions.</p>
      <p>Blend in and don't reveal yourself!</p>

      <button
        class="btn-secondary mt-4"
        onclick={() => showingLocationList = !showingLocationList}
      >
        {showingLocationList ? 'Hide' : 'Show'} Possible Locations
      </button>

      {#if showingLocationList}
        <div class="location-list">
          {#each allLocations as loc}
            <div class="location-item">{loc}</div>
          {/each}
        </div>
      {/if}
    </div>
  {:else}
    <!-- NON-SPY VIEW -->
    <div class="location-card">
      <div class="location-icon">{locationData?.icon}</div>
      <h1>{location}</h1>
      <div class="role-badge">Your Role: {role}</div>

      <div class="instructions">
        <p>Ask questions to find the spy without revealing the location!</p>
      </div>
    </div>
  {/if}

  {#if !acknowledged}
    <button
      class="btn-primary w-full mt-6"
      onclick={handleAcknowledge}
    >
      I'm Ready
    </button>
  {:else}
    <div class="waiting-text">
      Waiting for other players...
    </div>
  {/if}
</div>

<style>
  .role-reveal {
    max-width: 500px;
    margin: 0 auto;
    padding: 2rem;
  }

  .spy-card, .location-card {
    background: #282828;
    padding: 2rem;
    border-radius: 1rem;
    text-align: center;
  }

  .spy-icon {
    font-size: 4rem;
    margin-bottom: 1rem;
  }

  .spy-card h1 {
    color: #cc241d;
    margin-bottom: 1rem;
  }

  .location-icon {
    font-size: 5rem;
    margin-bottom: 1rem;
  }

  .location-card h1 {
    color: #d79921;
    margin-bottom: 1rem;
  }

  .role-badge {
    background: #458588;
    color: white;
    padding: 0.75rem 1.5rem;
    border-radius: 0.5rem;
    display: inline-block;
    margin-top: 1rem;
    font-weight: bold;
  }

  .instructions {
    margin-top: 1.5rem;
    color: #a89984;
    line-height: 1.6;
  }

  .location-list {
    margin-top: 1rem;
    max-height: 300px;
    overflow-y: auto;
    background: #3c3836;
    border-radius: 0.5rem;
    padding: 1rem;
  }

  .location-item {
    padding: 0.5rem;
    border-bottom: 1px solid #504945;
    color: #ebdbb2;
  }

  .location-item:last-child {
    border-bottom: none;
  }

  .waiting-text {
    text-align: center;
    color: #a89984;
    margin-top: 1rem;
  }
</style>
```

### Component: QuestionRound.svelte

```svelte
<script lang="ts">
  import TimerDisplay from './TimerDisplay.svelte';
  import QuestionLog from './QuestionLog.svelte';
  import LocationList from './LocationList.svelte';

  let {
    isSpy,
    location,
    role,
    timeRemaining,
    timerActive,
    currentTurn,
    isMyTurn,
    players,
    questionHistory,
    isHost,
    onAskQuestion,
    onAnswerQuestion,
    onCallAccusation,
    onPauseTimer,
    onResumeTimer,
    onExtendTimer
  } = $props<{
    isSpy: boolean;
    location?: string;
    role?: string;
    timeRemaining: number;
    timerActive: boolean;
    currentTurn: string;
    isMyTurn: boolean;
    players: Array<{ id: string; name: string }>;
    questionHistory: Array<{
      asker_id: string;
      answerer_id: string;
      question: string;
      answer?: string;
    }>;
    isHost: boolean;
    onAskQuestion: (answererId: string, question: string) => void;
    onAnswerQuestion: (answer: string) => void;
    onCallAccusation: () => void;
    onPauseTimer: () => void;
    onResumeTimer: () => void;
    onExtendTimer: (seconds: number) => void;
  }>();

  let showLocationList = $state(false);
  let selectedAnswerer = $state<string | null>(null);
  let questionText = $state('');
  let answerText = $state('');

  const currentTurnName = $derived(
    players.find(p => p.id === currentTurn)?.name || 'Unknown'
  );

  const lastQuestion = $derived(
    questionHistory.length > 0 ? questionHistory[questionHistory.length - 1] : null
  );

  const waitingForAnswer = $derived(
    lastQuestion && !lastQuestion.answer
  );

  const canAskQuestion = $derived(
    isMyTurn && !waitingForAnswer
  );

  const canAnswerQuestion = $derived(
    waitingForAnswer && lastQuestion?.answerer_id === currentPlayerId
  );

  function handleAskQuestion() {
    if (selectedAnswerer && questionText.trim()) {
      onAskQuestion(selectedAnswerer, questionText);
      selectedAnswerer = null;
      questionText = '';
    }
  }

  function handleAnswerQuestion() {
    if (answerText.trim()) {
      onAnswerQuestion(answerText);
      answerText = '';
    }
  }
</script>

<div class="question-round">
  <!-- Header with timer and role info -->
  <div class="header">
    <TimerDisplay
      {timeRemaining}
      {timerActive}
      {isHost}
      {onPauseTimer}
      {onResumeTimer}
      {onExtendTimer}
    />

    {#if !isSpy}
      <div class="role-info">
        <div class="location">📍 {location}</div>
        <div class="role">🎭 {role}</div>
      </div>
    {/if}
  </div>

  <!-- Turn indicator -->
  <div class="turn-indicator">
    {#if canAskQuestion}
      <strong>Your turn! Ask a question.</strong>
    {:else if canAnswerQuestion}
      <strong>Your turn! Answer the question.</strong>
    {:else}
      <span>{currentTurnName}'s turn</span>
    {/if}
  </div>

  <!-- Question history -->
  <QuestionLog {questionHistory} {players} />

  <!-- Ask question UI -->
  {#if canAskQuestion}
    <div class="ask-question">
      <h3>Ask a Question</h3>

      <label>Who do you want to ask?</label>
      <div class="player-select">
        {#each players as player}
          {#if player.id !== currentPlayerId}
            <button
              class="player-btn"
              class:selected={selectedAnswerer === player.id}
              onclick={() => selectedAnswerer = player.id}
            >
              {player.name}
            </button>
          {/if}
        {/each}
      </div>

      <label>Your Question</label>
      <textarea
        bind:value={questionText}
        placeholder="Type your question..."
        rows="3"
      />

      <button
        class="btn-primary w-full"
        disabled={!selectedAnswerer || !questionText.trim()}
        onclick={handleAskQuestion}
      >
        Ask Question
      </button>
    </div>
  {/if}

  <!-- Answer question UI -->
  {#if canAnswerQuestion}
    <div class="answer-question">
      <h3>Answer the Question</h3>
      <p class="question-text">"{lastQuestion?.question}"</p>

      <textarea
        bind:value={answerText}
        placeholder="Type your answer..."
        rows="3"
      />

      <button
        class="btn-primary w-full"
        disabled={!answerText.trim()}
        onclick={handleAnswerQuestion}
      >
        Submit Answer
      </button>
    </div>
  {/if}

  <!-- Action buttons -->
  <div class="actions">
    <button
      class="btn-secondary"
      onclick={() => showLocationList = !showLocationList}
    >
      📍 Location List
    </button>

    <button
      class="btn-accent"
      onclick={onCallAccusation}
    >
      🔍 Call for Accusation
    </button>
  </div>

  <!-- Location list modal -->
  {#if showLocationList}
    <LocationList onClose={() => showLocationList = false} />
  {/if}
</div>

<style>
  .question-round {
    max-width: 600px;
    margin: 0 auto;
    padding: 1rem;
  }

  .header {
    margin-bottom: 1.5rem;
  }

  .role-info {
    background: #282828;
    padding: 1rem;
    border-radius: 0.5rem;
    margin-top: 1rem;
    display: flex;
    justify-content: space-around;
  }

  .turn-indicator {
    text-align: center;
    padding: 1rem;
    background: #3c3836;
    border-radius: 0.5rem;
    margin-bottom: 1rem;
  }

  .ask-question, .answer-question {
    background: #282828;
    padding: 1.5rem;
    border-radius: 0.5rem;
    margin-bottom: 1rem;
  }

  .player-select {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin: 0.5rem 0 1rem 0;
  }

  .player-btn {
    background: #3c3836;
    color: #ebdbb2;
    padding: 0.5rem 1rem;
    border-radius: 0.25rem;
    border: 2px solid transparent;
    cursor: pointer;
    transition: all 0.2s;
  }

  .player-btn.selected {
    border-color: #d79921;
    background: #504945;
  }

  textarea {
    width: 100%;
    background: #3c3836;
    color: #ebdbb2;
    border: 1px solid #504945;
    border-radius: 0.25rem;
    padding: 0.75rem;
    margin-bottom: 1rem;
    font-family: inherit;
    resize: vertical;
  }

  .question-text {
    background: #3c3836;
    padding: 1rem;
    border-radius: 0.25rem;
    margin-bottom: 1rem;
    font-style: italic;
    color: #d79921;
  }

  .actions {
    display: flex;
    gap: 1rem;
    margin-top: 1.5rem;
  }

  .actions button {
    flex: 1;
  }
</style>
```

### Component: TimerDisplay.svelte

```svelte
<script lang="ts">
  let {
    timeRemaining,
    timerActive,
    isHost,
    onPauseTimer,
    onResumeTimer,
    onExtendTimer
  } = $props<{
    timeRemaining: number;
    timerActive: boolean;
    isHost: boolean;
    onPauseTimer: () => void;
    onResumeTimer: () => void;
    onExtendTimer: (seconds: number) => void;
  }>();

  let showExtendMenu = $state(false);

  const minutes = $derived(Math.floor(timeRemaining / 60));
  const seconds = $derived(timeRemaining % 60);

  const timeString = $derived(
    `${minutes}:${seconds.toString().padStart(2, '0')}`
  );

  const isLowTime = $derived(timeRemaining < 60);
  const isCritical = $derived(timeRemaining < 30);
</script>

<div class="timer-display" class:low-time={isLowTime} class:critical={isCritical}>
  <div class="timer-icon">⏱️</div>
  <div class="time">{timeString}</div>

  {#if isHost}
    <div class="timer-controls">
      {#if timerActive}
        <button class="btn-small" onclick={onPauseTimer}>⏸️ Pause</button>
      {:else}
        <button class="btn-small" onclick={onResumeTimer}>▶️ Resume</button>
      {/if}

      <button
        class="btn-small"
        onclick={() => showExtendMenu = !showExtendMenu}
      >
        ⏰ Extend
      </button>
    </div>

    {#if showExtendMenu}
      <div class="extend-menu">
        <button onclick={() => { onExtendTimer(30); showExtendMenu = false; }}>
          +30 sec
        </button>
        <button onclick={() => { onExtendTimer(60); showExtendMenu = false; }}>
          +1 min
        </button>
        <button onclick={() => { onExtendTimer(120); showExtendMenu = false; }}>
          +2 min
        </button>
      </div>
    {/if}
  {/if}
</div>

<style>
  .timer-display {
    background: #282828;
    padding: 1.5rem;
    border-radius: 0.5rem;
    text-align: center;
    border: 3px solid #458588;
    transition: border-color 0.3s;
  }

  .timer-display.low-time {
    border-color: #d79921;
  }

  .timer-display.critical {
    border-color: #cc241d;
    animation: pulse 1s infinite;
  }

  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.7; }
  }

  .timer-icon {
    font-size: 2rem;
  }

  .time {
    font-size: 3rem;
    font-weight: bold;
    color: #d79921;
    font-family: monospace;
  }

  .timer-controls {
    display: flex;
    gap: 0.5rem;
    justify-content: center;
    margin-top: 1rem;
  }

  .btn-small {
    background: #3c3836;
    color: #ebdbb2;
    padding: 0.5rem 1rem;
    border-radius: 0.25rem;
    border: none;
    cursor: pointer;
    font-size: 0.875rem;
  }

  .extend-menu {
    display: flex;
    gap: 0.5rem;
    justify-content: center;
    margin-top: 0.5rem;
  }

  .extend-menu button {
    background: #504945;
    color: #ebdbb2;
    padding: 0.5rem 0.75rem;
    border-radius: 0.25rem;
    border: none;
    cursor: pointer;
    font-size: 0.75rem;
  }
</style>
```

---

## State Management

### Game Store Integration

```typescript
// Add to game.svelte.ts
export interface SpyfallPlayerState {
    isSpy: boolean;
    location?: string;
    role?: string;
    hasAcknowledged: boolean;
    hasVoted: boolean;
}

export interface SpyfallPublicState {
    phase: string;
    timeRemaining: number;
    timerActive: boolean;
    currentTurn: string;
    questionHistory: QuestionAnswer[];
    votesSubmitted: number;
}

function handleSpyfallEvent(event: GameEvent) {
    switch (event.type) {
        case 'spy_assigned':
            playerState.isSpy = true;
            break;

        case 'location_assigned':
            playerState.isSpy = false;
            playerState.location = event.payload.location;
            playerState.role = event.payload.role;
            break;

        case 'question_asked':
            publicState.questionHistory.push(event.payload);
            break;

        case 'timer_updated':
            publicState.timeRemaining = event.payload.remaining;
            break;

        // ... handle all events
    }
}
```

---

## Implementation Checklist

### Backend Tasks
- [ ] Create `backend/internal/games/spyfall/` directory
- [ ] Implement state types (`state.go`)
- [ ] Implement location database (`locations.go`)
- [ ] Implement timer management (`timer.go`)
- [ ] Implement voting logic (`voting.go`)
- [ ] Implement config validation (`config.go`)
- [ ] Implement phase transitions (`phases.go`)
- [ ] Implement `Game` interface (`game.go`)
- [ ] Register game in `games/registry.go`
- [ ] Write unit tests for timer logic
- [ ] Write unit tests for voting
- [ ] Write unit tests for location guess

### Frontend Tasks
- [ ] Create `frontend/src/lib/games/spyfall/` directory
- [ ] Implement `locationConfig.ts`
- [ ] Implement `RoleReveal.svelte`
- [ ] Implement `TimerDisplay.svelte`
- [ ] Implement `QuestionLog.svelte`
- [ ] Implement `QuestionRound.svelte`
- [ ] Implement `Accusation.svelte`
- [ ] Implement `LocationGuess.svelte`
- [ ] Implement `LocationList.svelte`
- [ ] Implement `Results.svelte`
- [ ] Implement `SpyfallGame.svelte` main container
- [ ] Add event handlers to game store
- [ ] Add Spyfall to game selection UI

### Testing Tasks
- [ ] Test spy selection randomization
- [ ] Test location and role assignment
- [ ] Test timer (start/pause/resume/extend)
- [ ] Test timer expiration handling
- [ ] Test question/answer flow
- [ ] Test accusation voting
- [ ] Test location guessing
- [ ] Test all win conditions
- [ ] Multiplayer playtest with 3-8 players
- [ ] Mobile UI testing

### Polish Tasks
- [ ] Add timer animations
- [ ] Add sound effects for low time warning
- [ ] Add "How to Play" guide
- [ ] Add location reference sheet design
- [ ] Optimize for mobile
- [ ] Test reconnection handling
- [ ] Performance testing

---

## Estimated Timeline

- **Backend**: 2 days (location system, timer, turn management)
- **Frontend**: 2 days (timer UI, Q&A interface, location display)
- **Content**: 0.5 days (location database - 20-30 locations)
- **Testing**: 0.5 days (timer testing, multiplayer)
- **Polish**: 0.5 days (mobile optimization, animations)

**Total**: 4-5 days

---

## Notes

### Key Differences from Werewolf/Avalon
- **Real-time timer** instead of phase-based
- **Dynamic turn order** (answerer becomes next asker)
- **No voting rounds** (except accusation)
- **Location database** (content-heavy)
- **Spy vs Everyone** instead of teams

### Shared Patterns
- Role assignment with hidden information ✅
- Role reveal phase ✅
- Event sourcing architecture ✅
- Mobile-first UI ✅
- In-person social focus ✅

### Technical Challenges
1. **Timer synchronization** - must stay in sync across all clients
2. **Dynamic turn order** - answerer becomes next questioner
3. **Location database** - needs 30+ diverse locations with roles
4. **Q&A interface** - must be fast and mobile-friendly
5. **Spy guessing** - location list must be easy to browse

### Content Requirements
- Minimum 20 locations (30+ recommended)
- 6-8 roles per location
- Diverse themes (travel, work, leisure, etc.)
- Avoid ambiguous locations
- Roles should be distinctive and interesting

---

**Document Version**: 1.0
**Last Updated**: 2025-11-19
**Status**: Ready for Implementation
