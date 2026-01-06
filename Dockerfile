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

# Stage 2: Runtime
FROM alpine:latest

# Install Python, pip, and required system dependencies for matplotlib
RUN apk add --no-cache \
    python3 \
    py3-pip \
    sqlite \
    freetype \
    fontconfig \
    ttf-dejavu \
    && python3 -m pip install --no-cache-dir --upgrade pip --break-system-packages

# Copy Python requirements and install dependencies
COPY requirements.txt .
RUN pip3 install --no-cache-dir -r requirements.txt --break-system-packages

# Copy the Go binary from builder
COPY --from=builder /app/bot /app/bot

# Copy Python script
COPY main/wrapped/generate_personal_wrapped.py /app/main/wrapped/generate_personal_wrapped.py

# Copy emoji image
COPY utils/Poop_Emoji.webp /app/utils/Poop_Emoji.webp

# Set working directory
WORKDIR /app

# Command to run the application
CMD ["./bot"]