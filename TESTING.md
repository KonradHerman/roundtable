# Testing Guide for Cardless

This document outlines testing strategies for Cardless, from manual testing during development to potential automation approaches for the future.

## Testing Philosophy

Cardless is a **real-time multiplayer game platform** that requires:
- Multiple simultaneous connections (3-10 players)
- WebSocket communication
- Event synchronization across clients
- Mobile-first UX

This makes testing more complex than typical web apps. Our approach is:
1. **Start with manual testing** (sufficient for MVP)
2. **Add automation selectively** (when it speeds up development)
3. **Focus on critical paths** (lobby → game → results)

---

## Manual Testing

### Prerequisites

- Backend running on `localhost:8080`
- Frontend running on `localhost:5173`
- Multiple browser windows/devices for multi-player testing

### Test Setup Methods

#### Method 1: Incognito Windows (Fastest)

**Best for**: Quick iteration during development

1. Open app in normal browser window
2. Create a room (become host)
3. Open 2-7 incognito windows (`Ctrl+Shift+N` in Chrome)
4. Join room from each incognito window
5. Test game flow

**Pros:**
- Very fast to set up
- Can manage all windows on one screen
- Easy to see all player perspectives

**Cons:**
- Cramped on single monitor
- Not realistic mobile UX
- Can't test touch interactions

**Tips:**
- Use split screen or multiple monitors
- Label each window clearly in dev tools console
- Keep host window prominent for advancing phases

#### Method 2: Multiple Devices (Most Realistic)

**Best for**: Final validation before release

1. Find your local IP: `ipconfig` (Windows) or `ifconfig` (Mac/Linux)
2. Update frontend proxy in `vite.config.ts` to use your IP
3. Connect phones/tablets to same WiFi
4. Access `http://192.168.x.x:5173` from each device
5. Test game flow

**Pros:**
- Tests actual mobile UX
- Realistic touch interactions
- Can test with real friends
- Spot mobile-specific issues

**Cons:**
- Slower to set up
- Need multiple devices
- Harder to debug

**Tips:**
- Test on both iOS and Android
- Check different screen sizes
- Verify touch target sizes (min 48x48px)
- Test on slower devices

#### Method 3: Different Browsers (Alternative)

**Best for**: Quick cross-browser testing

Use different browsers as different players:
- Chrome (normal)
- Chrome (incognito)
- Firefox
- Edge
- Safari (Mac only)

**Pros:**
- Tests browser compatibility
- Don't need incognito
- Easy to debug

**Cons:**
- Limited to ~5 browsers
- Not realistic for mobile
- More resource intensive

---

## Test Scenarios

### Critical Path: Full Game Flow

**Frequency**: Test after every significant change

1. **Lobby Phase**
   - [ ] Host creates room, gets 6-character code
   - [ ] 2-3 players join with code
   - [ ] All players see each other in real-time
   - [ ] Player names display correctly
   - [ ] Host sees "Start Game" button

2. **Configuration Phase**
   - [ ] Host selects game (Werewolf)
   - [ ] Host configures roles (count = players + 3)
   - [ ] Invalid configs show error
   - [ ] "Start Game" button becomes enabled

3. **Role Reveal Phase**
   - [ ] All players receive their roles privately
   - [ ] Each role displays correctly
   - [ ] Players can acknowledge role
   - [ ] Acknowledgement count updates for all
   - [ ] Auto-advances to night when all acknowledged

