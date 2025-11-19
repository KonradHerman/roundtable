# Backend Remediation Plan
## Go Best Practices & Security Hardening

> **Status**: Most critical items completed. See [IMPLEMENTATION_STATUS.md](IMPLEMENTATION_STATUS.md) for progress tracking.

This document provides detailed guidance for addressing security, concurrency, and architectural gaps identified in the Roundtable backend audit.

---

## 1. Authentication & Authorization

### Current Issues
- No authentication middleware on REST endpoints
- WebSocket actions don't verify host privileges
- Session tokens validated only at connection time
- TODO comments in handlers indicate missing auth logic

### Recommended Solution

**Middleware-Based Authentication:**
```go
// internal/server/auth.go
func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        sessionToken := r.Header.Get("X-Session-Token")
        if sessionToken == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Extract room code from path
        roomCode := r.PathValue("code")
        room, err := s.store.GetRoom(roomCode)
        if err != nil {
            http.Error(w, "Room not found", http.StatusNotFound)
            return
        }

        player, err := room.GetPlayerByToken(sessionToken)
        if err != nil {
            http.Error(w, "Invalid session", http.StatusUnauthorized)
            return
        }

        // Inject player ID into request context
        ctx := context.WithValue(r.Context(), "playerID", player.ID)
        ctx = context.WithValue(ctx, "roomCode", roomCode)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}

// Usage in main.go:
mux.HandleFunc("POST /api/rooms/{code}/start", s.authMiddleware(s.HandleStartGame))
mux.HandleFunc("POST /api/rooms/{code}/reset", s.authMiddleware(s.HandleResetGame))
```

**Host-Only Action Validation:**
```go
// internal/server/handlers.go - in HandleStartGame
func (s *Server) HandleStartGame(w http.ResponseWriter, r *http.Request) {
    playerID := r.Context().Value("playerID").(string)
    roomCode := r.Context().Value("roomCode").(string)
    
    room, _ := s.store.GetRoom(roomCode)
    if !room.IsHost(playerID) {
        http.Error(w, "Only host can start game", http.StatusForbidden)
        return
    }
    // ... rest of handler
}
```

**WebSocket Action Authorization:**
```go
// internal/server/websocket.go - in handleClientMessage
case ClientMsgAction:
    // Check if action requires host privileges
    if requiresHost(actionPayload.Action.Type) && !room.IsHost(conn.PlayerID) {
        errMsg, _ := NewErrorMessage("Only host can perform this action")
        conn.Send <- errMsg
        return
    }
```

