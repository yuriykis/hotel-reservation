BINARY_NAME=hotel-reservation-api
ARGS="--listenAddr=:5001"

build:
	@go build -o bin/$(BINARY_NAME) -v

run: build
	@./bin/$(BINARY_NAME) $(ARGS)

seed:
	@go run scripts/seed.go

test:
	@go test -v ./... -count=1

test-race:
	@go test -v ./... --race
	
docker:
	echo "Building docker image"
	docker build -t api .
	echo "Running docker image"
	docker run -p 5001:5001 api