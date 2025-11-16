# Cardless: Party Game Platform Architecture

## Vision
Replace physical cards with phones for in-person multiplayer party games. The app handles role assignment, private information, and night actions while preserving the social, in-person nature of these games.

## Core Principles

### 1. Event Sourcing Foundation
Every game is modeled as:
```
Current State = Initial State + Sequence of Events
```

**Why this matters:**
- **Reconnection**: Replay events to rebuild player state
- **Board sync**: Single source of truth, multiple views
- **Spectators**: Subscribe to event stream, no special logic
- **Time travel**: Rewind/replay for debugging or "undo" features
- **Audit trail**: Full game history for stats/learning

### 2. Game Abstraction Layer
```go
type Game interface {
    // Lifecycle
    Initialize(config GameConfig, players []Player) ([]GameEvent, error)

    // Action validation & processing
    ValidateAction(playerID string, action Action) error
    ProcessAction(playerID string, action Action) ([]GameEvent, error)

    // State views (filtered by permission)
    GetPlayerState(playerID string) PlayerState
    GetPublicState() PublicState

    // Game flow
    GetPhase() GamePhase
    IsFinished() bool
    GetResults() GameResults
}
```

**Each game implements this interface.** Adding Avalon after Werewolf means:
1. Create `games/avalon/` package
2. Implement the `Game` interface
3. Register in `games/registry.go`
4. Done. Core infrastructure handles everything else.

### 3. Event-Driven Architecture

```
┌─────────────────────────────────────────────┐
│           Event Log (Source of Truth)       │
│  [PlayerJoined, RoleAssigned, NightAction...]│
└─────────────────────────────────────────────┘
                      │
         ┌────────────┼────────────┐
         ▼            ▼            ▼
   ┌─────────┐  ┌─────────┐  ┌─────────┐
   │ Phone 1 │  │ Phone 2 │  │ Phone 3 │
   │ (Alice) │  │  (Bob)  │  │ (Carol) │
   └─────────┘  └─────────┘  └─────────┘
     Private      Private      Private
```

**Views are derived from events:**
- **Player view**: Filtered to what that player can see (their role, night action results)
- Events can be public (everyone sees) or private (specific players only)
- No board/spectator view for MVP - focus is on player experience

### 4. Server-Authoritative Design
Clients **request actions**, server **validates and broadcasts results**.

```
Phone → WebSocket → Server validates → Process action → Emit events → Broadcast to affected clients
```

**Never trust client:**
- Server holds canonical state
- Clients send intent ("I vote for Alice"), not state changes
- Server decides validity and applies changes

---

## Data Model

### Core Entities

```go
// Room: Container for a game session
type Room struct {
    ID           string    // 6-char code: "XJ4K2P"
    CreatedAt    time.Time
    Status       RoomStatus // waiting, playing, finished
    GameType     string     // "werewolf", "avalon"

    Host         string     // PlayerID of host
    Players      []Player
    MaxPlayers   int

    EventLog     []GameEvent // Append-only event history
    Game         Game        // Game-specific state machine
}

// Player: Anonymous participant
type Player struct {
    ID           string    // UUID
    SessionToken string    // For reconnection (stored in localStorage)
    DisplayName  string    // "Alice"

    Connected    bool
    LastSeenAt   time.Time
    JoinedAt     time.Time
}

// GameEvent: Immutable fact about what happened
type GameEvent struct {
    ID        string          // Event UUID
    Timestamp time.Time
    Type      string          // "role_assigned", "vote_cast", "phase_changed"
    ActorID   string          // PlayerID who triggered it (or "system")
    Payload   json.RawMessage // Event-specific data
}

// Action: Player intent
type Action struct {
    Type    string          // "vote", "select_role", "pass_card"
    Payload json.RawMessage
}
```

### Game-Specific State (Werewolf Example)

```go
type WerewolfGame struct {
    Config WerewolfConfig
    State  WerewolfState
}

type WerewolfConfig struct {
    Roles          []RoleType // [Werewolf, Werewolf, Seer, Robber, Troublemaker, Villager...]
    NightDuration  time.Duration
    DayDuration    time.Duration
}

type WerewolfState struct {
    Phase         Phase // setup, night, day, results
    RoleAssignments map[string]RoleType // playerID → role
    OriginalRoles   map[string]RoleType // before night actions

    NightActions  []NightAction
    Votes         map[string]string // voterID → targetID

    StartedAt     time.Time
    PhaseEndsAt   time.Time
}
```

