
# LOAD - Load Balancer in Go

LOAD is a load balancer application developed in Go. It uses a Round Robin algorithm to distribute HTTP requests across multiple backend servers, with real-time health checks to ensure that only healthy instances are used. The project is designed to be extensible, easy to maintain, and deploy.

## 🚀 Features

- **Load Balancing**: Uses the Round Robin algorithm to distribute HTTP requests across backend servers.
- **Health Checking**: Continuous monitoring of backend health with automatic failover to healthy servers.
- **HTTP Proxy**: A reverse proxy that intercepts incoming requests and routes them to the selected backends.
- **Graceful Shutdown**: Handles server shutdown smoothly, ensuring active connections are closed gracefully.
- **Unit Tests**: Automated tests to ensure the proper functioning of the load balancer, proxy, and error handling.

## 📦 Prerequisites

- Go 1.20 or higher
- Docker (optional for backend configuration)

## 🏗️ Installation

### 1. Clone the repository

```bash
git clone https://github.com/0SansNom/load.git
cd load
```

### 2. Initialize the Go module

If you haven’t already, initialize your Go module:

```bash
go mod init github.com/0SansNom/load
```

### 3. Install dependencies

Install the required dependencies using `go get`:

```bash
go get github.com/stretchr/testify v1.8.0
go get go.uber.org/zap v1.22.0
```

### 4. Configure the backends

In `cmd/main.go`, configure the list of backends where the load balancer will send requests. For example:

```go
backends := []string{
    "http://localhost:9001",
    "http://localhost:9002",
}
```

You can use real backend servers, or you can use Docker to start mock servers for testing.

### 5. Run the application

Start the application by running:

```bash
go run cmd/main.go
```

This will start the HTTP server on port `8080`. Incoming HTTP requests will be balanced across the defined backends.

## 📜 Architecture

The architecture of LOAD is divided into several key components:

1. **`cmd/main.go`**: The entry point of the application. It initializes dependencies, configures the load balancer, and starts the HTTP server.
2. **`internal/server`**: Manages the HTTP server, including request handling and server start/stop.
3. **`internal/proxy`**: The reverse proxy that routes requests to the backends.
4. **`internal/balancer`**: Contains the load balancing algorithm (Round Robin) and manages backend servers.
5. **`internal/health`**: Performs health checks on the backends.
6. **`internal/logger`**: Provides centralized logging using `zap`.
7. **`tests`**: Contains unit and integration tests to validate the system's functionality.

## ⚙️ Tests

### Running Unit Tests

You can run the unit tests with the following command:

```bash
go test ./...
```

This will run all the tests defined in the `*_test.go` files in your project.

### Adding a Mock Backend for Tests

You can use the `setupMockBackend` function to create a mock HTTP server that responds to requests, which is useful for unit testing.

## 🛠️ Development

### Adding a Backend

If you want to add a new backend to your system, simply update the list of backends in `cmd/main.go`:

```go
backends := []string{
    "http://localhost:9001",
    "http://localhost:9002",
    "http://localhost:9003", // New backend
}
```

### Adding a New Load Balancing Algorithm

To add a new load balancing algorithm, create a new type that implements the `Balancer` interface and define it in `internal/balancer`.

### Contribution

If you’d like to contribute to this project, here are a few steps to get started:

1. Fork this repository.
2. Create a branch for your feature (`git checkout -b feature/my-feature`).
3. Make your changes and add tests to validate your feature.
4. Submit a Pull Request describing your changes.

## 📄 License

This project is licensed under the MIT License. See the `LICENSE` file for details.

---

### 📌 **Notes**:
- The README is designed to be clear and concise while providing all the necessary information to install, configure, and contribute to the project.
- If you have additional preferences (such as extra tools or specific configurations), feel free to add them to this template.