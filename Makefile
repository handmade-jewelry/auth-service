OPENAPI_RESOURCE = api/resourceapi.yaml
OPENAPI_AUTH = api/authapi.yaml

GEN_RESOURCE = pkg/api/resource/resourceapi.gen.go
GEN_AUTH = pkg/api/auth/authapi.gen.go

PACKAGE_RESOURCE = resource
PACKAGE_AUTH = auth

OAPI_CODEGEN = oapi-codegen

MIGRATIONS_DIR = ./migrations

DB_USER = auth_service_user
DB_NAME = auth_service_db
DB_SSLMODE = disable

.PHONY: generate-oapi new-migration migration-up migration-down

# Generate resourceapi.yaml
generate-resource:
	$(OAPI_CODEGEN) -generate types,chi-server -o $(GEN_RESOURCE) -package $(PACKAGE_RESOURCE) $(OPENAPI_RESOURCE)

# Generate authapi.yaml
generate-auth:
	$(OAPI_CODEGEN) -generate types,chi-server -o $(GEN_AUTH) -package $(PACKAGE_AUTH) $(OPENAPI_AUTH)

# Make new migration sql
new-migration:
ifndef NAME
	$(error Usage: make new-migration NAME=your_migration_name)
endif
	goose -dir $(MIGRATIONS_DIR) create $(NAME) sql

# Apply migrations
migration-up:
	goose -dir $(MIGRATIONS_DIR) postgres "user=$(DB_USER) dbname=$(DB_NAME) sslmode=$(DB_SSLMODE)" up

# Rollback migrations
migration-down:
	goose -dir $(MIGRATIONS_DIR) postgres "user=$(DB_USER) dbname=$(DB_NAME) sslmode=$(DB_SSLMODE)" down