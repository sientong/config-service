APP_NAME=config-service
DOCKER_IMAGE=config-service:latest
DOCKER_CONTAINER_NAME=config-service-container

# Build binary locally
build:
	@echo "Building Go binary..."
	CGO_ENABLED=1 go build -o $(APP_NAME) main.go

# Run locally
run:
	@echo "Running locally..."
	go run main.go

# Run tests
tests:
	@echo "Running tests..."
	go test ./test -v

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Run Docker container
docker-run:
	make docker-stop
	@echo "Running Docker container..."
	docker run --name $(DOCKER_CONTAINER_NAME) --rm -it -p 3000:3000 $(DOCKER_IMAGE)

# Stop Docker container
docker-stop:
	@echo "Stopping Docker container..."
	@docker rm -f $(DOCKER_CONTAINER_NAME) 2>/dev/null || true

# Clean local builds
clean:
	@echo "Cleaning build files..."
	rm -f $(APP_NAME)