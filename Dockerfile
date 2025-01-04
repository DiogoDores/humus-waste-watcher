# Stage 1: Build
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache sqlite

WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -o bot main/bot.go

# Command to run the application
CMD ["./bot"]