# syntax=docker/dockerfile:1

FROM golang:1.21.2-alpine

# Set destination for COPY
WORKDIR /backend

# Download Go modules
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the source code
COPY backend/ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /jobboard-back

# Clean
RUN rm -rf backend

# Run
CMD ["/jobboard-back"]
