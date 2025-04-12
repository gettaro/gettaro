# Database configuration
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= ems_dev

# Migration configuration
MIGRATION_DIR = database/migrations
MIGRATION_DSN = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATE_BIN = $(shell go env GOPATH)/bin/migrate

.PHONY: migrate-up migrate-down migrate-create migrate-force migrate-status db-create db-drop

# Database operations
db-create:
	@echo "Creating database $(DB_NAME)..."
	PGPASSWORD=$(DB_PASSWORD) createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

db-drop:
	@echo "Dropping database $(DB_NAME)..."
	PGPASSWORD=$(DB_PASSWORD) dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

# Migration operations
migrate-up:
	@echo "Running migrations up..."
	@$(MIGRATE_BIN) -path $(MIGRATION_DIR) -database "$(MIGRATION_DSN)" up

migrate-down:
	@echo "Running migrations down..."
	@$(MIGRATE_BIN) -path $(MIGRATION_DIR) -database "$(MIGRATION_DSN)" down

migrate-down-1:
	@echo "Running last migration down..."
	@$(MIGRATE_BIN) -path $(MIGRATION_DIR) -database "$(MIGRATION_DSN)" down 1

migrate-create:
	@echo "Creating new migration..."
	@read -p "Enter migration name: " name; \
	$(MIGRATE_BIN) create -ext sql -dir $(MIGRATION_DIR) -seq $$name

migrate-force:
	@echo "Forcing migration version..."
	@read -p "Enter version to force: " version; \
	$(MIGRATE_BIN) -path $(MIGRATION_DIR) -database "$(MIGRATION_DSN)" force $$version

migrate-status:
	@echo "Checking migration status..."
	@$(MIGRATE_BIN) -path $(MIGRATION_DIR) -database "$(MIGRATION_DSN)" version

# Development setup
setup: db-create migrate-up
	@echo "Development environment setup complete!"

# Cleanup
clean: db-drop
	@echo "Database cleanup complete!" 