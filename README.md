# Task Management Application

A comprehensive Jira-lite task management system built with React.js, Go, and PostgreSQL. Features team hierarchy, flexible workflows, timeline tracking, and calendar-centric visibility.

## 🏗️ Architecture

```
task-management/
├── backend/          # Go backend (Gin framework)
│   ├── config/       # Configuration management
│   ├── models/       # GORM database models
│   ├── repositories/ # Data access layer
│   ├── services/     # Business logic layer
│   ├── handlers/     # HTTP handlers (controllers)
│   └── middleware/   # Auth, CORS, Permission middleware
├── frontend/         # React frontend (Vite + TypeScript)
│   └── src/
└── database/         # SQL migrations & seeds
    ├── migrations/   # Database schema migrations
    └── seeds/        # Sample data
```

## ✨ Features

### Core Functionality
- **Organization & Team Management**: Hierarchical team structure with parent-child relationships
- **Flexible Workflows**: Customizable status workflows per team (not hardcoded)
- **Task Management**: Create, assign, and track issues with priorities (LOW, NORMAL, HIGH, URGENT)
- **Timeline Tracking**: Assign tasks with start/end dates for calendar visibility
- **Work Logging**: Track actual time spent vs planned time
- **Calendar View**: Visualize task assignments across teams and users
- **Hold/Block Management**: Track why tasks are blocked with detailed reasons
- **Activity Logging**: Complete audit trail of all issue events

### Security & Permissions
- **JWT Authentication**: Secure token-based authentication
- **Role-Based Access Control**: Per-team roles (Manager, Assistant, Member, Stakeholder)
- **Permission Hierarchy**: Fine-grained access control based on team membership

## 🚀 Getting Started

### Prerequisites

- **Go**: 1.24+ (will be auto-upgraded during dependency installation)
- **Node.js**: 18+ and npm
- **PostgreSQL**: 14+

### 1. Database Setup

```bash
# Create database
createdb taskmanagement

# Or using psql
psql -U postgres
CREATE DATABASE taskmanagement;
\q

# Run migrations
cd database/migrations
psql -U postgres -d taskmanagement -f run_migrations.sql

# (Optional) Load seed data
cd ../seeds
psql -U postgres -d taskmanagement -f seed_data.sql
```

### 2. Backend Setup

```bash
cd backend

# Install dependencies (already done if you ran go get earlier)
go mod download

# Copy environment file
cp .env.example .env

# Edit .env and configure your database credentials
# DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, JWT_SECRET

# Run the server
go run main.go

# Server will start on http://localhost:8080
```

### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Copy environment file
cp .env.example .env

# Run development server
npm run dev

# Frontend will start on http://localhost:5173
```

## 📊 Database Schema

### Core Tables

1. **organizations** - Top-level company entities
2. **users** - User accounts with timezone support
3. **teams** - Teams/departments with hierarchical support
4. **team_members** - User-team relationships with roles
5. **issue_statuses** - Flexible workflow statuses per team
6. **issues** - Core task/work items
7. **issue_assignments** - Task assignments with timeline (WHO + WHEN)
8. **issue_work_logs** - Actual work tracking
9. **issue_status_logs** - Status change audit trail
10. **issue_activities** - Comprehensive activity log
11. **issue_hold_reasons** - Blocked task tracking

### Key Relationships

```
Organization
  └── Teams (hierarchical)
      ├── Team Members (with roles)
      ├── Issue Statuses (workflow)
      └── Issues
          ├── Assignments (timeline)
          ├── Work Logs
          ├── Activities
          └── Hold Reasons
