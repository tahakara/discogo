# DiscoGo

DiscoGo is a service discovery and health monitoring system written in Go. It provides RESTful APIs for service registration, heartbeat, discovery, and deregistration, using Redis as a backend.

![DiscoGo Architecture](./shared/discogo-architecture.png)

## Features

- Service registration with metadata
- Heartbeat endpoint for health checks
- Service discovery with filtering and pagination
- Deregistration of services
- Health check for API and Redis
- Swagger/OpenAPI documentation

## Project Structure

```
.
├── cmd/                # Application entrypoint
│   └── mian.go
├── docs/               # Swagger/OpenAPI docs
├── internal/
│   ├── api/            # HTTP API handlers, DTOs, validators
│   ├── config/         # Configuration loading
│   ├── logger/         # Logging utilities
│   ├── redis/          # Redis client and helpers
│   ├── service/        # Service startup logic
│   └── utils/          # Utility functions
├── shared/             # Shared assets (architecture diagram, etc.)
├── conf.json           # Service types and providers config
├── .env                # Environment variables (development)
├── .env-prod           # Environment variables (production)
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.25+
- Redis server

### Setup

1. Clone the repository.
2. Copy `.env` or `.env-prod` and adjust environment variables as needed.
3. Install dependencies:

   ```sh
   go mod tidy
   ```

4. Run the application:

   ```sh
   go run ./cmd/mian.go
   ```

### API Documentation

Swagger UI is available at `/swagger/index.html` when the server is running.

## Main Endpoints

- `POST /disco/register` — Register a new service
- `POST /disco/heartbeat/{uuid}` — Send heartbeat for a service
- `GET  /disco/discover` — Discover services
- `POST /deregister` — Deregister a service
- `GET  /disco/health` — Health check
- `GET  /disco/version` — Version info

## Configuration

- Service types and providers are defined in [`conf.json`](conf.json).
- Environment variables are loaded from `.env` or `.env-prod`.

## License

MIT