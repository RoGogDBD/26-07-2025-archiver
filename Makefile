SERVER_DIR=cmd/server
APP_DIR=internal/app

COVERAGE=coverage.out
COVERAGE_HTML=coverage.html

.PHONY: all test build cover clean

all: test build

test: test-server test-app

test-server:
	@echo "--- Running tests in $(SERVER_DIR) ---"
	@go test -v -coverprofile=$(COVERAGE) ./$(SERVER_DIR)
	@go tool cover -html=$(COVERAGE) -o coverage.html
	@xdg-open coverage.html
	@echo "--- Completed ---"

test-app:
	@echo "--- Running tests in $(APP_DIR) ---"
	@go test -v -coverprofile=$(COVERAGE) -covermode=atomic ./$(APP_DIR)
	@go tool cover -html=$(COVERAGE) -o coverage.html
	@xdg-open coverage.html
	@echo "--- Completed ---"

build:
	@echo "--- Building the project ---"
	@go build -o bin/app ./cmd/server
	@echo "--- Completed ---"

clean:
	@echo "--- Cleaning build artifacts ---"
	@rm -f $(COVERAGE) $(COVERAGE_HTML) bin/app
	@rm -rf bin
	@echo "--- Completed ---"

.PHONY: all test clean build