# Makefile for GoCart API

# Variables
BINARY_NAME=server
BUILD_DIR=bin
CMD_PATH=./cmd/server

# help
help:
	@echo "Available targets:"
	@grep -E "^[a-zA-Z_-]+:" Makefile | grep -v "^help:" | sed "s/:.*//"	

# Build the application
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
run: build
	@echo "Running..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run in development mode (using go run)
dev:
	@echo "Running in dev mode..."
	@go run $(CMD_PATH)/main.go

# Lint the code
lint:
	@echo "Linting..."
	@golangci-lint run ./...

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Database Migrations (Placeholder commands)
migrate-up:
	@echo "Running migrations up..."
	# migrate -path db/migrations -database "postgresql://user:password@localhost:5432/gocart?sslmode=disable" up

migrate-down:
	@echo "Running migrations down..."
	# migrate -path db/migrations -database "postgresql://user:password@localhost:5432/gocart?sslmode=disable" down

# Docker Compose
docker-up:
	@echo "Starting docker services..."
	@docker-compose --env-file .env -f docker/docker-compose.yml up -d

docker-down:
	@echo "Stopping docker services..."
	@docker-compose --env-file .env -f docker/docker-compose.yml down

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

.PHONY: build run dev lint test migrate-up migrate-down clean docker-up docker-down
