# 🎲 Cardless

**A multiplayer party game platform that replaces physical cards and boards with phones + optional shared screens.**

Like Jackbox, but for tabletop games like One Night Werewolf, Avalon, Bohnanza, and more.

---

## ✨ Vision

- **In-person play**: Same room with friends, not remote
- **Phone as controller**: Each player's phone is their "hand" (private info, actions)
- **Optional shared board**: Public game state on TV/projector for games that need it
- **Anonymous-first**: No signup required, just join with a room code
- **Ultra-convenient**: Faster setup than shuffling cards, no lost pieces

---

## 🎯 MVP: One Night Werewolf

- Room creation with 6-character code (+ QR code generation)
- Anonymous join with display name
- Host configures roles and starts game
- Server randomly assigns roles (players see only theirs)
- Night phase: server wakes roles in sequence, players act on phone
- Day phase: discussion timer, voting on phones
- Results reveal with win condition
- "Play again" with same players
- Optional board view (architecture supports it from day 1)

---

## 🏗️ Architecture

### Tech Stack

| Layer | Technology | Why |
|-------|-----------|-----|
| **Backend** | Go 1.21+ | Goroutines for concurrency, fast iteration, excellent WebSocket support |
| **Frontend** | SvelteKit | You're comfortable with it, mobile-friendly, reactive |
| **Real-time** | Native WebSockets | `nhooyr.io/websocket` - modern, context-aware |
| **State** | In-memory → Redis | Start simple, scale later |
| **Architecture** | Event Sourcing | Reconnection, replay, audit trail, spectator mode |
| **Database** | SQLite (optional) | Game history/stats (write-after, not during gameplay) |
| **Styling** | Tailwind CSS | Mobile-first utilities, rapid iteration |
| **Deploy** | Docker | Single-container, homelab-ready |

### Core Concepts

**Event Sourcing**: Game state = initial state + sequence of events
```
Current State = Initial State + [Event1, Event2, Event3, ...]
```

**Game Abstraction**: Every game implements a common interface
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

**Multi-View Support**: Events can be public or private
- Public: Everyone sees (phase changes, votes revealed)
- Private: Specific players only (role assignments, seer views)

---

## 📁 Project Structure

```
roundtable/
├── backend/                    # Go backend
│   ├── cmd/server/            # Entry point
│   ├── internal/
│   │   ├── core/              # Platform core (game-agnostic)
│   │   │   ├── game.go        # Game interface
│   │   │   ├── room.go        # Room management
│   │   │   ├── player.go      # Player management
│   │   │   └── event.go       # Event types
│   │   ├── games/             # Game implementations
│   │   │   ├── registry.go    # Game factory
│   │   │   └── werewolf/      # One Night Werewolf
│   │   ├── server/            # HTTP & WebSocket
│   │   │   ├── handlers.go    # REST endpoints
│   │   │   ├── websocket.go   # WS connection manager
│   │   │   └── messages.go    # Message types
│   │   └── store/             # State persistence
│   │       ├── memory.go      # In-memory (MVP)
│   │       └── redis.go       # Redis (future)
│   └── Dockerfile
│
├── frontend/                  # SvelteKit frontend
│   ├── src/
│   │   ├── routes/
│   │   │   ├── +page.svelte           # Landing page
│   │   │   ├── create/                # Create room
│   │   │   ├── join/                  # Join room
│   │   │   └── room/[code]/           # Game room
│   │   ├── lib/
│   │   │   ├── stores/
│   │   │   │   ├── session.ts         # Player session
│   │   │   │   ├── websocket.ts       # WS connection
│   │   │   │   └── game.ts            # Game state
│   │   │   ├── components/            # Reusable UI
│   │   │   └── games/                 # Game-specific components
│   │   │       └── werewolf/
│   │   └── app.css                    # Tailwind styles
│   └── Dockerfile
│
├── docker-compose.yml
└── ARCHITECTURE.md            # Detailed design doc
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Node.js 20+

### Option 1: Railway Deployment (Recommended)

This project uses Railpack (Railway's build system) for deployment:

1. Push to GitHub
2. Create two services in Railway:
   - **Backend**: Set root directory to `backend/`
   - **Frontend**: Set root directory to `frontend/`
3. Railway will auto-detect and build both services
4. Set environment variables as needed

### Option 2: Local Development

**Backend:**
```bash
cd backend
go mod download
go run cmd/server/main.go

