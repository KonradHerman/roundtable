# Testing Guide for Cardless

This document provides comprehensive guidance on running and writing tests for the Cardless project.

## Table of Contents

1. [Overview](#overview)
2. [Backend Testing (Go)](#backend-testing-go)
3. [Frontend Testing (Svelte 5 + Vitest)](#frontend-testing-svelte-5--vitest)
4. [Running All Tests](#running-all-tests)
5. [Writing New Tests](#writing-new-tests)
6. [Test Coverage](#test-coverage)
7. [CI/CD Integration](#cicd-integration)
8. [Best Practices](#best-practices)

---

## Overview

The Cardless project uses modern testing frameworks following industry best practices:

- **Backend**: Go's built-in `testing` package with table-driven tests
- **Frontend**: Vitest + Testing Library for Svelte 5
- **Philosophy**: Test critical paths, maintain high coverage for core logic

### Test Structure

```
backend/
├── internal/
│   ├── core/
│   │   ├── room.go
│   │   └── room_test.go          # Room management tests
│   ├── games/werewolf/
│   │   ├── game.go
│   │   └── game_test.go          # Werewolf game logic tests
│   └── server/
│       ├── handlers.go
│       └── handlers_test.go      # API handler tests

frontend/
├── src/
│   ├── lib/api/
│   │   ├── client.ts
│   │   └── client.test.ts        # API client tests
│   └── test/
│       └── setup.ts               # Test configuration
├── vitest.config.ts
└── package.json
```

---

## Backend Testing (Go)

### Running Backend Tests

**Run all tests:**
```bash
cd backend
go test ./...
```

**Run with verbose output:**
```bash
go test ./... -v
```

**Run specific package:**
```bash
go test ./internal/core -v
go test ./internal/games/werewolf -v
go test ./internal/server -v
```

**Run with coverage:**
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out  # View coverage in browser
```

**Run specific test:**
```bash
go test ./internal/core -run TestRoom_AddPlayer -v
```

### Test Structure

Go tests follow the table-driven test pattern:

```go
func TestRoom_AddPlayer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupRoom   func() *Room
		playerToAdd *Player
		wantErr     bool
		errContains string
	}{
		{
			name: "successfully add player to waiting room",
			setupRoom: func() *Room {
				host := &Player{ID: "host", DisplayName: "Host"}
				return NewRoom("ABC123", "werewolf", host, 10)
			},
			playerToAdd: &Player{ID: "player1", DisplayName: "Player1"},
			wantErr:     false,
		},
		// More test cases...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Test implementation...
		})
	}
}
```

### What's Tested (Backend)

**✅ Core Package (`internal/core`)**
- Room creation and management
- Player operations (add, remove, get)
- Concurrent access safety
- Token-based player lookup

**✅ Werewolf Game (`internal/games/werewolf`)**
- Game initialization and role assignment
- Config validation
- Role team classifications
- Phase management
- State retrieval (player and public views)
- Role uniqueness across games

**✅ Server Handlers (`internal/server`)**
- Room creation endpoint
- Join room endpoint
- Get room endpoint
- Concurrent room creation
- Error handling for all endpoints

### Current Coverage

- **core**: 17.5%
- **werewolf**: 15.4%
- **server**: 29.5%

---

## Frontend Testing (Svelte 5 + Vitest)

### Running Frontend Tests

**Run tests in watch mode:**
```bash
cd frontend
npm test
```

**Run once (CI mode):**
```bash
npm run test:run
```

**Run with UI:**
```bash
npm run test:ui
# Opens interactive UI at http://localhost:51204
```

**Run with coverage:**
```bash
npm run test:coverage
```

### Test Configuration

**`vitest.config.ts`:**
```typescript
import { defineConfig } from 'vitest/config';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import { svelteTesting } from '@testing-library/svelte/vite';

export default defineConfig({
	plugins: [svelte(), svelteTesting()],
	test: {
		globals: true,
		environment: 'jsdom',
		setupFiles: ['./src/test/setup.ts'],
		include: ['src/**/*.{test,spec}.{js,ts}']
	}
});
```

### Test Structure

Frontend tests use Vitest with Testing Library:

```typescript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { api } from './client';

describe('API Client', () => {
	beforeEach(() => {
		vi.restoreAllMocks();
	});

	it('should create a room with the correct payload', async () => {
		const mockResponse = { roomCode: 'ABC123', ... };

		global.fetch = vi.fn().mockResolvedValue({
			ok: true,
			json: async () => mockResponse
		});

		const result = await api.createRoom({...});

		expect(result).toEqual(mockResponse);
		expect(global.fetch).toHaveBeenCalledWith(...);
	});
});
```

### What's Tested (Frontend)

**✅ API Client (`src/lib/api/client.ts`)**
- `createRoom()` - Room creation with validation
- `joinRoom()` - Joining existing rooms
- `getRoomState()` - Fetching room details
- `startGame()` - Starting games with config
- Error handling (network errors, JSON parse errors)
- HTTP status code handling (404, 400, etc.)

**📋 Future Tests (Recommended)**
- Component tests for critical UI flows
- WebSocket store tests
- Game state management tests

---

## Running All Tests

### Quick Verification

```bash
# From project root
cd backend && go test ./... && cd ../frontend && npm run test:run
```

### With Coverage

```bash
# Backend
cd backend
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Frontend
cd ../frontend
npm run test:coverage
```

---

## Writing New Tests

### Backend Tests

**1. Create test file:**
```bash
# Example: testing a new game feature
touch backend/internal/games/werewolf/actions_test.go
```

**2. Follow the pattern:**
```go
package werewolf

import "testing"

func TestNewFeature(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		input SomeInput
		want SomeOutput
		wantErr bool
	}{
		// Test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Test implementation
		})
	}
}
```

**3. Test considerations:**
- Use `t.Parallel()` for independent tests
- Test both success and error cases
- Use descriptive test names
- Mock external dependencies
- Verify thread safety for concurrent code

### Frontend Tests

**1. Create test file:**
```bash
# Place test next to the file being tested
touch frontend/src/lib/stores/game.test.ts
```

**2. Follow the pattern:**
```typescript
import { describe, it, expect, beforeEach, vi } from 'vitest';

describe('Feature Name', () => {
	beforeEach(() => {
		// Reset state before each test
	});

	it('should do something expected', () => {
		// Arrange
		// Act
		// Assert
	});
});
```

**3. Test considerations:**
- Mock `fetch` for API calls
- Use Testing Library for components
- Test user interactions, not implementation details
- Clean up after each test with `afterEach(cleanup)`

---

## Test Coverage

### Viewing Coverage

**Backend:**
```bash
cd backend
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Frontend:**
```bash
cd frontend
npm run test:coverage
# Coverage report: frontend/coverage/index.html
```

### Coverage Goals

- **Critical paths**: 80%+ (room management, game logic, API handlers)
- **UI components**: 60%+ (focus on user flows)
- **Utilities**: 90%+ (pure functions)

### What NOT to Test

- Third-party libraries
- Generated code
- Trivial getters/setters
- Configuration files

---

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      - name: Run backend tests
        working-directory: ./backend
        run: |
          go test ./... -v -coverprofile=coverage.out
          go tool cover -func=coverage.out

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install dependencies
        working-directory: ./frontend
        run: npm ci
      - name: Run frontend tests
        working-directory: ./frontend
        run: npm run test:run
```

### Pre-commit Hook (Optional)

```bash
#!/bin/sh
# .git/hooks/pre-commit

cd backend && go test ./... || exit 1
cd ../frontend && npm run test:run || exit 1
```

---

## Best Practices

### General Principles

1. **Test Behavior, Not Implementation**
   - Focus on what the code does, not how it does it
   - Test inputs and outputs, not internal state

2. **Keep Tests Simple and Focused**
   - One assertion per test (when reasonable)
   - Clear arrange-act-assert structure

3. **Use Descriptive Names**
   ```go
   // Good
   TestRoom_AddPlayer_FailsWhenRoomIsFull

   // Bad
   TestAddPlayer
   ```

4. **Maintain Test Independence**
   - Tests should not depend on execution order
   - Each test should set up its own state
   - Use `t.Parallel()` in Go

5. **Mock External Dependencies**
   - Database calls
   - API requests
   - WebSocket connections
   - Time-dependent operations

### Backend-Specific

- Use table-driven tests for multiple scenarios
- Test concurrent access with goroutines
- Verify error messages contain helpful context
- Test edge cases (empty inputs, boundary values)

### Frontend-Specific

- Mock `fetch` for API tests
- Use `vi.restoreAllMocks()` in `beforeEach`
- Test error states and loading states
- Prefer integration tests over unit tests for components

---

## Troubleshooting

### Common Issues

**Backend: "no test files"**
- Ensure files end with `_test.go`
- Run from correct directory

**Backend: Race detector warnings**
```bash
go test ./... -race
```

**Frontend: "Cannot find module"**
- Run `npx svelte-kit sync` first
- Check path aliases in `vitest.config.ts`

**Frontend: Tests timeout**
- Check for unmocked async operations
- Increase timeout in test config

### Getting Help

- Check existing tests for patterns
- Read Go testing docs: https://go.dev/doc/tutorial/add-a-test
- Read Vitest docs: https://vitest.dev
- Read Testing Library docs: https://testing-library.com

---

## Summary

- **Backend**: `go test ./...` (table-driven tests)
- **Frontend**: `npm run test:run` (Vitest + Testing Library)
- **Coverage**: Track with `-coverprofile` (Go) and `--coverage` (Vitest)
- **CI/CD**: Integrate both test suites in GitHub Actions
- **Best Practice**: Test critical paths, maintain independence, use descriptive names

**Current Status**: ✅ Basic test infrastructure complete with good coverage of core functionality.

---

*Last Updated: 2025-11-21*
