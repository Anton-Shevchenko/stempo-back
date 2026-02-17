.PHONY: build run migrate seed docker-up docker-down docker-build docker-logs docker-restart docker-migrate docker-seed docker-shell test test-coverage

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

migrate:
	go run ./cmd/server migrate

seed:
	go run ./cmd/server seed

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

docker-logs:
	docker-compose logs -f backend

docker-restart:
	docker-compose restart backend

docker-migrate:
	docker-compose exec backend ./server migrate

docker-seed:
	@docker-compose exec -T backend ./server seed 2>&1 | grep -v "SELECT\|FROM\|WHERE\|JOIN\|INFORMATION_SCHEMA\|pg_\|rows:\|bind:" | grep -E "(Starting|Seeded|completed|error|Error)" || true

docker-shell:
	docker-compose exec backend sh

docker-up-build:
	docker-compose up -d --build

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
	@echo "  make docker-migrate    - Run migrations in container"
	@echo "  make docker-seed       - Run seeders in container"
	@echo "  make docker-shell       - Open shell in backend container"
	@echo "  make docker-reset      - Reset and setup database (down, up, migrate, seed)"
	@echo "  make test               - Run tests"
	@echo "  make test-coverage      - Run tests with coverage"
