.PHONY: test test-coverage test-race test-all

# Run all tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -func=coverage.txt

# Run tests with race detection
test-race:
	go test -race ./...

# Run all tests with coverage and race detection
test-all:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -func=coverage.txt
