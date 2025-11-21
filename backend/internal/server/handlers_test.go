package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KonradHerman/roundtable/internal/store"
)

func TestHandleCreateRoom(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		method         string
		body           interface{}
		wantStatusCode int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:   "successfully create werewolf room",
			method: http.MethodPost,
			body: CreateRoomRequest{
				GameType:    "werewolf",
				DisplayName: "Alice",
				MaxPlayers:  10,
			},
			wantStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var resp CreateRoomResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.RoomCode == "" {
					t.Error("room code is empty")
				}
				if len(resp.RoomCode) != 6 {
					t.Errorf("room code should be 6 characters, got %d", len(resp.RoomCode))
				}
				if resp.SessionToken == "" {
					t.Error("session token is empty")
				}
				if resp.PlayerID == "" {
					t.Error("player ID is empty")
				}
			},
		},
		{
			name:   "successfully create room with default max players",
			method: http.MethodPost,
			body: CreateRoomRequest{
				GameType:    "werewolf",
				DisplayName: "Bob",
			},
			wantStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var resp CreateRoomResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.RoomCode == "" {
					t.Error("room code is empty")
				}
			},
		},
		{
			name:           "fail with wrong HTTP method",
			method:         http.MethodGet,
			body:           CreateRoomRequest{},
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:   "fail with unknown game type",
			method: http.MethodPost,
			body: CreateRoomRequest{
				GameType:    "unknown-game",
				DisplayName: "Charlie",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fail with malformed JSON",
			method:         http.MethodPost,
			body:           "not a json object",
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create server with in-memory store
			memStore := store.NewMemoryStore()
			server := NewServer(memStore)

			// Prepare request body
			var bodyReader io.Reader
			if bodyStr, ok := tt.body.(string); ok {
				bodyReader = bytes.NewBufferString(bodyStr)
			} else {
				bodyBytes, err := json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
				bodyReader = bytes.NewBuffer(bodyBytes)
			}

			// Create request
			req := httptest.NewRequest(tt.method, "/api/rooms", bodyReader)
			rec := httptest.NewRecorder()

			// Execute handler
			server.HandleCreateRoom(rec, req)

			// Check status code
			if rec.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, rec.Code)
			}

			// Check response if success
			if tt.wantStatusCode == http.StatusOK && tt.checkResponse != nil {
				tt.checkResponse(t, rec.Body.Bytes())
			}
		})
	}
}

func TestHandleJoinRoom(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupRoom      func(s *Server) string // Returns room code
		method         string
		roomCode       string
		body           interface{}
		wantStatusCode int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successfully join existing room",
			setupRoom: func(s *Server) string {
				req := CreateRoomRequest{
					GameType:    "werewolf",
					DisplayName: "Host",
					MaxPlayers:  10,
				}
				bodyBytes, _ := json.Marshal(req)
				httpReq := httptest.NewRequest(http.MethodPost, "/api/rooms", bytes.NewBuffer(bodyBytes))
				rec := httptest.NewRecorder()
				s.HandleCreateRoom(rec, httpReq)

				var resp CreateRoomResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)
				return resp.RoomCode
			},
			method: http.MethodPost,
			body: JoinRoomRequest{
				DisplayName: "Player1",
			},
			wantStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var resp JoinRoomResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.SessionToken == "" {
					t.Error("session token is empty")
				}
				if resp.PlayerID == "" {
					t.Error("player ID is empty")
				}
				if resp.RoomCode == "" {
					t.Error("room code is empty")
				}
			},
		},
		{
			name: "fail when joining non-existent room",
			setupRoom: func(s *Server) string {
				return "NOROOM"
			},
			method: http.MethodPost,
			body: JoinRoomRequest{
				DisplayName: "Player1",
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "fail with empty display name",
			setupRoom: func(s *Server) string {
				req := CreateRoomRequest{
					GameType:    "werewolf",
					DisplayName: "Host",
					MaxPlayers:  10,
				}
				bodyBytes, _ := json.Marshal(req)
				httpReq := httptest.NewRequest(http.MethodPost, "/api/rooms", bytes.NewBuffer(bodyBytes))
				rec := httptest.NewRecorder()
				s.HandleCreateRoom(rec, httpReq)

				var resp CreateRoomResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)
				return resp.RoomCode
			},
			method: http.MethodPost,
			body: JoinRoomRequest{
				DisplayName: "",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fail when joining full room",
			setupRoom: func(s *Server) string {
				// Create room with max 2 players (host + 1)
				req := CreateRoomRequest{
					GameType:    "werewolf",
					DisplayName: "Host",
					MaxPlayers:  2,
				}
				bodyBytes, _ := json.Marshal(req)
				httpReq := httptest.NewRequest(http.MethodPost, "/api/rooms", bytes.NewBuffer(bodyBytes))
				rec := httptest.NewRecorder()
				s.HandleCreateRoom(rec, httpReq)

				var resp CreateRoomResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				// Add first player to fill the room
				joinReq := JoinRoomRequest{DisplayName: "Player1"}
				joinBytes, _ := json.Marshal(joinReq)
				joinHttpReq := httptest.NewRequest(http.MethodPost, "/api/rooms/"+resp.RoomCode+"/join", bytes.NewBuffer(joinBytes))
				joinHttpReq.SetPathValue("code", resp.RoomCode)
				joinRec := httptest.NewRecorder()
				s.HandleJoinRoom(joinRec, joinHttpReq)

				return resp.RoomCode
			},
			method: http.MethodPost,
			body: JoinRoomRequest{
				DisplayName: "Player2",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fail with wrong HTTP method",
			setupRoom: func(s *Server) string {
				req := CreateRoomRequest{
					GameType:    "werewolf",
					DisplayName: "Host",
				}
				bodyBytes, _ := json.Marshal(req)
				httpReq := httptest.NewRequest(http.MethodPost, "/api/rooms", bytes.NewBuffer(bodyBytes))
				rec := httptest.NewRecorder()
				s.HandleCreateRoom(rec, httpReq)

				var resp CreateRoomResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)
				return resp.RoomCode
			},
			method:         http.MethodGet,
			body:           JoinRoomRequest{DisplayName: "Player1"},
			wantStatusCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create server with in-memory store
			memStore := store.NewMemoryStore()
			server := NewServer(memStore)

			// Setup room if needed
			roomCode := ""
			if tt.setupRoom != nil {
				roomCode = tt.setupRoom(server)
			}

			// Use explicit room code or the one from setup
			if tt.roomCode != "" {
				roomCode = tt.roomCode
			}

			// Prepare request body
			var bodyReader io.Reader
			if bodyStr, ok := tt.body.(string); ok {
				bodyReader = bytes.NewBufferString(bodyStr)
			} else {
				bodyBytes, err := json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
				bodyReader = bytes.NewBuffer(bodyBytes)
			}

			// Create request
			req := httptest.NewRequest(tt.method, "/api/rooms/"+roomCode+"/join", bodyReader)
			req.SetPathValue("code", roomCode)
			rec := httptest.NewRecorder()

			// Execute handler
			server.HandleJoinRoom(rec, req)

			// Check status code
			if rec.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d (body: %s)", tt.wantStatusCode, rec.Code, rec.Body.String())
			}

			// Check response if success
			if tt.wantStatusCode == http.StatusOK && tt.checkResponse != nil {
				tt.checkResponse(t, rec.Body.Bytes())
			}
		})
	}
}