# Server runs on :8080
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev

# Frontend runs on :5173 (dev) or :3000 (prod)
```

---

## 🎮 How to Play (MVP)

1. **Host creates room**: Go to app → "Host a Game" → Choose Werewolf → Enter name → Get room code `XJ4K2P`
2. **Players join**: Others go to app → "Join Game" → Enter code `XJ4K2P` + name
3. **Host starts**: Configure roles (3+ players) → "Start Game"
4. **Night phase**: Each player sees their role. Server wakes roles in sequence (Werewolves see each other, Seer views a role, etc.)
5. **Day phase**: Discuss (in person!), then vote on your phone for who to eliminate
6. **Results**: See who won, all roles revealed
7. **Play again**: Same players, new game

---

## 📋 Implementation Roadmap

### ✅ Phase 1: Foundation (COMPLETE)
- [x] Backend project structure
- [x] Core game abstraction
- [x] Event sourcing infrastructure
- [x] Room and player management
- [x] WebSocket connection handling
- [x] In-memory store
- [x] Frontend SvelteKit setup
- [x] Landing, create, join pages
- [x] WebSocket and session stores
- [x] Docker configuration

### 🔄 Phase 2: Werewolf MVP (NEXT - 5-6 days)
- [ ] Complete Werewolf game logic
  - [ ] Full night phase with all role actions
  - [ ] Day phase voting
  - [ ] Results calculation
- [ ] Game UI components
  - [ ] Lobby with player list
  - [ ] Role reveal animation
  - [ ] Night phase UI (role-specific)
  - [ ] Day phase voting UI
  - [ ] Results screen
- [ ] Host start game flow
- [ ] End-to-end playtest

### 📦 Phase 3: Polish (3-4 days)
- [ ] QR code generation for room
- [ ] Reconnection handling
- [ ] Timer animations
- [ ] Mobile UI polish (haptic feedback, sounds)
- [ ] "Play again" feature
- [ ] Error boundaries & loading states

### 🖥️ Phase 4: Board View (2-3 days)
- [ ] `/room/:code/board` route
- [ ] Public state display
- [ ] Werewolf board (timer, phase, vote status)
- [ ] QR code on board to join

### 🎯 Phase 5: Second Game - Avalon (4-5 days)
- [ ] Avalon game implementation
- [ ] Validate game abstraction
- [ ] Mission voting
- [ ] Team selection UI
- [ ] Merlin/Assassin reveal

### 🚢 Phase 6: Production (3-4 days)
- [ ] Redis store implementation
- [ ] Room expiry (cleanup)
- [ ] Rate limiting
- [ ] PWA manifest (installable)
- [ ] Homelab deployment guide

---

## 🎯 What to Build First (Day 1-5 Quick Win)

**Goal**: Playable prototype in one week

1. **Day 1-2**: Complete basic lobby
   - Room creation works
   - Players can join and see each other
   - WebSocket messages flow

2. **Day 3-4**: Simplest possible Werewolf
   - Assign random roles
   - Show role to each player
   - Skip night actions (just show roles)
   - Vote → show results

3. **Day 5**: Playtest with friends!
   - Gather feedback
   - Fix critical bugs
   - Iterate on UX

---

## 🏛️ Key Design Decisions

### 1. Why Event Sourcing?
- **Reconnection**: Replay events to rebuild state
- **Spectators**: Subscribe to event stream
- **Debugging**: Full audit trail
- **Time travel**: Rewind/replay for undo features

### 2. Why Server-Authoritative?
- **No cheating**: Server holds canonical state
- **Security**: Clients can't see hidden info
- **Consistency**: Single source of truth

### 3. Why Game Abstraction?
- **Extensibility**: Adding game #2 is easy
- **Isolation**: Game-specific logic doesn't leak into core
- **Testing**: Test games independently

### 4. Why Anonymous-First?
- **Instant play**: No signup friction
- **Privacy**: Just display names
- **Optional accounts**: Add later for stats

### 5. Why Board View from Day 1?
- **Future games need it**: Bohnanza, Avalon
- **Better UX**: Shared screen for public state
- **Spectators**: Easy to add

---

## 🔧 API Reference

### REST Endpoints

```bash
# Create room
POST /api/rooms
{
  "gameType": "werewolf",
  "displayName": "Alice",
  "maxPlayers": 10
}
→ { "roomCode": "XJ4K2P", "sessionToken": "...", "playerId": "..." }

