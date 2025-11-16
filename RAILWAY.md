# Railway Deployment Guide

This guide explains how to deploy Cardless to Railway.

## Architecture

Cardless requires **two separate services** on Railway:

1. **Backend Service** (Go) - API server and WebSocket handler
2. **Frontend Service** (Node.js) - SvelteKit application

## Quick Setup

### 1. Create Backend Service

1. In Railway, create a new service
2. Connect your GitHub repository
3. Set the **Root Directory** to `backend`
4. Railway will auto-detect the Go application

**Environment Variables:**
- `PORT` - Automatically set by Railway (usually 8080)

### 2. Create Frontend Service

1. Create another new service in the same project
2. Connect the same GitHub repository
3. Set the **Root Directory** to `frontend`
4. Railway will auto-detect the Node.js application

**Environment Variables (REQUIRED):**

```bash
# Use Railway's internal networking for best performance
VITE_API_URL=http://cardless.railway.internal:8080/api

# Alternative: Use the backend's public URL
# VITE_API_URL=https://your-backend-service.up.railway.app/api
```

**IMPORTANT:** The `http://` prefix is required! Without it, you'll get CORS errors.

### 3. Configure Networking

**Option A: Internal Networking (Recommended)**
- Use Railway's internal network: `http://cardless.railway.internal:8080`
- Faster and doesn't use public bandwidth
- Set in frontend's `VITE_API_URL`

**Option B: Public URL**
- Use the backend's public Railway URL
- Works but slower and uses public bandwidth

## Environment Variable Reference

### Backend Service

| Variable | Value | Required | Notes |
|----------|-------|----------|-------|
| `PORT` | Auto-set by Railway | Yes | Usually 8080 |

### Frontend Service

| Variable | Value | Required | Notes |
|----------|-------|----------|-------|
| `VITE_API_URL` | `http://roundtable.railway.internal:8080/api` | Yes | Must include `http://` prefix |
| `ORIGIN` | Auto-set by Railway | No | SvelteKit origin for CSRF protection |

## Common Issues

### CORS Error: "CORS request not http"

**Problem:** The `VITE_API_URL` is missing the `http://` or `https://` prefix.

**Bad:**
```bash
VITE_API_URL=cardless.railway.internal:8080/api  ❌
```

**Good:**
```bash
VITE_API_URL=http://cardless.railway.internal:8080/api  ✅
```

**Note:** The code now automatically adds `http://` if missing, but it's better to set it correctly.

### "Load failed" or Connection Errors

1. **Check backend is running:**
   - Go to backend service logs
   - Should see: "Server starting on port 8080"

2. **Check frontend environment variables:**
   - Verify `VITE_API_URL` is set correctly
   - Check it includes the `/api` suffix

3. **Test backend health:**
   ```bash
   curl https://your-backend.up.railway.app/health
   # Should return: OK
   ```

### Backend Not Accessible

- Make sure the backend service is deployed and running
- Check that the port matches (usually 8080)
- Verify CORS is enabled (it is by default in the code)

### Frontend Can't Connect to Backend

1. Check that both services are in the same Railway project
2. Verify `VITE_API_URL` points to the correct backend URL
3. Use Railway's internal networking for better reliability

## Railway Internal Networking

Railway provides private networking between services in the same project:

**Format:** `http://<service-name>.railway.internal:<port>`

**Example:**
```bash
# If your backend service is named "cardless" or detected as "backend"
VITE_API_URL=http://cardless.railway.internal:8080/api
```

**How to find your internal URL:**
1. Go to your backend service in Railway
2. Click "Settings" → "Networking"
3. Look for "Private Networking" section
4. Use the internal domain shown there

## Deployment Checklist

- [ ] Backend service created with root directory `backend/`
- [ ] Frontend service created with root directory `frontend/`
- [ ] Backend `PORT` environment variable set (or auto-set by Railway)
- [ ] Frontend `VITE_API_URL` set to backend URL with `http://` prefix
- [ ] Both services deployed successfully
- [ ] Backend health check passes: `/health` returns "OK"
- [ ] Frontend loads without errors
- [ ] Can create a game room (tests API connection)
- [ ] WebSocket connects (check browser console)

## Build Configuration

### Backend

Railway uses Railpack, which auto-detects Go:
- Builds with `go build cmd/server/main.go`
- Uses `backend/nixpacks.toml` if present
- See `backend/railway.json` for deployment config

### Frontend

Railway auto-detects Node.js:
- Runs `npm install`
- Runs `npm run build`
- Starts with `npm start`
- Uses SvelteKit's Node adapter

## Monitoring

### Backend Logs

Check for:
```
Server starting on port 8080
```

### Frontend Logs

Check for:
```
Listening on 0.0.0.0:3000
```

### Browser Console

Should see:
```
API Request: POST http://cardless.railway.internal:8080/api/rooms
```

Should NOT see:
- CORS errors
- Connection refused
- "cardless.railway.internal:8080" (without http://)

## Performance Tips

1. **Use internal networking** - Faster than public URLs
2. **Keep services in same region** - Reduces latency
3. **Monitor service health** - Check Railway metrics
4. **Enable auto-scaling** - If needed for high traffic

## Cost Optimization

- Both services can run on Railway's free tier for development
- Use internal networking to avoid egress fees
- Monitor usage in Railway dashboard

## Troubleshooting Commands

```bash
# Test backend health
curl https://your-backend.up.railway.app/health

# Test backend API (create room)
curl -X POST https://your-backend.up.railway.app/api/rooms \
  -H "Content-Type: application/json" \
  -d '{"gameType":"werewolf","displayName":"Test","maxPlayers":10}'

# Should return JSON with roomCode, sessionToken, playerId
```

## Support

If you encounter issues:

1. Check Railway service logs for both services
2. Verify environment variables are set correctly
3. Test backend health endpoint
4. Check browser console for detailed error messages
5. Review the `DEVELOPMENT.md` for local testing

## Security Notes

- Backend uses CORS with `Access-Control-Allow-Origin: *` (development mode)
- In production, restrict CORS origins to your frontend domain
- Session tokens are sent in response body (consider httpOnly cookies for production)
- WebSocket connections are authenticated with session tokens

## Future Improvements

- [ ] Restrict CORS origins in production
- [ ] Add Redis for session storage
- [ ] Implement rate limiting
- [ ] Add health check endpoints for Railway
- [ ] Set up monitoring and alerts
