.PHONY: build run migrate seed docker-up docker-down docker-build docker-rebuild-backend docker-logs docker-restart docker-migrate docker-seed docker-shell test test-coverage

# Docker Compose v2: `docker compose`. Override if needed: make docker-seed COMPOSE=docker-compose
COMPOSE ?= docker compose

# Host cannot resolve Docker DNS names — use docker-seed / docker-migrate, or DB_HOST=127.0.0.1 DB_PORT=5433
check-db-host-for-local:
	@if [ -f .env ]; then \
		if grep -qE '^[[:space:]]*DB_HOST[[:space:]]*=[[:space:]]*stempo_backend[[:space:]]*$$' .env; then \
			echo "DB_HOST=stempo_backend is wrong: that is the backend API container name, not PostgreSQL."; \
			echo "In .env set: DB_HOST=postgres  (the docker-compose service name for the database)."; \
			echo "Then: docker compose up -d --force-recreate backend  (to reload env)"; \
			echo "Seeds: make docker-seed"; \
			exit 1; \
		elif grep -qE '^[[:space:]]*DB_HOST[[:space:]]*=[[:space:]]*postgres[[:space:]]*$$' .env; then \
			echo "DB_HOST=postgres only works inside Docker. On the server run:"; \
			echo "  $(COMPOSE) exec backend ./server seed"; \
			echo "  (or: make docker-seed)"; \
			echo "Or from host with published DB port:"; \
			echo "  DB_HOST=127.0.0.1 DB_PORT=5433 make seed"; \
			exit 1; \
		fi; \
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

# Use when backend still runs old ./server after seeder changes (no --no-cache = stale layers)
docker-rebuild-backend:
	$(COMPOSE) build --no-cache backend

docker-logs:
	$(COMPOSE) logs -f backend

docker-restart:
	$(COMPOSE) restart backend

docker-migrate:
	$(COMPOSE) exec backend ./server migrate

docker-seed:
	@echo "Tip: if seed errors look like old code (SQL without owner_id), run: make docker-rebuild-backend && $(COMPOSE) up -d backend"
	$(COMPOSE) exec -T backend ./server seed

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
	@echo "  make docker-rebuild-backend - Rebuild backend image with --no-cache (after seeder/API code changes)"
	@echo "  make docker-restart     - Restart backend container"
	@echo "  make docker-migrate     - Run migrations in container (use on server)"
	@echo "  make docker-seed        - Run seeders in container (use on server)"
	@echo "  make migrate / seed     - Run on host; needs DB_HOST=127.0.0.1 DB_PORT=5433 if DB is in Docker"
	@echo "  make docker-shell       - Open shell in backend container"
	@echo "  make docker-reset       - Reset and setup database (down, up, migrate, seed)"
	@echo "  make test               - Run tests"
	@echo "  make test-coverage      - Run tests with coverage"
