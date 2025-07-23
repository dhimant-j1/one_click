# Makefile for Insurance Management API

.PHONY: help build up down logs clean restart shell db-shell

# Default target
help:
	@echo "Insurance Management API - Docker Commands"
	@echo ""
	@echo "Available commands:"
	@echo "  build     - Build all services"
	@echo "  up        - Start all services"
	@echo "  down      - Stop all services"
	@echo "  logs      - View logs from all services"
	@echo "  logs-app  - View application logs"
	@echo "  logs-db   - View database logs"
	@echo "  restart   - Restart all services"
	@echo "  clean     - Stop and remove all containers, networks, and volumes"
	@echo "  shell     - Open shell in application container"
	@echo "  db-shell  - Open PostgreSQL shell"
	@echo "  status    - Show status of all services"

# Build all services
build:
	docker-compose build

# Start all services
up:
	docker-compose up -d

# Start all services with build
up-build:
	docker-compose up --build -d

# Stop all services
down:
	docker-compose down

# View logs from all services
logs:
	docker-compose logs -f

# View application logs
logs-app:
	docker-compose logs -f app

# View database logs
logs-db:
	docker-compose logs -f postgres

# Restart all services
restart:
	docker-compose restart

# Clean everything (containers, networks, volumes)
clean:
	docker-compose down -v --remove-orphans
	docker system prune -f

# Open shell in application container
shell:
	docker-compose exec app sh

# Open PostgreSQL shell
db-shell:
	docker-compose exec postgres psql -U postgres -d insurance

# Show status of all services
status:
	docker-compose ps

# Development mode (with live reload)
dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

# Production deployment
prod:
	docker-compose up --build -d
	@echo "Services started in production mode"
	@echo "Application: http://localhost:8081"
	@echo "Database: localhost:5432"
