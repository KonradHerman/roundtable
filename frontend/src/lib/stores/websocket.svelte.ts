import { browser } from '$app/environment';

export interface ServerMessage {
	type: string;
	payload?: any;
}

export interface ClientMessage {
	type: string;
	payload?: any;
}

export type ConnectionStatus = 'disconnected' | 'connecting' | 'connected' | 'reconnecting';

class WebSocketStore {
	status = $state<ConnectionStatus>('disconnected');
	messages = $state<ServerMessage[]>([]);
	error = $state<string | null>(null);

	#ws: WebSocket | null = null;
	#reconnectAttempts = 0;
	#maxReconnectAttempts = 5;
	#reconnectTimeout: ReturnType<typeof setTimeout> | null = null;
	#roomCode: string;
	#sessionToken: string;
	#wsUrl: string;

	constructor(roomCode: string, sessionToken: string) {
		this.#roomCode = roomCode;
		this.#sessionToken = sessionToken;
		this.#wsUrl = this.#getWsUrl();

		// Auto-connect on creation
		if (browser) {
			this.connect();
		}
	}

	#getWsUrl(): string {
		if (!browser) return '';

		const apiBase = import.meta.env.VITE_API_URL;
		if (apiBase) {
			// Convert HTTP(S) URL to WS(S)
			const wsBase = apiBase.replace(/^http/, 'ws');
			return `${wsBase}/rooms/${this.#roomCode}/ws`;
		}

		// Fallback to relative URL (works with Vite proxy in dev)
		const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
		return `${protocol}//${window.location.host}/api/rooms/${this.#roomCode}/ws`;
	}

	connect() {
		if (!browser || !this.#wsUrl) return;

		this.status = 'connecting';

		this.#ws = new WebSocket(this.#wsUrl);

		this.#ws.onopen = () => {
			console.log('WebSocket connected');
			this.#reconnectAttempts = 0;

			// Send authentication message
			this.send({
				type: 'authenticate',
				payload: { sessionToken: this.#sessionToken }
			});

			this.status = 'connected';
			this.error = null;
		};

		this.#ws.onmessage = (event) => {
			try {
				const message: ServerMessage = JSON.parse(event.data);
				console.log('WebSocket message:', message);

				this.messages = [...this.messages, message];
			} catch (err) {
				console.error('Failed to parse WebSocket message:', err);
			}
		};

		this.#ws.onerror = (error) => {
			console.error('WebSocket error:', error);
			this.error = 'Connection error';
		};

		this.#ws.onclose = (event) => {
			console.log('WebSocket closed:', event.code, event.reason);
			this.#ws = null;

			if (event.code !== 1000 && this.#reconnectAttempts < this.#maxReconnectAttempts) {
				// Abnormal closure, attempt reconnect
				const delay = Math.min(1000 * Math.pow(2, this.#reconnectAttempts), 10000);
				this.#reconnectAttempts++;

				this.status = 'reconnecting';

				console.log(
					`Reconnecting in ${delay}ms (attempt ${this.#reconnectAttempts}/${this.#maxReconnectAttempts})`
				);

				this.#reconnectTimeout = setTimeout(() => {
					this.connect();
				}, delay);
			} else {
				this.status = 'disconnected';
			}
		};
	}

	disconnect() {
		if (this.#reconnectTimeout) {
			clearTimeout(this.#reconnectTimeout);
			this.#reconnectTimeout = null;
		}

		if (this.#ws) {
			this.#ws.close(1000, 'Client disconnect');
			this.#ws = null;
		}

		this.status = 'disconnected';
	}

	send(message: ClientMessage) {
		if (this.#ws && this.#ws.readyState === WebSocket.OPEN) {
			this.#ws.send(JSON.stringify(message));
		} else {
			console.error('WebSocket not connected');
		}
	}

	sendAction(action: any) {
		this.send({
			type: 'action',
			payload: { action }
		});
	}

	reconnect() {
		this.connect();
	}
}

// Export factory function
export function createWebSocket(roomCode: string, sessionToken: string) {
	return new WebSocketStore(roomCode, sessionToken);
}

