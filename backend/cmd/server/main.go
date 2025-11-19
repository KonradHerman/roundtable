package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/rs/cors"

	"github.com/yourusername/roundtable/internal/server"
	"github.com/yourusername/roundtable/internal/store"
)

// parseAllowedOrigins parses comma-separated origins from environment variable
// and trims whitespace from each origin
func parseAllowedOrigins(originsEnv string) []string {
	if originsEnv == "" {
		return nil
	}

	parts := strings.Split(originsEnv, ",")
	origins := make([]string, 0, len(parts))

	for _, origin := range parts {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	return origins
}

// isRailwayOrigin checks if an origin is from Railway (*.up.railway.app)
func isRailwayOrigin(origin string) bool {
	// Remove protocol
	withoutProtocol := strings.TrimPrefix(strings.TrimPrefix(origin, "https://"), "http://")
	// Check if it ends with .up.railway.app
	return strings.HasSuffix(withoutProtocol, ".up.railway.app")
}

// isRailwayEnvironment checks if we're running in Railway
func isRailwayEnvironment() bool {
	return os.Getenv("RAILWAY_ENVIRONMENT") != "" || os.Getenv("RAILWAY_PROJECT_ID") != ""
}

func main() {
	// Set up structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Create root context for shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create store
	memStore := store.NewMemoryStore()

	// Create server
	srv := server.NewServer(memStore)

	// Setup routes
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("POST /api/rooms", srv.HandleCreateRoom)
	mux.HandleFunc("GET /api/rooms/{code}", srv.HandleGetRoom)
	mux.HandleFunc("POST /api/rooms/{code}/join", srv.HandleJoinRoom)
	mux.HandleFunc("POST /api/rooms/{code}/start", srv.HandleStartGame)
	mux.HandleFunc("POST /api/rooms/{code}/reset", srv.HandleResetGame)

	// WebSocket route
	mux.HandleFunc("GET /api/rooms/{code}/ws", srv.HandleWebSocket)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// CORS middleware (rs/cors)
	allowedOrigins := parseAllowedOrigins(os.Getenv("ALLOWED_ORIGIN"))
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"http://localhost:5173"} // Dev default
	}

	// Check if wildcard is used
	hasWildcard := false
	for _, origin := range allowedOrigins {
		if origin == "*" {
			hasWildcard = true
			break
		}
	}

	inRailway := isRailwayEnvironment()

	var c *cors.Cors
	if hasWildcard {
		// If wildcard is used (not recommended for production with credentials), adjust options
		c = cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "X-Session-Token"},
			// AllowCredentials cannot be true with AllowedOrigins: []string{"*"}
			AllowCredentials: false,
		})
		slog.Info("CORS configured with wildcard (credentials disabled)")
	} else {
		// Custom origin validator: allow configured origins + Railway domains in Railway environment
		allowOriginFunc := func(origin string) bool {
			// Check if origin is in the allowed list
			for _, allowed := range allowedOrigins {
				if allowed == origin {
					return true
				}
			}

			// In Railway environment, also allow any *.up.railway.app origin
			if inRailway && isRailwayOrigin(origin) {
				slog.Info("allowing Railway origin", "origin", origin)
				return true
			}

			return false
		}

		c = cors.New(cors.Options{
			AllowOriginFunc:  allowOriginFunc,
			AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "X-Session-Token"},
			AllowCredentials: true,
			Debug:            os.Getenv("CORS_DEBUG") == "true",
		})

		if inRailway {
			slog.Info("CORS configured for Railway", "allowed_origins", allowedOrigins, "railway_wildcard", "*.up.railway.app")
		} else {
			slog.Info("CORS configured", "allowed_origins", allowedOrigins)
		}
	}

	handler := c.Handler(mux)

	// HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start cleanup goroutine with context
	go cleanupRoutine(ctx, memStore)

	// Start phase check routine for game timers
	go phaseCheckRoutine(ctx, memStore, srv)

	// Graceful shutdown
	go func() {
		slog.Info("server starting", "port", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Info("shutting down server")

	// Cancel background routines
	cancel()

	// Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}

// cleanupRoutine periodically cleans up stale rooms.
func cleanupRoutine(ctx context.Context, store store.Store) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("cleanup routine shutting down")
			return
		case <-ticker.C:
			if err := store.CleanupStaleRooms(); err != nil {
				slog.Error("cleanup error", "error", err)
			} else {
				slog.Info("cleanup completed")
			}
		}
	}
}

// phaseCheckRoutine periodically checks if game phases should advance
// Optimized: uses a priority queue (min-heap) via store.PopExpiredRooms to avoid O(N) polling.
func phaseCheckRoutine(ctx context.Context, store store.Store, srv *server.Server) {
	ticker := time.NewTicker(100 * time.Millisecond) // Check more frequently for better precision
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("phase check routine shutting down")
			return
		case <-ticker.C:
			// Get rooms that have expired phases
			rooms, err := store.PopExpiredRooms(time.Now())
			if err != nil {
				slog.Error("failed to pop expired rooms", "error", err)
				continue
			}

			for _, room := range rooms {
				// CheckAndAdvancePhase safely acquires the room's lock internally
				events, shouldUpdate, nextPhaseTime, err := room.CheckAndAdvancePhase()
				if err != nil {
					slog.Error("phase check error", "roomID", room.ID, "error", err)
					continue
				}

				// If there are events, broadcast them
				// AppendEvents will acquire its own lock
				if len(events) > 0 {
					room.AppendEvents(events)
					for _, event := range events {
						srv.ConnectionManager().BroadcastEvent(room.ID, event)
					}
				}

				// Update the timer for the next phase (if any)
				if shouldUpdate && nextPhaseTime.After(time.Now()) {
					store.UpdateRoomTimer(room.ID, nextPhaseTime)
				}
			}
		}
	}
}
