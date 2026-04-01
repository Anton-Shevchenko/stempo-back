.PHONY: build run migrate seed docker-up docker-down docker-build docker-logs docker-restart docker-migrate docker-seed docker-shell test test-coverage

# Docker Compose v2: `docker compose`. Override if needed: make docker-seed COMPOSE=docker-compose
COMPOSE ?= docker compose

# Host cannot resolve service name "postgres" — use docker-seed / docker-migrate, or DB_HOST=127.0.0.1 DB_PORT=5433
check-db-host-for-local:
	@if [ -f .env ] && grep -qE '^[[:space:]]*DB_HOST[[:space:]]*=[[:space:]]*postgres[[:space:]]*$$' .env; then \
		echo "DB_HOST=postgres only works inside Docker. On the server run:"; \
		echo "  $(COMPOSE) exec backend ./server seed"; \
		echo "  (or: make docker-seed)"; \
		echo "Or from this machine with published DB port:"; \
		echo "  DB_HOST=127.0.0.1 DB_PORT=5433 make seed"; \
		exit 1; \
	fi

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

migrate: check-db-host-for-local
	go run ./cmd/server migrate

seed: check-db-host-for-local
	go run ./cmd/server seed

docker-up:
	$(COMPOSE) up -d

docker-down:
	$(COMPOSE) down

docker-build:
	$(COMPOSE) build

docker-logs:
	$(COMPOSE) logs -f backend

docker-restart:
	$(COMPOSE) restart backend

docker-migrate:
	$(COMPOSE) exec backend ./server migrate

docker-seed:
	@$(COMPOSE) exec -T backend ./server seed 2>&1 | grep -v "SELECT\|FROM\|WHERE\|JOIN\|INFORMATION_SCHEMA\|pg_\|rows:\|bind:" | grep -E "(Starting|Seeded|completed|error|Error)" || true

docker-shell:
	$(COMPOSE) exec backend sh

docker-up-build:
	$(COMPOSE) up -d --build

docker-reset: docker-down docker-up-build
	@sleep 5
	@$(MAKE) docker-migrate
	@$(MAKE) docker-seed

test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

help:
	@echo "Available commands:"
	@echo "  make docker-up          - Start Docker containers"
	@echo "  make docker-down        - Stop Docker containers"
	@echo "  make docker-build       - Build Docker images"
	@echo "  make docker-restart     - Restart backend container"
	@echo "  make docker-migrate     - Run migrations in container (use on server)"
	@echo "  make docker-seed        - Run seeders in container (use on server)"
	@echo "  make migrate / seed     - Run on host; needs DB_HOST=127.0.0.1 DB_PORT=5433 if DB is in Docker"
	@echo "  make docker-shell       - Open shell in backend container"
	@echo "  make docker-reset       - Reset and setup database (down, up, migrate, seed)"
	@echo "  make test               - Run tests"
	@echo "  make test-coverage      - Run tests with coverage"
