# API Documentation

Base URL: `http://localhost:8080/api`

## Authentication

All endpoints except `/auth/*` require JWT token in Authorization header:
```
Authorization: Bearer <token>
```

### Register
**POST** `/auth/register`

Request:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "full_name": "John Doe",
  "organization_id": 1,
  "timezone": "Asia/Jakarta"
}
```

### Login
**POST** `/auth/login`

Request:
```json
{
  "email": "admin@demo.com",
  "password": "password123"
}
```

### Change Password
**POST** `/auth/change-password`

Request:
```json
{
  "current_password": "oldpassword",
  "new_password": "newpassword123"
}
```

---

## User & Permissions

### Get My Permissions
**GET** `/users/me/permissions`

Returns current user's team roles and capabilities.

Response:
```json
{
  "user_id": 1,
  "email": "admin@demo.com",
  "full_name": "Admin User",
  "teams": [
    {
      "team_id": 1,
      "team_name": "Engineering",
      "role": "manager",
      "can_edit": true,
      "can_delete": true,
      "can_manage": true
    }
  ],
  "is_org_admin": false
}
```

---

## Organizations

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/organizations` | List organizations |
| POST | `/organizations` | Create organization |
| GET | `/organizations/:id` | Get organization |
| PUT | `/organizations/:id` | Update organization |
| DELETE | `/organizations/:id` | Delete organization |

---

## Teams

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/teams` | List teams (user's organization) |
| POST | `/teams` | Create team |
| GET | `/teams/:id` | Get team details |
| PUT | `/teams/:id` | Update team |
| DELETE | `/teams/:id` | Delete team |
| GET | `/teams/:id/members` | Get team members |
| POST | `/teams/:id/members` | Add member |
| DELETE | `/teams/:id/members/:userId` | Remove member |

**Roles:** `stakeholder`, `member`, `assistant`, `manager`

---

## Issue Statuses

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/statuses` | Get organization statuses |
| POST | `/statuses` | Create status |
| PUT | `/statuses/:id` | Update status |
| DELETE | `/statuses/:id` | Delete status |

---

## Issues

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/issues?team_id=1` | List issues for team |
| POST | `/issues` | Create issue |
| GET | `/issues/:id` | Get issue details |
| PUT | `/issues/:id` | Update issue (with deadline) |
| DELETE | `/issues/:id` | Delete issue |
| POST | `/issues/:id/assign` | Assign to user |
| POST | `/issues/:id/status` | Update status |
| POST | `/issues/:id/hold` | Put on hold |
| POST | `/issues/:id/resume` | Resume from hold |
| GET | `/issues/:id/activities` | Get activity log |
| POST | `/issues/:id/worklog` | Log work |

**Priority:** `LOW`, `NORMAL`, `HIGH`, `URGENT`

### Create/Update Issue
```json
{
  "team_id": 1,
  "status_id": 1,
  "title": "Issue title",
  "description": "Description",
  "priority": "HIGH",
  "deadline": "2025-12-31"
}
```

---

## Comments

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/issues/:id/comments` | List comments |
| POST | `/issues/:id/comments` | Add comment |
| PUT | `/issues/:id/comments/:commentId` | Update comment |
| DELETE | `/issues/:id/comments/:commentId` | Delete comment |

Request:
```json
{
  "content": "This is a comment"
}
```

---

## Attachments

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/issues/:id/attachments` | List attachments |
| POST | `/issues/:id/attachments` | Upload file (multipart) |
| GET | `/attachments/:id/download` | Download file |
| DELETE | `/attachments/:id` | Delete attachment |

**Allowed types:** Images, Documents (PDF, DOC, XLS), Text files

---

## Meetings

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/meetings?team_id=1` | List team meetings |
| POST | `/meetings` | Create meeting |
| GET | `/meetings/:id` | Get meeting details |
| PUT | `/meetings/:id` | Update meeting |
| DELETE | `/meetings/:id` | Delete meeting |
| POST | `/meetings/:id/attendees` | Add attendee |
| POST | `/meetings/:id/respond` | Respond (accept/decline) |

### Create Meeting
```json
{
  "team_id": 1,
  "title": "Sprint Planning",
  "description": "Weekly planning meeting",
  "meeting_date": "2025-12-27",
  "start_time": "09:00",
  "end_time": "10:00",
  "location": "Conference Room A",
  "is_recurring": false,
  "attendee_ids": [2, 3, 4]
}
```

---

## Calendar

**GET** `/calendar`

Query params:
- `start_date` (required): YYYY-MM-DD
- `end_date` (required): YYYY-MM-DD
- `team_id` (optional): Filter by team
- `user_id` (optional): Filter by user

---

## Analytics

### Dashboard Analytics
**GET** `/analytics/dashboard`

Response:
```json
{
  "total_tasks": 10,
  "completed_tasks": 5,
  "in_progress_tasks": 3,
  "on_hold_tasks": 1,
  "overdue_tasks": 1,
  "tasks_by_status": [...],
  "tasks_by_priority": [...],
  "weekly_activity": [...],
  "team_stats": [...]
}
```

---

## Error Responses

```json
{
  "error": "Error message description"
}
```

| Code | Description |
|------|-------------|
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 500 | Server Error |
