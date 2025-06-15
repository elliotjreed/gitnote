.PHONY: build test clean install run-tests lint

# Build the application
build:
	go build -o bin/gitnote .

# Install the application
install:
	go install .

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run the application
run:
	go run .

# Download dependencies
deps:
	go mod download
	go mod tidy

# Verify dependencies
verify:
	go mod verify

# Format code
fmt:
	go fmt ./...

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/gitnote-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o bin/gitnote-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/gitnote-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o bin/gitnote-windows-amd64.exe .