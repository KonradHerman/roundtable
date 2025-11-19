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

class SessionStore {
	value = $state<Session | null>(loadSession());

	set(session: Session | null) {
		this.value = session;
		if (browser) {
			if (session) {
				localStorage.setItem(SESSION_KEY, JSON.stringify(session));
			} else {
				localStorage.removeItem(SESSION_KEY);
			}
		}
	}

	update(updater: (session: Session | null) => Session | null) {
		const updated = updater(this.value);
		this.set(updated);
	}

	clear() {
		this.set(null);
	}
}

export const session = new SessionStore();

