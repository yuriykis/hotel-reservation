BINARY_NAME=hotel-reservation-api
ARGS="--listenAddr=:5001"

build:
	@go build -o bin/$(BINARY_NAME) -v

run: build
	@./bin/$(BINARY_NAME) $(ARGS)

seed:
	@go run scripts/seed.go

test:
	@go test -v ./...

test-race:
	@go test -v ./... --race
	