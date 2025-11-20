# Game Prototypes - Summary & Validation

> **Created**: 2025-11-19
> **Status**: Ready for Implementation
> **Games**: Avalon (The Resistance), Spyfall

---

## Overview

This document summarizes the two Phase 2 game prototypes (**Avalon** and **Spyfall**) and validates them against Roundtable's core design principles. Both games have been fully planned with detailed rules, event flows, state management, and implementation checklists.

---

## Design Principles Validation

### 1. ✅ Enhancing In-Person Interactions

Both prototypes strictly adhere to the principle of **enhancing physical play**, not replacing it.

#### Avalon
| Digital (App Handles) | Physical (Players Handle) |
|----------------------|---------------------------|
| Role assignment with hidden knowledge | Team discussions and deliberations |
| Quest card submissions (prevents tells) | Leader proposals (verbal) |
| Vote counting and revelation | Social reads and body language |
| Quest history tracking | Accusations and defenses |
| Assassination selection | Victory celebration |

#### Spyfall
| Digital (App Handles) | Physical (Players Handle) |
|----------------------|---------------------------|
| Location/role assignment | Question asking and answering |
| Timer management | Reading body language |
| Location database | Social deduction through conversation |
| Voting collection | Bluffing and misdirection |
| Spy location guess | Dramatic accusations |

**Verdict**: ✅ **Both games preserve the social core while automating mechanics.**

---

### 2. ✅ Tracking and Automating Card States

Both prototypes eliminate the need for physical components while preventing cheating and information leaks.

