FROM golang:1.21.1-alpine3.17

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# This container exposes port 5001 to the outside world
EXPOSE 5001

# Run the executable
CMD ["./main"]