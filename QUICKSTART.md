# Quick Start Guide

## Prerequisites
- PostgreSQL 14+
- Go 1.24+
- Node.js 18+

## Setup (5 minutes)

### 1. Database
```bash
createdb taskmanagement
cd database/migrations
psql -U postgres -d taskmanagement -f run_migrations.sql
cd ../seeds
psql -U postgres -d taskmanagement -f seed_data.sql
```

### 2. Backend
```bash
cd backend
go run main.go
```

### 3. Frontend
```bash
cd frontend
npm install
npm run dev
```

### 4. Login
- Open http://localhost:5173
- Email: `admin@demo.com`
- Password: `password123`

## What's Working

✅ **Backend API** (100% functional)
- Authentication (JWT)
- Organizations, Teams, Users
- Issues with workflow
- Assignments with timeline
- Calendar view
- Work logging
- Activity tracking

✅ **Frontend** (Basic)
- Login/logout
- Protected routes
- Dashboard placeholder

## Next Steps

The backend is production-ready. Frontend needs:
- Kanban board UI
- Calendar component
- Issue forms
- Team management UI

See [README.md](file:///Users/sonyadriko/Projects/task-management/README.md) for full documentation.
