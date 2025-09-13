# Makefile for C2 Server

# Variables
BINARY_NAME=c2server
BUILD_DIR=build
MAIN_PATH=./cmd/server
GO_FILES=$(shell find . -name "*.go")

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
.PHONY: run
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	golint ./...

# Vet code
.PHONY: vet
vet:
	go vet ./...

# Build Docker image
.PHONY: docker-build
docker-build:
	docker build -t hermes-c2-server .

# Run Docker container
.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 hermes-c2-server

# Start Docker Compose
.PHONY: docker-compose-up
docker-compose-up:
	docker-compose up -d

# Stop Docker Compose
.PHONY: docker-compose-down
docker-compose-down:
	docker-compose down