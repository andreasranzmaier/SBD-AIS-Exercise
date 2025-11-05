# syntax=docker/dockerfile:1
# build runner 
FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN chmod +x scripts/build-application.sh && ./scripts/build-application.sh

# minimal image for ordersystem
FROM alpine:latest
COPY --from=builder /app/ordersystem /ordersystem
USER nonroot

ENTRYPOINT ["/ordersystem"]
