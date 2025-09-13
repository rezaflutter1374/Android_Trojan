FROM golang:1.16-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o c2server ./cmd/server

# Use a minimal alpine image for the final container
FROM alpine:latest

# Add ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/c2server .

# Create necessary directories
RUN mkdir -p /root/logs

# Expose the port the server listens on
EXPOSE 8080

# Command to run the executable
CMD ["./c2server"]