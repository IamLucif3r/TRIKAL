.PHONY: run build dev
build:
	go build -o bin/trikal ./cmd/trikal

run:
	TRIKAL_CONFIG=./configs/config.example.yaml go run ./cmd/trikal

dev:
	docker compose up --build