```

## 🔐 API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login and get JWT token
- `POST /api/auth/logout` - Logout

### Organizations
- `GET /api/organizations` - List all organizations
- `POST /api/organizations` - Create organization
- `GET /api/organizations/:id` - Get organization details
- `PUT /api/organizations/:id` - Update organization
- `DELETE /api/organizations/:id` - Delete organization

### Teams
- `GET /api/teams` - List teams (filtered by organization)
- `POST /api/teams` - Create team
- `GET /api/teams/:id` - Get team details
- `PUT /api/teams/:id` - Update team
- `DELETE /api/teams/:id` - Delete team
- `GET /api/teams/:id/members` - Get team members
- `POST /api/teams/:id/members` - Add team member
- `DELETE /api/teams/:id/members/:userId` - Remove team member

### Issue Statuses
- `GET /api/statuses/team/:teamId` - Get workflow statuses for team
- `POST /api/statuses` - Create status
- `PUT /api/statuses/:id` - Update status
- `DELETE /api/statuses/:id` - Delete status

### Issues
- `GET /api/issues?team_id=X` - List issues for team
- `POST /api/issues` - Create issue
- `GET /api/issues/:id` - Get issue details
- `PUT /api/issues/:id` - Update issue
- `DELETE /api/issues/:id` - Delete issue
- `POST /api/issues/:id/assign` - Assign issue to user
- `POST /api/issues/:id/status` - Update issue status
- `POST /api/issues/:id/hold` - Put issue on hold
- `POST /api/issues/:id/resume` - Resume issue from hold
- `GET /api/issues/:id/activities` - Get issue activity log
- `POST /api/issues/:id/worklog` - Log work on issue

### Calendar
- `GET /api/calendar?team_id=X&user_id=Y&start_date=YYYY-MM-DD&end_date=YYYY-MM-DD` - Get calendar events

## 👥 Team Roles & Permissions

| Role | View | Create | Edit Own | Edit All | Manage Team |
|------|------|--------|----------|----------|-------------|
| **Stakeholder** | ✅ | ❌ | ❌ | ❌ | ❌ |
| **Member** | ✅ | ✅ | ✅ | ❌ | ❌ |
| **Assistant** | ✅ | ✅ | ✅ | ✅ | ❌ |
| **Manager** | ✅ | ✅ | ✅ | ✅ | ✅ |

## 🔄 Workflow Example

Default workflow for Engineering team (from seed data):

```
WAITING → IN_PROGRESS → QA → READY_TO_DEPLOY → DONE
              ↓
            HOLD (with reason tracking)
```

Each team can customize their own workflow by creating/modifying statuses.

## 📅 Calendar View

The calendar view displays:
- Task title and description
- Assigned user (PIC)
- Current status with color coding
- Start and end dates
- Team information

Filter by:
- Team (show all tasks for a team)
- User (show all tasks assigned to a user)
- Date range (week, month, custom)

## 🧪 Testing with Seed Data

The seed data creates:
- 1 Demo organization
- 2 Teams (Engineering + Frontend sub-team)
- 4 Users (admin, manager, 2 developers)
- 6 Workflow statuses
- 5 Sample issues with assignments

**Login credentials** (password for all: `password123`):
- `admin@demo.com` - Manager of Engineering
- `manager@demo.com` - Manager of Frontend Team
- `developer1@demo.com` - Member
- `developer2@demo.com` - Member

## 🛠️ Development

### Backend Hot Reload (Optional)

```bash
# Install Air for hot reload
go install github.com/air-verse/air@latest

# Run with hot reload
air
```

### Frontend Development

```bash
cd frontend
npm run dev
```

### Database Migrations

To create a new migration:

```bash
cd database/migrations
# Create new file: XXX_description.sql
# Add your SQL statements
```

## 📝 Environment Variables

### Backend (.env)

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=taskmanagement
JWT_SECRET=your-super-secret-jwt-key
PORT=8080
ALLOWED_ORIGINS=http://localhost:5173
```

### Frontend (.env)

```env
VITE_API_URL=http://localhost:8080/api
```

## 🏗️ Tech Stack

### Backend
- **Go** 1.24+ - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM library
- **PostgreSQL** - Database
- **JWT** - Authentication
- **bcrypt** - Password hashing

### Frontend
- **React** 18+ - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool
- **React Router** - Routing
- **Axios** - HTTP client

## 📚 Project Structure Details

### Backend Layers

1. **Models**: GORM models matching database schema
2. **Repositories**: Data access layer with database operations
3. **Services**: Business logic layer
4. **Handlers**: HTTP request handlers
5. **Middleware**: Authentication, CORS, Permission checks

### Frontend Structure (To be implemented)

1. **API Client**: Axios setup with interceptors
2. **Contexts**: Auth and Team state management
3. **Pages**: Main application pages
4. **Components**: Reusable UI components
5. **Hooks**: Custom React hooks

## 🚧 Next Steps

- [ ] Implement frontend UI components
- [ ] Add drag-and-drop Kanban board
- [ ] Implement calendar component
- [ ] Add real-time notifications
- [ ] Add file attachments to issues
- [ ] Implement comments on issues
- [ ] Add dashboard analytics
- [ ] Deploy to production

## 📄 License

This project is created for demonstration purposes.

## 🤝 Contributing

This is a custom project. Feel free to fork and modify for your needs.

---

**Built with ❤️ using React.js, Go, and PostgreSQL**
# task-managemet