### References
- [OWASP API Security Top 10](https://owasp.org/www-project-api-security/)
- [Go net/http Context Values](https://pkg.go.dev/context#WithValue)

---

## 2. Concurrency & Data Integrity

### Current Issues
- `room.EventLog` accessed without holding `room.mu`
- `CleanupStaleRooms` iterates `room.Players` map without room lock
- Race conditions in broadcast paths reading room state
- Slices copied unsafely for JSON marshaling

### Recommended Solution

**Encapsulate Room State:**
```go
// internal/core/room.go
// Add safe accessor methods:

// GetEventLogLength returns the current event log length safely
func (r *Room) GetEventLogLength() int {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return len(r.EventLog)
}

// GetEventsSince returns events from a given index
func (r *Room) GetEventsSince(startIndex int) []GameEvent {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    if startIndex >= len(r.EventLog) {
        return nil
    }
    
    // Copy events to avoid races
    events := make([]GameEvent, len(r.EventLog)-startIndex)
    copy(events, r.EventLog[startIndex:])
    return events
}

// IsAnyPlayerConnected safely checks connection status
func (r *Room) IsAnyPlayerConnected() bool {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    for _, player := range r.Players {
        if player.Connected {
            return true
        }
    }
    return false
}
```

**Fix Handler Access:**
```go
// internal/server/handlers.go - HandleStartGame
startIndex := room.GetEventLogLength()
if err := room.StartGame(game, config); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
}

// Now safe to read new events
newEvents := room.GetEventsSince(startIndex)
for _, event := range newEvents {
    s.connMgr.BroadcastEvent(roomCode, event)
}
```

**Fix Store Cleanup:**
```go
// internal/store/memory.go
func (s *MemoryStore) CleanupStaleRooms() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    const staleTimeout = 24 * time.Hour
    toDelete := make([]string, 0)

    for roomCode, room := range s.rooms {
        // Lock room before inspecting
        room.mu.RLock()
        
        if room.Status == core.RoomStatusFinished && time.Since(room.CreatedAt) > 1*time.Hour {
            room.mu.RUnlock()
            toDelete = append(toDelete, roomCode)
            continue
        }

        allDisconnected := true
        for _, player := range room.Players {
            if player.Connected {
                allDisconnected = false
                break
            }
        }
        room.mu.RUnlock()

        if allDisconnected && time.Since(room.CreatedAt) > staleTimeout {
            toDelete = append(toDelete, roomCode)
        }
    }

    for _, roomCode := range toDelete {
        delete(s.rooms, roomCode)
    }
    return nil
}
```

**Alternative: Use Helper Method:**
```go
// internal/store/memory.go
func (s *MemoryStore) CleanupStaleRooms() error {
    // ... existing lock code ...
    
    for roomCode, room := range s.rooms {
        if room.Status == core.RoomStatusFinished && time.Since(room.CreatedAt) > 1*time.Hour {
            toDelete = append(toDelete, roomCode)
            continue
        }

        // Use safe accessor instead of direct access
        if !room.IsAnyPlayerConnected() && time.Since(room.CreatedAt) > staleTimeout {
            toDelete = append(toDelete, roomCode)
        }
    }
    // ... cleanup ...
}
```

### Testing for Races
```bash
# Run with race detector
go test -race ./...

# For specific packages
go test -race ./internal/core
go test -race ./internal/store
```

### References
- [Go Data Race Detector](https://go.dev/blog/race-detector)
- [Go sync package docs](https://pkg.go.dev/sync)
- [Effective Go: Concurrency](https://go.dev/doc/effective_go#concurrency)

---

## 3. Network Resource Limiting

### Current Issues
- CORS allows all origins permanently
- WebSocket accepts all origins with TODO comment
- No read limits on WebSocket messages
- No read/write deadlines on connections
- HTTP handlers lack request size limits

### Recommended Solution

**Environment-Based CORS:**
```go
// cmd/server/main.go
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
        if allowedOrigin == "" {
            allowedOrigin = "http://localhost:5173" // Dev default
        }

        origin := r.Header.Get("Origin")
        if origin == allowedOrigin || allowedOrigin == "*" {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }
        
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-Token")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

**HTTP Request Size Limits:**
```go
// internal/server/handlers.go
func (s *Server) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
    // Limit request body to 1MB
    r.Body = http.MaxBytesReader(w, r.Body, 1*1024*1024)
    
    var req CreateRoomRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Request too large or malformed", http.StatusBadRequest)
        return
    }
    // ... rest of handler
}
```

**WebSocket Limits:**
```go
// internal/server/websocket.go
func (cm *ConnectionManager) HandleConnection(ctx context.Context, conn *websocket.Conn, roomCode string) {
    // Set read limit to 1MB
    conn.SetReadLimit(1 * 1024 * 1024)
    
    // Set read deadline for auth message (10 seconds)
    conn.SetReadDeadline(time.Now().Add(10 * time.Second))
    
    var authMsg ClientMessage
    if err := wsjson.Read(ctx, conn, &authMsg); err != nil {
        log.Printf("Failed to read auth message: %v", err)
        conn.Close(websocket.StatusPolicyViolation, "authentication timeout")
        return
    }
    
    // Clear deadline after auth
    conn.SetReadDeadline(time.Time{})
    
    // ... rest of auth logic ...
    
    // Restrict origin patterns
    allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
    if len(allowedOrigins) == 0 {
        allowedOrigins = []string{"localhost:*"}
    }
}

// Update websocket accept options in handlers.go:
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    // ... room verification ...
    
    allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
    if len(allowedOrigins) == 0 {
        allowedOrigins = []string{"localhost:*"}
    }
    
    conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
        OriginPatterns: allowedOrigins,
    })
    // ... rest ...
}
```

**Read with Context Timeout:**
```go
// internal/server/websocket.go - in readPump
func (c *Connection) readPump(cm *ConnectionManager, room *core.Room) {
    defer c.cancel()

    for {
        // Set 60-second read deadline per message
        c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        
        var msg ClientMessage
        if err := wsjson.Read(c.ctx, c.Conn, &msg); err != nil {
            if c.ctx.Err() != nil {
                return
            }
            log.Printf("Read error for player %s: %v", c.PlayerID, err)
            return
        }

        cm.handleClientMessage(c, room, msg)
    }
}
```

### References
- [nhooyr.io/websocket docs](https://pkg.go.dev/nhooyr.io/websocket)
- [Go net/http Security Best Practices](https://go.dev/doc/articles/wiki/#tmp_6)

---

## 4. Cryptographic Randomness

### Current Issues
- `math/rand.Shuffle` used without seeding in `werewolf/game.go`
- Role assignments are deterministic across server restarts
- Predictable shuffles undermine game fairness

### Recommended Solution

**Secure Shuffle Using crypto/rand:**
```go
// internal/games/werewolf/game.go
import (
    "crypto/rand"
    "math/big"
)

