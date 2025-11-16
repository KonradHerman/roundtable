# Development Guide

## Quick Start

To run the application locally, you need to start **both** the backend and frontend servers:

### 1. Start the Backend Server

```bash
cd backend
go run cmd/server/main.go
```

The backend will start on **http://localhost:8080**

You should see: `Server starting on port 8080`

### 2. Start the Frontend Server

In a **new terminal window**:

```bash
cd frontend
npm install  # First time only
npm run dev
```

The frontend will start on **http://localhost:5173**

You should see: `Local: http://localhost:5173/`

### 3. Access the Application

Open your browser and go to: **http://localhost:5173**

⚠️ **Important:** Do NOT access `http://localhost:8080` directly - that's the backend API server, not the frontend!

## How It Works

- **Backend (port 8080)**: Handles API requests and WebSocket connections
- **Frontend (port 5173)**: Serves the web UI and proxies API requests to the backend
- **Vite Proxy**: Automatically forwards `/api/*` requests from port 5173 to port 8080

## Troubleshooting

### "Load failed" or "Cannot connect to backend"

**Cause:** The backend server is not running, or you're accessing the wrong port.

**Solution:**
1. Make sure the backend is running on port 8080
2. Make sure you're accessing the frontend at `http://localhost:5173` (not 8080)
3. Check that both servers are running in separate terminal windows

### Ports already in use

If port 8080 or 5173 is already in use:

**For backend:** Set a different port
```bash
PORT=3001 go run cmd/server/main.go
```

Then update `frontend/vite.config.ts`:
```typescript
target: 'http://localhost:3001'  // Update this line
```

**For frontend:** Vite will automatically use the next available port (5174, 5175, etc.)

### CORS errors

The backend has CORS enabled for development. If you still see CORS errors:
1. Check that you're accessing through the Vite dev server (port 5173)
2. The Vite proxy handles CORS automatically

## Production Build

To build for production:

```bash
# Build frontend
cd frontend
npm run build
npm start  # Runs the production server on port 3000

# In production, set the API URL:
VITE_API_URL=https://your-backend-url.com/api npm run build
```

## Environment Variables

### Frontend

Create `frontend/.env` based on `frontend/.env.example`:

```bash
# Only needed for production or if backend is on a different host
VITE_API_URL=http://localhost:8080/api
```

### Backend

```bash
PORT=8080  # Server port (default: 8080)
```

## Development Workflow

1. Start both servers (backend on 8080, frontend on 5173)
2. Access **http://localhost:5173** in your browser
3. Make changes to code - both servers support hot reload:
   - Frontend: Changes reload automatically
   - Backend: Restart with Ctrl+C and `go run cmd/server/main.go`

## Testing Multi-Player Flows

### Testing Locally with Multiple Sessions

To test game flows that require multiple players, you'll need to simulate multiple users joining the same room:

#### Option 1: Multiple Browser Windows (Quick & Easy)

1. Start both backend and frontend servers
2. Open the app in a normal browser window (`http://localhost:5173`)
3. Create a room and note the room code
4. Open **incognito/private windows** for each additional player
5. Join with the room code in each window

**Tips:**
- Use `Ctrl+Shift+N` (Chrome) or `Ctrl+Shift+P` (Firefox) for incognito windows
- Each incognito window is a separate session
- Can open 5-8 windows on a single monitor for testing
- Name windows clearly (Player 1, Player 2, etc.) using browser dev tools console

#### Option 2: Multiple Devices (Most Realistic)

1. Start backend on your computer
2. Find your local IP address:
   ```bash
   # Windows
   ipconfig
   # Look for "IPv4 Address" (usually 192.168.x.x)
   
   # Mac/Linux
   ifconfig
   # Look for "inet" under your active connection
   ```
3. Update frontend to point to your computer's IP:
   - In `frontend/vite.config.ts`, update the proxy target to your IP
   - Or set `VITE_API_URL=http://192.168.x.x:8080/api`
4. Connect phones/tablets to same WiFi network
5. Access `http://192.168.x.x:5173` from each device

**Tips:**
- Most realistic testing environment
- Tests mobile UX properly
- Can test with actual friends
- Easier to spot UI issues

#### Option 3: Different Browsers (Alternative)

Use different browsers for different players:
- Chrome
- Firefox
- Edge
- Safari (if on Mac)

Each browser has separate session storage, so you can test multiple players.

### Testing Typical Scenarios

#### 3-Player Game (Minimum)
1. Create room in Window 1 (Host)
2. Join in Windows 2 and 3 (Players)
3. Configure 6 roles (3 players + 3 center)
4. Test role reveal phase (all must acknowledge)
5. Test night phase (verify role-specific UIs appear)
6. Test day phase (discussion timer)
7. Test role reveal and play again

#### 6-8 Player Game (Typical)
Same as above but with more players. This is the sweet spot for One Night Werewolf.

**Key things to test:**
- All players see real-time updates
- Night actions work for each role
- Timer synchronizes across all clients
- Role swaps tracked correctly
- Reconnection works if someone refreshes

### Testing Edge Cases

- **Late join**: Try joining after game starts (should fail gracefully)
- **Disconnect mid-game**: Refresh one player's window during night phase
- **Host disconnect**: Close host window, see if game continues
- **Invalid role config**: Try starting with wrong number of roles
- **Spam clicking**: Click actions multiple times rapidly
- **Network delays**: Use Chrome DevTools Network throttling

## Debugging

### Check if servers are running

```bash
# Check backend
curl http://localhost:8080/health
# Should return: OK

# Check frontend
curl http://localhost:5173
# Should return: HTML content
```

### View API requests

Check the browser console - API requests are logged with the full URL.

### View backend logs

The backend logs all requests to the console where it's running.

### Debug WebSocket connections

Open browser DevTools → Network tab → WS filter to see WebSocket messages in real-time.

### View game state

The frontend game store logs events to console. Check `gameStore` in Redux DevTools or browser console.