4. **Night Phase**
   - [ ] Phase changes for all players
   - [ ] Host sees narration script
   - [ ] Each role sees appropriate UI:
     - Werewolves see other werewolves
     - Seer can view player or center cards
     - Robber can swap and see new role
     - Troublemaker can swap two players
     - Drunk must swap with center
     - Insomniac sees final role
     - Villager/Tanner see "no action" message
   - [ ] Actions are validated (can't act twice)
   - [ ] Private results sent only to acting player
   - [ ] Host can advance to day phase

5. **Day Phase**
   - [ ] Phase changes for all players
   - [ ] Timer controls visible to host
   - [ ] Timer can start/pause/extend
   - [ ] Timer syncs across all clients
   - [ ] Discussion prompt visible
   - [ ] No voting UI present

6. **Role Reveal**
   - [ ] After day phase, roles revealed
   - [ ] Each player sees their FINAL role
   - [ ] Role swaps are reflected correctly
   - [ ] Clear display of what happened

7. **Play Again**
   - [ ] Host sees "Play Again" button
   - [ ] Clicking returns to lobby
   - [ ] Players remain in room
   - [ ] Can reconfigure and start new game

### Edge Cases

**Connection Issues:**
- [ ] Player disconnects during lobby → can rejoin
- [ ] Player disconnects during game → can reconnect
- [ ] Host disconnects → game continues or handles gracefully
- [ ] All players disconnect → room cleans up

**Invalid Actions:**
- [ ] Can't join room that doesn't exist
- [ ] Can't join game after it started
- [ ] Can't start game with wrong role count
- [ ] Can't perform night action twice
- [ ] Can't perform night action as wrong role

**Concurrent Actions:**
- [ ] Multiple players acknowledging roles simultaneously
- [ ] Multiple players performing night actions at once
- [ ] Spam clicking buttons doesn't break state

**Browser Compatibility:**
- [ ] Works in Chrome (desktop & mobile)
- [ ] Works in Firefox
- [ ] Works in Safari (desktop & mobile)
- [ ] Works in Edge

**Mobile Specific:**
- [ ] Touch targets are easily tappable
- [ ] No horizontal scrolling
- [ ] Keyboard doesn't obscure UI
- [ ] Phone rotation works
- [ ] Works in both portrait and landscape

---

## Performance Testing

### Load Testing (Manual)

**Goal**: Verify app works with maximum players

1. Open 10 incognito windows
2. All join same room
3. Start game with 13 roles (10 players + 3 center)
4. Verify all phases work smoothly
5. Check for lag or synchronization issues

**Expected Results:**
- No noticeable lag
- Events arrive within 100ms
- UI remains responsive
- WebSocket connections stable

### Stress Testing

**When to do**: Before major releases

1. Create multiple rooms simultaneously
2. Run multiple games at once
3. Rapidly create and join rooms
4. Test server handles load

**Look for:**
- Memory leaks (check backend console)
- Connection drops
- Event processing delays
- Server errors

---

## Future: Automation Testing

### When to Implement

Consider automation testing when:
- ✅ Manual testing becomes too time-consuming
- ✅ Working on game #2 (need regression tests)
- ✅ Team grows beyond solo developer
- ✅ Frequent bugs introduced by changes

### Automation Approaches

#### 1. Unit Tests (Backend)

**Priority**: HIGH (easy wins)

Test game logic in isolation:

```go
// Example: Test Seer action
func TestSeerViewPlayer(t *testing.T) {
    game := werewolf.NewGame()
    // ... setup game state
    
    action := core.Action{
        Type: "seer_view",
        Payload: json.RawMessage(`{"targetId":"player-2"}`),
    }
    
    events, err := game.ProcessAction("player-1", action)
    
    assert.NoError(t, err)
    assert.Equal(t, 1, len(events))
    assert.Equal(t, "seer_viewed", events[0].Type)
}
```

**What to test:**
- Role assignment logic
- Night action processing
- Vote counting
- Results calculation
- Edge cases and error handling

**Tools:**
- Go's built-in `testing` package
- `testify` for assertions

**Estimated effort**: 2-3 days

#### 2. Integration Tests (Backend)

**Priority**: MEDIUM

Test API endpoints and WebSocket flows:

```go
func TestRoomCreationFlow(t *testing.T) {
    // Create room via HTTP
    room := createRoom(t, "werewolf", "Alice")
    
    // Join via HTTP
    player := joinRoom(t, room.Code, "Bob")
    
    // Connect via WebSocket
    ws := connectWebSocket(t, player.SessionToken)
    
    // Verify events received
    event := waitForEvent(ws, "player_joined")
    assert.Equal(t, "Bob", event.DisplayName)
}
```

**What to test:**
- REST API endpoints
- WebSocket authentication
- Event broadcasting
- Room cleanup

**Tools:**
- `httptest` for HTTP testing
- `gorilla/websocket` test clients

**Estimated effort**: 2-3 days

#### 3. E2E Tests (Full Stack)

**Priority**: MEDIUM (high value, but complex)

Simulate multiple real players:

```typescript
// Example: Playwright multi-session test
test('complete game flow', async () => {
    // Create 3 browser contexts (3 players)
    const host = await browser.newContext();
    const player2 = await browser.newContext();
    const player3 = await browser.newContext();
    
    // Host creates room
    const hostPage = await host.newPage();
    await hostPage.goto('http://localhost:5173');
    await hostPage.click('text=Host a Game');
    const roomCode = await hostPage.locator('.room-code').textContent();
    
    // Players join
    const p2Page = await player2.newPage();
    await p2Page.goto('http://localhost:5173');
    await p2Page.fill('input[placeholder="Enter room code"]', roomCode);
    await p2Page.click('text=Join Game');
    
    // ... continue test
});
```

**What to test:**
- Full lobby → game → results flow
- Multi-player synchronization
- Critical user journeys
- Regression tests for bugs

**Tools:**
- **Playwright** (recommended)
  - Supports multiple browser contexts
  - Great for real-time apps
  - Good documentation
- **Puppeteer** (alternative)
- **Cypress** (not ideal for multi-session)

**Estimated effort**: 3-4 days initial setup, then add tests incrementally

#### 4. WebSocket Load Testing

**Priority**: LOW (only if scaling)

Simulate many concurrent connections:

```python
# Example: Using websocket-client
import websocket
import threading

def create_client(room_code):
    ws = websocket.WebSocket()
    ws.connect("ws://localhost:8080/ws")
    ws.send(json.dumps({
        "type": "authenticate",
        "token": get_session_token()
    }))
    # ... interact with server

# Spawn 100 clients
threads = [threading.Thread(target=create_client, args=("ABC123",)) 
           for _ in range(100)]
for t in threads:
    t.start()
```

**What to test:**
- Server handles many connections
- Event broadcasting performance
- Memory usage under load
- Connection limits

**Tools:**
- Custom WebSocket clients (Go, Python)
- `k6` for load testing
- `artillery` for WebSocket testing

**Estimated effort**: 2-3 days

---

## Test Coverage Goals

### Phase 2 (Current): Manual Testing Only
- ✅ 100% manual coverage of critical paths
- ✅ Multi-player testing with 3-8 players
- ✅ Edge case testing (disconnects, errors)
- ❌ No automation yet

### Phase 3 (Polish): Add Backend Unit Tests
- ✅ 80%+ coverage of game logic
- ✅ Unit tests for all night actions
- ✅ Edge case tests for vote counting
- ✅ Manual E2E testing continues

### Phase 4 (Game #2): Add E2E Tests
- ✅ Critical path E2E tests (Werewolf & Avalon)
- ✅ Regression tests for known bugs
- ✅ Manual testing for new features
- ❌ Not aiming for 100% E2E coverage

### Phase 5 (Production): Full Coverage
- ✅ Backend unit tests (80%+)
- ✅ Backend integration tests (key flows)
- ✅ E2E tests (critical paths)
- ✅ Load testing before major releases

---

## CI/CD Integration

### When to add CI/CD

After implementing automation tests, integrate with CI:

```yaml
# Example: GitHub Actions workflow
name: Tests
on: [push, pull_request]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
      - run: cd backend && go test ./...
  
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-node@v2
      - run: npm install
      - run: npx playwright test
```

**Benefits:**
- Catch regressions before merge
- Faster feedback on PRs
- Confidence in deployments

---

## Testing Checklist

Before each release:

- [ ] Manual test full game flow with 6+ players
- [ ] Test on mobile devices (iOS + Android)
- [ ] Test reconnection (refresh mid-game)
- [ ] Test edge cases (see list above)
- [ ] Load test with max players (10)
- [ ] Run all automated tests (when available)
- [ ] Check for console errors
- [ ] Verify WebSocket connections stable

---

## Debugging Tips

### Backend
- Check logs in terminal where `go run` is running
- Add debug logging: `log.Printf("Debug: %+v", data)`
- Use `dlv` debugger for complex issues

### Frontend
- Open browser DevTools (F12)
- Check Console for errors
- Network tab → WS to see WebSocket messages
- Use Vue/React DevTools to inspect component state

### WebSocket Issues
- Use DevTools Network → WS filter
- Check connection state (Open/Closed)
- Verify authentication message sent
- Check event payloads

---

## Questions & Next Steps

### Current Approach (Phase 2)
- ✅ Manual testing sufficient for MVP
- ✅ Focus on real device testing
- ❌ No automation yet

### Future Considerations
- ❓ Add backend unit tests in Phase 3?
- ❓ Add Playwright E2E tests before Game #2?
- ❓ When does manual testing become too slow?

---

Last updated: November 2025

