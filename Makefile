.PHONY: run fmt tidy lint test migrate-test e2e migrate-up migrate-down compose-up compose-down

MIGRATIONS_DIR := ./migrations
MIGRATE        ?= migrate
COMPOSE        ?= docker compose
E2E_COMPOSE    ?= docker compose -p pr-reviewer-service-e2e -f docker-compose.e2e.yml

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
	go test ./internal/repository/postgres/integration -count=1

migrate-test:
	$(COMPOSE) up migrate-test

e2e:
	@set -e; \
		trap '$(E2E_COMPOSE) down -v --remove-orphans' EXIT; \
		$(E2E_COMPOSE) up -d --build --remove-orphans; \
		go test ./test/e2e -count=1

compose-up:
	$(COMPOSE) up -d --build

compose-down:
	$(COMPOSE) down -v

migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_DSN)" up

migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_DSN)" down
