# syntax=docker/dockerfile:1

# ---- build stage ----
# Pick a real Go version (1.25 is current at time of writing)
FROM golang:1.25 AS builder

# Build for the docker target architecture, disable CGO for static binary if possible
ENV CGO_ENABLED=0 GOOS=linux
# TARGETARCH is provided by BuildKit; GOARCH inherits automatically if not overridden

WORKDIR /app

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest and build
COPY . .

RUN go build -o /app/ordersystem .

# ---- minimal runtime image ----
FROM alpine:3.20

# Create a non-root user
RUN adduser -D -H -s /sbin/nologin appuser

# Copy the statically-linked binary
COPY --from=builder /app/ordersystem /ordersystem

# Make sure it’s executable
RUN chmod +x /ordersystem

USER appuser

# If your service listens on 3000, expose it (purely documentation)
EXPOSE 3000

ENTRYPOINT ["/ordersystem"]
