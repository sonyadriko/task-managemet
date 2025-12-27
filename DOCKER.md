# Docker Deployment Guide

## Quick Start

### 1. Build and Run
```bash
# Build and start all services
docker-compose up -d --build

# View logs
docker-compose logs -f
```

### 2. Access the Application
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Database**: localhost:5432

### 3. Initial Setup

Run database migrations and seed data:
```bash
# Connect to postgres container
docker exec -it taskmanagement-db psql -U postgres -d taskmanagement

# Or run migrations manually
docker exec -it taskmanagement-db psql -U postgres -d taskmanagement -f /docker-entrypoint-initdb.d/migrations/001_create_organizations.sql
```

### 4. Default Login
- **Email**: admin@demo.com
- **Password**: password123

## Commands

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# Rebuild after code changes
docker-compose up -d --build

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Access database
docker exec -it taskmanagement-db psql -U postgres -d taskmanagement

# Reset database (delete volume)
docker-compose down -v
docker-compose up -d --build
```

## Environment Variables

Create `.env` file in root directory for production:

```env
# Database
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your-secure-password
POSTGRES_DB=taskmanagement

# JWT
JWT_SECRET=your-super-secret-key

# Cloudflare R2 (optional)
R2_ACCOUNT_ID=
R2_ACCESS_KEY_ID=
R2_SECRET_ACCESS_KEY=
R2_BUCKET_NAME=
R2_PUBLIC_URL=
```

## Production Tips

1. **Change default passwords** in docker-compose.yml
2. **Use secrets management** for sensitive data
3. **Add SSL/TLS** with reverse proxy (Traefik/Caddy)
4. **Enable database backups** with volume mounts
5. **Set resource limits** for containers

## Troubleshooting

### Backend can't connect to database
```bash
# Check if postgres is healthy
docker-compose ps

# Check postgres logs
docker-compose logs postgres
```

### Frontend shows blank page
```bash
# Rebuild frontend
docker-compose up -d --build frontend

# Check nginx logs
docker-compose logs frontend
```
