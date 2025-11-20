# 🎲 Cardless


[![Play Live](https://img.shields.io/badge/Play_Live-cardless.games-orange?style=for-the-badge&logo=gamepad)](https://cardless.games)

**Play party games without the cards.**

Cardless replaces physical cards with your phone for in-person party games. Create a room, share the code, and start playing—no shuffling, no lost pieces, no cheating.

---

## What is Cardless?

Cardless brings the best tabletop party games to your phone while keeping the social, in-person experience that makes them fun. Think of it like Jackbox, but for social deduction and bluffing games.

**Currently available:**
- **One Night Werewolf** - Social deduction with secret roles and night actions

**Coming soon:**
- Avalon, Spyfall, Skull, Wavelength, and more!

### Why Cardless?

- **✨ Instant setup**: No shuffling, dealing, or organizing components
- **🎭 Perfect information**: Digital night actions prevent accidental peeks and track role swaps
- **📱 Anonymous play**: No account required—just share a room code
- **🏠 In-person focused**: Designed for same-room play with friends
- **🔒 Privacy-first**: No personal data collection, rooms expire automatically

---

## 🎮 How It Works

1. **Host creates a room** - Choose a game and get a 6-character room code
2. **Friends join** - Enter the code and a display name (no account needed)
3. **Play together** - Each phone becomes your private "hand" of cards
4. **Discuss in person** - The real magic happens face-to-face

### One Night Werewolf Features

- 🎭 **Secret role assignment** - Everyone gets a hidden role
- 🌙 **Digital night actions** - Seer, Robber, Troublemaker, and more
- ⏱️ **Discussion timer** - Keep conversations focused
- 👉 **Physical voting** - Point at suspects together (preserves the drama!)
- 🔄 **Play again** - Same room, new game instantly

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

### Play Online

Visit **[cardless.games](https://cardless.games)** to start playing immediately!

### Run Locally

Want to self-host or contribute? See [DEVELOPMENT.md](DEVELOPMENT.md) for detailed setup instructions.

**Quick version:**
```bash
# Backend (Go 1.21+)
cd backend
go run cmd/server/main.go

# Frontend (Node.js 20+)
cd frontend
npm install
npm run dev
```

Then visit `http://localhost:5173`

---

## 🎯 Current Status

**✅ Working Now:**
- Full One Night Werewolf implementation
- Real-time multiplayer with WebSockets
- Room creation and joining with 6-character codes
- Role assignment and reveal
- Night phase with host narration script
- Day phase with discussion timer
- Event sourcing architecture for game state

**🔨 In Progress:**
- Digital night actions for all roles (Seer, Robber, Troublemaker, Drunk, Insomniac)
- Role reveal screen showing final roles
- Play again feature

**📋 Planned:**
- QR code room sharing
- Reconnection handling
- Additional games (Avalon, Spyfall, Skull, Wavelength)

See [ROADMAP.md](ROADMAP.md) for detailed development plan.

---

## 🗺️ Roadmap

We're building Cardless incrementally, validating the platform with each new game.

### Phase 1: Foundation ✅
Core platform with room management, WebSocket real-time communication, and event sourcing architecture.

### Phase 2: One Night Werewolf 🔄
First complete game implementation to validate the platform. Currently completing digital night actions.

### Phase 3: Polish & Stability 📋
QR code sharing, reconnection handling, mobile UI polish, and production-ready features.

### Phase 4: Additional Games 🎲
- **Avalon** - Quest voting and team selection
- **Spyfall** - Location deduction with question rounds
- **Skull** - Bluffing and bidding mechanics
- **Wavelength** - Team-based spectrum guessing

See [ROADMAP.md](ROADMAP.md) and [GAMES_ROADMAP.md](GAMES_ROADMAP.md) for detailed plans.

---

## 🏗️ Architecture Highlights

Cardless is built with scalability and extensibility in mind:

### Event Sourcing
Game state is derived from a sequence of events, enabling:
- Seamless reconnection (replay events to rebuild state)
- Full audit trail for debugging
- Easy addition of spectator mode or game replays

### Game Abstraction Layer
Each game implements a common interface, making it straightforward to add new games without changing the core platform.

### Server-Authoritative Design
The server validates all actions and maintains the canonical game state—clients can't cheat or see hidden information.

### Anonymous-First
No accounts required. Players join with a room code and display name, keeping the friction low and privacy high.

For technical details, see [ARCHITECTURE.md](ARCHITECTURE.md).

---

## 👥 For Developers

### Contributing

Want to add a new game or improve the platform?

1. Fork the repository
2. Check out [DEVELOPMENT.md](DEVELOPMENT.md) for setup
3. Review [ARCHITECTURE.md](ARCHITECTURE.md) to understand the design
4. See [GAMES_ROADMAP.md](GAMES_ROADMAP.md) for game ideas
5. Submit a pull request

### Tech Stack

- **Backend**: Go 1.21+ with goroutines for concurrency
- **Frontend**: SvelteKit with Svelte 5 (migrating)
- **Real-time**: Native WebSockets (nhooyr.io/websocket)
- **State**: Event sourcing with in-memory store (Redis later)
- **Styling**: Tailwind CSS for mobile-first design

### Project Structure

```
roundtable/
├── backend/              # Go backend
│   ├── cmd/server/      # Entry point
│   ├── internal/
│   │   ├── core/        # Platform core (game-agnostic)
│   │   ├── games/       # Game implementations
│   │   ├── server/      # HTTP & WebSocket
│   │   └── store/       # State persistence
│
├── frontend/            # SvelteKit frontend
│   ├── src/
│   │   ├── routes/      # Pages
│   │   ├── lib/
│   │   │   ├── stores/  # State management
│   │   │   ├── games/   # Game-specific components
│   │   │   └── components/  # Reusable UI
```

---

## 📜 License

[CC BY-NC 4.0](LICENSE) - Creative Commons Attribution-NonCommercial 4.0 International

Free for personal use. Contact for commercial licensing.

---

## 📞 Get in Touch

- 🐛 Found a bug? [Open an issue](https://github.com/yourusername/cardless/issues)
- 💡 Have a game idea? Check [GAMES_ROADMAP.md](GAMES_ROADMAP.md)
- 🎮 Want to contribute? See [DEVELOPMENT.md](DEVELOPMENT.md)

---

**🎲 [Start Playing Now](https://cardless.games)**