---

## Project Structure

```
roundtable/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go                 # Entry point
│   │
│   ├── internal/
│   │   ├── core/                       # Platform core (game-agnostic)
│   │   │   ├── game.go                 # Game interface
│   │   │   ├── room.go                 # Room management
│   │   │   ├── player.go               # Player management
│   │   │   ├── event.go                # Event types
│   │   │   └── errors.go
│   │   │
│   │   ├── games/                      # Game implementations
│   │   │   ├── registry.go             # Game factory
│   │   │   ├── werewolf/
│   │   │   │   ├── game.go             # Implements core.Game
│   │   │   │   ├── state.go            # Game state
│   │   │   │   ├── events.go           # Event types
│   │   │   │   ├── actions.go          # Action types
│   │   │   │   ├── roles.go            # Role logic
│   │   │   │   └── phases.go           # Phase transitions
│   │   │   │
│   │   │   └── avalon/                 # Future game
│   │   │       └── ...
│   │   │
│   │   ├── server/                     # HTTP & WebSocket server
│   │   │   ├── server.go               # HTTP server setup
│   │   │   ├── handlers.go             # REST endpoints
│   │   │   ├── websocket.go            # WS connection manager
│   │   │   ├── messages.go             # WS message types
│   │   │   └── middleware.go           # Auth, logging
│   │   │
│   │   ├── store/                      # State persistence
│   │   │   ├── store.go                # Store interface
│   │   │   ├── memory.go               # In-memory (MVP)
│   │   │   └── redis.go                # Redis (production)
│   │   │
│   │   └── util/
│   │       ├── codegen.go              # 6-char room codes
│   │       └── token.go                # Session tokens
│   │
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── routes/
│   │   │   ├── +page.svelte            # Landing: Create or Join
│   │   │   │
│   │   │   ├── create/
│   │   │   │   └── +page.svelte        # Choose game type, create room
│   │   │   │
│   │   │   ├── join/
│   │   │   │   └── +page.svelte        # Enter code + name
│   │   │   │
│   │   │   ├── room/
│   │   │   │   └── [code]/
│   │   │   │       ├── +page.svelte    # Router: lobby → game → results
│   │   │   │       ├── +layout.svelte  # Room shell, WS connection
│   │   │   │       │
│   │   │   │       └── board/
│   │   │   │           └── +page.svelte # Optional board view (QR code to join)
│   │   │   │
│   │   │   └── api/                    # API proxy routes (optional)
│   │   │
│   │   ├── lib/
│   │   │   ├── components/
│   │   │   │   ├── ui/
│   │   │   │   │   ├── Button.svelte
│   │   │   │   │   ├── Card.svelte
│   │   │   │   │   └── Timer.svelte
│   │   │   │   │
│   │   │   │   ├── room/
│   │   │   │   │   ├── Lobby.svelte        # Player list, config, start
│   │   │   │   │   ├── PlayerList.svelte
│   │   │   │   │   ├── QRCode.svelte       # Room code QR
│   │   │   │   │   └── Results.svelte      # Game over screen
│   │   │   │   │
│   │   │   │   └── board/
│   │   │   │       └── BoardView.svelte    # Game-agnostic board
│   │   │   │
│   │   │   ├── stores/
│   │   │   │   ├── room.ts                 # Room state (players, status)
│   │   │   │   ├── game.ts                 # Game state (derived from events)
│   │   │   │   ├── websocket.ts            # WS connection manager
│   │   │   │   └── session.ts              # Player session (token, ID)
│   │   │   │
│   │   │   ├── games/                      # Game-specific components
│   │   │   │   ├── werewolf/
│   │   │   │   │   ├── WerewolfGame.svelte # Main game component
│   │   │   │   │   ├── RoleCard.svelte     # Show player's role
│   │   │   │   │   ├── NightPhase.svelte   # Night action UI
│   │   │   │   │   ├── DayPhase.svelte     # Voting UI
│   │   │   │   │   └── WerewolfBoard.svelte # Public board (minimal for Werewolf)
│   │   │   │   │
│   │   │   │   └── avalon/
│   │   │   │       └── ...
│   │   │   │
│   │   │   └── api/
│   │   │       └── client.ts               # API client (fetch wrapper)
│   │   │
│   │   └── app.html
│   │
│   ├── static/
│   │   └── favicon.png
│   │
│   ├── svelte.config.js
│   ├── vite.config.js
│   ├── package.json
│   └── tailwind.config.js              # Mobile-first styling
│
├── docker-compose.yml
├── Dockerfile
└── README.md
```

