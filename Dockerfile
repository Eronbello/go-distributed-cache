# ----------------------------
# Stage 1: Build the Go binary
# ----------------------------
FROM golang:1.23-bullseye AS builder

WORKDIR /app

# Copy go.mod and go.sum first to cache module downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the cache-node binary
RUN CGO_ENABLED=0 go build -o /bin/cache-node ./cmd/cache-node/main.go

# -----------------------------------
# Stage 2: Create the final runtime
# -----------------------------------
FROM debian:bullseye-slim

# Ensure we have CA certificates, and clean up apt data
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the compiled binary from the builder stage
COPY --from=builder /bin/cache-node /usr/local/bin/cache-node

# Expose the application port
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/cache-node"]
