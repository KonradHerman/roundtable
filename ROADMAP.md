# Cardless Development Roadmap

This is the focused development roadmap for Cardless, a party game platform that replaces physical cards with phones for in-person play.

For detailed technical tasks, see [TODO.md](TODO.md).
For future game concepts, see [GAMES.md](GAMES.md).

## Design Philosophy

Cardless is a **card replacement**, not a full digital game:
- ✅ Handles role assignment and private information
- ✅ Provides digital night actions (prevents cheating, tracks swaps accurately)
- ✅ Timer for discussion phases
- ❌ **NOT** digital voting - voting remains physical (everyone points simultaneously)
- ❌ **NOT** automatic winner calculation - players determine this together

This preserves the joy and social nature of in-person party games while solving the problems of physical cards (shuffling, lost pieces, accidental peeking).

---

## Current Status

### ✅ What's Working (Phase 1 Complete)
- Full lobby system (create room, join with code, real-time player list)
- WebSocket real-time communication with event sourcing
- Role assignment following One Night Werewolf rules (players + 3 center cards)
- Role reveal phase with player acknowledgements
- Night phase with host narration script (tabletop-style)
- Day phase with timer controls (start/pause/extend)
- Game abstraction layer ready for multiple games
- Results calculation logic (village vs werewolves, tanner edge case)
- **Night Actions**: Werewolf, Seer, Robber, Troublemaker, Drunk, Insomniac
- **Physical Cues**: Minion and Mason actions handled via UI instructions
- **Play Again**: Host can reset game keeping players in room

### 🔨 What's In Progress (Phase 2)
- Removing phone voting (replacing with physical voting)
- Role reveal screen (show final roles after discussion)
- UI Polish (mobile touch targets)

---

## Phase 2: Werewolf MVP - Digital Night Actions

**Goal**: Complete, playable One Night Werewolf with digital night actions

**Priority: HIGH - Current Sprint**

### Backend Tasks
- [x] **Implement Werewolf night action** (see other werewolves)
- [x] **Implement Seer action** (view player or 2 center cards)
- [x] **Implement Robber action** (swap with player, view new role)
- [x] **Implement Troublemaker action** (swap two other players)
- [x] **Implement Drunk action** (swap with center card)
- [x] **Implement Insomniac action** (view final role)
- [x] **Minion & Mason actions** (Physical cues implemented)
- [x] **Host tracking** (Verified)
- [x] **Implement play again endpoint**

### Frontend Tasks
- [x] **Create role-specific night phase UIs**
- [ ] **Remove voting UI from day phase**
  - Remove vote selection interface
  - Remove vote submission button
  - Keep timer and discussion prompt
- [ ] **Create role reveal screen**
  - After day phase ends (timer or manual advance)
  - Show each player their FINAL role
  - Clear display of any swaps that occurred
  - Option to reveal to table for discussion
- [x] **Implement play again UI**

### Testing & Validation
- [ ] Test with 3 players (minimum viable game)
- [ ] Test with 6-8 players (typical game night)
- [ ] Verify all night actions work correctly
- [ ] Test physical voting flow feels natural
- [ ] Verify role swaps tracked accurately
- [ ] Test play again feature works smoothly

**Success Criteria**: Can play complete games with digital night actions, physical voting flows naturally, play again works

**Estimated Time**: 2-3 days of focused development

---

## Phase 3: Polish & Stability

**Goal**: Stable, production-ready Werewolf

**Priority: MEDIUM - After Phase 2 Complete**

### Backend
- [ ] Reconnection handling (Event replay)
- [ ] Room cleanup (Auto-expire inactive rooms)
- [ ] WebSocket improvements (Error handling)

### Frontend
- [ ] QR code generation for room sharing
- [ ] Reconnection UI (Banner, auto-retry)
- [ ] Error boundaries and loading states
- [ ] Mobile polish (Touch targets, screen sizes)

---

## Phase 4: Second Game - Avalon

**Goal**: Validate architecture, prove it's multi-game

**Priority: MEDIUM - After Werewolf is stable**

### Implementation
- [ ] Avalon game logic (Quests, Voting, Teams)
- [ ] Avalon UI components
- [ ] Test and validate

---

## Phase 5: Production Ready

**Goal**: Deploy to production, stable for regular use

**Priority: LOW - After core games are solid**

- [ ] Redis store implementation
- [ ] Docker optimization
- [ ] PWA manifest
- [ ] SEO & Analytics

---

## Phase 6: Planned Games

We are prioritizing "no box required" games that fit the "phones as hands" model.

### 1. Spyfall (The Spy)
- **Concept**: Location-based social deduction.
- **Technical Challenge**: Non-card based state; requires location database.

### 2. Skull (Roses & Skulls)
- **Concept**: Bluffing and bidding.
- **Technical Challenge**: Bidding mechanic and state tracking for "flipped" discs.

### 3. Love Letter (Courtship)
- **Concept**: Micro-deck deduction game.
- **Technical Challenge**: Deck management and discard pile tracking.

---

## Technical Debt & Optimization

See [TODO.md](TODO.md) for the full list.

### Backend
- [ ] **Optimization**: Implement priority queue (min-heap) for phase timeouts to replace O(N) polling.
- [ ] **Security**: Add `rs/cors` for robust CORS handling.

### Frontend
- [ ] **Migration**: Pilot Svelte 5 Runes (`$state`, `$derived`) in a single component.
- [ ] **Type Safety**: Strengthen TypeScript definitions for WebSocket payloads.

---

Last updated: November 2025
