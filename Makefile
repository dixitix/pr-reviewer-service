.PHONY: run lint test migrate-up migrate-down

MIGRATIONS_DIR := ./migrations
MIGRATE        ?= migrate

run:
	go run ./cmd/pr-reviewer-service

lint:
	go vet ./...
	golangci-lint run

test:
	go test ./...

migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_DSN)" up

migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_DSN)" down
