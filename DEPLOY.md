# 🚀 Quick Deployment Reference

## Backend on Render

### 1. Push to GitHub
```bash
git add .
git commit -m "Ready for deployment"
git push origin main
```

### 2. On Render (https://dashboard.render.com)
- New → Web Service
- Connect repository → Select `safeware`
- Configure:
  - Root Directory: `backend`
  - Build: `go build -o bin/server ./cmd/server`
  - Start: `./bin/server`

### 3. Environment Variables (on Render)
```
MONGODB_URI=your-mongodb-atlas-uri
JWT_SECRET=your-jwt-secret
JWT_REFRESH_SECRET=your-refresh-secret
GIN_MODE=release
ALLOWED_ORIGIN=https://your-app.vercel.app
PORT=8080
```

### 4. Save Backend URL
```
https://your-backend-name.onrender.com
```

---

## Frontend on Vercel

### 1. Update Frontend ENV
Create `frontend/.env.production`:
```
VITE_API_URL=https://your-backend-name.onrender.com/api/v1
```

### 2. Commit and Push
```bash
git add frontend/.env.production
git commit -m "Add production config"
git push origin main
```

### 3. On Vercel (https://vercel.com/dashboard)
- New Project → Import `safeware`
- Configure:
  - Root Directory: `frontend`
  - Framework: Vite
  - Build: `npm run build`
  - Output: `dist`

### 4. Environment Variables (on Vercel)
```
VITE_API_URL=https://your-backend-name.onrender.com/api/v1
```

### 5. Deploy!

---

## Post-Deployment

### Update CORS on Render
Add to environment variables:
```
ALLOWED_ORIGIN=https://safeware-weld.vercel.app
```

### Test
1. ✅ Backend: `https://your-backend.onrender.com/health`
2. ✅ Frontend: `https://your-app.vercel.app`
3. ✅ Try logging in

---

## Auto-Deploy
Both services auto-deploy on push to `main` branch! 🎉
