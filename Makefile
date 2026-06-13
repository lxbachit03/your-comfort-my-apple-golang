include .env.local
export

IDENTITY_MIGRATION_DIR := "internal/services/identity/internal/infrastructure/db/migrations"
IDENTITY_DB_CONN_STRING = postgresql://$(IDENTITY_POSTGRES_DB_USER):$(IDENTITY_POSTGRES_DB_PASSWORD)@$(IDENTITY_POSTGRES_DB_HOST):$(IDENTITY_POSTGRES_DB_PORT)/$(IDENTITY_POSTGRES_DB_NAME)?sslmode=$(IDENTITY_POSTGRES_DB_SSL_MODE)

start-identity:
	go run ./internal/services/identity/cmd/server/main.go ENV=$(ENV)

identity-sqlc:
	sqlc generate

identity-migrate-create:
	migrate create -ext sql -dir $(IDENTITY_MIGRATION_DIR) -seq $(NAME)

identity-migrate-up:
	migrate -path $(IDENTITY_MIGRATION_DIR) -database "$(IDENTITY_DB_CONN_STRING)" up

identity-migrate-down:
	migrate -path $(IDENTITY_MIGRATION_DIR) -database "$(IDENTITY_DB_CONN_STRING)" down 1

identity-migrate-force:
	migrate -path $(IDENTITY_MIGRATION_DIR) -database "$(IDENTITY_DB_CONN_STRING)" force $(VERSION)

identity-migrate-drop:
	migrate -path $(IDENTITY_MIGRATION_DIR) -database "$(IDENTITY_DB_CONN_STRING)" drop

identity-migrate-goto:
	migrate -path $(IDENTITY_MIGRATION_DIR) -database "$(IDENTITY_DB_CONN_STRING)" goto $(VERSION)

docker-local:
	docker compose --env-file=.env.local --file=provisions/docker-compose.noapp.yaml up

.PHONY: start-identity identity-sqlc identity-migrate-create identity-migrate-up identity-migrate-down docker-local