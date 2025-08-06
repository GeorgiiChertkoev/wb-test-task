# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /order-service ./cmd/server

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /order-service .
COPY web ./web
# COPY config ./config

EXPOSE 8080
CMD ["./order-service"]