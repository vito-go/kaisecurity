# === Makefile for kai_security ===

APP_NAME = kai_security
DB_PATH = kai_security.db
COVER_FILE = cover.out
BUILD_DIR = build
PORT = 8080
MAIN_PKG = cmd/server/main.go

# Run all unit tests with coverage
.PHONY: test

test:
	@echo          tests with coverage...
	go test -coverprofile=$(COVER_FILE) ./...

# Show HTML coverage report
.PHONY: cover

cover: test
	@echo  Opening coverage report...
	go tool cover -html=$(COVER_FILE)

# Run the app with default args
.PHONY: run

run:
	@echo        Starting $(APP_NAME)...
	go run $(MAIN_PKG) -db=$(DB_PATH) -port=$(PORT)

# Build binary
.PHONY: build

build:
	@echo      Compiling binary to $(BUILD_DIR)/$(APP_NAME)...
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PKG)

# Clean up all generated files
.PHONY: clean

clean:
	@echo      Cleaning coverage and build artifacts...
	rm -f $(COVER_FILE) $(DB_PATH)
	rm -rf $(BUILD_DIR)

# Build Docker image
.PHONY: docker

docker:
	@echo      Building Docker image $(APP_NAME):latest...
	docker build -t $(APP_NAME):latest .

# Run tests + report only (for CI)
.PHONY: ci

ci:
	@echo      Running tests in CI mode...
	go test -v ./...

# Help
.PHONY: help

help:
	@echo "Available targets:"
	@echo "  make run       - Run the app"
	@echo "  make build     - Build binary"
	@echo "  make test      - Run unit tests with coverage"
	@echo "  make cover     - Open coverage report in browser"
	@echo "  make docker    - Build Docker image"
	@echo "  make clean     - Clean build and test files"
	@echo "  make ci        - Run test for CI"
	@echo "  make help      - Show this help"
