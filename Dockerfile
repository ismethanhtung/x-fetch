# Build stage
FROM golang:1.21-alpine AS builder

# Install git và ca-certificates
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o twitter-backend main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/twitter-backend .

# Expose port
EXPOSE 8080

# Set environment variables để đảm bảo server bind vào 0.0.0.0
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8080
ENV CONTAINER=true

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run
CMD ["./twitter-backend"]

