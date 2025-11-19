# Future Improvements to Reimplement

This document records infrastructure improvements that were reverted on 2025-11-19 due to deployment issues. These should be reimplemented later with more oversight.

## Reverted Commits

The following commits were rolled back from commit c1fadaa1fc1786c168184e8abaf285cd2a99a771:

- **18dc401** - Documentation cleanup
- **cbd636d** - go.mod update
- **49e6781** - Svelte 5 migration + backend optimizations
- **e86e4ce** - Svelte 5 connection fixes
- **265a08a** - Railway deployment improvements
- **da1d62d** - Merge commit

## Key Improvements to Reimplement

### 1. Svelte 5 Migration

**Status**: Reverted, needs reimplementation with testing

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

**Files affected:**
```
frontend/src/lib/stores/game.svelte.ts (new)
frontend/src/lib/stores/game.ts (deleted)
frontend/src/lib/stores/websocket.ts
frontend/src/routes/+layout.svelte
frontend/src/routes/room/[code]/+page.svelte
frontend/src/lib/games/werewolf/CenterCardSelect.svelte
frontend/src/lib/games/werewolf/DayPhase.svelte
frontend/src/lib/games/werewolf/NightPhase.svelte
frontend/src/lib/games/werewolf/PlayerCardSelect.svelte
frontend/src/lib/games/werewolf/Results.svelte
frontend/src/lib/games/werewolf/WerewolfGame.svelte
```

**Why important:**
- Svelte 4 syntax is deprecated and will be removed in future versions
- Svelte 5 provides better performance and type safety
- The new Runes system is more predictable and easier to reason about

**Reimplement with:**
- Comprehensive testing on all routes before deployment
- Gradual rollout, one component at a time
- Verify WebSocket connection states work correctly on all pages
- Test room creation and join flows thoroughly

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
   - Deploy to staging environment first
   - Monitor production logs after deployment
   - Have rollback plan ready
   - Document breaking changes and configuration requirements

## Notes

- The security hardening from commit c1fadaa (structured logging, concurrency fixes, crypto/rand, request limits) was kept and is working correctly
- Focus on incremental improvements with proper testing
- Consider feature flags for gradual rollout
- Keep documentation updated as features are reimplemented

