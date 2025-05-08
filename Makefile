OPENAPI_SPEC = api/openapi.yaml
GEN_FILE = pkg/api/openapi.gen.go
OAPI_CODEGEN = oapi-codegen
PACKAGE = api
MIGRATIONS_DIR = ./migrations

DB_USER = auth_service_user
DB_NAME = auth_service_db
DB_SSLMODE = disable

.PHONY: generate-oapi new-migration migration-up migration-down

generate-oapi:
	$(OAPI_CODEGEN) -generate types,chi-server -o $(GEN_FILE) -package $(PACKAGE) $(OPENAPI_SPEC)

# Make new migration sql
ifndef NAME
	$(error Usage: make new-migration NAME=your_migration_name)
endif

new-migration:
	goose -dir $(MIGRATIONS_DIR) create $(NAME) sql

# Применить миграции
migration-up:
	goose -dir $(MIGRATIONS_DIR) postgres "user=$(DB_USER) dbname=$(DB_NAME) sslmode=$(DB_SSLMODE)" up

# Откатить миграции
migration-down:
	goose -dir $(MIGRATIONS_DIR) postgres "user=$(DB_USER) dbname=$(DB_NAME) sslmode=$(DB_SSLMODE)" down