import { describe, it, expect, vi, beforeEach } from 'vitest';
import { api } from './client';

describe('API Client', () => {
	beforeEach(() => {
		// Reset fetch mock before each test
		vi.restoreAllMocks();
	});

	describe('createRoom', () => {
		it('should create a room with the correct payload', async () => {
			const mockResponse = {
				roomCode: 'ABC123',
				sessionToken: 'token-123',
				playerId: 'player-123'
			};

			global.fetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => mockResponse
			});

			const result = await api.createRoom({
				gameType: 'werewolf',
				displayName: 'TestPlayer',
				maxPlayers: 10
			});

			expect(global.fetch).toHaveBeenCalledWith(
				'/api/rooms',
				expect.objectContaining({
					method: 'POST',
					headers: expect.objectContaining({
						'Content-Type': 'application/json'
					}),
					body: JSON.stringify({
						gameType: 'werewolf',
						displayName: 'TestPlayer',
						maxPlayers: 10
					})
				})
			);

			expect(result).toEqual(mockResponse);
		});

		it('should throw error on failed request', async () => {
			global.fetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 400,
				statusText: 'Bad Request',
				text: async () => 'Invalid game type'
			});

			await expect(
				api.createRoom({
					gameType: 'invalid',
					displayName: 'TestPlayer',
					maxPlayers: 10
				})
			).rejects.toThrow();
		});
	});

	describe('joinRoom', () => {
		it('should join a room with correct parameters', async () => {
			const mockResponse = {
				sessionToken: 'token-456',
				playerId: 'player-456',
				roomCode: 'ABC123'
			};

			global.fetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => mockResponse
			});

			const result = await api.joinRoom('ABC123', { displayName: 'Player2' });

			expect(global.fetch).toHaveBeenCalledWith(
				'/api/rooms/ABC123/join',
				expect.objectContaining({
					method: 'POST',
					headers: expect.objectContaining({
						'Content-Type': 'application/json'
					}),
					body: JSON.stringify({
						displayName: 'Player2'
					})
				})
			);

			expect(result).toEqual(mockResponse);
		});

		it('should throw error when room not found', async () => {
			global.fetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 404,
				statusText: 'Not Found',
				text: async () => 'Room not found'
			});

			await expect(api.joinRoom('NOROOM', { displayName: 'Player2' })).rejects.toThrow();
		});
	});

	describe('getRoomState', () => {
		it('should get room details', async () => {
			const mockRoom = {
				id: 'ABC123',
				gameType: 'werewolf',
				status: 'waiting',
				players: []
			};

			global.fetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => mockRoom
			});

			const result = await api.getRoomState('ABC123');

			expect(global.fetch).toHaveBeenCalledWith(
				'/api/rooms/ABC123',
				expect.objectContaining({
					method: 'GET'
				})
			);

			expect(result).toEqual(mockRoom);
		});
	});

	describe('startGame', () => {
		it('should start a game with config', async () => {
			const config = {
				roles: ['werewolf', 'seer', 'villager']
			};

			global.fetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ success: true })
			});

			await api.startGame('ABC123', { config });

			expect(global.fetch).toHaveBeenCalledWith(
				'/api/rooms/ABC123/start',
				expect.objectContaining({
					method: 'POST',
					headers: expect.objectContaining({
						'Content-Type': 'application/json'
					}),
					body: JSON.stringify({ config })
				})
			);
		});
	});

	describe('error handling', () => {
		it('should handle network errors', async () => {
			global.fetch = vi.fn().mockRejectedValue(new Error('Network error'));

			await expect(
				api.createRoom({
					gameType: 'werewolf',
					displayName: 'TestPlayer',
					maxPlayers: 10
				})
			).rejects.toThrow();
		});

		it('should handle JSON parse errors gracefully', async () => {
			global.fetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => {
					throw new Error('Invalid JSON');
				}
			});

			await expect(
				api.createRoom({
					gameType: 'werewolf',
					displayName: 'TestPlayer',
					maxPlayers: 10
				})
			).rejects.toThrow('Invalid JSON');
		});
	});
});
