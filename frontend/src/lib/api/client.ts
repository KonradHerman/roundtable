// API client for backend communication

// Normalize API URL - ensure it has a protocol if it's an absolute URL
function normalizeApiUrl(url: string): string {
	if (!url) return '/api';

	// If it's a relative URL (starts with /), use it as-is
	if (url.startsWith('/')) return url;

	// If it already has a protocol, use it as-is
	if (url.startsWith('http://') || url.startsWith('https://')) return url;

	// If it's missing a protocol, add http:// (Railway internal URLs use http)
	return `http://${url}`;
}

const API_BASE = normalizeApiUrl(import.meta.env.VITE_API_URL || '/api');

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
	try {
		const url = `${API_BASE}${endpoint}`;
		console.log(`API Request: ${options?.method || 'GET'} ${url}`);

		const response = await fetch(url, {
			...options,
			headers: {
				'Content-Type': 'application/json',
				...options?.headers
			}
		});

		if (!response.ok) {
			const contentType = response.headers.get('content-type');
			let error: string;

			// If response is HTML (likely an error page), provide a better message
			if (contentType && contentType.includes('text/html')) {
				error = `Backend connection failed. Make sure the backend server is running on ${API_BASE}`;
			} else {
				error = await response.text();
			}

			throw new APIError(response.status, error || response.statusText);
		}

		return response.json();
	} catch (err) {
		if (err instanceof APIError) {
			throw err;
		}
		// Network errors or other fetch failures
		if (err instanceof TypeError && err.message.includes('fetch')) {
			throw new APIError(0, `Cannot connect to backend server at ${API_BASE}. Please ensure:\n1. Backend is running (go run cmd/server/main.go in backend/)\n2. Backend is accessible at ${API_BASE}\n3. You're accessing the frontend at http://localhost:5173`);
		}
		throw new APIError(0, `Request failed: ${err instanceof Error ? err.message : String(err)}`);
	}
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
