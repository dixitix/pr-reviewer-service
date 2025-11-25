.PHONY: run fmt tidy lint test migrate-test e2e migrate-up migrate-down compose-up compose-down

MIGRATIONS_DIR := ./migrations
MIGRATE        ?= migrate

run:
	go run ./cmd/pr-reviewer-service

fmt:
	go fmt ./...

tidy:
	go mod tidy

lint:
	go vet ./...
	golangci-lint run

test: compose-up migrate-test
	go test ./...

migrate-test:
	docker compose up migrate-test

e2e:
	@set -e; \
		trap 'docker compose -f docker-compose.e2e.yml down -v' EXIT; \
		docker compose -f docker-compose.e2e.yml up -d --build; \
		go test ./test/e2e -count=1

compose-up:
	docker compose up -d --build

compose-down:
	docker compose down -v

migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_DSN)" up

migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_DSN)" down
