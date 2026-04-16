.PHONY: build run test clean install-deps

# Build the server binary
build:
	go build -o bin/crochetbot cmd/server/main.go

# Run the server
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run integration tests (requires test STL files)
test-integration:
	go test -v -tags=integration ./internal/pattern/

# Run all tests
test-all: test test-integration

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Install Go dependencies
install-deps:
	go mod download
	go mod tidy

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/crochetbot-linux cmd/server/main.go

# Docker build
docker-build:
	docker build -t crochetbot:latest .

# Docker run
docker-run:
	docker run -p 8080:8080 crochetbot:latest
