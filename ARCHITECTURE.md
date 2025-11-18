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
- **Audit trail**: Full game history for stats/learning
- **Spectators**: Subscribe to event stream

### 2. Game Abstraction Layer
```go
type Game interface {
    Initialize(config, players) ([]Event, error)
    ValidateAction(playerID, action) error
    ProcessAction(playerID, action) ([]Event, error)
    GetPlayerState(playerID) PlayerState
    GetPublicState() PublicState
    IsFinished() bool
    GetResults() GameResults
}
```

### 3. Server-Authoritative Design
Clients **request actions**, server **validates and broadcasts results**.
- Server holds canonical state
- Clients send intent ("I vote for Alice"), not state changes

---

## Data Model

### Core Entities

```go
// Room: Container for a game session
type Room struct {
    ID           string    // 6-char code: "XJ4K2P"
    Status       RoomStatus // waiting, playing, finished
    GameType     string     // "werewolf", "avalon"
    Host         string     // PlayerID of host
    Players      []Player
    EventLog     []GameEvent // Append-only event history
    Game         Game        // Game-specific state machine
}

// Player: Anonymous participant
type Player struct {
    ID           string
    SessionToken string
    DisplayName  string
}

// GameEvent: Immutable fact about what happened
type GameEvent struct {
    Type      string          // "role_assigned", "vote_cast"
    ActorID   string
    Payload   json.RawMessage
}
```

---

## Project Structure

```
roundtable/
├── backend/
│   ├── cmd/server/main.go      # Entry point
│   ├── internal/
│   │   ├── core/               # Platform core (Room, Player, Game interface)
│   │   ├── games/              # Game implementations (Werewolf)
│   │   ├── server/             # HTTP & WebSocket handlers
│   │   └── store/              # State persistence (Memory/Redis)
│
├── frontend/
│   ├── src/
│   │   ├── routes/             # Pages
│   │   ├── lib/
│   │   │   ├── components/     # Reusable UI
│   │   │   ├── games/          # Game-specific components
│   │   │   └── stores/         # State management
```

---

## WebSocket Protocol

### Client → Server
```json
{ "type": "authenticate", "payload": { "sessionToken": "..." } }
{ "type": "action", "payload": { "action": { "type": "vote", "payload": {...} } } }
```

### Server → Client
```json
{ "type": "authenticated", "payload": { "roomState": {...} } }
{ "type": "event", "payload": { "event": {...} } }
{ "type": "events", "payload": { "events": [...] } }
```

---

## Architecture Decisions

### 1. Event Sourcing = Future-Proof
Adding spectator mode, replays, or time-travel debugging is trivial. Event log is your foundation.

### 2. Game Interface = Easy Expansion
Werewolf and Avalon are drastically different games. If the interface works for both, it'll work for Bohnanza, Coup, Secret Hitler.

### 3. Server-Authoritative = No Cheating
Never trust the client. Server validates everything. Clients are just dumb views.

### 4. In-Person First
The app replaces cards, not in-person interaction. Physical voting and social deduction remain central to the experience.
