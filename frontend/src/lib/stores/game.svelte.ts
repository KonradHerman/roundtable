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

class GameStore {
	room = $state<RoomState | null>(null);
	events = $state<GameEvent[]>([]);
	playerState = $state<any>(null);
	publicState = $state<any>(null);

	setRoomState(room: RoomState) {
		this.room = room;
	}

	appendEvent(event: GameEvent) {
		this.events = [...this.events, event];
		// Process event to update game state
		this.processEvent(event);
	}

	appendEvents(events: GameEvent[]) {
		this.events = [...this.events, ...events];
		// Process all events
		events.forEach(event => this.processEvent(event));
	}

	setPlayerState(playerState: any) {
		this.playerState = playerState;
	}

	setPublicState(publicState: any) {
		this.publicState = publicState;
	}

	reset() {
		this.room = null;
		this.events = [];
		this.playerState = null;
		this.publicState = null;
	}

	private processEvent(event: GameEvent) {
		// Game-specific event processing
		// This is where you'd update playerState and publicState based on events
		console.log('Processing event:', event.type, event);

		// TODO: Implement game-specific state updates
		// For now, this is a placeholder
	}
}

export const gameStore = new GameStore();

