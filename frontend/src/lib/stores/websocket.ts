import { writable, derived, get } from 'svelte/store';
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

interface WebSocketStore {
	status: ConnectionStatus;
	messages: ServerMessage[];
	error: string | null;
}

function createWebSocketStore(roomCode: string, sessionToken: string) {
	const { subscribe, set, update } = writable<WebSocketStore>({
		status: 'disconnected',
		messages: [],
		error: null
	});

	let ws: WebSocket | null = null;
	let reconnectAttempts = 0;
	const maxReconnectAttempts = 5;
	let reconnectTimeout: ReturnType<typeof setTimeout> | null = null;

	const wsUrl = browser
		? `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/rooms/${roomCode}/ws`
		: '';

	function connect() {
		if (!browser || !wsUrl) return;

		update((state) => ({ ...state, status: 'connecting' }));

		ws = new WebSocket(wsUrl);

		ws.onopen = () => {
			console.log('WebSocket connected');
			reconnectAttempts = 0;

			// Send authentication message
			send({
				type: 'authenticate',
				payload: { sessionToken }
			});

			update((state) => ({ ...state, status: 'connected', error: null }));
		};

		ws.onmessage = (event) => {
			try {
				const message: ServerMessage = JSON.parse(event.data);
				console.log('WebSocket message:', message);

				update((state) => ({
					...state,
					messages: [...state.messages, message]
				}));
			} catch (err) {
				console.error('Failed to parse WebSocket message:', err);
			}
		};

		ws.onerror = (error) => {
			console.error('WebSocket error:', error);
			update((state) => ({ ...state, error: 'Connection error' }));
		};

		ws.onclose = (event) => {
			console.log('WebSocket closed:', event.code, event.reason);
			ws = null;

			if (event.code !== 1000 && reconnectAttempts < maxReconnectAttempts) {
				// Abnormal closure, attempt reconnect
				const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 10000);
				reconnectAttempts++;

				update((state) => ({ ...state, status: 'reconnecting' }));

				console.log(`Reconnecting in ${delay}ms (attempt ${reconnectAttempts}/${maxReconnectAttempts})`);

				reconnectTimeout = setTimeout(() => {
					connect();
				}, delay);
			} else {
				update((state) => ({ ...state, status: 'disconnected' }));
			}
		};
	}

	function disconnect() {
		if (reconnectTimeout) {
			clearTimeout(reconnectTimeout);
			reconnectTimeout = null;
		}

		if (ws) {
			ws.close(1000, 'Client disconnect');
			ws = null;
		}

		update((state) => ({ ...state, status: 'disconnected' }));
	}

	function send(message: ClientMessage) {
		if (ws && ws.readyState === WebSocket.OPEN) {
			ws.send(JSON.stringify(message));
		} else {
			console.error('WebSocket not connected');
		}
	}

	function sendAction(action: any) {
		send({
			type: 'action',
			payload: { action }
		});
	}

	// Auto-connect on creation
	if (browser) {
		connect();
	}

	return {
		subscribe,
		send,
		sendAction,
		disconnect,
		reconnect: connect
	};
}

// Export factory function
export function createWebSocket(roomCode: string, sessionToken: string) {
	return createWebSocketStore(roomCode, sessionToken);
}

// Derived store for connection status
export function createStatusStore(wsStore: ReturnType<typeof createWebSocketStore>) {
	return derived(wsStore, ($ws) => $ws.status);
}

// Derived store for latest message
export function createLatestMessageStore(wsStore: ReturnType<typeof createWebSocketStore>) {
	return derived(wsStore, ($ws) => $ws.messages[$ws.messages.length - 1]);
}
