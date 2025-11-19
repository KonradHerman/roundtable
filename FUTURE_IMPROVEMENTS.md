# Future Improvements to Reimplement

This document records infrastructure improvements that were reverted on 2025-11-19 due to deployment issues. These have been successfully reimplemented on 2025-11-19 following the incremental migration plan.

## ✅ COMPLETED: Svelte 5 Migration + Go Module Path Fix (2025-11-19)

The following improvements have been successfully implemented:
- **Go Module Path Fix**: Updated all imports from `github.com/yourusername/roundtable` to `github.com/KonradHerman/roundtable`
- **Phase A - Store Migration**: Migrated game.ts, session.ts, and websocket.ts to Svelte 5 Runes (.svelte.ts files)
- **Phase B - Layout & Pages**: Updated all route pages to use new store system
- **Phase C - Werewolf Components**: Migrated all game components to Svelte 5 props and reactivity
- **Phase D - UI Components**: Updated Button, Card, and Badge components to use Svelte 5 children rendering

## Reverted Commits (Now Reimplemented)

The following commits were rolled back from commit c1fadaa1fc1786c168184e8abaf285cd2a99a771:

- **18dc401** - Documentation cleanup
- **cbd636d** - go.mod update
- **49e6781** - Svelte 5 migration + backend optimizations
- **e86e4ce** - Svelte 5 connection fixes
- **265a08a** - Railway deployment improvements
- **da1d62d** - Merge commit

## Key Improvements to Reimplement

### 1. Svelte 5 Migration ✅

**Status**: COMPLETED (2025-11-19) - Successfully migrated incrementally with testing at each phase

**What was done:**
- Migrated from Svelte 4 `$:` reactive syntax to Svelte 5 Runes system
- Updated reactive declarations: `$:` → `$derived`
- Updated state management: `let` → `$state`
- Updated props: destructuring → `$props()`
- Fixed store subscription patterns in `+layout.svelte`
- Converted `lib/stores/game.ts` to `lib/stores/game.svelte.ts` for Svelte 5 compatibility
- Updated all werewolf components:
  - `CenterCardSelect.svelte`
  - `DayPhase.svelte`
  - `NightPhase.svelte`
  - `PlayerCardSelect.svelte`
  - `Results.svelte`
  - `WerewolfGame.svelte`
- Fixed props type syntax: `let { children }: Props = $props()` pattern
- Resolved event processing bugs and deprecation warnings

**Files updated:**
```
frontend/src/lib/stores/game.svelte.ts (created)
frontend/src/lib/stores/session.svelte.ts (created)
frontend/src/lib/stores/websocket.svelte.ts (created)
frontend/src/lib/stores/game.ts (kept for now, can be deleted after verification)
frontend/src/lib/stores/session.ts (kept for now, can be deleted after verification)
frontend/src/lib/stores/websocket.ts (kept for now, can be deleted after verification)
frontend/src/routes/+page.svelte
frontend/src/routes/create/+page.svelte
frontend/src/routes/join/+page.svelte
frontend/src/routes/room/[code]/+page.svelte
frontend/src/lib/games/werewolf/WerewolfGame.svelte
frontend/src/lib/games/werewolf/DayPhase.svelte
frontend/src/lib/games/werewolf/NightPhase.svelte
frontend/src/lib/games/werewolf/RoleReveal.svelte
frontend/src/lib/games/werewolf/CenterCardSelect.svelte
frontend/src/lib/games/werewolf/PlayerCardSelect.svelte
frontend/src/lib/games/werewolf/Results.svelte
frontend/src/lib/components/ui/button.svelte
frontend/src/lib/components/ui/card.svelte
frontend/src/lib/components/ui/badge.svelte
backend/go.mod
backend/cmd/server/main.go
backend/internal/**/*.go (all Go files with import statements)
```

**Why important:**
- Svelte 4 syntax is deprecated and will be removed in future versions
- Svelte 5 provides better performance and type safety
- The new Runes system is more predictable and easier to reason about

**Implementation approach used:**
- ✅ Incremental migration following Phases A→B→C→D
- ✅ Tested imports and basic functionality after each phase
- ✅ Verified store reactivity with $effect and $derived
- ✅ Updated all components systematically
- ✅ Fixed Go module paths first to ensure backend stability

### 2. Go Backend Optimizations

**Status**: Reverted, needs reimplementation following Go best practices

#### 2a. Priority Queue for Phase Timeouts

**What was done:**
- Replaced O(N) polling with priority queue implementation
- Added heap-based data structure for efficient timeout management
- Implemented in `backend/internal/store/memory.go`
- Added new methods to store interface: `GetNextPhaseCheck()`, `UpdatePhaseCheckTime()`

**Why important:**
- Current implementation polls all rooms every second (inefficient at scale)
- Priority queue reduces complexity from O(N) to O(log N)
- Better CPU usage and scalability

