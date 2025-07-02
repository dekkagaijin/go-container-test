FROM golang:1.24.1-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY main.go .

# Build the Weather service
RUN go build -a -o main .

# Use a nice, vulnerable-ridden base image.
FROM ubuntu:jammy-20211029

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the port your application listens on (if applicable)
EXPOSE 8080

# Command to run the application when the container starts
CMD ["./main"]