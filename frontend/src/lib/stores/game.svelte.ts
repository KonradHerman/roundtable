import type { RoomState } from '$lib/api/client';

export interface GameEvent {
	id: string;
	timestamp: string;
	type: string;
	actorId: string;
	payload: any;
}

export class GameStore {
	room = $state<RoomState | null>(null);
	events = $state<GameEvent[]>([]);
	playerState = $state<any>(null); // Game-specific player state
	publicState = $state<any>(null); // Game-specific public state

	setRoomState(room: RoomState) {
		this.room = room;
	}

	appendEvent(event: GameEvent) {
		this.events = [...this.events, event];
		this.processEvent(event);
	}

	appendEvents(events: GameEvent[]) {
		this.events = [...this.events, ...events];
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