---

## Critical User Flows

### Flow 1: Create Room

```
[Landing Page]
    ↓
User clicks "Host Game"
    ↓
[Create Page] - Select game type (Werewolf)
    ↓
POST /api/rooms { gameType: "werewolf" }
    ↓
Server:
  - Generate 6-char code (XJ4K2P)
  - Create Room in store
  - Generate session token for host
  - Return { roomCode, sessionToken, playerID }
    ↓
Client:
  - Save token to localStorage
  - Redirect to /room/XJ4K2P
    ↓
[Room Page - Lobby]
  - Open WebSocket connection (sends token in handshake)
  - Server validates token → associates connection with playerID
  - Server sends initial room state
  - Show QR code for others to join
```

### Flow 2: Join Room

```
[Landing Page]
    ↓
User clicks "Join Game" or scans QR
    ↓
[Join Page] - Enter code + display name
    ↓
POST /api/rooms/XJ4K2P/join { displayName: "Alice" }
    ↓
Server:
  - Validate room exists and is joinable
  - Create Player
  - Generate session token
  - Emit event: PlayerJoined
  - Return { sessionToken, playerID, roomCode }
    ↓
Client:
  - Save token to localStorage
  - Redirect to /room/XJ4K2P
    ↓
WebSocket connection established
    ↓
Server broadcasts to all clients:
  - Event: PlayerJoined { playerID, displayName }
  - Updated room state
```

### Flow 3: Start Game (Werewolf)

```
[Lobby] - Host configures roles
    ↓
Host clicks "Start Game"
    ↓
POST /api/rooms/XJ4K2P/start { config: { roles: [...] } }
    ↓
Server:
  - Validate config (enough roles for players)
  - Initialize game: game.Initialize(config, players)
  - Game emits events:
      - GameStarted
      - RoleAssigned (per player)
      - PhaseChanged { phase: "night" }
  - Append events to room.EventLog
    ↓
Server broadcasts events to clients:
  - All clients: GameStarted, PhaseChanged
  - Individual clients: RoleAssigned (only to that player)
    ↓
Clients:
  - Receive events
  - Update game store (apply events to state)
  - UI transitions to game view
  - Each player sees their role card
```

### Flow 4: Night Phase (Werewolf)

```
[Night Phase UI]
    ↓
Server broadcasts: PhaseChanged { phase: "night" }
    ↓
Each role's phone shows their specific action UI:

Werewolves:
  - See list of other werewolves (or can view center card if alone)
    ↓
Seer:
  - Choose: view one player's role OR two center cards
  - Taps a player or center cards
    ↓
WS → Server: { type: "action", action: { type: "seer_view", targetID: "bob-uuid" } }
    ↓
Server validates and processes:
  - Checks Seer hasn't acted yet
  - Records the view action
  - Sends private result back to Seer
    ↓
Seer sees: "Bob is a Villager" (private to seer only)
    ↓
Robber:
  - Chooses another player to swap with
  - Server processes swap, updates roleAssignments
  - Robber sees their NEW role
    ↓
Troublemaker:
  - Chooses two OTHER players to swap
  - Server swaps them (Troublemaker doesn't see what they swapped)
    ↓
Drunk:
  - Chooses a center card to swap with (forced action)
  - Server swaps but Drunk doesn't see new role
    ↓
Insomniac:
  - Automatically shown their FINAL role (after all swaps)
    ↓
Host manually advances when everyone is done
    ↓
Server emits: PhaseChanged { phase: "day" }
```

### Flow 5: Day Phase (Werewolf)

