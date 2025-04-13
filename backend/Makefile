run:
	go run cmd/server/main.go

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@db:5432/chatapp?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@db:5432/chatapp?sslmode=disable" down

docker-up:
	docker compose up --build

docker-down:
	docker compose down
