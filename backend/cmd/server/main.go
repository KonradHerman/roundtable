package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/KonradHerman/roundtable/internal/server"
	"github.com/KonradHerman/roundtable/internal/store"
)

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

	// Debug endpoint to check env vars
	mux.HandleFunc("GET /debug/env", func(w http.ResponseWriter, r *http.Request) {
		allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "(not set - using default)"
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ALLOWED_ORIGIN=" + allowedOrigin))
	})

	// CORS middleware (for development)
	handler := corsMiddleware(mux)

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

// corsMiddleware adds CORS headers with environment-based origin restrictions.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "http://localhost:5173" // Dev default
		}

		origin := r.Header.Get("Origin")

		// Check if origin is allowed and set appropriate CORS headers
		if allowedOrigin == "*" {
			// Wildcard: allow any origin
			// Note: When using credentials, we must echo the specific origin, not "*"
			// For now, we'll use "*" and note that credentials won't work with wildcard
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-Token")
			// Cannot use credentials with "*", so we don't set Allow-Credentials
		} else if origin == allowedOrigin {
			// Specific origin match
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-Token")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else if origin != "" {
			// Origin provided but doesn't match - log warning and don't set CORS headers
			slog.Warn("rejected CORS request", "origin", origin, "allowed", allowedOrigin)
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
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
func phaseCheckRoutine(ctx context.Context, store store.Store, srv *server.Server) {
	ticker := time.NewTicker(1 * time.Second) // Check every second
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("phase check routine shutting down")
			return
		case <-ticker.C:
			rooms, err := store.ListRooms()
			if err != nil {
				continue
			}

			for _, room := range rooms {
				// Only check rooms with active games
				if room.Status != "playing" || room.Game == nil {
					continue
				}

				// Check if phase should advance
				events, err := room.Game.CheckPhaseTimeout()
				if err != nil {
					slog.Error("phase check error", "roomID", room.ID, "error", err)
					continue
				}

				// If there are events, broadcast them
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
