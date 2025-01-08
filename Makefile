# Define project name
PROJECT_NAME := student-api

# Define image name
IMAGE_NAME := $(PROJECT_NAME)-api

# Define database container name
DB_CONTAINER := $(PROJECT_NAME)_db_1

# Define database URL
DATABASE_URL := ${DATABASE_URL}

# Docker Compose file
COMPOSE_FILE := docker-compose.yaml

# Default target
.DEFAULT_GOAL := up

# Target to start the DB container
db-up:
	docker-compose -f $(COMPOSE_FILE) up -d db

# Target to build the REST API Docker image
build:
	docker-compose -f $(COMPOSE_FILE) build --no-cache app1 app2

# Target to build nginx image
nginx-build:
	docker-compose -f $(COMPOSE_FILE) build --no-cache nginx

# Target to run the REST API Docker container (including dependencies and migrations)
up: db-up build
	docker-compose -f $(COMPOSE_FILE) up -d app1 app2 nginx

# Target to stop all containers
down:
	docker-compose -f $(COMPOSE_FILE) down

# Target to stop and remove all containers and volumes
down-v:
	docker-compose -f $(COMPOSE_FILE) down -v

# Target to check if DB is running
db-check:
	@docker ps --filter "name=$(DB_CONTAINER)" --format '{{.State}}' | grep -q 'running'
	@if [ $$? -eq 0 ]; then \
		echo "DB is running"; \
	else \
		echo "DB is NOT running"; \
	fi

logs:
	docker-compose logs


logs-%: # Pattern rule for service logs
	docker-compose logs $*

# Target to check if migrations are applied (basic check - improve if needed)
migrations-check: db-up
	@docker-compose -f $(COMPOSE_FILE) exec db psql -U devasheesh -d student_db -c '\dt goose_db_version' > /dev/null 2>&1
	@if [ $$? -eq 0 ]; then \
 		echo "Migrations are applied (goose_db_version table exists)"; \
	else \
		echo "Migrations are NOT applied (goose_db_version table does not exist)"; \
	fi

.PHONY: db-up build nginx-build up down down-v db-check logs migrations-check