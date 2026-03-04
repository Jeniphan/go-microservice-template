# =============================================================================
# Stage 1: Build Stage
# =============================================================================
# Use Go 1.21 for building the application
FROM golang:1.25-alpine AS builder

# Set environment variables for reproducible builds
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
# Build context should be the project root (where Dockerfile is located)
COPY cmd ./cmd
COPY internal ./internal
COPY configs ./configs

# Build the application
# CGO_ENABLED=0 for static binary (no c dependencies)
# -a forces rebuild of all packages
# -installsuffix cgo omits the stdlib's cgo packages from the binary
# -o main specifies output binary name
RUN go build -a -installsuffix cgo -ldflags="-s -w" -o main ./cmd/main.go

# =============================================================================
# Stage 2: Production Stage
# =============================================================================
# Use Alpine Linux for minimal production image
FROM alpine:latest

# Install CA certificates for HTTPS connections and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Create non-root user for security
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

# Copy binary from builder stage
COPY --from=builder /app/main .

# Change ownership of the binary to non-root user
RUN chown appuser:appgroup /app/main

# Expose application port
EXPOSE 8080

# Switch to non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/access-control/healthcheck || exit 1

# Run the application
CMD ["./main"]
