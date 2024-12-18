# Start from an official Golang image
FROM golang:1.23.4-alpine3.20

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory
WORKDIR /app


# Copy go.mod and go.sum files first for dependency caching
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN go build -o server

# Expose the application's port
EXPOSE 8080

# Run the binary
CMD ["./server"]
