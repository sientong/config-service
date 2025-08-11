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
	go test ./... -v

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Run Docker container
docker-run:
	make docker-stop
	@echo "Running Docker container..."
	docker run --name $(DOCKER_CONTAINER_NAME) --rm -d -p 3000:3000 $(DOCKER_IMAGE)

# Stop Docker container
docker-stop:
	@echo "Stopping Docker container..."
	docker stop $(DOCKER_CONTAINER_NAME)

# Clean local builds
clean:
	@echo "Cleaning build files..."
	rm -f $(APP_NAME)