func TestHandleGetRoom(t *testing.T) {
	t.Parallel()

	// Create server with in-memory store
	memStore := store.NewMemoryStore()
	server := NewServer(memStore)

	// Create a room first
	createReq := CreateRoomRequest{
		GameType:    "werewolf",
		DisplayName: "Host",
		MaxPlayers:  10,
	}
	bodyBytes, _ := json.Marshal(createReq)
	createHttpReq := httptest.NewRequest(http.MethodPost, "/api/rooms", bytes.NewBuffer(bodyBytes))
	createRec := httptest.NewRecorder()
	server.HandleCreateRoom(createRec, createHttpReq)

	var createResp CreateRoomResponse
	json.Unmarshal(createRec.Body.Bytes(), &createResp)

	tests := []struct {
		name           string
		roomCode       string
		wantStatusCode int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:           "successfully get room",
			roomCode:       createResp.RoomCode,
			wantStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var room map[string]interface{}
				if err := json.Unmarshal(body, &room); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if room["id"] != createResp.RoomCode {
					t.Errorf("expected room code '%s', got '%s'", createResp.RoomCode, room["id"])
				}
				if room["gameType"] != "werewolf" {
					t.Errorf("expected game type 'werewolf', got '%s'", room["gameType"])
				}
			},
		},
		{
			name:           "fail when room not found",
			roomCode:       "NOROOM",
			wantStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/rooms/"+tt.roomCode, nil)
			req.SetPathValue("code", tt.roomCode)
			rec := httptest.NewRecorder()

			// Execute handler
			server.HandleGetRoom(rec, req)

			// Check status code
			if rec.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, rec.Code)
			}

			// Check response if success
			if tt.wantStatusCode == http.StatusOK && tt.checkResponse != nil {
				tt.checkResponse(t, rec.Body.Bytes())
			}
		})
	}
}

func TestConcurrentRoomCreation(t *testing.T) {
	t.Parallel()

	memStore := store.NewMemoryStore()
	server := NewServer(memStore)

	// Create multiple rooms concurrently
	numRooms := 10
	done := make(chan string, numRooms)

	for i := 0; i < numRooms; i++ {
		go func(id int) {
			req := CreateRoomRequest{
				GameType:    "werewolf",
				DisplayName: "Host",
				MaxPlayers:  10,
			}
			bodyBytes, _ := json.Marshal(req)
			httpReq := httptest.NewRequest(http.MethodPost, "/api/rooms", bytes.NewBuffer(bodyBytes))
			rec := httptest.NewRecorder()

			server.HandleCreateRoom(rec, httpReq)

			if rec.Code == http.StatusOK {
				var resp CreateRoomResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)
				done <- resp.RoomCode
			} else {
				done <- ""
			}
		}(i)
	}

	// Collect room codes
	roomCodes := make(map[string]bool)
	for i := 0; i < numRooms; i++ {
		code := <-done
		if code != "" {
			if roomCodes[code] {
				t.Errorf("duplicate room code generated: %s", code)
			}
			roomCodes[code] = true
		}
	}

	// Verify all rooms were created successfully
	if len(roomCodes) != numRooms {
		t.Errorf("expected %d unique rooms, got %d", numRooms, len(roomCodes))
	}
}
