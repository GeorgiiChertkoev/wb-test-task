# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY scripts/kafka_producer/go.mod .
COPY scripts/kafka_producer/go.sum .
COPY . .

WORKDIR /app/scripts/kafka_producer
RUN go build -o /kafka-producer

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /kafka-producer .

EXPOSE 8082
CMD ["./kafka-producer"]


