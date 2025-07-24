# Insurance Management API - Docker Deployment

This project provides a complete Docker Compose setup for deploying the Insurance Management API with PostgreSQL database.

## Quick Start

### Prerequisites

- Docker Desktop installed on your machine
- Docker Compose (included with Docker Desktop)

### 1. Clone and Navigate to Project

```bash
git clone <your-repo-url>
cd one_click_2
```

### 2. Build and Start Services

```bash
docker-compose up --build
```

This command will:

- Build the Go application Docker image
- Start PostgreSQL database container
- Start the application container
- Automatically create the database schema
- Create a default MasterAdmin account

### 3. Access the Application

- **API Base URL**: http://localhost:8081
- **Database**: localhost:5432 (accessible from host machine)

### Default Admin Account

```
Email: admin@insurance.com
Password: admin123
```

**⚠️ IMPORTANT: Change this password after first login!**

## Services

### Application Service (`app`)

- **Container Name**: insurance_app
- **Port**: 8081
- **Environment**: Production mode (GIN_MODE=release)
- **Health**: Depends on PostgreSQL being healthy

### Database Service (`postgres`)

- **Container Name**: insurance_db
- **Port**: 5432
- **Database**: PostgreSQL 15 Alpine
- **Default Database**: `insurance` (created automatically)
- **Persistent Storage**: Docker volume `postgres_data`

## Environment Variables

The application supports the following environment variables:

| Variable           | Default             | Description              |
| ------------------ | ------------------- | ------------------------ |
| `DB_HOST`          | localhost           | Database host            |
| `DB_PORT`          | 5432                | Database port            |
| `DB_USER`          | postgres            | Database user            |
| `DB_PASSWORD`      | example             | Database password        |
| `DB_NAME`          | insurance           | Database name            |
| `GIN_MODE`         | release             | Gin framework mode       |
| `ADMIN_EMAIL`      | admin@insurance.com | Default admin email      |
| `ADMIN_PASSWORD`   | admin123            | Default admin password   |
| `ADMIN_FIRST_NAME` | System              | Default admin first name |
| `ADMIN_LAST_NAME`  | Administrator       | Default admin last name  |

### Using Custom Environment Variables

You can override any of these variables by:

1. **Creating a `.env` file** (copy from `.env.example`):

```bash
cp .env.example .env
# Edit .env with your values
```

2. **Setting environment variables in docker-compose.yml**:

```yaml
environment:
  - ADMIN_EMAIL=your-admin@company.com
  - ADMIN_PASSWORD=your-secure-password
```

3. **Using environment variables in your shell**:

```bash
export ADMIN_EMAIL=your-admin@company.com
export ADMIN_PASSWORD=your-secure-password
docker-compose up
```

## Docker Commands

### Start Services (detached mode)

```bash
docker-compose up -d
```

### Stop Services

```bash
docker-compose down
```

### View Logs

```bash
# All services
docker-compose logs

# Specific service
docker-compose logs app
docker-compose logs postgres
```

### Rebuild Application

```bash
docker-compose up --build app
```

### Remove Everything (including volumes)

```bash
docker-compose down -v
```

## Development vs Production

### Development

For development with hot reload, you can mount the source code:

```yaml
volumes:
  - .:/app
  - /app/vendor
```

### Production

The current setup is optimized for production with:

- Multi-stage Docker build for smaller image size
- Non-root user execution
- Health checks
- Restart policies
- Persistent data volumes

## API Documentation

Once running, you can test the API endpoints:

### Authentication

```bash
POST http://localhost:8081/api/auth/login
POST http://localhost:8081/api/auth/reset-password
PUT http://localhost:8081/api/auth/reset-password/:token
```

### Profile Management (All Roles)

```bash
# Get current user's profile
GET http://localhost:8081/api/profile/me

# Get specific user's profile (role-based access)
GET http://localhost:8081/api/profile/user/:id

# Get users by role (role-based access)
GET http://localhost:8081/api/profile/users?role=Agent
```

### Master Admin Endpoints

```bash
GET http://localhost:8081/api/master-admin/agencies
GET http://localhost:8081/api/master-admin/providers
# Profile management
GET http://localhost:8081/api/master-admin/profiles/users?role=Agent
GET http://localhost:8081/api/master-admin/profiles/user/:id
```

### Agency Admin Endpoints

```bash
GET http://localhost:8081/api/agency-admin/locations
GET http://localhost:8081/api/agency-admin/agents
# Profile management (agency-scoped)
GET http://localhost:8081/api/agency-admin/profiles/users?role=Agent
GET http://localhost:8081/api/agency-admin/profiles/user/:id
```

### Agent Endpoints

```bash
GET http://localhost:8081/api/agent/customers
GET http://localhost:8081/api/agent/quotes
# Profile management (agency-scoped)
GET http://localhost:8081/api/agent/profiles/users?role=Customer
GET http://localhost:8081/api/agent/profiles/user/:id
```

### Profile API Access Control

**Master Admin:**

- Can view profiles of all users (MasterAdmin, AgencyAdmin, LocationAdmin, Agent, Customer)
- Has access to all profile endpoints

**Agency Admin:**

- Can view profiles of Agent, AgencyAdmin, LocationAdmin, and Customer within their agency
- Cannot view MasterAdmin profiles or users from other agencies

**Agent:**

- Can view profiles of Customer and other Agents within their agency
- Cannot view admin profiles or users from other agencies

**All Users:**

- Can always view their own profile using `/api/profile/me`
  GET http://localhost:8081/api/master-admin/providers

````

### Database Access

If you need direct database access:

```bash
docker exec -it insurance_db psql -U postgres -d insurance
````

## Troubleshooting

### Port Already in Use

If port 8081 or 5432 is already in use, modify the ports in `docker-compose.yml`:

```yaml
ports:
  - "8082:8081" # Change host port
```

### Database Connection Issues

Check if PostgreSQL is healthy:

```bash
docker-compose ps
docker-compose logs postgres
```

### Application Logs

```bash
docker-compose logs app
```

### Reset Database

To reset the database completely:

```bash
docker-compose down -v
docker-compose up --build
```

## File Structure in Container

```
/root/
├── main (Go binary)
└── database/
    └── schema.sql
```

## Network

Services communicate on the `insurance_network` bridge network:

- App can reach database at `postgres:5432`
- Database is isolated from external access except mapped ports

## Volumes

- `postgres_data`: Persistent PostgreSQL data storage
- `./database:/app/database:ro`: Read-only mount of schema files
