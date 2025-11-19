// API client for backend communication

// Normalize API URL - ensure it has a protocol if it's an absolute URL
function normalizeApiUrl(url: string): string {
	if (!url) return '/api';

	// If it's a relative URL (starts with /), use it as-is
	if (url.startsWith('/')) return url;

	// If it already has a protocol, use it as-is
	if (url.startsWith('http://') || url.startsWith('https://')) return url;

	// If it's missing a protocol, detect based on current page protocol
	// In production (HTTPS), use HTTPS for backend too to avoid mixed content errors
	// In development (HTTP), use HTTP
	const protocol = typeof window !== 'undefined' && window.location.protocol === 'https:'
		? 'https://'
		: 'http://';

	return `${protocol}${url}`;
}

const API_BASE = normalizeApiUrl(import.meta.env.VITE_API_URL || '/api');

// Log API configuration on startup (helps debugging Railway deployments)
if (typeof window !== 'undefined') {
	console.log('[API Client] Configuration:', {
		VITE_API_URL: import.meta.env.VITE_API_URL || '(not set)',
		API_BASE,
		mode: import.meta.env.MODE,
		isDev: import.meta.env.DEV
	});
}

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
		const errorMessage = err instanceof Error ? err.message : String(err);

		// Detect if we're in production (Railway) without VITE_API_URL set
		const isProduction = import.meta.env.PROD;
		const hasApiUrl = !!import.meta.env.VITE_API_URL;

		if (!hasApiUrl && isProduction && API_BASE === '/api') {
			throw new APIError(0,
				`Backend configuration error: VITE_API_URL environment variable is not set.\n\n` +
				`In Railway, you need to:\n` +
				`1. Go to your frontend service settings\n` +
				`2. Add environment variable: VITE_API_URL\n` +
				`3. Set it to your backend URL (e.g., https://your-backend.up.railway.app/api)\n` +
				`4. Redeploy the frontend\n\n` +
				`Original error: ${errorMessage}`
			);
		}

		if (err instanceof TypeError && errorMessage.includes('fetch')) {
			throw new APIError(0,
				`Cannot connect to backend server at ${API_BASE}.\n\n` +
				`Development: Ensure backend is running (go run cmd/server/main.go in backend/)\n` +
				`Production: Verify VITE_API_URL points to your backend service\n\n` +
				`Error: ${errorMessage}`
			);
		}

		throw new APIError(0, `Request failed: ${errorMessage}`);
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
		}),

	resetGame: (roomCode: string) =>
		request<void>(`/rooms/${roomCode}/reset`, {
			method: 'POST'
		})
};