```
[Day Phase UI]
    ↓
All players see:
  - Discussion timer (optional, can start/pause/extend)
  - "Discuss who you think is a Werewolf"
    ↓
Players discuss in person (no in-game chat needed)
    ↓
Timer ends (or host advances manually)
    ↓
Everyone votes PHYSICALLY by pointing at the same time (traditional ONUW)
    ↓
Host or anyone taps "Reveal Roles"
    ↓
All players see their FINAL roles on their phones
  (shows which swaps happened during night)
    ↓
Players determine winner together based on:
  - Who got eliminated (physical vote)
  - What their final roles were (shown on phones)
    ↓
Host can tap "Play Again" to start new game with same players
```

### Flow 6: Reconnection

```
Player loses connection (phone sleeps, WiFi drops)
    ↓
Client detects disconnect (WS close event)
    ↓
Client automatically attempts reconnect:
  - Read sessionToken from localStorage
  - Open new WebSocket connection
  - Send in handshake: { token: "abc123..." }
    ↓
Server:
  - Validate token → playerID
  - Mark player as connected again
  - Send full state reconstruction:
      - Room state
      - All events since game started (client replays to rebuild state)
    ↓
Client:
  - Replay events in order
  - Rebuild game state
  - Resume from current phase
    ↓
Player sees current game state, can continue playing
```

---

## WebSocket Protocol

### Client → Server Messages

```typescript
type ClientMessage =
  | { type: "authenticate", token: string }
  | { type: "action", action: Action }
  | { type: "ping" }

type Action =
  | { type: "vote", targetID: string }
  | { type: "seer_view", targetID: string }
  | { type: "robber_swap", targetID: string }
  | ... (game-specific)
```

### Server → Client Messages

```typescript
type ServerMessage =
  | { type: "authenticated", playerID: string, roomState: RoomState }
  | { type: "event", event: GameEvent }
  | { type: "events", events: GameEvent[] } // Batch for reconnection
  | { type: "error", error: string }
  | { type: "pong" }

type GameEvent =
  | { type: "player_joined", playerID: string, displayName: string }
  | { type: "player_left", playerID: string }
  | { type: "game_started", config: GameConfig }
  | { type: "role_assigned", role: string } // Sent only to that player
  | { type: "phase_changed", phase: string, endsAt: string }
  | { type: "vote_cast", voterID: string }
  | { type: "votes_revealed", votes: Record<string, string> }
  | { type: "game_finished", results: GameResults }
  | ... (game-specific)
```

---

## Game Abstraction in Detail

### Interface

```go
package core

type Game interface {
    // Initialize game with players and config
    // Returns initial events (GameStarted, RoleAssigned, etc.)
    Initialize(config GameConfig, players []Player) ([]GameEvent, error)

    // Validate if player can perform this action now
    ValidateAction(playerID string, action Action) error

    // Process valid action, return resulting events
    ProcessAction(playerID string, action Action) ([]GameEvent, error)

    // Get state visible to this player (filtered)
    GetPlayerState(playerID string) PlayerState

    // Get public state (board view, spectators)
    GetPublicState() PublicState

    // Current phase for UI rendering
    GetPhase() GamePhase

    // Check if game is over
    IsFinished() bool

    // Get final results (only valid if IsFinished)
    GetResults() GameResults
}
```

### Werewolf Implementation Example

