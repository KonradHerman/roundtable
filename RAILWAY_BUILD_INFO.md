# Railway Build Configuration

This project uses **Railpack** (not Docker).

## Build Settings Required

### Backend Service
- **Root Directory**: `backend`
- **Builder**: Railpack (auto-detected)
- **Build Command**: (leave empty - auto-detected)
- **Start Command**: (leave empty - uses ./server)

### Frontend Service
- **Root Directory**: `frontend`
- **Builder**: Railpack (auto-detected)
- **Build Command**: (leave empty - auto-detected)
- **Start Command**: (leave empty - uses node build)

## If Railway Still Uses Docker

1. **Delete the service** completely from Railway
2. **Create a NEW service** from GitHub
3. Set the **Root Directory** FIRST (before any build)
4. Railway will auto-detect Railpack

OR

1. Go to Service Settings
2. Scroll to Build section
3. Click "Clear Build Cache"
4. Redeploy

## Configuration Files

- `backend/railway.json` - Forces Railpack builder
- `frontend/railway.json` - Forces Railpack builder
- `backend/nixpacks.toml` - Build configuration
- `frontend/nixpacks.toml` - Build configuration

NO Dockerfiles should exist in this repo!
