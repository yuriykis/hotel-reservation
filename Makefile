BINARY_NAME=hotel-reservation
build:
	@go build -o bin/$(BINARY_NAME) -v

run: build
	@./bin/$(BINARY_NAME)

test:
	@go test -v ./...

test-race:
	@go test -v ./... --race
	