```go
package werewolf

type WerewolfGame struct {
    config  Config
    state   State
    players map[string]*Player // playerID → player
}

func (g *WerewolfGame) Initialize(config GameConfig, players []Player) ([]GameEvent, error) {
    // Validate config
    wc := config.(*WerewolfConfig)
    if len(wc.Roles) != len(players) {
        return nil, errors.New("role count must match player count")
    }

    // Shuffle roles
    roles := shuffle(wc.Roles)

    // Create events
    events := []GameEvent{
        {Type: "game_started", Payload: wc},
    }

    // Assign roles
    g.state.roleAssignments = make(map[string]RoleType)
    for i, player := range players {
        role := roles[i]
        g.state.roleAssignments[player.ID] = role
        events = append(events, GameEvent{
            Type: "role_assigned",
            ActorID: player.ID,
            Payload: RoleAssignedPayload{Role: role},
        })
    }

    // Start night phase
    g.state.phase = PhaseNight
    g.state.phaseEndsAt = time.Now().Add(wc.NightDuration)
    events = append(events, GameEvent{
        Type: "phase_changed",
        Payload: PhaseChangedPayload{Phase: PhaseNight, EndsAt: g.state.phaseEndsAt},
    })

    return events, nil
}

func (g *WerewolfGame) ValidateAction(playerID string, action Action) error {
    switch action.Type {
    case "vote":
        if g.state.phase != PhaseDay {
            return errors.New("can only vote during day phase")
        }
        if g.state.votes[playerID] != "" {
            return errors.New("already voted")
        }
        // ... more validation
        return nil

    case "seer_view":
        if g.state.roleAssignments[playerID] != RoleSeer {
            return errors.New("only seer can view roles")
        }
        if g.state.phase != PhaseNight {
            return errors.New("can only view during night")
        }
        // ... more validation
        return nil

    default:
        return errors.New("unknown action type")
    }
}

func (g *WerewolfGame) ProcessAction(playerID string, action Action) ([]GameEvent, error) {
    switch action.Type {
    case "vote":
        targetID := action.Payload["targetID"].(string)
        g.state.votes[playerID] = targetID

        events := []GameEvent{
            {Type: "vote_cast", ActorID: playerID},
        }

        // Check if all voted
        if len(g.state.votes) == len(g.players) {
            events = append(events, GameEvent{
                Type: "votes_revealed",
                Payload: VotesRevealedPayload{Votes: g.state.votes},
            })

            // Determine winner
            results := g.calculateResults()
            events = append(events, GameEvent{
                Type: "game_finished",
                Payload: results,
            })
        }

        return events, nil

    // ... other actions
    }

    return nil, nil
}

func (g *WerewolfGame) GetPlayerState(playerID string) PlayerState {
    state := PlayerState{
        Phase: g.state.phase,
        YourRole: g.state.roleAssignments[playerID],
    }

    // Add role-specific info
    switch state.YourRole {
    case RoleWerewolf:
        // Show other werewolves
        for pid, role := range g.state.roleAssignments {
            if role == RoleWerewolf && pid != playerID {
                state.OtherWerewolves = append(state.OtherWerewolves, pid)
            }
        }
    }

    if g.state.phase == PhaseDay {
        state.HasVoted = g.state.votes[playerID] != ""
    }

    return state
}

func (g *WerewolfGame) GetPublicState() PublicState {
    voteCount := len(g.state.votes)
    totalPlayers := len(g.players)

    return PublicState{
        Phase: g.state.phase,
        PhaseEndsAt: g.state.phaseEndsAt,
        VotesSubmitted: voteCount,
        TotalPlayers: totalPlayers,
    }
}
```

---

## Phased Implementation Plan

### **Phase 1: Foundation (3-4 days) - "It works!" moment**

**Goal**: Create room, join, see each other's names

**Backend:**
- [ ] Project setup: `go mod init`, basic server
- [ ] Room creation endpoint: `POST /api/rooms`
- [ ] Join endpoint: `POST /api/rooms/:code/join`
- [ ] In-memory store for rooms
- [ ] WebSocket connection handling
- [ ] Broadcast player joined/left events

**Frontend:**
- [ ] SvelteKit project setup
- [ ] Landing page (Create / Join buttons)
- [ ] Create room flow
- [ ] Join room flow (enter code + name)
- [ ] Room lobby: player list updates in real-time
- [ ] WebSocket store with auto-reconnect

**Success criteria**: 5 friends can join room XJ4K2P from their phones, see each other's names appear live

---

### **Phase 2: Werewolf MVP (5-6 days)**

**Goal**: Play full game of One Night Werewolf

**Completed:**
- [x] Game interface definition
- [x] Event sourcing infrastructure
- [x] Werewolf game implementation (basic flow)
- [x] Role assignment (players + 3 center cards)
- [x] Start game endpoint
- [x] Game state store (event-driven)
- [x] Role reveal with acknowledgements
- [x] Night phase with host narration
- [x] Day phase with timer controls
- [x] Results calculation logic

**In Progress (Current Sprint):**
- [ ] Digital night actions for each role:
  - [ ] Werewolf: See other werewolves
  - [ ] Seer: View player or center cards
  - [ ] Robber: Swap and view new role
  - [ ] Troublemaker: Swap two others
  - [ ] Drunk: Swap with center (blind)
  - [ ] Insomniac: View final role
- [ ] Remove phone voting (use physical voting)
- [ ] Role reveal screen (show final roles)
- [ ] Play again feature
- [ ] Fix host tracking

**Success criteria**: Digital night actions work, physical voting flows naturally, can play multiple rounds

