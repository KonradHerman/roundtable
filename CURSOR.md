# CURSOR.md - AI Assistant Guide for Cardless

> **Last Updated**: 2025-11-20
> **Status**: One Night Werewolf MVP in progress (Phase 2)
> **Live at**: [cardless.games](https://cardless.games)

This document provides Cursor IDE's AI assistant with essential context about the Cardless codebase, including architecture, conventions, and workflows.

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Quick Reference](#quick-reference)
3. [Codebase Structure](#codebase-structure)
4. [Architecture & Design Patterns](#architecture--design-patterns)
5. [Technology Stack](#technology-stack)
6. [Development Workflow](#development-workflow)
7. [Key Conventions & Guidelines](#key-conventions--guidelines)
8. [Common Tasks](#common-tasks)
9. [Testing Strategy](#testing-strategy)
10. [Deployment](#deployment)
11. [Important Context](#important-context)

---

## Project Overview

**Cardless** is a real-time multiplayer platform for in-person party games. It replaces physical cards with phones while preserving the social, face-to-face nature of games like One Night Werewolf, Avalon, and others.

### Core Philosophy
- **In-person first**: App replaces cards, NOT social interaction
- **Physical voting preserved**: No digital voting—players point simultaneously in real life
- **Anonymous play**: No accounts required for MVP
- **Privacy-first**: No data collection, rooms auto-expire after 24 hours
- **Mobile-first**: Designed primarily for phone touchscreens

### Current Status
- ✅ **Phase 1 Complete**: Lobby, WebSocket communication, role assignment, narration
- 🔨 **Phase 2 In Progress**: Digital night actions, role reveal, play again
- 📋 **Phase 3 Planned**: Polish, Svelte 5 migration, additional games

---

## Quick Reference

### File Locations

| What | Where |
|------|-------|
| Backend entry point | `backend/cmd/server/main.go` |
| Game interface | `backend/internal/core/game.go` |
| Werewolf game logic | `backend/internal/games/werewolf/game.go` |
| HTTP/WebSocket handlers | `backend/internal/server/handlers.go`, `websocket.go` |
| Frontend routes | `frontend/src/routes/` |
| WebSocket client | `frontend/src/lib/stores/websocket.svelte.ts` |
| Game state store | `frontend/src/lib/stores/game.svelte.ts` |
| Werewolf UI components | `frontend/src/lib/games/werewolf/` |
| API client | `frontend/src/lib/api/client.ts` |
| Style guide | `STYLE_GUIDE.md` |
| Architecture docs | `ARCHITECTURE.md` |
| Development setup | `DEVELOPMENT.md` |

### Common Commands

```bash
# Start backend (port 8080)
cd backend && go run cmd/server/main.go

# Start frontend (port 5173)
cd frontend && npm run dev

# Run backend tests (when they exist)
cd backend && go test ./...

# Build for production
cd frontend && npm run build
cd backend && CGO_ENABLED=0 GOOS=linux go build -a -o server ./cmd/server
```

### Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `PORT` | 8080 | Backend server port |
| `ALLOWED_ORIGIN` | http://localhost:5173 | CORS allowed origin |
| `VITE_API_URL` | /api | Frontend API URL (proxied in dev) |

---

## Codebase Structure

```
/home/user/roundtable/
├── backend/                      # Go backend service
│   ├── cmd/server/              # Application entry point
│   │   └── main.go              # HTTP server, routing, background tasks
│   ├── internal/
│   │   ├── core/                # Game-agnostic platform core
│   │   │   ├── game.go          # Game interface (80 lines)
│   │   │   ├── room.go          # Room management (366 lines)
│   │   │   ├── player.go        # Player entity
│   │   │   └── event.go         # Event types and visibility
│   │   ├── games/               # Game implementations
│   │   │   ├── registry.go      # Game factory/registration
│   │   │   ├── werewolf/        # One Night Werewolf (1,240 lines)
│   │   │   │   ├── game.go      # Main game logic (715 lines)
│   │   │   │   ├── state.go     # State structures (133 lines)
│   │   │   │   ├── config.go    # Configuration (99 lines)
│   │   │   │   ├── phases.go    # Phase transitions (188 lines)
│   │   │   │   └── narration.go # Host narration (105 lines)
│   │   │   └── avalon/          # Avalon (in progress)
│   │   ├── server/              # HTTP & WebSocket layer
│   │   │   ├── handlers.go      # REST API endpoints
│   │   │   ├── websocket.go     # WebSocket connection manager
│   │   │   └── messages.go      # Message types
│   │   ├── store/               # State persistence
│   │   │   ├── store.go         # Interface definition
│   │   │   └── memory.go        # In-memory implementation
│   │   └── util/                # Utilities
│   │       └── codegen.go       # 6-character room codes
│   ├── go.mod                   # Go dependencies
│   └── Dockerfile               # Multi-stage Docker build
│
├── frontend/                     # SvelteKit frontend
│   ├── src/
│   │   ├── routes/              # File-based routing
│   │   │   ├── +page.svelte     # Landing page
│   │   │   ├── +layout.svelte   # Root layout
│   │   │   ├── create/          # Room creation
│   │   │   ├── join/            # Join pages
│   │   │   └── room/[code]/     # Dynamic game room
│   │   ├── lib/
│   │   │   ├── api/
│   │   │   │   └── client.ts    # REST API client (146 lines)
│   │   │   ├── components/ui/   # Reusable UI components
│   │   │   │   ├── button.svelte
│   │   │   │   ├── card.svelte
│   │   │   │   ├── badge.svelte
│   │   │   │   └── CardBack.svelte
│   │   │   ├── games/           # Game-specific components
│   │   │   │   ├── werewolf/    # Werewolf UI (large files)
│   │   │   │   │   ├── WerewolfGame.svelte (3,709 lines)
│   │   │   │   │   ├── RoleReveal.svelte (4,076 lines)
│   │   │   │   │   ├── NightPhase.svelte (17,068 lines)
│   │   │   │   │   ├── DayPhase.svelte (4,888 lines)
│   │   │   │   │   ├── Results.svelte (5,486 lines)
│   │   │   │   │   ├── RoleCard.svelte (1,157 lines)
│   │   │   │   │   └── roleConfig.ts (1,891 lines)
│   │   │   │   └── avalon/      # Avalon (early stage)
│   │   │   └── stores/          # State management
│   │   │       ├── game.svelte.ts (Svelte 5 version)
│   │   │       ├── session.svelte.ts
│   │   │       └── websocket.svelte.ts (3,582 lines)
│   │   └── static/              # Static assets, favicons
│   ├── package.json
│   ├── vite.config.ts           # Vite config with API proxy
│   ├── svelte.config.js         # Path aliases, adapters
│   ├── tailwind.config.js       # Gruvbox theme
│   └── tsconfig.json
│
└── [Documentation Files]         # 10 major .md files
    ├── ARCHITECTURE.md          # Detailed technical architecture (830 lines)
    ├── DEVELOPMENT.md           # Developer setup guide
    ├── TESTING.md               # Testing strategies
    ├── ROADMAP.md               # Development roadmap
    ├── GAMES_ROADMAP.md         # Game implementation plans
    ├── STYLE_GUIDE.md           # Gruvbox color palette
    ├── WEREWOLF_RULES.md        # One Night Werewolf rules
    ├── RAILWAY.md               # Railway deployment guide
    ├── FUTURE_IMPROVEMENTS.md   # Reverted features
    ├── CLAUDE.md                # AI assistant guide
    └── CURSOR.md                # This file - Cursor IDE guide
```

---

## Architecture & Design Patterns

### 1. Event Sourcing (Core Principle)

**Game state is derived from an append-only log of events:**

```
Current State = Initial State + [Event1, Event2, Event3, ...]
```

**Key Benefits:**
- **Reconnection**: Replay events to rebuild player state seamlessly
- **Audit trail**: Full game history for debugging and stats
- **Spectator mode**: Subscribe to event stream (future feature)
- **Time travel**: Rewind/replay for debugging
- **Multiple views**: Same events, different perspectives per player

**Implementation:**
- Room maintains `EventLog []GameEvent` (append-only)
- Events have visibility controls (public vs. private)
- Events filtered per player based on permissions
- WebSocket delivers events to connected clients

**File references:**
- Backend: `backend/internal/core/event.go`
- Frontend: `frontend/src/lib/stores/websocket.svelte.ts`

---

### 2. Game Abstraction Layer

**Every game implements the `core.Game` interface:**

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
    CheckPhaseTimeout() ([]GameEvent, error)
}
```

**Adding a new game requires:**
1. Create `internal/games/<gamename>/` package
2. Implement the `Game` interface
3. Register in `games/registry.go`
4. Create frontend components in `frontend/src/lib/games/<gamename>/`

**No changes to platform core needed.**

**File references:**
- Interface: `backend/internal/core/game.go:80`
- Werewolf implementation: `backend/internal/games/werewolf/game.go:715`
- Registry: `backend/internal/games/registry.go`

---

### 3. Server-Authoritative Design

**Clients send actions (intent), server validates and broadcasts events (results).**

```
Phone → WebSocket → Server validates → Process action → Emit events → Broadcast to clients
```

**Never trust the client:**
- Server holds canonical state
- Clients send "I want to view Alice's card" (intent)
- Server validates: "Is this player the Seer? Have they already used their action?"
- Server processes and broadcasts: "Seer viewed Alice: she is a Werewolf" (result)

**Prevents cheating and ensures consistent game state.**

**File references:**
- Server handlers: `backend/internal/server/handlers.go`
- WebSocket manager: `backend/internal/server/websocket.go`
- Action processing: `backend/internal/games/werewolf/game.go` (ProcessAction method)

---

### 4. Anonymous Session Management

**No accounts required:**
- Session tokens generated on first visit
- Stored in `localStorage` (frontend)
- Sent via HTTP headers and WebSocket handshake
- Enables instant play without signup friction

**File references:**
- Session store: `frontend/src/lib/stores/session.svelte.ts`
- API client: `frontend/src/lib/api/client.ts:146`

---

### 5. Concurrent State Management

**Go backend uses `sync.RWMutex` for safe concurrent access:**
- Room and Player structs have embedded mutexes
- Read locks for queries, write locks for mutations
- Safe access from multiple goroutines (WebSocket connections, background tasks)

**File references:**
- Room concurrency: `backend/internal/core/room.go:366`

---

### 6. Background Processing

**Two main background goroutines:**
1. **`cleanupRoutine`**: Removes stale rooms (1 hour interval)
2. **`phaseCheckRoutine`**: Checks game timers (1 second interval)

Both respect context cancellation for graceful shutdown.

**File references:**
- Background tasks: `backend/cmd/server/main.go`

---

## Technology Stack

### Backend
- **Language**: Go 1.22
- **HTTP Server**: `net/http` (standard library)
- **WebSocket**: `nhooyr.io/websocket` v1.8.10 (modern, context-aware)
- **Logging**: `log/slog` (structured JSON logging)
- **Concurrency**: Goroutines, channels, `sync.RWMutex`
- **UUID**: `github.com/google/uuid` v1.6.0
- **Deployment**: Multi-stage Docker build (alpine-based)

### Frontend
- **Framework**: SvelteKit 2.5.0 with Svelte 5.0.0 (migration in progress)
- **Language**: TypeScript 5.3.3
- **Build Tool**: Vite 5.1.0
- **Adapter**: `@sveltejs/adapter-node` 5.0.1
- **Styling**: Tailwind CSS 3.4.1 with Gruvbox Dark theme
- **UI Libraries**:
  - `bits-ui` 0.21.13 (headless components)
  - `lucide-svelte` 0.447.0 (icons)
  - `qr-code-styling` 1.9.2 (QR code generation)
  - `@neoconfetti/svelte` 2.2.1 (confetti effects)

### Development Tools
- **Type Checking**: `svelte-check` 3.6.4
- **PostCSS**: 8.4.35 with autoprefixer
- **Preprocessing**: vitePreprocess

### Production Infrastructure
- **Deployment**: Railway (two separate services)
- **Builder**: Railpack (auto-detected)
- **Future**: Redis for state persistence (currently in-memory)

---

## Development Workflow

### Local Development Setup

**Prerequisites:**
- Go 1.22+
- Node.js 18+
- Two terminal windows

**Step 1: Start Backend**
```bash
cd backend
go run cmd/server/main.go
# Runs on http://localhost:8080
# Look for: "Server starting on port 8080"
```

**Step 2: Start Frontend**
```bash
cd frontend
npm install  # First time only
npm run dev
# Runs on http://localhost:5173
# Look for: "Local: http://localhost:5173/"
```

**Step 3: Access Application**
- Open browser: `http://localhost:5173`
- Vite proxies `/api/*` to `http://localhost:8080`
- WebSocket proxying enabled automatically

**⚠️ Important**: Do NOT access `localhost:8080` directly—that's the backend API, not the frontend.

---

### Testing Multi-Player Flows

**Method 1: Incognito Windows (Fastest)**
1. Open app in normal browser window
2. Create a room (become host)
3. Open 2-7 incognito windows (`Ctrl+Shift+N` in Chrome)
4. Join room from each window with the room code
5. Test game flow

**Method 2: Multiple Devices (Most Realistic)**
1. Find local IP: `ipconfig` (Windows) or `ifconfig` (Mac/Linux)
2. Ensure backend is running on `localhost:8080`
3. Update `frontend/vite.config.ts` proxy target to your IP (if needed)
4. Connect phones/tablets to same WiFi
5. Access `http://192.168.x.x:5173` from each device

**Method 3: Different Browsers**
- Chrome, Firefox, Edge, Safari (Mac)
- Each browser has separate session storage

**Reference**: See `DEVELOPMENT.md` and `TESTING.md` for detailed guidance.

---

### Hot Reload Behavior

| Layer | Hot Reload? | How to Restart |
|-------|-------------|----------------|
| Frontend | ✅ Yes (Vite) | Automatic on file save |
| Backend | ❌ No | `Ctrl+C` and re-run `go run cmd/server/main.go` |

---

### Debugging Tools

- **Backend logs**: Structured JSON via `log/slog` in terminal
- **Frontend logs**: Browser DevTools console
- **API requests**: Logged with full URLs in browser console
- **WebSocket**: Browser DevTools → Network tab → WS filter
- **Game state**: Check `gameStore` in console
- **Health check**: `curl http://localhost:8080/health` → "OK"

---

## Key Conventions & Guidelines

### Code Style

**Backend (Go):**
- Follow standard Go conventions (`gofmt`, `golint`)
- Use structured logging with `slog` (JSON format)
- Mutex-based concurrency (read/write locks)
- Error handling: return errors, don't panic (except in main setup)
- Context-aware: respect `context.Context` for cancellation

**Frontend (TypeScript/Svelte):**
- **In transition**: Svelte 4 → Svelte 5 (dual implementation exists)
- Use TypeScript strictly (no `any` types)
- Path aliases: `$components`, `$stores`, `$games`, `$api`
- Component naming: PascalCase (e.g., `RoleCard.svelte`)
- State management: Svelte runes (`$state`, `$derived`, `$effect`) in new code

---

### Design System (Gruvbox Dark)

**See `STYLE_GUIDE.md` for full details.**

**Key Colors (HSL):**
- **Background**: `#282828` (HSL: 30 9% 15.7%)
- **Foreground**: `#ebdbb2` (HSL: 42 19% 85.3%)
- **Primary**: `#d79921` (gold/orange) - CTAs, focus states
- **Secondary**: `#458588` (teal/blue) - secondary actions
- **Destructive**: `#cc241d` (red) - errors
- **Input Background**: `#f9f5d7` (light cream)
- **Input Text**: `#1a1a1a` (nearly black for contrast)

**Typography:**
- Body: `-apple-system, BlinkMacSystemFont, 'Segoe UI', ...`
- Monospace (room codes): `font-mono`

**Spacing:**
- Minimum touch target: **48x48px** (mobile accessibility)
- Card padding: `1.5rem` (24px)
- Border radius: `0.75rem` (cards: `1rem`)

**Accessibility:**
- All text combinations meet WCAG AAA contrast standards
- Focus rings: 2px `#d79921` (primary color)
- Touch-friendly: 48px minimum targets

---

### File Naming

| Type | Convention | Example |
|------|-----------|---------|
| Go packages | lowercase | `internal/games/werewolf/` |
| Go files | snake_case | `room.go`, `websocket.go` |
| Svelte components | PascalCase | `RoleCard.svelte` |
| TypeScript files | camelCase | `client.ts`, `roleConfig.ts` |
| Routes | kebab-case (SvelteKit convention) | `+page.svelte`, `[code]/+page.svelte` |

---

### Git Workflow

**Branches:**
- `main` - production branch
- `claude/*` - AI assistant feature branches (e.g., `claude/add-avalon-game-01ABC123`)

**Commits:**
- Use conventional commits style: `feat:`, `fix:`, `docs:`, `refactor:`, etc.
- Example: `feat: implement Seer night action for Werewolf`
- Keep commits focused and atomic

**Pull Requests:**
- Include clear description of changes
- Reference related issues if applicable
- Test manually before creating PR

---

### API Conventions

**Backend REST Endpoints:**
- `POST /api/rooms` - Create room
- `GET /api/rooms/:code` - Get room details
- `POST /api/rooms/:code/join` - Join room
- `POST /api/rooms/:code/start` - Start game
- `GET /health` - Health check

**WebSocket Messages:**
- Client → Server: JSON with `{ type: "action", payload: {...} }`
- Server → Client: JSON with `{ type: "event", payload: {...} }`
- Events are game-specific (defined per game implementation)

**File references:**
- REST handlers: `backend/internal/server/handlers.go`
- WebSocket manager: `backend/internal/server/websocket.go`
- Message types: `backend/internal/server/messages.go`

---

### Error Handling

**Backend:**
- Return errors up the call stack
- Log errors with context using `slog`
- Send appropriate HTTP status codes (400, 404, 500)
- Include error messages in API responses

**Frontend:**
- API client wraps errors with enhanced messages
- Display user-friendly error messages in UI
- Log errors to console for debugging
- Graceful degradation on connection loss

**File references:**
- API error handling: `frontend/src/lib/api/client.ts:146`

---

## Common Tasks

### Adding a New Game

**Backend:**
1. Create `internal/games/<gamename>/` directory
2. Implement `core.Game` interface in `game.go`
3. Add supporting files: `state.go`, `config.go`, etc.
4. Register in `internal/games/registry.go`

**Frontend:**
1. Create `src/lib/games/<gamename>/` directory
2. Create main game component (e.g., `AvalongGame.svelte`)
3. Create phase-specific components
4. Add role configuration file (if applicable)
5. Update route handler in `src/routes/room/[code]/+page.svelte`

**Reference**: See `ARCHITECTURE.md` for detailed game implementation guide.

---

### Adding a New API Endpoint

**Backend:**
1. Add handler function in `internal/server/handlers.go`
2. Register route in `cmd/server/main.go`
3. Add request/response types if needed

**Frontend:**
1. Add method to `src/lib/api/client.ts`
2. Use TypeScript interfaces for type safety

---

### Modifying Werewolf Game Logic

**Role-specific actions:**
- Edit `backend/internal/games/werewolf/game.go` (ProcessAction method)
- Update state in `state.go` if needed
- Modify corresponding UI in `frontend/src/lib/games/werewolf/NightPhase.svelte`

**Phase transitions:**
- Edit `backend/internal/games/werewolf/phases.go`
- Update narration in `narration.go` if needed

**Configuration:**
- Edit `backend/internal/games/werewolf/config.go` for new settings

---

### Updating Styles

**Global styles:**
- Edit `frontend/tailwind.config.js` for theme changes
- Follow Gruvbox palette in `STYLE_GUIDE.md`

**Component-specific:**
- Use Tailwind utility classes
- Extract to component classes if needed
- Minimum touch target: 48x48px

---

### Adding a WebSocket Event

**Backend:**
1. Define event type in game's `state.go`
2. Emit event in `ProcessAction` or phase transition
3. Set visibility (public or private per player)

**Frontend:**
1. Handle event in `src/lib/stores/websocket.svelte.ts`
2. Update UI based on event type
3. Test reconnection (events should replay correctly)

---

## Testing Strategy

### Current Approach: Manual Testing Only

**No automated tests exist yet** (intentional for MVP).

**Testing philosophy:**
1. Manual testing sufficient for MVP
2. Add automation when it speeds up development
3. Focus on critical paths

### Manual Testing Checklist

**Critical Path:**
- ✅ Create room → Join → Configure roles → Start game
- ✅ Role reveal phase → All players acknowledge
- ✅ Night phase → Role-specific actions work
- ✅ Day phase → Timer functions correctly
- ✅ Role reveal → Final roles displayed
- ✅ Play again → Room resets properly

**Edge Cases:**
- Disconnect/reconnect during game
- Invalid actions (wrong role, wrong phase)
- Concurrent actions from multiple players
- Browser compatibility (Chrome, Firefox, Safari)
- Mobile-specific issues (touch targets, viewport)

### Future Automation (Planned)

**Phase 3: Backend Unit Tests**
- Use Go's built-in `testing` package
- Test game logic in isolation
- Target: 80%+ coverage

**Phase 4: E2E Tests**
- Playwright (recommended)
- Multi-session testing (simulate multiple players)
- Critical user journeys
- Regression tests

**Load Testing:**
- WebSocket client simulations
- Tools: k6, artillery

**Reference**: See `TESTING.md` for detailed testing guide.

---

## Deployment

### Railway (Production)

**Two separate services:**

**Backend Service:**
- Root directory: `backend/`
- Builder: Railpack (auto-detected)
- Build command: Auto (`go build`)
- Start command: `./server`
- Environment: `PORT` (set by Railway)

**Frontend Service:**
- Root directory: `frontend/`
- Builder: Railpack
- Build command: `npm run build`
- Start command: `npm start`
- Environment: `VITE_API_URL` (must include `http://` prefix)
  - Example: `http://cardless.railway.internal:8080/api`

**Common Issues:**
- Missing `http://` prefix in `VITE_API_URL` causes CORS errors
- Must set root directory BEFORE first build
- Railway may default to Docker (requires service recreation)

**Reference**: See `RAILWAY.md` for deployment guide.

---

### Docker Build

**Backend:**
```dockerfile
# Multi-stage build
# Builder: golang:1.22-alpine
# Final: alpine:latest
```

Build command:
```bash
cd backend
docker build -t cardless-backend .
```

**Frontend:**
```bash
cd frontend
npm run build
# Output: ./build directory
```

---

### Environment Configuration

**Production:**
- Backend: Set `PORT` via environment variable
- Frontend: Set `VITE_API_URL` to backend URL with `/api` suffix
- CORS: Restrict `ALLOWED_ORIGIN` to production domain

---

## Important Context

### What This App Does (and Doesn't Do)

**✅ Does:**
- Replace physical cards with phones
- Handle role assignment and private information
- Provide digital night actions (prevents cheating)
- Track role swaps accurately
- Provide discussion timer
- Enable instant room creation and joining

**❌ Doesn't:**
- Digital voting (voting remains physical—players point simultaneously)
- Automatic winner calculation (players discuss and determine)
- Replace in-person social interaction (designed for same-room play)

**This preserves the social joy of party games while solving card problems.**

---

### Current Development Phase

**Phase 2: Digital Night Actions (In Progress)**

**Backend tasks remaining:**
- [ ] Implement night actions for all roles (Werewolf, Seer, Robber, Troublemaker, Drunk, Insomniac)
- [ ] Fix host tracking (use actual host from room creation)
- [ ] Implement play again endpoint

**Frontend tasks remaining:**
- [ ] Create role-specific night phase UIs for each role
- [ ] Create role reveal screen (show final roles)
- [ ] Add play again button and flow
- [ ] Remove voting UI (switch to physical voting)

**Reference**: See `ROADMAP.md` for full Phase 2 details.

---

### Known Limitations / Future Plans

**State Persistence:**
- Currently in-memory (rooms lost on server restart)
- Planned: Redis for production persistence

**Scalability:**
- Single server instance
- Planned: Horizontal scaling with Redis

**Authentication:**
- Anonymous sessions only
- Planned: Optional accounts for stats tracking

**Games:**
- Only Werewolf currently playable
- Avalon in progress
- More games planned (Spyfall, Skull, Wavelength)

**Reference**: See `FUTURE_IMPROVEMENTS.md` for detailed plans.

---

### Migration Notes: Svelte 4 → Svelte 5

**Current state:**
- Dual implementation exists for stores (`.ts` and `.svelte.ts`)
- New code should use Svelte 5 runes (`$state`, `$derived`, `$effect`)
- Migration planned for Phase 3

**When working with state:**
- Prefer `.svelte.ts` files (Svelte 5)
- Use runes for reactivity
- Avoid `.ts` versions (legacy Svelte 4)

---

### Project Naming

- **Code name**: "roundtable" (repo name, internal references)
- **Public name**: "Cardless"
- **Domain**: cardless.games
- **Tagline**: "Play party games without the cards"

---

### Documentation Files

| File | Purpose |
|------|---------|
| `ARCHITECTURE.md` | Detailed technical architecture (830 lines) |
| `DEVELOPMENT.md` | Local development setup guide |
| `TESTING.md` | Manual and automated testing strategies |
| `ROADMAP.md` | Development roadmap with phases |
| `GAMES_ROADMAP.md` | Game implementation plans |
| `STYLE_GUIDE.md` | Gruvbox color palette and design system |
| `WEREWOLF_RULES.md` | One Night Werewolf rules reference |
| `RAILWAY.md` | Railway deployment guide |
| `FUTURE_IMPROVEMENTS.md` | Reverted features to reimplement |
| `CLAUDE.md` | AI assistant guide (Claude) |
| `CURSOR.md` | **This file** - Cursor IDE guide |

---

### License

**CC-BY-NC-4.0** (Creative Commons Attribution-NonCommercial 4.0)

- ✅ Free for personal, educational, non-commercial use
- ❌ Commercial use requires permission
- ✅ Attribution required

---

## Final Notes for AI Assistants

### When Making Changes

1. **Always read relevant documentation first** (e.g., `ARCHITECTURE.md`, `STYLE_GUIDE.md`)
2. **Follow existing patterns** (event sourcing, game interface, server-authoritative design)
3. **Test manually with multiple sessions** (incognito windows or devices)
4. **Respect the design philosophy** (in-person first, physical voting, privacy-first)
5. **Use TypeScript strictly** (no `any` types)
6. **Follow Gruvbox color palette** (see `STYLE_GUIDE.md`)
7. **Maintain mobile-first approach** (48px touch targets)
8. **Add context to commits** (conventional commits style)

### When Asking for Clarification

**If uncertain about:**
- **Architecture decisions**: Reference `ARCHITECTURE.md`
- **Game rules**: Reference `WEREWOLF_RULES.md` or ask user
- **Design choices**: Reference `STYLE_GUIDE.md`
- **Current roadmap**: Reference `ROADMAP.md`
- **Testing approach**: Reference `TESTING.md`

### When Adding New Features

1. Check `ROADMAP.md` to see if it's planned
2. Check `FUTURE_IMPROVEMENTS.md` for reverted features
3. Follow the game abstraction layer pattern
4. Maintain event sourcing architecture
5. Test reconnection scenarios
6. Ensure mobile responsiveness

---

**Questions?** Check the documentation files first, then ask the user for clarification.

**Happy coding! 🎲**
