
start-identity:
	go run ./internal/services/identity/cmd/server/main.go

identity-sqlc:
	sqlc generate

identity-migrate-create:
	migrate create -ext sql -dir internal/services/identity/internal/infrastructure/db/migrations -seq $(NAME)

identity-migrate-up:
	migrate -path internal/services/identity/internal/infrastructure/db/migrations -database "postgresql://admin:adminpassword@localhost:5432/identity?sslmode=disable" up

identity-migrate-down:
	migrate -path internal/services/identity/internal/infrastructure/db/migrations -database "postgresql://admin:adminpassword@localhost:5432/identity?sslmode=disable" down

docker-local:
	docker compose --env-file=.env.local --file=provisions/docker-compose.noapp.yaml up

.PHONY: start-identity identity-sqlc docker-local