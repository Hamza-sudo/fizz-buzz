APP_NAME := fizz-buzz

.PHONY: run test build docker-build

run:
	go run ./cmd/server

test:
	go test ./...

build:
	go build -o bin/$(APP_NAME) ./cmd/server

docker-build:
	docker build -t $(APP_NAME):latest .
