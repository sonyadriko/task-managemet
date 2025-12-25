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

Response:
```json
{
  "token": "eyJhbGc...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe",
    "organization_id": 1,
    "timezone": "Asia/Jakarta"
  }
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

Response: Same as register

## Organizations

### List Organizations
**GET** `/organizations`

### Create Organization
**POST** `/organizations`

Request:
```json
{
  "name": "My Company",
  "description": "Company description"
}
```

### Get Organization
**GET** `/organizations/:id`

### Update Organization
**PUT** `/organizations/:id`

### Delete Organization
**DELETE** `/organizations/:id`

## Teams

### List Teams
**GET** `/teams`

Returns teams for authenticated user's organization.

### Create Team
**POST** `/teams`

Request:
```json
{
  "name": "Engineering",
  "description": "Software development team",
  "parent_team_id": null
}
```

### Get Team
**GET** `/teams/:id`

### Update Team
**PUT** `/teams/:id`

### Delete Team
**DELETE** `/teams/:id`

### Get Team Members
**GET** `/teams/:id/members`

### Add Team Member
**POST** `/teams/:id/members`

Request:
```json
{
  "user_id": 2,
  "role": "member"
}
```

Roles: `stakeholder`, `member`, `assistant`, `manager`

### Remove Team Member
**DELETE** `/teams/:id/members/:userId`

## Issue Statuses

### Get Team Statuses
**GET** `/statuses/team/:teamId`

### Create Status
**POST** `/statuses`

Request:
```json
{
  "team_id": 1,
  "name": "IN_PROGRESS",
  "position": 2,
  "is_final": false,
  "color": "#3B82F6"
}
```

### Update Status
**PUT** `/statuses/:id`

### Delete Status
**DELETE** `/statuses/:id`

## Issues

### List Issues
**GET** `/issues?team_id=1`

Query params:
- `team_id` (required): Filter by team

### Create Issue
**POST** `/issues`

Request:
```json
{
  "team_id": 1,
  "status_id": 1,
  "title": "Implement user authentication",
  "description": "Add JWT-based authentication",
  "priority": "HIGH"
}
```

Priority: `LOW`, `NORMAL`, `HIGH`, `URGENT`

### Get Issue
**GET** `/issues/:id`

### Update Issue
**PUT** `/issues/:id`

### Delete Issue
**DELETE** `/issues/:id`

### Assign Issue
**POST** `/issues/:id/assign`

Request:
```json
{
  "user_id": 3,
  "start_date": "2025-12-26T00:00:00Z",
  "end_date": "2025-12-30T00:00:00Z"
}
```

### Update Issue Status
**POST** `/issues/:id/status`

Request:
```json
{
  "status_id": 2
}
```

### Put Issue on Hold
**POST** `/issues/:id/hold`

Request:
```json
{
  "reason": "Waiting for client approval"
}
```

### Resume Issue
**POST** `/issues/:id/resume`

### Get Issue Activities
**GET** `/issues/:id/activities`

### Log Work
**POST** `/issues/:id/worklog`

Request:
```json
{
  "work_date": "2025-12-25T00:00:00Z",
  "minutes_spent": 120,
  "notes": "Implemented authentication logic"
}
```

## Calendar

### Get Calendar Events
**GET** `/calendar`

Query params:
- `start_date` (required): YYYY-MM-DD
- `end_date` (required): YYYY-MM-DD
- `team_id` (optional): Filter by team
- `user_id` (optional): Filter by user

Example:
```
GET /calendar?start_date=2025-12-01&end_date=2025-12-31&team_id=1
```

Response:
```json
[
  {
    "issue_id": 1,
    "issue_title": "Setup CI/CD",
    "priority": "HIGH",
    "status_id": 2,
    "status_name": "IN_PROGRESS",
    "status_color": "#3B82F6",
    "user_id": 3,
    "user_name": "Developer One",
    "team_id": 1,
    "team_name": "Engineering",
    "start_date": "2025-12-25T00:00:00Z",
    "end_date": "2025-12-30T00:00:00Z"
  }
]
```

## Error Responses

All endpoints return errors in this format:

```json
{
  "error": "Error message description"
}
```

Common status codes:
- `400` - Bad Request (validation error)
- `401` - Unauthorized (missing/invalid token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `500` - Internal Server Error
