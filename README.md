# DiscoGo Service Discovery Tool

DiscoGo is a service discovery tool built using Go and Memcached. It aims to provide a reliable way to discover services in a distributed environment. This README provides an overview of the project, setup instructions, and usage guidelines.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Features

- Service discovery using Memcached
- Environment variable configuration
- Logging support for monitoring and debugging
- Simple HTTP API for service interactions

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/microsoft/vscode-remote-try-go.git
   cd discogo
   ```

2. Ensure you have Go installed on your machine. You can download it from [golang.org](https://golang.org/dl/).

3. Install the necessary dependencies:
   ```
   go mod tidy
   ```

## Configuration

Before running the application, you need to set up the environment variables. Create a `.env` file in the root directory with the following content:

```
SERVER_IP=<your_server_ip>
```

Replace `<your_server_ip>` with the actual IP address of your server.

## Usage

To start the service discovery tool, run the following command:

```
go run cmd/main.go
```

The application will attempt to read the server IP address from the `.env` file. If it fails to retrieve the IP address, it will log an error message and exit.

The tool exposes an HTTP API for service discovery. You can interact with the API using tools like `curl` or Postman.

## Contributing

Contributions are welcome! If you have suggestions for improvements or new features, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.