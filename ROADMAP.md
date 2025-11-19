# Cardless Development Roadmap

> **Play at [cardless.games](https://cardless.games)**

This is the focused development roadmap for Cardless, a party game platform that replaces physical cards with phones for in-person play.

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

### 🔨 What's In Progress (Phase 2)
- Digital night actions for each role
- Removing phone voting (replacing with physical voting)
- Role reveal screen (show final roles after discussion)
- Play again feature
- Host tracking fix

---

## Phase 2: Werewolf MVP - Digital Night Actions

**Goal**: Complete, playable One Night Werewolf with digital night actions

**Priority: HIGH - Current Sprint**

### Backend Tasks
- [ ] **Implement Werewolf night action** (see other werewolves)
  - Private event to show other werewolf players
  - Handle solo werewolf viewing center card
- [ ] **Implement Seer action** (view player or 2 center cards)
  - Validate action (Seer only, once per game)
  - Process view request, send private result
- [ ] **Implement Robber action** (swap with player, view new role)
  - Validate action (Robber only, once per game)
  - Update roleAssignments
  - Send private result with new role
- [ ] **Implement Troublemaker action** (swap two other players)
  - Validate action (Troublemaker only, once per game)
  - Update roleAssignments for both players
  - No result sent (Troublemaker doesn't see what they swapped)
- [ ] **Implement Drunk action** (swap with center card)
  - Validate action (Drunk only, once per game)
  - Update roleAssignments
  - No result sent (Drunk doesn't see new role)
- [ ] **Implement Insomniac action** (view final role)
  - Automatically triggered at end of night
  - Send private result with final role
- [ ] **Fix host tracking** (use actual host from room creation, not first player)
- [ ] **Implement play again endpoint**
  - Keep players in room
  - Reset game state to lobby/configuration
  - Clear old events/state

### Frontend Tasks
- [ ] **Create role-specific night phase UIs**
  - Werewolf: Show other werewolves or center card option
  - Seer: Player selection OR center cards selection
  - Robber: Player selection, show result
  - Troublemaker: Select two other players
  - Drunk: Center card selection
  - Insomniac: Auto-display final role
  - Villager/Tanner: Simple "no action" message
- [ ] **Remove voting UI from day phase**
  - Remove vote selection interface
  - Remove vote submission button
  - Keep timer and discussion prompt
- [ ] **Create role reveal screen**
  - After day phase ends (timer or manual advance)
  - Show each player their FINAL role
  - Clear display of any swaps that occurred
  - Option to reveal to table for discussion
- [ ] **Implement play again UI**
  - Button for host at end of game
  - Return to lobby/configuration screen
  - Keep same players and room code
- [ ] **Update results component**
  - Remove automatic winner calculation
  - Show all final roles
  - Let players determine winner based on physical votes

### Testing & Validation
- [ ] Test with 3 players (minimum viable game)
- [ ] Test with 6-8 players (typical game night)
- [ ] Verify all night actions work correctly
- [ ] Test physical voting flow feels natural
- [ ] Verify role swaps tracked accurately
- [ ] Test play again feature works smoothly

**Success Criteria**: Can play complete games with digital night actions, physical voting flows naturally, play again works

**Estimated Time**: 3-4 days of focused development

---

## Phase 3: Polish & Stability

**Goal**: Stable, production-ready Werewolf

**Priority: MEDIUM - After Phase 2 Complete**

### Backend
- [ ] Reconnection handling
  - Event replay for reconnecting clients
  - Handle mid-game disconnects gracefully
- [ ] Room cleanup
  - Auto-expire inactive rooms
  - Clear old game data
- [ ] WebSocket improvements
  - Better error handling
  - Connection state management

### Frontend
- [ ] QR code generation for room sharing
  - Generate QR code for room URL
  - Display prominently in lobby
  - Easy for friends to scan and join
- [ ] Reconnection UI
  - Detect disconnects automatically
  - Show reconnecting banner
  - Auto-retry with backoff
- [ ] Error boundaries
  - Graceful error handling
  - User-friendly error messages
  - Fallback UI for crashed components
- [ ] Loading states
  - Better loading indicators
  - Skeleton screens
  - Smooth transitions
- [ ] Mobile polish
  - Ensure all touch targets are 48x48px minimum
  - Test on various screen sizes
  - Improve one-handed usability

### Testing
- [ ] Extended playtesting with multiple groups
- [ ] Test on various devices (iOS, Android, different screen sizes)
- [ ] Test edge cases (disconnects, late joins, etc.)
- [ ] Performance testing with 10 players

**Success Criteria**: Stable enough for regular game nights, handles edge cases gracefully

**Estimated Time**: 3-4 days

---

## Phase 4: Frontend Framework Upgrade

**Goal**: Migrate to Svelte 5 for improved performance and developer experience

**Priority: MEDIUM - After Phase 2 complete and tested**

### Background
The application currently uses Svelte 4 syntax. Svelte 5 introduces Runes, a more modern and predictable reactivity system. An earlier migration attempt was reverted due to networking issues (see `FUTURE_IMPROVEMENTS.md`).

### Migration Strategy
Incremental migration in phases:
1. **Stores** (game.ts → game.svelte.ts, session.ts, websocket.ts)
2. **Layouts and Core Pages** (+layout.svelte, landing, create, join, room)
3. **Werewolf Components** (WerewolfGame, NightPhase, DayPhase, etc.)
4. **UI Components** (button, card, badge, etc.)

### Key Changes
- `$:` reactive declarations → `$derived` rune
- `let` state → `$state` rune
- `export let` props → `$props()` rune
- Writable stores → `.svelte.ts` files with runes
- Store subscriptions → Direct reactivity

### Testing Requirements
- Test each phase independently
- Verify WebSocket connections on all pages
- Test full game flow after each phase
- Test on multiple devices/browsers
- Have rollback plan ready

**Estimated Time**: 4-5 days (with careful testing)

See detailed migration plan in `FUTURE_IMPROVEMENTS.md`

---

## Phase 5: Second Game - Avalon

**Goal**: Validate game abstraction, prove platform supports multiple games

**Priority: MEDIUM - After Werewolf is stable**

### Pre-Work
- [ ] Extract reusable patterns from Werewolf implementation
  - Lobby flow
  - Role assignment
  - Phase transitions
  - Event handling
- [ ] Document game implementation guide
  - Required interfaces to implement
  - Event sourcing patterns
  - Frontend component structure
- [ ] Identify platform vs game-specific code
  - What should move to core?
  - What stays game-specific?

### Implementation
- [ ] Avalon game logic
  - Quest voting mechanics
  - Team selection
  - Role reveals (Merlin sees evil, etc.)
  - Assassin endgame
- [ ] Avalon UI components
  - Quest tracking
  - Team selection interface
  - Mission voting
  - Role information screens
- [ ] Test and validate

**Success Criteria**: Avalon works with minimal changes to core platform, validates game abstraction is working

**Estimated Time**: 4-5 days

---

## Phase 6: Production Ready

**Goal**: Deploy to production, stable for regular use

**Priority: LOW - After core games are solid**

### Backend
- [ ] Redis store implementation (replace in-memory)
- [ ] Docker deployment optimization
- [ ] Rate limiting
- [ ] Health monitoring and alerts
- [ ] Backup and recovery

### Frontend
- [ ] PWA manifest (make installable)
- [ ] Performance optimization
- [ ] Analytics (optional, privacy-focused)
- [ ] SEO optimization

### Infrastructure
- [ ] CI/CD pipeline
- [ ] Monitoring and logging
- [ ] Deployment automation
- [ ] Documentation for self-hosting

**Success Criteria**: Deployed and stable for regular game nights

**Estimated Time**: 3-4 days

---

## Future Considerations

## Phase 7: Additional Games

**Goal**: Expand game library to prove platform versatility

**Priority: MEDIUM - After Avalon validates architecture**

### Planned Games

#### High Priority
- **Spyfall** - Location deduction with question/answer rounds
  - Unique spy/location mechanic
  - Timer-based rounds
  - Question tracking UI
  - Estimated time: 3-4 days

- **Skull** - Pure bluffing and bidding
  - Card stacking mechanics
  - Bidding system
  - Simple rules, deep strategy
  - Estimated time: 2-3 days

#### Medium Priority
- **Wavelength** - Team spectrum guessing
  - Dial/spectrum UI component
  - Team-based scoring
  - Clue-giving system
  - Estimated time: 3-4 days

- **Coup** - Bluffing and deduction
  - Character abilities
  - Challenge system
  - Coins and actions
  - Estimated time: 3-4 days

#### Future Consideration
- **Secret Hitler** - Policy and social deduction
- **Love Letter** - Card drafting
- **Cockroach Poker** - Pure bluffing
- **Bohnanza** - Trading mechanics

See `GAMES_ROADMAP.md` for detailed implementation plans.

---

## Phase 8: Advanced Features

### Features (Low Priority)
- Game history/statistics
- User accounts (optional, for stats only)
- Custom role configurations
- Sound effects and haptic feedback
- Multiple languages
- Spectator mode
- Game replays

---

## Questions & Decisions

### Resolved ✅
- ✅ **Night actions**: Digital (prevents cheating, tracks swaps)
- ✅ **Voting**: Physical (preserves in-person excitement)
- ✅ **Winner calculation**: Manual (players determine together)
- ✅ **Board/spectator view**: Not a priority for MVP
- ✅ **Play again**: Must implement (keep players in room)
- ✅ **Deployment**: Railway for production (cardless.games)
- ✅ **Next games**: Avalon, then Spyfall, Skull, Wavelength

### Open Questions ❓
- ❓ When to migrate to Svelte 5? (After Phase 2 complete)
- ❓ When to implement Redis? (When scaling beyond 100s of concurrent games)
- ❓ PWA features priority? (After Phase 3)

---

## How to Use This Roadmap

1. **Current focus**: Phase 2 (Digital Night Actions)
2. **Next**: Phase 3 (Polish & Stability)
3. **Then**: Phase 4 (Svelte 5 Migration)
4. **After that**: Phase 5 (Avalon) to validate multi-game architecture
5. **Long-term**: Phases 6-8 (Production ready, more games, advanced features)

This roadmap is a living document. Update as priorities shift and new information emerges.

**Last updated**: November 2024

