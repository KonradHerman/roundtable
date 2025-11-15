package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/yourusername/roundtable/internal/server"
	"github.com/yourusername/roundtable/internal/store"
)

func main() {
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

	// WebSocket route
	mux.HandleFunc("GET /api/rooms/{code}/ws", srv.HandleWebSocket)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
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

	// Start cleanup goroutine
	go cleanupRoutine(memStore)

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

// corsMiddleware adds CORS headers for development.
// TODO: Restrict origins in production.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// cleanupRoutine periodically cleans up stale rooms.
func cleanupRoutine(store store.Store) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		if err := store.CleanupStaleRooms(); err != nil {
			log.Printf("Cleanup error: %v", err)
		} else {
			log.Println("Cleanup completed")
		}
	}
}
