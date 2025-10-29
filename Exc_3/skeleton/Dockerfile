# building exe on machine and copying to minimal image 

# syntax=docker/dockerfile:1
FROM golang:1.25 AS builder 

# container dir
WORKDIR /app
COPY . .
RUN scripts/build-application.sh

# minimal image
FROM alpine:latest
COPY --from=builder /app/ordersystem /app/ordersystem
CMD ["/app/ordersystem"]
