# syntax=docker/dockerfile:1

FROM golang:1.21.2

# Set destination for COPY
WORKDIR /backend

# Download Go modules
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the source code
COPY backend/*.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /jobboard-back

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

# Run
CMD ["/jobboard-back"]