// secureShuffleRoles shuffles roles using cryptographically secure randomness
func secureShuffleRoles(roles []RoleType) {
    n := len(roles)
    for i := n - 1; i > 0; i-- {
        // Generate random index using crypto/rand
        jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
        if err != nil {
            // Fallback: this should never happen with crypto/rand
            panic(fmt.Sprintf("failed to generate random number: %v", err))
        }
        j := int(jBig.Int64())
        roles[i], roles[j] = roles[j], roles[i]
    }
}

// Update Initialize method:
func (g *Game) Initialize(config core.GameConfig, players []*core.Player) ([]core.GameEvent, error) {
    // ... validation ...

    // Shuffle and assign roles
    shuffledRoles := make([]RoleType, len(wConfig.Roles))
    copy(shuffledRoles, wConfig.Roles)
    
    // Use secure shuffle instead of math/rand
    secureShuffleRoles(shuffledRoles)

    // ... rest of initialization ...
}
```

**Alternative: Seeded math/rand (Less Secure):**
```go
// If crypto/rand is too heavy, at minimum seed math/rand globally once
import (
    "math/rand"
    "time"
)

// In init() or main():
func init() {
    rand.Seed(time.Now().UnixNano())
}

// Then use math/rand.Shuffle as before
```

### Why crypto/rand?
- `math/rand` is deterministic (same seed → same sequence)
- Server restarts with similar timestamps produce predictable sequences
- `crypto/rand` uses OS entropy sources for true randomness
- Fisher-Yates shuffle with crypto/rand is the gold standard for card games

### References
- [Go crypto/rand package](https://pkg.go.dev/crypto/rand)
- [Fisher-Yates Shuffle](https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle)

---

## 5. Background Routine Lifecycle

### Current Issues
- `cleanupRoutine` and `phaseCheckRoutine` run forever with no shutdown signal
- Goroutines leak on server restart/shutdown
- No context cancellation propagation

### Recommended Solution

**Context-Aware Background Routines:**
```go
// cmd/server/main.go
func main() {
    // ... setup ...

    // Create root context for shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Start cleanup goroutine with context
    go cleanupRoutine(ctx, memStore)
    go phaseCheckRoutine(ctx, memStore, srv)

    // ... HTTP server setup ...

    // Graceful shutdown
    go func() {
        log.Printf("Server starting on port %s", port)
        if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    // Wait for interrupt
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit

    log.Println("Shutting down server...")

    // Cancel background routines
    cancel()

    // Shutdown HTTP server
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownCancel()

    if err := httpServer.Shutdown(shutdownCtx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Server stopped")
}

// Updated cleanup routine
func cleanupRoutine(ctx context.Context, store store.Store) {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Println("Cleanup routine shutting down")
            return
        case <-ticker.C:
            if err := store.CleanupStaleRooms(); err != nil {
                log.Printf("Cleanup error: %v", err)
            } else {
                log.Println("Cleanup completed")
            }
        }
    }
}

