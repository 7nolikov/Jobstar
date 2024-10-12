# Use official Golang image as base
FROM golang:1.23-alpine

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main ./cmd/app

# Expose port
EXPOSE 8080

# Start the application
CMD ["./main"]
