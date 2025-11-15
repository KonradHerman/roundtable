// API client for backend communication

const API_BASE = import.meta.env.VITE_API_URL || '/api';

export interface CreateRoomRequest {
	gameType: string;
	displayName: string;
	maxPlayers?: number;
}

export interface CreateRoomResponse {
	roomCode: string;
	sessionToken: string;
	playerId: string;
}

export interface JoinRoomRequest {
	displayName: string;
}

export interface JoinRoomResponse {
	sessionToken: string;
	playerId: string;
	roomCode: string;
}

export interface StartGameRequest {
	config: any; // Game-specific config
}

export interface RoomState {
	id: string;
	status: 'waiting' | 'playing' | 'finished';
	gameType: string;
	maxPlayers: number;
	hostId: string;
	players: Player[];
}

export interface Player {
	id: string;
	displayName: string;
	connected: boolean;
	joinedAt: string;
	lastSeenAt: string;
}

export class APIError extends Error {
	constructor(public status: number, message: string) {
		super(message);
	}
}

async function request<T>(
	endpoint: string,
	options?: RequestInit
): Promise<T> {
	const response = await fetch(`${API_BASE}${endpoint}`, {
		...options,
		headers: {
			'Content-Type': 'application/json',
			...options?.headers
		}
	});

	if (!response.ok) {
		const error = await response.text();
		throw new APIError(response.status, error || response.statusText);
	}

	return response.json();
}

export const api = {
	createRoom: (req: CreateRoomRequest) =>
		request<CreateRoomResponse>('/rooms', {
			method: 'POST',
			body: JSON.stringify(req)
		}),

	joinRoom: (roomCode: string, req: JoinRoomRequest) =>
		request<JoinRoomResponse>(`/rooms/${roomCode}/join`, {
			method: 'POST',
			body: JSON.stringify(req)
		}),

	getRoomState: (roomCode: string) =>
		request<RoomState>(`/rooms/${roomCode}`, {
			method: 'GET'
		}),

	startGame: (roomCode: string, req: StartGameRequest) =>
		request<void>(`/rooms/${roomCode}/start`, {
			method: 'POST',
			body: JSON.stringify(req)
		})
};