// Updated phase check routine
func phaseCheckRoutine(ctx context.Context, store store.Store, srv *server.Server) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Println("Phase check routine shutting down")
            return
        case <-ticker.C:
            rooms, err := store.ListRooms()
            if err != nil {
                continue
            }

            for _, room := range rooms {
                if room.Status != "playing" || room.Game == nil {
                    continue
                }

                events, err := room.Game.CheckPhaseTimeout()
                if err != nil {
                    log.Printf("Phase check error for room %s: %v", room.ID, err)
                    continue
                }

                if len(events) > 0 {
                    room.AppendEvents(events)
                    for _, event := range events {
                        srv.ConnectionManager().BroadcastEvent(room.ID, event)
                    }
                }
            }
        }
    }
}
```

### References
- [Go Blog: Concurrency Patterns: Context](https://go.dev/blog/context)
- [Go context package](https://pkg.go.dev/context)

---

## 6. Structured Logging & Observability

### Current Issues
- Using global `log` package with unstructured messages
- No trace IDs or request correlation
- Difficult to query/filter logs in production

### Recommended Solution

**Migrate to log/slog (Go 1.21+):**
```go
// cmd/server/main.go
import (
    "log/slog"
    "os"
)

func main() {
    // Set up structured logger
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    slog.SetDefault(logger)

    // ... rest of main ...
}

// Usage in handlers:
func (s *Server) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
    // ... decode request ...

    slog.Info("creating room",
        "gameType", req.GameType,
        "hostName", req.DisplayName,
        "maxPlayers", req.MaxPlayers,
    )

    // ... create room ...

    slog.Info("room created",
        "roomCode", roomCode,
        "hostID", hostPlayer.ID,
    )
    // ... response ...
}

// In websocket:
func (cm *ConnectionManager) HandleConnection(ctx context.Context, conn *websocket.Conn, roomCode string) {
    // ... auth ...

    slog.Info("player connected",
        "playerID", player.ID,
        "playerName", player.DisplayName,
        "roomCode", roomCode,
    )
    // ... rest ...
}
```

**Add Request IDs:**
```go
// internal/server/middleware.go
func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        ctx := context.WithValue(r.Context(), "requestID", requestID)
        w.Header().Set("X-Request-ID", requestID)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Use in handlers:
func (s *Server) HandleStartGame(w http.ResponseWriter, r *http.Request) {
    requestID := r.Context().Value("requestID").(string)
    
    slog.Info("starting game",
        "requestID", requestID,
        "roomCode", r.PathValue("code"),
    )
    // ... handler logic ...
}
```

### References
- [Go log/slog package](https://pkg.go.dev/log/slog)
- [Structured Logging in Go](https://go.dev/blog/slog)

---

## 7. Testing Strategy

### Current Issues
- No unit tests for any package
- No integration tests for handlers or WebSocket flows
- No way to validate concurrency fixes or security changes

### Why Tests Are Essential (Not Optional)

Testing is **absolutely critical** for this backend because:

1. **Concurrency Safety**: Race conditions only appear under load; `-race` detector requires test execution
2. **Game Logic Correctness**: Complex rules (Tanner wins, role swaps, vote counting) must work perfectly
3. **Security Validation**: Auth/authorization bugs can't be manually verified in production
4. **Refactoring Confidence**: Can't safely fix concurrency issues or add auth without tests
5. **Go Community Standard**: Go projects typically have 60-80%+ test coverage; it's idiomatic

### Recommended Test Coverage

**Critical Priority (Write First):**
```
internal/core/room_test.go         # Room state, concurrency, security
internal/core/event_test.go        # Event visibility rules
internal/store/memory_test.go      # Store operations & cleanup
internal/games/werewolf/game_test.go   # Role assignment, voting, results
internal/games/werewolf/phases_test.go # Phase transitions
internal/server/handlers_test.go    # Auth middleware, host checks
```

**High Priority (Write Soon):**
```
internal/server/websocket_test.go  # Connection handling, broadcasts
go test -race ./...                # Must pass on all packages
```

**Medium Priority (Nice to Have):**
```
test/integration/game_flow_test.go     # Full game simulation
test/integration/websocket_test.go     # Multi-client scenarios
```

### Recommended Solution

#### **1. Core Domain Tests**

**Room Concurrency & Security:**
```go
// internal/core/room_test.go
package core

import (
    "sync"
    "testing"
)

