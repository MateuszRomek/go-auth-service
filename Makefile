# Load environment variables from .env file
include .env
# Export only specific variables as needed
export POSTGRES_USER POSTGRES_PASSWORD POSTGRES_DB POSTGRES_HOST POSTGRES_PORT MIGRATIONS_DIR

# Database connection string
DB_URL=user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) sslmode=disable

# Goose binary
GOOSE=goose

# Goose commands
.PHONY: migrate-up migrate-down migrate-reset migrate-status migrate-create

migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

migrate-reset:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" reset

migrate-status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

migrate-create:
	@read -p "Enter migration name: " name; \
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $${name} sql
