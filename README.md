
# Distributed Cache with Consistent Hashing (DDD, Clean Code, SOLID)

## Project Overview

This project demonstrates a scalable distributed cache system using consistent hashing, written in Go (Golang). The architecture follows principles from **Domain-Driven Design (DDD)**, **SOLID**, and **Clean Code**, making it maintainable and extensible. The system can be containerized using Docker and orchestrated with Docker Compose, allowing you to easily run multiple cache nodes for a small local cluster.

## Features

- **Consistent Hashing** via a custom hash ring for balanced key distribution.
- **In-memory caching** with automatic key expiration (TTL).
- **HTTP API** for basic cache operations: `GET` and `POST`.
- **Containerization** with Docker and Docker Compose for multi-node clusters.
- **Clean Architecture** (Domain, Application, Infrastructure layers).
- Easily extensible (e.g., replication, advanced discovery, etc.).

## Technologies

- **Go (Golang)**
- **Docker, Docker Compose**
- **Debian or Alpine-based Docker images** (configurable)
- **Gorilla Mux** (HTTP routing)

## Directory Structure (High-Level)

- `cmd/cache-node`: Main entry point for each cache node.
- `internal/domain`: Core domain logic (cache entities, consistent hashing).
- `internal/application`: Orchestration and use cases (`DistributedCacheService`).
- `internal/infrastructure`: HTTP handlers, logging, etc.
- `docker-compose.yml` and `Dockerfile`: Container build and orchestration.

## Pre-Requisites

- **Go 1.20+** (if building locally)
- **Docker & Docker Compose** (for running multi-node setups)
- **Git** (for cloning the repository)

## Installation & Build (Local)

1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/distributed-cache.git
    ```

2. Navigate to the project directory:
    ```bash
    cd distributed-cache
    ```

3. Initialize Go modules (if not already done):
    ```bash
    go mod tidy
    ```

4. Build the binary (optional if you only plan to use Docker):
    ```bash
    go build -o cache-node ./cmd/cache-node/main.go
    ```

## Docker & Docker Compose Usage

1. In the project root, ensure you have Docker and Docker Compose installed.

2. Build and start the multi-node cluster:
    ```bash
    docker-compose up --build
    ```

3. You should see logs for each node (`node1`, `node2`, `node3`). Each node listens on port 8080 internally, mapped to different ports locally (e.g., 8081, 8082, 8083).

4. Test a `POST` (store) operation on `node1`:
    ```bash
    curl -X POST http://localhost:8081/cache \
         -H "Content-Type: application/json" \
         -d '{"key":"gopher","value":"Hello, World!","ttl_seconds":60}'
    ```

5. Retrieve the same key from `node2` or `node3`:
    ```bash
    curl "http://localhost:8082/cache?key=gopher"
    ```

## Contributing (Open Source Guidelines)

1. Fork this repository.
2. Create a feature branch for your changes:
    ```bash
    git checkout -b feature/my-awesome-improvement
    ```
3. Write clear, concise commit messages and maintain code style (`Go fmt`, lint).
4. Include tests and documentation for any new functionality.
5. Submit a pull request (PR) to the `main` branch with a detailed description of your changes.

## License

This project is released under the **MIT License** (or the license you prefer to use). Refer to the [LICENSE](./LICENSE) file in the repository for details.

## Support & Contact

- File any issues on the project's issue tracker with as much detail as possible.
- For questions or clarifications, please open a discussion or send a PR with documentation improvements.