func TestRoom_AddPlayer(t *testing.T) {
    tests := []struct {
        name        string
        maxPlayers  int
        players     int
        expectError bool
    }{
        {"add to empty room", 5, 1, false},
        {"add to full room", 2, 1, true},
        {"add to room at capacity", 3, 2, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            host := NewPlayer("host")
            room := NewRoom("TEST01", "werewolf", host, tt.maxPlayers)

            // Add additional players
            for i := 0; i < tt.players-1; i++ {
                p := NewPlayer("player")
                if err := room.AddPlayer(p); err != nil {
                    t.Fatalf("failed to add player: %v", err)
                }
            }

            // Try adding one more
            newPlayer := NewPlayer("new")
            err := room.AddPlayer(newPlayer)

            if tt.expectError && err == nil {
                t.Error("expected error but got none")
            }
            if !tt.expectError && err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}

// Critical: Test concurrent access (catches race conditions)
func TestRoom_ConcurrentAccess(t *testing.T) {
    host := NewPlayer("host")
    room := NewRoom("TEST01", "werewolf", host, 10)

    // Spawn multiple goroutines accessing room
    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            p := NewPlayer("player")
            room.AddPlayer(p)
        }()
    }
    wg.Wait()

    if len(room.Players) != 6 { // host + 5
        t.Errorf("expected 6 players, got %d", len(room.Players))
    }
}

// Critical: Test concurrent event appending (catches races)
func TestRoom_ConcurrentEventAppend(t *testing.T) {
    host := NewPlayer("host")
    room := NewRoom("TEST01", "werewolf", host, 10)
    
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            event, _ := NewPublicEvent("test", "system", nil)
            room.AppendEvent(event)
        }()
    }
    wg.Wait()
    
    if len(room.EventLog) != 10 {
        t.Errorf("expected 10 events, got %d (race condition)", len(room.EventLog))
    }
}

// Security: Test host-only operations
func TestRoom_IsHost(t *testing.T) {
    host := NewPlayer("host")
    other := NewPlayer("other")
    room := NewRoom("TEST01", "werewolf", host, 10)
    room.AddPlayer(other)

    if !room.IsHost(host.ID) {
        t.Error("host should be identified as host")
    }
    if room.IsHost(other.ID) {
        t.Error("non-host should not be identified as host")
    }
}
```

**Event Visibility Tests:**
```go
// internal/core/event_test.go
package core

import "testing"

func TestGameEvent_CanPlayerSee(t *testing.T) {
    tests := []struct {
        name       string
        visibility EventVisibility
        playerID   string
        canSee     bool
    }{
        {
            name:       "public event visible to all",
            visibility: EventVisibility{Public: true},
            playerID:   "player1",
            canSee:     true,
        },
        {
            name:       "private event visible to target",
            visibility: EventVisibility{Public: false, PlayerIDs: []string{"player1"}},
            playerID:   "player1",
            canSee:     true,
        },
        {
            name:       "private event not visible to others",
            visibility: EventVisibility{Public: false, PlayerIDs: []string{"player1"}},
            playerID:   "player2",
            canSee:     false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            event := GameEvent{Visibility: tt.visibility}
            if got := event.CanPlayerSee(tt.playerID); got != tt.canSee {
                t.Errorf("CanPlayerSee() = %v, want %v", got, tt.canSee)
            }
        })
    }
}
```

#### **2. Game Logic Tests**

**Werewolf Win Conditions:**
```go
// internal/games/werewolf/game_test.go
package werewolf

import (
    "testing"
    "github.com/yourusername/roundtable/internal/core"
)

func TestGame_TannerWinCondition(t *testing.T) {
    game := NewGame().(*Game)
    
    // Setup: 3 players, one is Tanner
    players := []*core.Player{
        core.NewPlayer("Alice"),
        core.NewPlayer("Bob"),
        core.NewPlayer("Charlie"),
    }
    
    config := &Config{
        Roles: []RoleType{
            RoleTanner,      // Alice gets Tanner
            RoleVillager,    // Bob
            RoleVillager,    // Charlie
            RoleWerewolf,    // Center
            RoleVillager,    // Center
            RoleVillager,    // Center
        },
    }
    
    game.Initialize(config, players)
    
    // Alice (Tanner) is voted out
    game.votes = map[string]string{
        players[0].ID: players[0].ID, // Alice votes Alice
        players[1].ID: players[0].ID, // Bob votes Alice
        players[2].ID: players[0].ID, // Charlie votes Alice
    }
    
    results := game.calculateResults()
    
    // Tanner should win alone
    if len(results.Winners) != 1 {
        t.Errorf("expected 1 winner (Tanner), got %d", len(results.Winners))
    }
    if results.Winners[0] != players[0].ID {
        t.Error("Tanner should be the winner")
    }
    if results.WinReason != "Tanner wins by getting eliminated!" {
        t.Errorf("unexpected win reason: %s", results.WinReason)
    }
}

