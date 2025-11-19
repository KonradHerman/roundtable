package server

import (
	"encoding/json"

	"github.com/KonradHerman/roundtable/internal/core"
)

// ClientMessage represents messages sent from client to server.
type ClientMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// Client message types
const (
	ClientMsgAuthenticate = "authenticate"
	ClientMsgAction       = "action"
	ClientMsgPing         = "ping"
)

// AuthenticatePayload is sent when a client connects or reconnects.
type AuthenticatePayload struct {
	SessionToken string `json:"sessionToken"`
}

// ActionPayload wraps a game action.
type ActionPayload struct {
	Action core.Action `json:"action"`
}

// ServerMessage represents messages sent from server to client.
type ServerMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// Server message types
const (
	ServerMsgAuthenticated = "authenticated"
	ServerMsgRoomState     = "room_state"
	ServerMsgEvent         = "event"
	ServerMsgEvents        = "events" // Batch for reconnection
	ServerMsgError         = "error"
	ServerMsgPong          = "pong"
)

// AuthenticatedPayload confirms successful authentication.
type AuthenticatedPayload struct {
	PlayerID  string         `json:"playerId"`
	RoomState core.RoomState `json:"roomState"`
}

// RoomStatePayload contains current room state.
type RoomStatePayload struct {
	RoomState core.RoomState `json:"roomState"`
}

// EventPayload contains a single game event.
type EventPayload struct {
	Event core.GameEvent `json:"event"`
}

// EventsPayload contains multiple events (for reconnection).
type EventsPayload struct {
	Events []core.GameEvent `json:"events"`
}

// ErrorPayload contains error information.
type ErrorPayload struct {
	Message string `json:"message"`
}

// Helper functions to create server messages

func NewServerMessage(msgType string, payload interface{}) (ServerMessage, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return ServerMessage{}, err
	}

	return ServerMessage{
		Type:    msgType,
		Payload: payloadBytes,
	}, nil
}

func NewAuthenticatedMessage(playerID string, roomState core.RoomState) (ServerMessage, error) {
	return NewServerMessage(ServerMsgAuthenticated, AuthenticatedPayload{
		PlayerID:  playerID,
		RoomState: roomState,
	})
}

func NewRoomStateMessage(roomState core.RoomState) (ServerMessage, error) {
	return NewServerMessage(ServerMsgRoomState, RoomStatePayload{
		RoomState: roomState,
	})
}

func NewEventMessage(event core.GameEvent) (ServerMessage, error) {
	return NewServerMessage(ServerMsgEvent, EventPayload{
		Event: event,
	})
}

func NewEventsMessage(events []core.GameEvent) (ServerMessage, error) {
	return NewServerMessage(ServerMsgEvents, EventsPayload{
		Events: events,
	})
}

func NewErrorMessage(errMsg string) (ServerMessage, error) {
	return NewServerMessage(ServerMsgError, ErrorPayload{
		Message: errMsg,
	})
}

func NewPongMessage() (ServerMessage, error) {
	return ServerMessage{Type: ServerMsgPong}, nil
}
