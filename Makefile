
start-identity:
	go run ./internal/services/identity/cmd/server/main.go

docker-local:
	docker compose --env-file=.env.local --file=provisions/docker-compose.noapp.yaml up

