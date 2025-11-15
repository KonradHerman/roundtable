import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export interface Session {
	playerId: string;
	sessionToken: string;
	roomCode: string;
	displayName?: string;
}

const SESSION_KEY = 'roundtable_session';

// Load session from localStorage
function loadSession(): Session | null {
	if (!browser) return null;

	const stored = localStorage.getItem(SESSION_KEY);
	if (!stored) return null;

	try {
		return JSON.parse(stored);
	} catch {
		return null;
	}
}

// Create the session store
function createSessionStore() {
	const { subscribe, set, update } = writable<Session | null>(loadSession());

	return {
		subscribe,
		set: (session: Session | null) => {
			if (browser) {
				if (session) {
					localStorage.setItem(SESSION_KEY, JSON.stringify(session));
				} else {
					localStorage.removeItem(SESSION_KEY);
				}
			}
			set(session);
		},
		update: (updater: (session: Session | null) => Session | null) => {
			update((current) => {
				const updated = updater(current);
				if (browser) {
					if (updated) {
						localStorage.setItem(SESSION_KEY, JSON.stringify(updated));
					} else {
						localStorage.removeItem(SESSION_KEY);
					}
				}
				return updated;
			});
		},
		clear: () => {
			if (browser) {
				localStorage.removeItem(SESSION_KEY);
			}
			set(null);
		}
	};
}

export const session = createSessionStore();
