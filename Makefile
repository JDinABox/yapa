build: generate
	go build -ldflags="-s -w" -o ./cmd/yapa/yapa.out ./cmd/yapa

generate:
	go generate ./...
	cd ./internal/sqlc && sqlc generate

make run: generate
	go run ./cmd/yapa