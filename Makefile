SERVER_DIR=cmd/server
INTERNAL_DIR=internal/...

COVERAGE=coverage.out
COVERAGE_HTML=coverage.html

.PHONY: all test build clean test-server test-internal

all: test build

test: test-server test-internal

test-server:
	@echo "--- Running tests in $(SERVER_DIR) ---"
	@go test -v -coverprofile=$(COVERAGE) ./$(SERVER_DIR)
	@go tool cover -html=$(COVERAGE) -o coverage.html
	@xdg-open coverage.html
	@echo "--- Completed ---"

test-internal:
	@echo "--- Running tests in $(INTERNAL_DIR) ---"
	@go test -v -coverprofile=$(COVERAGE) -covermode=atomic ./$(INTERNAL_DIR)
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