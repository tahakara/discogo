# DiscoGo Service Discovery Tool

<p align="center">
  <img src="discoGo.png" alt="DiscoGo Logo" width="200"/>
</p>

DiscoGo is a service discovery tool built with Go and Memcached. It provides a simple, reliable, and configurable HTTP API for service registration, deregistration, and discovery in distributed environments.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)
- [License](#license)

## Features

- Memcached-based service discovery
- Easy configuration via environment variables (`.env` support)
- Colored and detailed logging (configurable)
- Standardized JSON responses for all endpoints
- Health, version, registration, deregistration, and discovery HTTP API

## Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/your-username/discogo.git
   cd discogo
   ```

2. **Install Go** (if not already installed):  
   [Download Go](https://golang.org/dl/)

3. **Install dependencies:**
   ```sh
   go mod tidy
   ```

## Configuration

Create a `.env` file in the root directory with the following content:

```
DISCOGO_NAME=discoGo
DISCOGO_VERSION=1.0.0
DISCOGO_VERSION_NAME=Astrid
DISCOGO_LOG_COLOR=true

DISCOGO_HTTP_HOST=localhost
DISCOGO_HTTP_PORT=8080

MEMCACHED_HOST=127.0.0.1
MEMCACHED_PORT=11211
```

Adjust the values as needed for your environment.

## Usage

To start the service:

```sh
go run cmd/main.go
```

The application will read configuration from the `.env` file. If required variables are missing, it will log an error and exit.

## API Endpoints

| Method | Endpoint      | Description                |
|--------|--------------|----------------------------|
| GET    | /health      | Health check               |
| GET    | /version     | Returns version info       |
| GET    | /heartbeat   | Service uptime             |
| GET    | /discover    | Discover registered service|
| POST   | /register    | Register a new service     |
| POST   | /deregister  | Deregister a service       |

All responses follow a standard JSON structure.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for improvements or new features.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.