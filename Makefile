# Variables
BINARY_NAME=dns-resolver
GO=go
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GORUN=$(GO) run

# Targets
all: build

# Build the project
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o bin/$(BINARY_NAME) main.go

# Run the project
run:
	@echo "Running $(BINARY_NAME)..."
	$(GORUN) main.go

# Clean the build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/

# Test the project
test:
	@echo "Testing..."
	$(GOTEST) -v ./...

# Help (list all targets)
help:
	@echo "Available targets:"
	@echo "  build    - Build the project"
	@echo "  run      - Run the project"
	@echo "  clean    - Clean the build artifacts"
	@echo "  test     - Run tests"

.PHONY: all build run clean test deps fmt lint check help