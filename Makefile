.PHONY: run lint test

run:
	go run ./cmd/pr-reviewer-service

lint:
	go vet ./...
	golangci-lint run

test:
	go test ./...
