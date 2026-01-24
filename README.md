# Assignly

A Go-based REST API backend for task management.

## Prerequisites

- Docker
- Docker Compose

## Quick Start

1. Clone the repository

2. Create a `.env` file in the root directory:
```env
DB_CONNECTION_STRING="postgresql://user:password@postgres:5432/assignly?sslmode=disable"
DB_USER=user
DB_PASSWORD=password
DB_NAME=assignly
JWT_SECRET=your-secret-key-here
```

3. Run with Docker Compose:
```bash
docker-compose up --build
```

4. The API will be available at `http://localhost:8080`

## Stopping the Application

```bash
docker-compose down
```

To also remove the database volume:
```bash
docker-compose down -v
```

## API Endpoints

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/users/create` | Create a new user |
| POST | `/users/login` | Login with email |
| GET | `/users/get` | Get user by ID |
| GET | `/users/employees-list` | Get list of employees (auth required) |

### Tasks
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/tasks/assign-task` | Assign a task (managers only) |
| GET | `/tasks/get-assigned-task` | Get assigned tasks |
| PUT | `/tasks/status/:taskId` | Update task status |
| DELETE | `/tasks/delete/:taskId` | Delete a task (managers only) |

## Project Structure

```
.
├── api/
│   ├── handler/v1/     # HTTP handlers
│   └── middleware/     # Auth middleware
├── app/
│   ├── task/          # Task module
│   └── user/          # User module
├── cmd/api/           # Application entry point
├── domain/            # Domain entities
├── pkg/               # Shared packages
├── postgres/init/     # Database initialization scripts
├── Dockerfile
└── docker-compose.yml
```