**Reimplement with:**
- Use standard library `container/heap` package
- Add comprehensive unit tests for edge cases
- Verify no goroutine leaks
- Load test with many concurrent rooms

#### 2b. Production-Grade CORS with rs/cors

**What was done:**
- Replaced custom CORS middleware with `rs/cors` library
- Added proper environment variable configuration
- Implemented spec-compliant CORS behavior
- Updated `backend/cmd/server/main.go`
- Added dependency in `backend/go.mod`: `github.com/rs/cors v1.10.1`

**Why important:**
- Custom CORS implementation had bugs (empty string issues)
- `rs/cors` is battle-tested and maintained
- Follows CORS specification correctly
- Better security with proper origin validation

**Reimplement with:**
- Careful testing of CORS headers in development and production
- Verify preflight requests work correctly
- Test with actual Railway deployment URLs
- Document environment variable configuration

### 3. UX Improvements

**Status**: Reverted, consider reimplementing

**What was done:**
- Added visual reconnection status indicators in UI
- Updated `frontend/src/lib/stores/websocket.ts` with connection state tracking
- Modified `frontend/src/routes/+layout.svelte` to display connection status banner
- Only shows status when in room (not on home page)

**Why important:**
- Users need feedback when connection is lost
- Prevents confusion during network issues
- Better user experience during reconnections

**Reimplement with:**
- Test reconnection scenarios thoroughly
- Ensure banner doesn't appear on pages without WebSocket connections
- Consider toast notifications instead of persistent banner
- Add retry logic with exponential backoff

### 4. Railway Deployment Improvements

**Status**: Reverted, may reimplement

**What was done:**
- Enhanced API client configuration logging
- Added startup logs showing VITE_API_URL and API_BASE
- Railway-specific error messages for missing configuration
- Step-by-step instructions in console for environment variable setup
- Updated `frontend/src/lib/api/client.ts`

**Why important:**
- Helps diagnose configuration issues quickly
- Reduces debugging time for deployment problems
- Provides clear instructions for fixing common issues

**Reimplement with:**
- Consider making this development-only (strip in production builds)
- Add similar logging for WebSocket connection URL
- Document Railway deployment process in RAILWAY.md

## Implementation Priorities

When reimplementing these improvements:

1. **High Priority**:
   - Svelte 5 migration (framework future-proofing)
   - Priority queue optimization (scalability)
   - rs/cors library (security)

2. **Medium Priority**:
   - Reconnection status UX (user experience)
   - Railway deployment logging (developer experience)

3. **Process Changes**:
   - Test all changes locally before deployment
   - Deploy to staging environment first (or use feature flags)
   - Monitor production logs after deployment
   - Have rollback plan ready (git revert)
   - Document breaking changes and configuration requirements
   - Commit after each phase completion

## Detailed Svelte 5 Migration Guide

### Prerequisites
- Current: Svelte 4 syntax with `$:`, `export let`, writable stores
- Target: Svelte 5 with Runes system (`$state`, `$derived`, `$props()`)
- Package: `svelte: ^5.0.0` (already in package.json)

### Phase A: Store Migration (game.svelte.ts, session.svelte.ts, websocket.svelte.ts)

#### 1. game.ts → game.svelte.ts

**Before** (Svelte 4 writable store):
```typescript
import { writable } from 'svelte/store';

export const gameStore = writable({
  room: null,
  events: [],
  playerState: null,
  publicState: null
});
```

**After** (Svelte 5 runes):
```typescript
class GameStore {
  room = $state(null);
  events = $state([]);
  playerState = $state(null);
  publicState = $state(null);

  setRoomState(room) {
    this.room = room;
  }

  appendEvent(event) {
    this.events = [...this.events, event];
  }

  reset() {
    this.room = null;
    this.events = [];
    this.playerState = null;
    this.publicState = null;
  }
}

export const gameStore = new GameStore();
```

**Testing**: Verify room state updates, event appending, and reset functionality

#### 2. session.ts → session.svelte.ts

**Before**:
```typescript
import { writable } from 'svelte/store';

export const session = writable(null);
```

**After**:
```typescript
class SessionStore {
  value = $state(null);

  set(session) {
    this.value = session;
    if (browser) {
      localStorage.setItem('session', JSON.stringify(session));
    }
  }

  clear() {
    this.value = null;
    if (browser) {
      localStorage.removeItem('session');
    }
  }
}

export const session = new SessionStore();
```

**Testing**: Verify localStorage sync, session persistence across page refreshes

#### 3. websocket.ts → websocket.svelte.ts

**Before** (factory function returning store):
```typescript
function createWebSocketStore(roomCode, sessionToken) {
  const { subscribe, set, update } = writable({
    status: 'disconnected',
    messages: [],
    error: null
  });
  // ... connection logic
  return { subscribe, send, sendAction, disconnect };
}
```