#### Avalon - Card Automation
- **Role cards**: Digital assignment with role-specific knowledge filtering
  - Merlin sees Evil (except Mordred)
  - Percival sees Merlin + Morgana (can't distinguish)
  - Evil sees each other (except Oberon)
- **Quest cards**: Digital submission prevents:
  - Card orientation tells
  - Timing tells (submit order)
  - Physical stacking patterns
- **Vote tokens**: Simultaneous digital voting
- **Quest track**: Automated 5-quest tracking with fail requirements

#### Spyfall - Card Automation
- **Location cards**: Digital assignment prevents peeking
- **Role cards**: Unique roles per location with no physical tells
- **Spy card**: Secret assignment with no distribution pattern
- **Timer**: Automated countdown (no phone timer needed)
- **Reference sheet**: Digital location list (no printed cards)

**Verdict**: ✅ **Both games fully automate card mechanics while maintaining game integrity.**

---

### 3. ✅ Event Sourcing Architecture

Both prototypes leverage the existing event sourcing foundation.

#### Event Sourcing Benefits Applied
1. **Seamless Reconnection**: Replay events to restore state
2. **Player-Specific Views**: Private events for hidden information
3. **Public State Sync**: Shared game state (quest results, timer, etc.)
4. **Audit Trail**: Full game history for debugging/replays

#### Avalon Events
- **Private**: `role_assigned`, `role_knowledge`
- **Public**: `team_proposed`, `team_approved`, `quest_completed`, `assassination_result`

#### Spyfall Events
- **Private**: `spy_assigned`, `location_assigned`, `location_list`
- **Public**: `timer_started`, `question_asked`, `accusation_resolved`, `spy_guesses_location`

**Verdict**: ✅ **Both games integrate seamlessly with event sourcing platform.**

---

### 4. ✅ Game Interface Compliance

Both prototypes implement the standard `Game` interface without platform changes.

```go
type Game interface {
    Initialize(config GameConfig, players []*Player) ([]GameEvent, error)
    ValidateAction(playerID string, action Action) error
    ProcessAction(playerID string, action Action) ([]GameEvent, error)
    GetPlayerState(playerID string) PlayerState
    GetPublicState() PublicState
    GetPhase() GamePhase
    IsFinished() bool
    GetResults() GameResults
    CheckPhaseTimeout() ([]GameEvent, error)
}
```

**Avalon Implementation**:
- ✅ State management via `AvalonState`
- ✅ Phase transitions: Setup → RoleReveal → TeamBuilding → TeamVoting → QuestExecution → QuestResults → Assassination → Finished
- ✅ Action processing: `propose_team`, `vote_team`, `play_quest_card`, `assassinate`
- ✅ Timeout handling: None required (manual phase advances)

**Spyfall Implementation**:
- ✅ State management via `SpyfallState`
- ✅ Phase transitions: Setup → RoleReveal → QuestionRound → Accusation → LocationGuess → Finished
- ✅ Action processing: `ask_question`, `answer_question`, `call_accusation`, `vote_accusation`, `guess_location`
- ✅ Timeout handling: Timer expiration triggers `LocationGuess` phase

**Verdict**: ✅ **Both games use standard interface with zero platform modifications.**

---

### 5. ✅ Mobile-First UX

Both prototypes prioritize mobile usability with Gruvbox Dark theme.

#### Avalon Mobile Considerations
- **Touch Targets**: 48x48px minimum for player selection
- **Team Selection**: Tap-to-toggle player grid
- **Quest Board**: Horizontal scrollable 5-quest track
- **Vote Buttons**: Large APPROVE/REJECT buttons
- **Assassination UI**: Full-screen player grid for selection
- **One-Handed**: All primary actions accessible with thumb

#### Spyfall Mobile Considerations
- **Always-Visible Location**: Sticky header for non-spies
- **Timer Display**: Large, prominent countdown
- **Question Input**: Full-width text areas with soft keyboard support
- **Player Selection**: Touch-friendly grid layout
- **Location List**: Scrollable modal for spy reference
- **One-Handed**: Timer controls and accusation button thumb-reachable

**Verdict**: ✅ **Both games designed for mobile-first, one-handed operation.**

---

## Architecture Validation

### Backend Structure Compliance

Both games follow the established pattern:

```
backend/internal/games/
├── werewolf/        # ✅ Existing
│   ├── game.go
│   ├── state.go
│   ├── phases.go
│   └── ...
├── avalon/          # 🆕 New
│   ├── game.go
│   ├── state.go
│   ├── phases.go
│   ├── roles.go
│   ├── quests.go
│   └── voting.go
└── spyfall/         # 🆕 New
    ├── game.go
    ├── state.go
    ├── phases.go
    ├── locations.go
    ├── timer.go
    └── voting.go
```

**Key Patterns Used**:
1. ✅ Secure role shuffling with `crypto/rand`
2. ✅ Mutex-protected state access
3. ✅ Event creation with visibility control
4. ✅ Explicit phase state machines
5. ✅ Configuration validation

---

### Frontend Structure Compliance

Both games follow the component pattern:

```
frontend/src/lib/games/
├── werewolf/          # ✅ Existing
│   ├── WerewolfGame.svelte
│   ├── RoleReveal.svelte
│   ├── NightPhase.svelte
│   └── roleConfig.ts
├── avalon/            # 🆕 New
│   ├── AvalonGame.svelte
│   ├── RoleReveal.svelte
│   ├── TeamBuilding.svelte
│   ├── TeamVoting.svelte
│   ├── QuestExecution.svelte
│   ├── QuestBoard.svelte
│   └── roleConfig.ts
└── spyfall/           # 🆕 New
    ├── SpyfallGame.svelte
    ├── RoleReveal.svelte
    ├── QuestionRound.svelte
    ├── TimerDisplay.svelte
    ├── LocationList.svelte
    └── locationConfig.ts
```

**Key Patterns Used**:
1. ✅ Svelte 5 runes (`$state`, `$derived`, `$props`)
2. ✅ Phase-specific components
3. ✅ Reusable UI elements (player grids, vote buttons, timers)
4. ✅ Role configuration files
5. ✅ Gruvbox Dark theming

---

## Shared Component Opportunities

### Reusable Components

Both prototypes identify components that can be extracted for reuse:

#### Timer Component (Used in Spyfall, Wavelength, future games)
```svelte
TimerDisplay.svelte
- Countdown display
- Pause/resume controls
- Extend time options
- Low-time warnings
- Critical time animations
```

#### Voting System (Used in Werewolf, Avalon, Spyfall)
```svelte
VotingUI.svelte
- Player selection grid
- Vote submission
- Vote count display
- Result reveal animation
```

#### Player Selector (Used in Avalon, Spyfall, Wavelength)
```svelte
PlayerSelector.svelte
- Touch-friendly player grid
- Single or multi-select modes
- Selected state visualization
- Disabled state support
```

#### Role Reveal (Used in all hidden-role games)
```svelte
RoleReveal.svelte
- Role card display
- Team banner
- Knowledge section (who you see)
- Acknowledgment button
```

**Recommendation**: Extract these during implementation to reduce code duplication.

---

## Implementation Complexity Analysis

### Avalon Complexity: **MEDIUM**

| Component | Complexity | Reason |
|-----------|-----------|--------|
| Role Knowledge System | Medium | Complex visibility rules (Merlin, Percival, Mordred, Oberon) |
| Quest Mechanics | Medium | Variable team sizes, 2-fail requirement for quest 4 |
| Voting System | Low | Simple approve/reject |
| Assassination Phase | Low | Simple player selection |
| UI Components | Medium | Quest board visualization, team selection |

**Estimated Time**: 5-6 days

**Risks**:
- Role knowledge filtering (must prevent information leaks)
- 5 consecutive rejections edge case
- Assassination phase only triggers conditionally

---

### Spyfall Complexity: **MEDIUM-LOW**

| Component | Complexity | Reason |
|-----------|-----------|--------|
| Timer System | Medium | Real-time countdown with pause/resume/extend |
| Location Database | Low | Content creation (30+ locations) |
| Turn Management | Low | Simple rotation (answerer becomes asker) |
| Voting System | Low | Simple accusation vote |
| UI Components | Medium | Timer display, Q&A interface, location list |

**Estimated Time**: 4-5 days

**Risks**:
- Timer synchronization across clients
- Dynamic turn order tracking
- Content quality (locations must be diverse and balanced)

---

## Testing Strategy

### Unit Testing (Backend)

#### Avalon
- [ ] Role knowledge filtering (Merlin, Percival, Evil)
- [ ] Quest resolution (1 fail vs 2 fails for quest 4)
- [ ] Team vote majority calculation
- [ ] 5 consecutive rejections auto-loss
- [ ] Assassination win condition reversal
- [ ] Configuration validation (team sizes, role requirements)

#### Spyfall
- [ ] Timer calculations (start/pause/resume/extend)
- [ ] Spy selection randomization
- [ ] Role assignment (all non-spies get unique roles)
- [ ] Accusation vote majority
- [ ] Location guess validation
- [ ] Win condition logic

### Integration Testing

#### Avalon
- [ ] Full game flow (5 quests)
- [ ] Team approval/rejection cycles
- [ ] Quest results update quest board
- [ ] Assassination phase trigger
- [ ] Reconnection mid-game

#### Spyfall
- [ ] Timer expiration triggers location guess
- [ ] Question/answer turn rotation
- [ ] Accusation interrupts question round
- [ ] Spy guess after incorrect accusation
- [ ] Reconnection with timer sync

### Multiplayer Playtesting

#### Avalon
- [ ] 5 players (minimum)
- [ ] 7 players (2-fail quest 4 test)
- [ ] 10 players (maximum)
- [ ] All role combinations
- [ ] Edge cases (5 rejections, assassination)

#### Spyfall
- [ ] 3 players (minimum)
- [ ] 5 players (typical)
- [ ] 8 players (maximum)
- [ ] Multiple rounds (new location each time)
- [ ] Timer edge cases (pause/extend during low time)

---

## Implementation Roadmap

### Phase 2A: Avalon (Priority 1)

**Week 1**: Backend + Frontend Core
- Days 1-2: Backend (roles, quests, voting)
- Days 3-4: Frontend (components, quest board)
- Day 5: Integration and testing
- Day 6: Multiplayer playtesting and polish

**Deliverables**:
- ✅ Avalon fully playable
- ✅ All roles implemented
- ✅ Assassination phase working
- ✅ Mobile-optimized UI

---

### Phase 2B: Spyfall (Priority 2)

**Week 2**: Backend + Frontend + Content
- Days 1-2: Backend (timer, locations, voting)
- Days 2-3: Frontend (timer UI, Q&A interface)
- Day 3: Location database creation (30+ locations)
- Day 4: Integration and testing
- Day 5: Multiplayer playtesting and polish

**Deliverables**:
- ✅ Spyfall fully playable
- ✅ 30+ diverse locations
- ✅ Timer system robust
- ✅ Mobile-optimized UI

---

### Phase 2C: Shared Component Extraction (Optional)

**Days 1-2**: Refactor common patterns
- Extract `TimerDisplay` component
- Extract `VotingUI` component
- Extract `PlayerSelector` component
- Extract `RoleReveal` base component
- Update Werewolf/Avalon/Spyfall to use shared components

**Deliverables**:
- ✅ Reduced code duplication
- ✅ Consistent UX across games
- ✅ Easier maintenance

---

## Success Criteria

### Avalon Launch Criteria
- [ ] All 8 roles implemented and tested
- [ ] Quest board displays 5 quests correctly
- [ ] Team voting works (approve/reject)
- [ ] Quest card submissions work (Good can only play Success)
- [ ] Assassination phase triggers after Good wins 3 quests
- [ ] All win conditions work correctly
- [ ] Mobile UI is smooth and one-handed friendly
- [ ] No information leaks (role knowledge is secure)
- [ ] Playtested with 5, 7, and 10 players
- [ ] Documentation complete (How to Play guide)

### Spyfall Launch Criteria
- [ ] Timer system works (start/pause/resume/extend)
- [ ] 30+ locations with diverse themes
- [ ] Location and role assignment works
- [ ] Spy sees location list, non-spies see location
- [ ] Question/answer flow works
- [ ] Accusation voting works
- [ ] Spy location guess works
- [ ] All win conditions work correctly
- [ ] Timer syncs across clients (no drift)
- [ ] Mobile UI is smooth and one-handed friendly
- [ ] Playtested with 3, 5, and 8 players
- [ ] Documentation complete (How to Play guide)

---

## Risk Assessment

### Avalon Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Information leaks (role knowledge) | Medium | High | Extensive unit tests for visibility filtering |
| 5 rejection rule edge case | Low | Medium | Dedicated integration test |
| Assassination UX confusion | Medium | Low | Clear UI prompts and tutorial |
| Complex role interactions | Low | Medium | Thorough documentation and playtesting |

### Spyfall Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Timer drift across clients | Medium | Medium | Server-authoritative time with client sync |
| Insufficient location variety | Medium | Medium | 30+ locations minimum, diverse themes |
| Q&A interface clunky on mobile | Medium | High | Mobile-first design, large touch targets |
| Turn order confusion | Low | Low | Clear turn indicator UI |

---

## Documentation Deliverables

### For Each Game

1. **Prototype Plan** ✅ (Created)
   - Full rules reference
   - Event flow breakdown
   - Implementation checklist

2. **How to Play Guide** 📋 (To be created during implementation)
   - Beginner-friendly tutorial
   - Role explanations
   - Strategy tips

3. **Developer Documentation** 📋 (To be created during implementation)
   - API documentation
   - State management guide
   - Testing guide

4. **Component Documentation** 📋 (To be created during implementation)
   - Component props and usage
   - Styling guide
   - Integration examples

---

## Conclusion

Both **Avalon** and **Spyfall** prototypes are **ready for implementation**. They:

✅ Adhere to all design principles (in-person enhancement, card automation)
✅ Integrate with existing platform architecture (event sourcing, Game interface)
✅ Follow established patterns (Svelte 5, Gruvbox Dark, mobile-first)
✅ Have detailed implementation plans (backend, frontend, testing)
✅ Are estimated accurately (Avalon: 5-6 days, Spyfall: 4-5 days)

### Recommended Sequence

1. **Avalon First** - Validates multi-round, team-based mechanics
2. **Spyfall Second** - Validates real-time timer and dynamic turn systems
3. **Extract Shared Components** - Reduce duplication and prepare for Phase 3 games

### Total Timeline

- **Avalon**: 5-6 days
- **Spyfall**: 4-5 days
- **Component Extraction**: 2 days (optional)
- **Total**: 9-11 days for Phase 2 completion

---

**Status**: ✅ **READY TO IMPLEMENT**

**Next Steps**:
1. Review prototypes with team (if applicable)
2. Create implementation branch: `feature/avalon`
3. Begin Avalon backend development
4. Follow implementation checklist in `AVALON_PROTOTYPE.md`

---

**Document Version**: 1.0
**Last Updated**: 2025-11-19
**Author**: Claude (Roundtable AI Assistant)
