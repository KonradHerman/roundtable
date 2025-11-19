# Railway Deployment Guide

## Environment Variables

### Backend Service

The backend requires the following environment variables to be set in Railway:

#### Required

- **`ALLOWED_ORIGIN`**: Comma-separated list of allowed frontend origins for CORS
  - Format: `https://your-frontend.up.railway.app,https://another-domain.com`
  - Example: `https://roundtable-frontend-production.up.railway.app`
  - **IMPORTANT**: Set this to your frontend's Railway URL
  - For multiple origins, separate with commas: `https://domain1.com,https://domain2.com`
  - For development: `http://localhost:5173`
  - For wildcard (not recommended with credentials): `*`

#### Optional

- **`PORT`**: Port to run the server on (default: `8080`)
  - Railway automatically sets this
- **`CORS_DEBUG`**: Enable CORS debugging logs (default: `false`)
  - Set to `true` for debugging CORS issues

### Frontend Service

The frontend requires the following environment variables:

#### Required

- **`VITE_API_URL`**: Backend API URL
  - Format: `https://your-backend.up.railway.app/api`
  - Railway template: `https://${{roundtable-backend.RAILWAY_PUBLIC_DOMAIN}}/api`
  - Example: `https://roundtable-backend-production-8923.up.railway.app/api`

## Setting Environment Variables in Railway

### Via Railway Dashboard

1. Go to your Railway project
2. Select the service (Backend or Frontend)
3. Go to the "Variables" tab
4. Click "New Variable"
5. Add the variable name and value
6. Click "Add"
7. The service will automatically redeploy

### Via Railway CLI

```bash
# Set backend ALLOWED_ORIGIN
railway variables --set ALLOWED_ORIGIN=https://your-frontend.up.railway.app

# Set frontend API URL
railway variables --set VITE_API_URL=https://your-backend.up.railway.app/api
```

## CORS Configuration

The backend now supports:
- ✅ Multiple comma-separated origins
- ✅ Automatic whitespace trimming
- ✅ Development default (`http://localhost:5173`)
- ✅ Wildcard support (not recommended for production with credentials)
- ✅ Debug logging with `CORS_DEBUG=true`

### Example CORS Configurations

#### Single Origin (Production)
```bash
ALLOWED_ORIGIN=https://roundtable-frontend-production.up.railway.app
```

#### Multiple Origins
```bash
ALLOWED_ORIGIN=https://roundtable.com,https://www.roundtable.com,https://app.roundtable.com
```

#### Development + Production
```bash
ALLOWED_ORIGIN=http://localhost:5173,https://roundtable-frontend-production.up.railway.app
```

#### Wildcard (Development Only)
```bash
ALLOWED_ORIGIN=*
```
⚠️ Note: Using `*` disables credentials (cookies/auth headers)

## Troubleshooting CORS Issues

### Issue: "CORS header 'Access-Control-Allow-Origin' missing"

**Cause**: The frontend origin is not in the backend's `ALLOWED_ORIGIN` list.

**Solution**:
1. Check your frontend URL (the domain making the request)
2. Add it to the backend's `ALLOWED_ORIGIN` environment variable
3. Redeploy the backend service

**Example**:
```bash
# If your frontend is at: https://roundtable-frontend-xyz.up.railway.app
# Set backend ALLOWED_ORIGIN to:
ALLOWED_ORIGIN=https://roundtable-frontend-xyz.up.railway.app
```

### Issue: "CORS request did not succeed"

**Cause**: Backend is not reachable or not responding.

**Solution**:
1. Check backend deployment logs in Railway
2. Verify backend service is running
3. Test backend health endpoint: `https://your-backend.up.railway.app/health`

### Debugging CORS

Enable CORS debug logging:
```bash
CORS_DEBUG=true
```

Check backend logs in Railway to see:
- Configured allowed origins
- Incoming request origins
- CORS validation results

## Deployment Checklist

- [ ] Backend `ALLOWED_ORIGIN` is set to frontend URL
- [ ] Frontend `VITE_API_URL` is set to backend URL
- [ ] Both services are deployed and running
- [ ] Health check passes: `curl https://your-backend.up.railway.app/health`
- [ ] CORS test: Try creating a room from the frontend
- [ ] Check browser console for any CORS errors
- [ ] Check Railway logs for any errors

## Getting Your Railway URLs

### Frontend URL
1. Go to Railway dashboard
2. Select frontend service
3. Go to "Settings" tab
4. Find "Public URL" or "Domains" section
5. Copy the URL (e.g., `https://roundtable-frontend-xyz.up.railway.app`)

### Backend URL
1. Go to Railway dashboard
2. Select backend service
3. Go to "Settings" tab
4. Find "Public URL" or "Domains" section
5. Copy the URL (e.g., `https://roundtable-backend-xyz.up.railway.app`)
6. Add `/api` to the end for frontend's `VITE_API_URL`

## Quick Setup Script

After getting both URLs from Railway:

```bash
# Set backend CORS (replace with your actual frontend URL)
railway link roundtable-backend
railway variables --set ALLOWED_ORIGIN=https://roundtable-frontend-xyz.up.railway.app

# Set frontend API URL (replace with your actual backend URL)
railway link roundtable-frontend
railway variables --set VITE_API_URL=https://roundtable-backend-xyz.up.railway.app/api
```
