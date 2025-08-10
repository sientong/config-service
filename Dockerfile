# Use ARM64 Debian-based image for both build and runtime
FROM golang:1.22-bookworm AS builder

WORKDIR /app

# Install dependencies for CGO (SQLite3, GCC, etc.)
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libc6-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

# Copy go mod/sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build with CGO enabled
ENV CGO_ENABLED=1
RUN go build -o config-service ./main.go

# Final runtime stage — same base as build to avoid GLIBC mismatch
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime deps for SQLite
RUN apt-get update && apt-get install -y --no-install-recommends \
    libc6 \
    libsqlite3-0 \
    && rm -rf /var/lib/apt/lists/*

# Copy the built binary from builder
COPY --from=builder /app/config-service .

EXPOSE 3000
CMD ["./config-service"]