---

### **Phase 3: Polish & Stability (3-4 days)**

**Goal**: Stable, reconnection-proof Werewolf

**Backend:**
- [ ] Session token reconnection
- [ ] Event replay for reconnecting clients
- [ ] Room expiry and cleanup

**Frontend:**
- [ ] QR code generation for room code
- [ ] Reconnection handling (detect disconnect, auto-retry)
- [ ] Error boundaries and loading states
- [ ] Mobile UI polish (better touch targets)

**Success criteria**: Handles edge cases, reconnection works smoothly, ready for extended playtesting

---

### **Phase 4: Second Game - Avalon (4-5 days)**

**Goal**: Validate architecture, prove it's multi-game

**Before starting:**
- [ ] Extract reusable patterns from Werewolf
- [ ] Document game implementation guide
- [ ] Identify what's game-specific vs platform-specific

**Backend:**
- [ ] Avalon game implementation
- [ ] Quest voting mechanics
- [ ] Team selection logic
- [ ] Merlin/Assassin reveal

**Frontend:**
- [ ] Avalon-specific components
- [ ] Quest tracking UI
- [ ] Team selection UI

**Success criteria**: Avalon works with minimal changes to core platform, validates game abstraction

---

### **Phase 5: Production Ready (3-4 days)**

**Backend:**
- [ ] Redis store implementation
- [ ] Docker deployment optimization
- [ ] Rate limiting
- [ ] Health monitoring

**Frontend:**
- [ ] PWA manifest (installable)
- [ ] Performance optimization
- [ ] Analytics (optional)

**Success criteria**: Deploy to homelab or cloud, stable for regular game nights

---

## Architecture Decisions That Matter

### 1. **Event Sourcing = Future-Proof**
Adding spectator mode, replays, or time-travel debugging is trivial. Event log is your foundation.

### 2. **Game Interface = Easy Expansion**
Werewolf and Avalon are drastically different games. If the interface works for both, it'll work for Bohnanza, Coup, Secret Hitler.

### 3. **Server-Authoritative = No Cheating**
Never trust the client. Server validates everything. Clients are just dumb views.

### 4. **Anonymous Sessions = Instant Play**
localStorage session token. No signup friction. Add accounts later for stats.

### 5. **Mobile-First = Touch Targets**
Buttons min 44x44px. One-handed operation. Test on real phones early.

### 6. **In-Person First**
The app replaces cards, not in-person interaction. Physical voting and social deduction remain central to the experience.

---

## What to Build First: Quick Win Path

**Day 1-2: Skeleton**
1. Go server with health endpoint
2. SvelteKit landing page
3. Room creation → returns code → redirect to lobby
4. Join room → add name → redirect to lobby
5. WebSocket connection, broadcast "player joined" messages

**Day 3-4: Werewolf Roles**
6. Start game → assign random roles → show role to each player
7. "Night phase" (just display role, no actions)
8. "Day phase" → tap to vote → show results

**Day 5: Playtest!**
9. Invite friends, play a game, collect feedback

**This gets you to a playable prototype in one week.**

---

## Tech Stack Summary

| Layer | Technology | Rationale |
|-------|-----------|-----------|
| Backend | **Go 1.21+** | Goroutines for concurrency, fast compile, excellent WS support |
| Frontend | **SvelteKit** | You know it, mobile-friendly, reactive stores |
| Real-time | **nhooyr.io/websocket** | Modern Go WS library, context-aware |
| State | **In-memory → Redis** | Start simple, scale horizontally later |
| Events | **Event Sourcing** | Reconnection, replay, audit, spectators |
| Database | **SQLite** (optional) | Game history, stats (write-after, not during) |
| Styling | **Tailwind CSS** | Mobile-first utilities, rapid iteration |
| Deploy | **Docker** | Single container, easy homelab deployment |

---

## Next Steps

I'll now create the full project structure with:
1. Backend scaffolding (core interfaces, room management, WebSocket server)
2. Frontend scaffolding (routes, stores, WebSocket client)
3. Docker setup
4. Implementation guide

This architecture is **opinionated** about separation of concerns:
- **Core** = game-agnostic platform
- **Games** = isolated, pluggable
- **Server** = thin transport layer
- **Store** = swappable persistence

Adding game #2 and #3 will be **easy**. Let's build it.