func TestGame_VillageTeamWins(t *testing.T) {
    game := NewGame().(*Game)
    
    players := []*core.Player{
        core.NewPlayer("Alice"),
        core.NewPlayer("Bob"),
        core.NewPlayer("Charlie"),
    }
    
    config := &Config{
        Roles: []RoleType{
            RoleWerewolf,    // Alice is werewolf
            RoleVillager,    // Bob
            RoleVillager,    // Charlie
            RoleSeer,        // Center
            RoleVillager,    // Center
            RoleVillager,    // Center
        },
    }
    
    game.Initialize(config, players)
    
    // Werewolf (Alice) is voted out
    game.votes = map[string]string{
        players[0].ID: players[0].ID,
        players[1].ID: players[0].ID,
        players[2].ID: players[0].ID,
    }
    
    results := game.calculateResults()
    
    // Village team (Bob and Charlie) should win
    if len(results.Winners) != 2 {
        t.Errorf("expected 2 winners (village team), got %d", len(results.Winners))
    }
    if results.WinReason != "Village team eliminated a werewolf!" {
        t.Errorf("unexpected win reason: %s", results.WinReason)
    }
}

func TestGame_RobberSwapChangesRole(t *testing.T) {
    game := NewGame().(*Game)
    
    players := []*core.Player{
        core.NewPlayer("Robber"),
        core.NewPlayer("Target"),
    }
    
    config := &Config{
        Roles: []RoleType{
            RoleRobber,      // Player 0
            RoleWerewolf,    // Player 1
            RoleVillager,    // Center
            RoleVillager,    // Center
            RoleVillager,    // Center
        },
    }
    
    game.Initialize(config, players)
    game.phase = PhaseNight
    
    // Robber swaps with werewolf
    action := core.Action{
        Type:    "robber_swap",
        Payload: mustMarshal(RobberSwapPayload{TargetID: players[1].ID}),
    }
    
    events, err := game.ProcessAction(players[0].ID, action)
    if err != nil {
        t.Fatalf("robber swap failed: %v", err)
    }
    
    // Check that robber now has werewolf role
    if game.roleAssignments[players[0].ID] != RoleWerewolf {
        t.Error("robber should now be werewolf")
    }
    if game.roleAssignments[players[1].ID] != RoleRobber {
        t.Error("target should now be robber")
    }
    
    // Check that event was emitted
    if len(events) != 1 {
        t.Errorf("expected 1 event, got %d", len(events))
    }
}
```

#### **3. Handler & Security Tests**

**Authentication Tests:**
```go
// internal/server/handlers_test.go
package server

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/yourusername/roundtable/internal/store"
)

func TestHandleCreateRoom(t *testing.T) {
    memStore := store.NewMemoryStore()
    srv := NewServer(memStore)

    req := CreateRoomRequest{
        GameType:    "werewolf",
        DisplayName: "TestHost",
        MaxPlayers:  8,
    }
    body, _ := json.Marshal(req)

    r := httptest.NewRequest("POST", "/api/rooms", bytes.NewReader(body))
    w := httptest.NewRecorder()

    srv.HandleCreateRoom(w, r)

    if w.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", w.Code)
    }

    var resp CreateRoomResponse
    if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
        t.Fatalf("failed to decode response: %v", err)
    }

    if resp.RoomCode == "" {
        t.Error("expected room code, got empty string")
    }
    if resp.SessionToken == "" {
        t.Error("expected session token, got empty string")
    }
}

func TestHandleCreateRoom_InvalidGameType(t *testing.T) {
    memStore := store.NewMemoryStore()
    srv := NewServer(memStore)

    req := CreateRoomRequest{
        GameType:    "invalid",
        DisplayName: "TestHost",
    }
    body, _ := json.Marshal(req)

    r := httptest.NewRequest("POST", "/api/rooms", bytes.NewReader(body))
    w := httptest.NewRecorder()

    srv.HandleCreateRoom(w, r)

    if w.Code != http.StatusBadRequest {
        t.Errorf("expected status 400, got %d", w.Code)
    }
}

