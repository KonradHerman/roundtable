# 🛠️ Technical Task List

This document tracks specific technical debt and optimization tasks.

## 🟢 Backend (Go)

### Optimization
- [ ] **Phase Check Routine**: Currently, `phaseCheckRoutine` in `main.go` iterates through ALL rooms every second (O(N)).
  - **Task**: Implement a `PriorityQueue` (min-heap) for rooms based on `NextPhaseTime`.
  - **Benefit**: Changes complexity to O(1) for checking top item.
  - **Ref**: `backend/cmd/server/main.go`

### Code Quality
- [ ] **Error Handling**: Replace `slog.Error` with proper error return propagation where applicable in `ProcessAction`.
- [ ] **Testing**: Add unit tests for `werewolf/game.go` logic (specifically role swaps).

### Security
- [ ] **CORS**: Replace custom middleware in `main.go` with `rs/cors` library for production-grade handling.
- [ ] **Session Tokens**: Ensure tokens are high-entropy and stored securely (consider hashing if persistent).

---

## 🔵 Frontend (SvelteKit)

### Migration to Svelte 5 Runes
The project currently uses Svelte 4 syntax (legacy mode in Svelte 5). We should migrate components gradually.

- [ ] **Stores**: Convert `lib/stores/game.ts` to use `$state` instead of `writable`.
- [ ] **Components**: Update `WerewolfGame.svelte` to use `$props()` and `$state()` instead of `export let` and top-level variables.
  - *Example Pattern*:
    ```typescript
    // Before
    export let roomCode;
    let count = 0;
    $: double = count * 2;

    // After
    let { roomCode } = $props();
    let count = $state(0);
    let double = $derived(count * 2);
    ```

### UX/UI
- [ ] **Touch Targets**: Audit all buttons in `app.css` to ensure `min-height: 48px` (partially done).
- [ ] **Reconnection**: Add visual indicator for "Reconnecting..." in `+layout.svelte` using the `websocket.ts` store status.

---

## 🏗️ Infrastructure

- [ ] **Redis**: Replace `store/memory.go` with `store/redis.go` implementation.
  - Interface `Store` is already defined, just need implementation.
- [ ] **Cleanup**: Verify `cleanupRoutine` handles edge cases (e.g., game in progress but no players connected for > 1 hour).