**After** (class with runes):
```typescript
class WebSocketStore {
  status = $state('disconnected');
  messages = $state([]);
  error = $state(null);
  
  #ws = null;
  #roomCode;
  #sessionToken;

  constructor(roomCode, sessionToken) {
    this.#roomCode = roomCode;
    this.#sessionToken = sessionToken;
    this.connect();
  }

  connect() {
    this.status = 'connecting';
    // ... WebSocket setup
  }

  send(message) {
    if (this.#ws?.readyState === WebSocket.OPEN) {
      this.#ws.send(JSON.stringify(message));
    }
  }

  disconnect() {
    this.#ws?.close();
    this.status = 'disconnected';
  }
}

export function createWebSocket(roomCode, sessionToken) {
  return new WebSocketStore(roomCode, sessionToken);
}
```

**Testing**: Verify connection, message sending/receiving, reconnection logic

### Phase B: Layout & Core Pages

#### 4. +layout.svelte

**Before**:
```svelte
<script>
  import { session } from '$lib/stores/session';
  $: isLoggedIn = $session !== null;
</script>
```

**After**:
```svelte
<script>
  import { session } from '$lib/stores/session.svelte';
  let isLoggedIn = $derived(session.value !== null);
</script>
```

**Testing**: Verify reactive updates when session changes

#### 5. room/[code]/+page.svelte

**Before**:
```svelte
<script>
  import { gameStore } from '$lib/stores/game';
  $: roomState = $gameStore.room;
  $: isHost = $session?.playerId === roomState?.hostId;
</script>
```

**After**:
```svelte
<script>
  import { gameStore } from '$lib/stores/game.svelte';
  import { session } from '$lib/stores/session.svelte';
  
  let roomState = $derived(gameStore.room);
  let isHost = $derived(session.value?.playerId === roomState?.hostId);
</script>
```

**Testing**: Verify room state reactivity, host detection, WebSocket connection

### Phase C: Werewolf Game Components

#### 6. WerewolfGame.svelte

**Before**:
```svelte
<script lang="ts">
  export let roomCode: string;
  export let roomState: any;
  export let wsStore: any;

  $: myRole = /* derived from gameStore */;
  $: currentPhase = /* derived from gameStore */;
</script>
```

**After**:
```svelte
<script lang="ts">
  let { roomCode, roomState, wsStore } = $props();

  let myRole = $derived(/* derived from gameStore */);
  let currentPhase = $derived(/* derived from gameStore */);
</script>
```

**Testing**: Verify all props are received, reactive computations work

#### 7. DayPhase.svelte

**Before**:
```svelte
<script lang="ts">
  export let roomState: any;
  export let wsStore: any;
  export let timerActive: boolean = false;
  export let phaseEndsAt: Date | null = null;

  $: isHost = $session?.playerId === roomState?.hostId;
</script>
```

**After**:
```svelte
<script lang="ts">
  let { roomState, wsStore, timerActive = false, phaseEndsAt = null } = $props();

  let isHost = $derived($session?.playerId === roomState?.hostId);
</script>
```

**Testing**: Verify timer reactivity, host controls, action sending

### Phase D: UI Components

#### 8. button.svelte

**Before**:
```svelte
<script lang="ts">
  export let variant: 'default' | 'outline' = 'default';
  export let disabled: boolean = false;
  export let className: string = '';
</script>
```

**After**:
```svelte
<script lang="ts">
  let { 
    variant = 'default',
    disabled = false,
    class: className = '',
    children
  } = $props();
</script>

<button {disabled} class={cn(buttonVariants({ variant }), className)}>
  {@render children?.()}
</button>
```

**Testing**: Verify all variants, disabled state, className merging

### Testing Checklist for Each Phase

- [ ] Local dev server runs without errors
- [ ] All pages load correctly
- [ ] WebSocket connections establish
- [ ] Room creation works
- [ ] Room joining works
- [ ] Real-time updates (player list) work
- [ ] Game start flow works
- [ ] Role reveal works
- [ ] Night phase displays correctly
- [ ] Day phase timer works
- [ ] No console errors or warnings
- [ ] Test on Chrome, Firefox, Safari
- [ ] Test on mobile device

### Rollback Procedure

If issues arise:
1. Note the specific issue (networking, reactivity, etc.)
2. `git log` to find the migration commit
3. `git revert [commit-hash]`
4. Document the issue in this file
5. Deploy the revert
6. Analyze the issue before reattempting

### Known Gotchas
- Store subscriptions: `$store` auto-subscription doesn't work with runes - use `$derived(store.value)` instead
- Props destructuring: Must use `$props()`, not `export let`
- Derived state: `$derived` re-runs when dependencies change - be careful with expensive computations
- Children: Use `{@render children?.()}` for slot content
- Event handlers: Still use `on:click`, no changes there

## Notes

- The security hardening from commit c1fadaa (structured logging, concurrency fixes, crypto/rand, request limits) was kept and is working correctly
- Focus on incremental improvements with proper testing
- Consider feature flags for gradual rollout
- Keep documentation updated as features are reimplemented
- Test WebSocket connections thoroughly after each migration phase