# Join room
POST /api/rooms/:code/join
{ "displayName": "Bob" }
→ { "sessionToken": "...", "playerId": "...", "roomCode": "XJ4K2P" }

# Get room state
GET /api/rooms/:code
→ { "id": "XJ4K2P", "status": "waiting", "players": [...], ... }

# Start game
POST /api/rooms/:code/start
{ "config": { "roles": ["werewolf", "seer", ...] } }
→ 200 OK
```

### WebSocket Protocol

**Client → Server:**
```json
{ "type": "authenticate", "payload": { "sessionToken": "..." } }
{ "type": "action", "payload": { "action": { "type": "vote", "payload": {...} } } }
{ "type": "ping" }
```

**Server → Client:**
```json
{ "type": "authenticated", "payload": { "playerId": "...", "roomState": {...} } }
{ "type": "event", "payload": { "event": {...} } }
{ "type": "events", "payload": { "events": [...] } }
{ "type": "error", "payload": { "message": "..." } }
{ "type": "pong" }
```

---

## 🧪 Testing Strategy

### Manual Testing (MVP)

1. **Single device**: Create room, join in incognito, test voting
2. **Multiple phones**: You + 2-3 friends, full game
3. **Reconnection**: Turn off WiFi mid-game, reconnect
4. **Edge cases**: Leave room, rejoin, host disconnects

### Future: Automated Tests

- Unit tests for game logic
- Integration tests for WebSocket flow
- E2E tests with Playwright

---

## 📝 Development Notes

### Backend

- **Module name**: Update `go.mod` with your actual repo URL
- **CORS**: Currently allows all origins (dev mode). Restrict in production.
- **Cleanup**: Runs hourly to delete stale rooms (24h old, no active players)

### Frontend

- **API proxy**: Vite proxies `/api` to backend (see `vite.config.ts`)
- **Mobile targets**: All touch targets 48x48px minimum
- **Offline**: Detects disconnect, shows banner, auto-reconnects

### Deployment

- **Railpack**: Railway's build system with automatic language detection
- **Monorepo**: Separate services for backend and frontend
- **Health checks**: Backend has `/health` endpoint
- **Environment**: Configure API_URL for frontend to connect to backend

---

## 🐛 Known Issues / TODOs

- [ ] Session token security (use httpOnly cookies in production)
- [ ] WebSocket authentication race condition handling
- [ ] Mobile Safari PWA installation flow
- [ ] Accessibility (ARIA labels, keyboard nav)
- [ ] Game config validation on frontend before sending
- [ ] Proper error messages (not just console.error)

---

## 🤝 Contributing

This is your personal project, but if you want to collaborate:

1. Fork the repo
2. Create a feature branch
3. Make your changes
4. Open a PR with clear description

---

## 📜 License

MIT (or your preferred license)

---

## 🎉 Credits

Built by a UX Engineer who loves party games and hates losing game pieces.

**Tech inspiration:**
- Jackbox Games (party game UX)
- Among Us (mobile-first social deduction)
- Netcode.io (real-time networking)

**Game inspiration:**
- One Night Werewolf (Bezier Games)
- Avalon (Don Eskridge)
- Bohnanza (Uwe Rosenberg)

---

## 📞 Support

Questions? Check `ARCHITECTURE.md` for detailed design.

Found a bug? Open an issue with:
- Steps to reproduce
- Expected vs actual behavior
- Browser/device info

---

**Happy gaming! 🎲**
