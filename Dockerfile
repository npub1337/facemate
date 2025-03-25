# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc g++ make cmake pkgconfig git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o face-recognition-processor ./cmd/api

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache dlib-dev

# Copy the binary from builder
COPY --from=builder /app/face-recognition-processor /app/face-recognition-processor
COPY --from=builder /app/models /app/models

# Create data directory
RUN mkdir -p /app/data

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 8080

# Run the application
CMD ["./face-recognition-processor"]
