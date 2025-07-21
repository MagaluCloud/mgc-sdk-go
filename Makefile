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

readthedocs:
	cd docs && rm -rf output && rm -rf source && mkdir source && touch source/.keep
	sphinx-build -b html ./docs/source ./docs/output/html -c ./docs
