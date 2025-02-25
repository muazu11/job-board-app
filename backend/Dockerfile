# Build stage: compile the Go application
FROM golang:1.24.0-alpine AS builder

WORKDIR /backend

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . ./

# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /jobboard-back

# Final stage: create a minimal, secure runtime image
FROM alpine:3.18

# Install CA certificates (if your app needs to make HTTPS requests)
RUN apk add --no-cache ca-certificates

# Create a non-root user to run the application
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the compiled binary from the builder stage
COPY --from=builder /jobboard-back /jobboard-back

# Change ownership of the binary to the non-root user
RUN chown appuser:appgroup /jobboard-back

# Switch to non-root user for enhanced security
USER appuser

# Command to run the binary
CMD ["/jobboard-back"]

