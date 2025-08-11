# ===== Builder Stage =====
FROM golang:1.23-bookworm AS builder

WORKDIR /app

# Install dependencies for CGO (SQLite3, GCC, etc.)
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libc6-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

# Install swag CLI for Swagger doc generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod/sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Generate Swagger docs (requires docs to be in ./docs or ./swagger)
RUN swag init -g main.go -o ./docs

# Build with CGO enabled
ENV CGO_ENABLED=1
RUN go build -o config-service ./main.go

# ===== Runtime Stage =====
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime deps for SQLite
RUN apt-get update && apt-get install -y --no-install-recommends \
    libc6 \
    libsqlite3-0 \
    && rm -rf /var/lib/apt/lists/*

# Copy the built binary and Swagger docs
COPY --from=builder /app/config-service .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/schemas ./schemas

EXPOSE 3000
CMD ["./config-service"]