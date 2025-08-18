# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/server/main.go

RUN go build -o /order-service ./cmd/server

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /order-service .
COPY web ./web
COPY docs ./docs

EXPOSE 8081
CMD ["./order-service"]