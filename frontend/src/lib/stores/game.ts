import { writable, derived } from 'svelte/store';
import type { RoomState } from '$lib/api/client';

export interface GameEvent {
	id: string;
	timestamp: string;
	type: string;
	actorId: string;
	payload: any;
}

export interface GameState {
	room: RoomState | null;
	events: GameEvent[];
	playerState: any; // Game-specific player state
	publicState: any; // Game-specific public state
}

function createGameStore() {
	const { subscribe, set, update } = writable<GameState>({
		room: null,
		events: [],
		playerState: null,
		publicState: null
	});

	return {
		subscribe,

		setRoomState: (room: RoomState) => {
			update((state) => ({ ...state, room }));
		},

		appendEvent: (event: GameEvent) => {
			update((state) => ({
				...state,
				events: [...state.events, event]
			}));

			// Process event to update game state
			// This is game-specific logic
			processEvent(event);
		},

		appendEvents: (events: GameEvent[]) => {
			update((state) => ({
				...state,
				events: [...state.events, ...events]
			}));

			// Process all events
			events.forEach(processEvent);
		},

		setPlayerState: (playerState: any) => {
			update((state) => ({ ...state, playerState }));
		},

		setPublicState: (publicState: any) => {
			update((state) => ({ ...state, publicState }));
		},

		reset: () => {
			set({
				room: null,
				events: [],
				playerState: null,
				publicState: null
			});
		}
	};
}

function processEvent(event: GameEvent) {
	// Game-specific event processing
	// This is where you'd update playerState and publicState based on events
	console.log('Processing event:', event.type, event);

	// TODO: Implement game-specific state updates
	// For now, this is a placeholder
}

export const gameStore = createGameStore();

// Derived stores for convenience
export const roomState = derived(gameStore, ($game) => $game.room);
export const playerState = derived(gameStore, ($game) => $game.playerState);
export const publicState = derived(gameStore, ($game) => $game.publicState);
export const events = derived(gameStore, ($game) => $game.events);
