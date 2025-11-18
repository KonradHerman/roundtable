# 🎲 Cardless

[![Play Live](https://img.shields.io/badge/Play_Live-cardless.games-orange?style=for-the-badge&logo=gamepad)](https://cardless.games)

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

## 🎯 Current Status

**Working Features:**
- ✅ **Lobby System**: Create/Join rooms, real-time player list
- ✅ **Role Assignment**: Works with Werewolf rules (players + 3 center cards)
- ✅ **Game Loop**: Start game, Night phase (narration), Day phase (timer), Results
- ✅ **Architecture**: Event sourcing and game abstraction fully implemented
- ✅ **Night Actions**: Werewolf, Seer, Robber, Troublemaker, Drunk, Insomniac
- ✅ **Physical Cues**: Minion and Mason guided by physical instructions (no digital action needed)
- ✅ **Flow**: Play Again mechanics fully implemented

**In Progress:**
- 🚧 **UI Polish**: Better mobile touch targets and role reveal screens
- 🚧 **Physical Voting**: Transitioning from phone voting to physical pointing (traditional style)

---

## 🎲 Planned Games

We aim to add games that are perfect for "I wish we had the box right now" moments:

### 1. Spyfall (The Spy)
- **Concept**: Location-based social deduction.
- **Why**: Pure conversation, minimal UI (just a location card), 5-minute rounds.
- **Mechanic**: Everyone sees a location except the Spy. Questions/Answers to find the spy without revealing the location.

### 2. Skull (Roses & Skulls)
- **Concept**: Bluffing and bidding.
- **Why**: Replaces physical coasters/cards perfectly. High tension, simple rules.
- **Mechanic**: Each player has 3 Roses and 1 Skull. Bid how many you can flip without hitting a skull.

### 3. Love Letter (Courtship)
- **Concept**: Micro-deck deduction game.
- **Why**: Very few "cards" to manage, fast rotation.
- **Mechanic**: Draw one, play one. Knock out other players to get your letter to the Princess.

---

## 🗺️ Roadmap

For detailed technical tasks, see [TODO.md](TODO.md).
For future game concepts, see [GAMES.md](GAMES.md).

### 🛠️ Phase 1: Polish & Stability (Current)
- [ ] **Tech Debt**: Migrate Svelte 4 syntax (`$:`) to Svelte 5 Runes (`$state`, `$derived`) for future proofing.
- [ ] **Backend Optimization**: Implement priority queue for phase timeouts (currently O(N) polling).
- [ ] **Security**: Add robust CORS handling (`rs/cors`) for production.
- [ ] **UX**: Add better reconnection handling.

### 📦 Phase 2: New Games
- [ ] Implement Spyfall (validates game abstraction for non-card games)
- [ ] Implement Skull (validates betting/bidding mechanics)

### 🚢 Phase 3: Accounts & Stats
- [ ] Optional user accounts to track wins/losses
- [ ] Global leaderboards (optional)
- [ ] Replay system (visualize game history)

---

## 📁 Project Structure

```
roundtable/
├── backend/                    # Go backend
│   ├── cmd/server/            # Entry point
│   ├── internal/
│   │   ├── core/              # Platform core (game-agnostic)
│   │   ├── games/             # Game implementations (werewolf/)
│   │   ├── server/            # HTTP & WebSocket handlers
│   │   └── store/             # State persistence (memory/redis)
│   └── Dockerfile
│
├── frontend/                  # SvelteKit frontend
│   ├── src/
│   │   ├── routes/            # Pages (create, join, room)
│   │   ├── lib/               # Components, stores, API client
│   │   └── app.css            # Tailwind styles
│   └── Dockerfile
│
├── docker-compose.yml
└── ARCHITECTURE.md            # Detailed design doc
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.22+
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

## 🤝 Contributing

This is your personal project, but if you want to collaborate:

1. Fork the repo
2. Create a feature branch
3. Make your changes
4. Open a PR with clear description

---

## 📜 License

MIT

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
