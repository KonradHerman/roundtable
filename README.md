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

## 🎮 How to Play

1. **Host creates room**: Go to app → "Host a Game" → Choose Werewolf → Enter name → Get room code `XJ4K2P`
2. **Players join**: Others go to app → "Join Game" → Enter code `XJ4K2P` + name
3. **Host starts**: Configure roles (must be player count + 3) → "Start Game"
4. **Role reveal**: Each player privately views their role on their phone
5. **Night phase**: Roles wake in sequence and perform actions on their phones:
   - Werewolves see each other
   - Seer views another player's role or two center cards
   - Robber swaps cards with another player
   - Troublemaker swaps two other players' cards
   - Drunk swaps with a center card (doesn't see new role)
   - Insomniac sees their final role
6. **Day phase**: Discuss (in person!) with optional timer. Use your knowledge to figure out who the werewolves are!
7. **Vote**: Everyone simultaneously points at who to eliminate (traditional ONUW style)
8. **Reveal**: All players reveal their final roles on their phones to determine winner
9. **Play again**: Host can start a new game with same players

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

### 🔄 Phase 2: Werewolf MVP (IN PROGRESS)
**Completed:**
- [x] Lobby with player list
- [x] Role assignment (players + 3 center cards)
- [x] Role reveal with acknowledgements
- [x] Night phase with host narration script
- [x] Day phase with timer (start/pause/extend)
- [x] Results calculation logic
- [x] Game abstraction layer
- [x] Host start game flow

**Current Sprint:**
- [ ] Digital night actions (Seer, Robber, Troublemaker, Drunk, Insomniac)
- [ ] Remove phone voting (use physical voting instead)
- [ ] Role reveal screen (show all roles after discussion)
- [ ] Play again feature
- [ ] Fix host tracking

### 📦 Phase 3: Polish & Stability (3-4 days)
- [ ] QR code generation for room
- [ ] Reconnection handling
- [ ] Mobile UI polish (better touch targets)
- [ ] Error boundaries & loading states
- [ ] Edge case handling

### 🎯 Phase 4: Second Game - Avalon (4-5 days)
- [ ] Extract reusable game patterns
- [ ] Avalon game implementation
- [ ] Validate game abstraction works for different game types
- [ ] Mission voting UI
- [ ] Team selection UI
- [ ] Merlin/Assassin reveal

### 🚢 Phase 5: Production Ready (3-4 days)
- [ ] Redis store implementation
- [ ] Room expiry (cleanup)
- [ ] Rate limiting
- [ ] PWA manifest (installable)
- [ ] Homelab deployment guide

---

## 🎯 Current Status & Next Steps

**✅ Working Now:**
- Full lobby system with real-time player list
- Role assignment following ONUW rules (players + 3 center cards)
- Role reveal phase with acknowledgements
- Night phase with host narration script
- Day phase with timer controls
- Event sourcing architecture ready for multiple games

**🔨 Next Sprint (Digital Night Actions):**
1. **Backend**: Implement action handlers for each role (Seer, Robber, Troublemaker, Drunk, Insomniac)
2. **Frontend**: Create role-specific night phase UIs for each player
3. **Backend**: Track role swaps and send private results back to players
4. **Frontend**: Remove voting UI, replace with "reveal roles" screen
5. **Backend + Frontend**: Implement play again feature

**🧪 Then: Playtest with 6-8 people!**
- Test all night actions work correctly
- Verify physical voting flows naturally
- Check reconnection edge cases
- Polish based on feedback

---

## 🏛️ Key Design Decisions

### 1. Card Replacement, Not Full Digital Game
This app **replaces physical cards**, not in-person interaction:
- ✅ Role assignment and private viewing
- ✅ Digital night actions (prevents cheating, tracks swaps)
- ✅ Discussion timer
- ❌ **NOT** digital voting - voting is physical (everyone points)
- ❌ **NOT** automatic winner calculation - players determine this together

### 2. Why Event Sourcing?
- **Reconnection**: Replay events to rebuild state
- **Debugging**: Full audit trail
- **Multiple games**: Easy to add game #2, #3, etc.
- **Time travel**: Rewind/replay for undo features

### 3. Why Server-Authoritative?
- **No cheating**: Server holds canonical state
- **Security**: Clients can't see hidden info
- **Consistency**: Single source of truth

### 4. Why Game Abstraction?
- **Extensibility**: Adding game #2 is easy
- **Isolation**: Game-specific logic doesn't leak into core
- **Testing**: Test games independently

### 5. Why Anonymous-First?
- **Instant play**: No signup friction
- **Privacy**: Just display names
- **Optional accounts**: Add later for stats

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

**Critical for MVP:**
- [ ] Digital night actions not yet implemented (Seer, Robber, Troublemaker, Drunk, Insomniac)
- [ ] Phone voting needs to be removed (use physical voting instead)
- [ ] Play again feature not implemented
- [ ] Host tracking uses first player instead of actual host
- [ ] Reconnection handling not tested

**Polish Items:**
- [ ] QR code generation for room sharing
- [ ] Session token security (use httpOnly cookies in production)
- [ ] Mobile Safari PWA installation flow
- [ ] Error boundaries and loading states
- [ ] Accessibility (ARIA labels, keyboard nav)
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
