# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

COPY . .

# Copy module files and download dependencies
COPY src/go.mod src/go.sum ./src/
WORKDIR /app/src
RUN go mod download

# Copy the database file
# Copy the application source code


# Build the Go binary
RUN go build -o bot main/bot.go

# Stage 2: Run
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy the built binary from the build stage
COPY --from=builder /app/src/bot .

# Expose the port (optional for documentation)
EXPOSE 8080

# Set environment variables for Fly.io
ENV PORT=8080

# Command to run the application
CMD ["./bot"]