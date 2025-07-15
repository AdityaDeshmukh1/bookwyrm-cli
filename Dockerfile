# Stage 1 — Build
FROM golang:1.24.4-alpine AS builder

# Set working directory inside container
WORKDIR /app

# Copy Go mod files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy entire project (source code, internal packages, cmd, etc.)
COPY . .

# Build the CLI binary
RUN go build -o bookwyrm-cli main.go

# Stage 2 — Minimal runtime image
FROM alpine:latest

# Working dir inside runtime container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/bookwyrm-cli .

# Set permissions just in case
RUN chmod +x bookwyrm-cli

# Default command when container starts
ENTRYPOINT ["./bookwyrm-cli"]

