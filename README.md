# Task Management Application

A comprehensive Jira-lite task management system built with React.js, Go, and PostgreSQL. Features team hierarchy, flexible workflows, timeline tracking, calendar-centric visibility, and role-based access control.

## ✨ Features

### Core Functionality
- **📋 Kanban Board** - Drag-and-drop task management with customizable workflows
- **📅 Calendar View** - Visualize task assignments with day detail modals
- **👥 Team Management** - Hierarchical teams with role-based access
- **📊 Dashboard Analytics** - Completion rate, status charts, weekly activity
- **📎 File Attachments** - Upload files to tasks (Cloudflare R2 storage)
- **💬 Comments** - Discussion threads on tasks
- **📆 Meetings** - Schedule and manage team meetings with attendees
- **🌙 Dark/Light Theme** - User preference with localStorage persistence

### Security & Permissions
- **JWT Authentication** - Secure token-based authentication
- **Role-Based Access Control** - Per-team roles with UI enforcement

| Role | View | Create | Edit | Delete | Manage Team |
|------|------|--------|------|--------|-------------|
| **Manager** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Assistant** | ✅ | ✅ | ✅ | ✅ | ❌ |
| **Member** | ✅ | ✅ | ✅ (assigned) | ❌ | ❌ |
| **Stakeholder** | ✅ | ❌ | ❌ | ❌ | ❌ |

## 🚀 Quick Start

### Prerequisites
- **Go** 1.21+
- **Node.js** 18+ and npm
- **PostgreSQL** 14+

### 1. Database Setup

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
cp .env.example .env
# Edit .env with your database credentials
go run main.go
# Server: http://localhost:8080
```

### 3. Frontend

```bash
cd frontend
npm install
npm run dev
# App: http://localhost:5173
```

## 🧪 Test Accounts

| Email | Password | Role |
|-------|----------|------|
| `admin@demo.com` | `password123` | Manager (Engineering) |
| `manager@demo.com` | `password123` | Manager (Frontend) |
| `developer1@demo.com` | `password123` | Member |
| `developer2@demo.com` | `password123` | Member |

## 🔗 API Endpoints

See [API.md](./API.md) for full documentation.

### Key Endpoints
- `POST /api/auth/login` - Authentication
- `GET /api/users/me/permissions` - User roles & permissions
- `GET /api/issues?team_id=X` - List tasks
- `GET /api/analytics/dashboard` - Dashboard stats
- `GET /api/meetings?team_id=X` - Team meetings
- `GET /api/calendar` - Calendar events

## 🏗️ Architecture

```
task-management/
├── backend/           # Go (Gin + GORM)
│   ├── handlers/      # HTTP handlers
│   ├── services/      # Business logic
│   ├── repositories/  # Data access
│   ├── models/        # Database models
│   └── middleware/    # Auth, CORS, Permissions
├── frontend/          # React (Vite + TypeScript)
│   ├── pages/         # Dashboard, Board, Calendar, Teams, Meetings
│   ├── components/    # Sidebar, shared components
│   ├── contexts/      # Auth, Theme, Permissions
│   └── api/           # API client
└── database/          # SQL migrations & seeds
```

## 📊 Database Schema

**Core Tables:**
- `organizations` - Companies
- `users` - User accounts
- `teams` - Team hierarchy
- `team_members` - Roles per team
- `issue_statuses` - Workflow statuses
- `issues` - Tasks with deadlines
- `issue_assignments` - Task assignments
- `issue_comments` - Task comments
- `attachments` - File uploads
- `meetings` - Team meetings

## 📝 Environment Variables

### Backend (.env)
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=taskmanagement
JWT_SECRET=your-secret-key
PORT=8080
ALLOWED_ORIGINS=http://localhost:5173

# Optional: Cloudflare R2 for attachments
R2_ACCOUNT_ID=
R2_ACCESS_KEY_ID=
R2_SECRET_ACCESS_KEY=
R2_BUCKET_NAME=
R2_PUBLIC_URL=
```

### Frontend (.env)
```env
VITE_API_URL=http://localhost:8080/api
```

## 🛠️ Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go, Gin, GORM, PostgreSQL |
| Frontend | React 18, TypeScript, Vite |
| Auth | JWT, bcrypt |
| Storage | Cloudflare R2 (S3-compatible) |
| Drag & Drop | @hello-pangea/dnd |

## 📄 License

MIT License - Free to use and modify.

---

**Built with ❤️ using React.js, Go, and PostgreSQL**
