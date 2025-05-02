# todo fix

# Переменные
OPENAPI_SPEC = api/openapi.yaml
GEN_FILE = pkg/api/openapi.gen.go
PACKAGE = api
# Команды
OAPI_CODEGEN = oapi-codegen
MIGRATIONS_DIR = ./migrations

.PHONY: generate-oapi new-migration migration-up migration-down

generate-oapi:
	$(OAPI_CODEGEN) -generate types,chi-server -o $(GEN_FILE) -package $(PACKAGE) $(OPENAPI_SPEC)

# Make new migration sql
#ifndef NAME
#	$(error Usage: make new-migration NAME=your_migration_name)
#endif

new-migration:
	goose -dir $(MIGRATIONS_DIR) create $(NAME) sql

migration-up:
	goose -dir $(MIGRATIONS_DIR) postgres "user=auth_service_user dbname=auth_service_db sslmode=disable" up

migration-down:
	goose -dir $(MIGRATIONS_DIR) postgres "user=auth_service_user dbname=auth_service_db sslmode=disable" down

