FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ddalab-manager .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates docker-cli docker-compose

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/ddalab-manager .

# Set environment variable for socket path
ENV SOCKET_PATH=/run/guest/backend.sock

# Run the binary
CMD ["./ddalab-manager"]