// Security: Test auth middleware rejects missing token
func TestAuthMiddleware_MissingToken(t *testing.T) {
    memStore := store.NewMemoryStore()
    srv := NewServer(memStore)
    
    handler := srv.authMiddleware(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })
    
    r := httptest.NewRequest("POST", "/api/rooms/TEST01/start", nil)
    w := httptest.NewRecorder()
    
    handler(w, r)
    
    if w.Code != http.StatusUnauthorized {
        t.Errorf("expected 401, got %d", w.Code)
    }
}

// Security: Test host-only operations
func TestHandleStartGame_RequiresHost(t *testing.T) {
    memStore := store.NewMemoryStore()
    srv := NewServer(memStore)
    
    // Create room with host
    host := core.NewPlayer("Host")
    room := core.NewRoom("TEST01", "werewolf", host, 10)
    memStore.CreateRoom(room)
    
    // Add non-host player
    nonHost := core.NewPlayer("Player")
    room.AddPlayer(nonHost)
    
    // Non-host tries to start game
    req := StartGameRequest{Config: json.RawMessage(`{}`)}
    body, _ := json.Marshal(req)
    
    r := httptest.NewRequest("POST", "/api/rooms/TEST01/start", bytes.NewReader(body))
    r.Header.Set("X-Session-Token", nonHost.SessionToken)
    w := httptest.NewRecorder()
    
    // Apply auth middleware
    handler := srv.authMiddleware(srv.HandleStartGame)
    handler(w, r)
    
    if w.Code != http.StatusForbidden {
        t.Errorf("expected 403 Forbidden, got %d", w.Code)
    }
}
```

#### **4. Store Tests**

**Cleanup & Concurrency:**
```go
// internal/store/memory_test.go
package store

import (
    "sync"
    "testing"
    "time"
    
    "github.com/yourusername/roundtable/internal/core"
)

func TestMemoryStore_CleanupStaleRooms(t *testing.T) {
    store := NewMemoryStore()
    
    // Create finished room (old)
    oldFinished := core.NewRoom("OLD01", "werewolf", core.NewPlayer("host"), 10)
    oldFinished.SetStatus(core.RoomStatusFinished)
    oldFinished.CreatedAt = time.Now().Add(-2 * time.Hour)
    store.CreateRoom(oldFinished)
    
    // Create finished room (recent)
    recentFinished := core.NewRoom("NEW01", "werewolf", core.NewPlayer("host"), 10)
    recentFinished.SetStatus(core.RoomStatusFinished)
    store.CreateRoom(recentFinished)
    
    // Run cleanup
    if err := store.CleanupStaleRooms(); err != nil {
        t.Fatalf("cleanup failed: %v", err)
    }
    
    // Old finished room should be deleted
    if _, err := store.GetRoom("OLD01"); err != ErrRoomNotFound {
        t.Error("old finished room should be cleaned up")
    }
    
    // Recent finished room should remain
    if _, err := store.GetRoom("NEW01"); err != nil {
        t.Error("recent finished room should not be cleaned up")
    }
}

