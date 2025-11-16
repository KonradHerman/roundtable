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
