# 🛠️ Technical Task List

This document tracks specific technical debt and optimization tasks.

## 🟢 Backend (Go)

### Optimization
- [x] **Phase Check Routine**: Implemented priority queue optimization and fixed race conditions.
  - Implemented `PriorityQueue` (min-heap) for rooms based on `NextPhaseTime`.
  - Changes complexity from O(N) polling to O(1) for checking top item.
  - Fixed race conditions by adding `CheckAndAdvancePhase()` method to `Room` that safely acquires locks.
  - **Ref**: `backend/cmd/server/main.go`, `backend/internal/core/room.go`

### Code Quality
- [ ] **Error Handling**: Replace `slog.Error` with proper error return propagation where applicable in `ProcessAction`.
- [ ] **Testing**: Add unit tests for `werewolf/game.go` logic (specifically role swaps).

### Security
- [x] **CORS**: Replace custom middleware in `main.go` with `rs/cors` library for production-grade handling.
- [ ] **Session Tokens**: Ensure tokens are high-entropy and stored securely (consider hashing if persistent).

---

## 🔵 Frontend (SvelteKit)

### Migration to Svelte 5 Runes
The project currently uses Svelte 4 syntax (legacy mode in Svelte 5). We should migrate components gradually.

- [x] **Stores**: Convert `lib/stores/game.ts` to use `$state` instead of `writable`.
- [x] **Components**: Update `WerewolfGame.svelte` to use `$props()` and `$state()` instead of `export let` and top-level variables.
  - *Example Pattern*:
    ```typescript
    // Before
    export let roomCode;
    let count = 0;
    $: double = count * 2;

    // After
    let { roomCode} = $props();
    let count = $state(0);
    let double = $derived(count * 2);
    ```
- [x] **Complete Runes Migration**: Fixed all Svelte 5 deprecation warnings and event processing bugs
  - Replaced `<slot />` with `{@render children()}` in `+layout.svelte`
  - Replaced `<svelte:component>` with direct component usage in `WerewolfGame.svelte`
  - Replaced all `on:click` event handlers with `onclick` attributes across all components
  - Replaced `onMount`/`onDestroy` with `$effect` in `DayPhase.svelte` and `+page.svelte`
  - Converted `export let` to `$props()` in `CenterCardSelect.svelte` and `PlayerCardSelect.svelte`
  - Replaced `$:` reactive statements with `$derived` in `+page.svelte` and `DayPhase.svelte`
  - **Fixed critical event reprocessing bugs**: Added `lastProcessedEventIndex` tracking to prevent `$effect` from processing all events repeatedly in `WerewolfGame.svelte`, `NightPhase.svelte`, `Results.svelte`, and WebSocket messages in `+page.svelte`
  - **Fixed reactivity bugs**: Changed `$derived` to use `$session` instead of `get(session)` for proper reactive dependencies in `Results.svelte` and `NightPhase.svelte`

### UX/UI
- [ ] **Touch Targets**: Audit all buttons in `app.css` to ensure `min-height: 48px` (partially done).
- [x] **Reconnection**: Add visual indicator for "Reconnecting..." in `+layout.svelte` using the `websocket.ts` store status.

---

## 🏗️ Infrastructure

- [ ] **Redis**: Replace `store/memory.go` with `store/redis.go` implementation.
  - Interface `Store` is already defined, just need implementation.
- [ ] **Cleanup**: Verify `cleanupRoutine` handles edge cases (e.g., game in progress but no players connected for > 1 hour).