func TestMemoryStore_ConcurrentAccess(t *testing.T) {
    store := NewMemoryStore()
    
    // Spawn multiple goroutines creating rooms
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            roomCode := "TEST" + string(rune('0'+id))
            room := core.NewRoom(roomCode, "werewolf", core.NewPlayer("host"), 10)
            store.CreateRoom(room)
        }(i)
    }
    wg.Wait()
    
    rooms, err := store.ListRooms()
    if err != nil {
        t.Fatalf("failed to list rooms: %v", err)
    }
    
    if len(rooms) != 10 {
        t.Errorf("expected 10 rooms, got %d", len(rooms))
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# With race detector (CRITICAL - run this regularly)
go test -race ./...

# With coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Specific packages
go test ./internal/core
go test ./internal/server
go test ./internal/games/werewolf

# Run only specific tests
go test -run TestRoom_ConcurrentAccess ./internal/core
go test -run TestGame_TannerWinCondition ./internal/games/werewolf

# Verbose output
go test -v ./...

# With timeout
go test -timeout 30s ./...
```

### CI/CD Integration

**.github/workflows/test.yml:**
```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      
      - name: Run tests
        run: |
          cd backend
          go test -v ./...
      
      - name: Race detector
        run: |
          cd backend
          go test -race -timeout 30s ./...
      
      - name: Coverage
        run: |
          cd backend
          go test -coverprofile=coverage.out ./...
          go tool cover -func=coverage.out
```

### Minimal Test Suite (Start Here)

If you only have time for **one test file**, make it this:

```go
// internal/core/room_test.go
// Contains: AddPlayer, ConcurrentAccess, ConcurrentEventAppend, IsHost tests
```

If you can write **three test files**, prioritize:
1. `internal/core/room_test.go` (concurrency safety)
2. `internal/games/werewolf/game_test.go` (game logic correctness)
3. `internal/server/handlers_test.go` (security validation)

### Helper Utilities

```go
// internal/core/testing.go
package core

import "encoding/json"

// Helper for tests
func mustMarshal(v interface{}) json.RawMessage {
    b, err := json.Marshal(v)
    if err != nil {
        panic(err)
    }
    return b
}
```

### References
- [Go testing package](https://pkg.go.dev/testing)
- [Table Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [Go net/http/httptest](https://pkg.go.dev/net/http/httptest)
- [Advanced Testing with Go](https://go.dev/blog/subtests)

---

## 8. Implementation Roadmap

### Phase 1: Foundation (Week 1)
1. Add structured logging (slog)
2. Implement context cancellation for background routines
3. Write tests for `internal/core` (room, player, event)

### Phase 2: Security (Week 2)
4. Add authentication middleware
5. Implement host-only authorization checks
6. Replace math/rand with crypto/rand for shuffles
7. Add HTTP/WebSocket resource limits

### Phase 3: Concurrency (Week 3)
8. Refactor room to use accessor methods
9. Fix store cleanup to lock rooms
10. Run race detector and fix remaining issues
11. Write concurrency tests

### Phase 4: Network Hardening (Week 4)
12. Implement environment-based CORS
13. Add WebSocket origin restrictions
14. Set read/write deadlines on connections
15. Add request size limits

### Phase 5: Verification (Week 5)
16. Run full test suite with `-race`
17. Load test WebSocket connections
18. Audit logs for security events
19. Document deployment configuration

### Rollout Strategy
- Deploy changes incrementally to staging
- Use feature flags for new auth middleware
- Monitor error rates and latency
- Keep rollback plan ready for each phase

### Verification Checklist
- [ ] All tests pass with `go test -race ./...`
- [ ] No panics in concurrent scenarios
- [ ] Auth middleware rejects invalid tokens
- [ ] Host-only actions return 403 for non-hosts
- [ ] WebSocket connections respect origin restrictions
- [ ] Background routines shut down cleanly on SIGINT
- [ ] Role shuffles produce different results each game
- [ ] Structured logs include trace IDs and context

---

## Additional Resources

### Go Official Documentation
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Security Best Practices](https://go.dev/doc/articles/wiki/)

### Security
- [OWASP Go Secure Coding Practices](https://cheatsheetseries.owasp.org/cheatsheets/Go_Secure_Coding_Practices_Cheat_Sheet.html)
- [Gosec - Go Security Checker](https://github.com/securego/gosec)

### Concurrency
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Go Memory Model](https://go.dev/ref/mem)

### Testing
- [Testing Best Practices](https://go.dev/doc/effective_go#testing)
- [Advanced Testing with Go](https://go.dev/blog/subtests)

### WebSockets
- [nhooyr.io/websocket Examples](https://github.com/nhooyr/websocket/tree/master/examples)
- [WebSocket Security](https://owasp.org/www-community/vulnerabilities/WebSocket_Security)

---

## Conclusion

This remediation plan addresses critical security, concurrency, and architectural gaps in the Roundtable backend. By following these recommendations and implementing changes incrementally with thorough testing, the codebase will align with Go community best practices and provide a secure, reliable foundation for the multiplayer game platform.

**Priority Order:**
1. **Critical**: Authentication, concurrency fixes, crypto/rand
2. **High**: Context cancellation, resource limits, testing
3. **Medium**: Structured logging, CORS restrictions
4. **Low**: Documentation, monitoring enhancements

