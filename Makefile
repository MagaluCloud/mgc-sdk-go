.PHONY: test test-coverage test-race test-all go-fmt

# Run all tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -func=coverage.txt
	go tool cover -html=coverage.txt
	
# Run tests with race detection
test-race:
	go test -race ./...

# Run all tests with coverage and race detection
test-all:
	@go install gotest.tools/gotestsum@latest
	go run gotest.tools/gotestsum@latest -- -covermode=atomic -coverprofile=cover.out ./... 

go-fmt:
	gofmt -s -l -w .

go-vet:
	go vet ./...

# Generate documentation from Go code
docs:
	@echo "Generating documentation from Go code..."
	@chmod +x scripts/generate-docs.sh
	@./scripts/generate-docs.sh

# Generate documentation using Python script
docs-python:
	@echo "Generating documentation using Python script..."
	@python3 scripts/generate_docs.py . --html

# Generate documentation and build HTML
docs-html:
	@echo "Generating documentation and building HTML..."
	@cd docs && make html
