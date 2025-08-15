module order_producer

go 1.23.0

toolchain go1.23.11

require (
	github.com/gorilla/mux v1.8.1
	github.com/segmentio/kafka-go v0.4.48
	github.com/GeorgiiChertkoev/wb-test-task v0.0.0-00010101000000-000000000000
)

require (
	github.com/brianvoe/gofakeit/v7 v7.3.0 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
)

replace github.com/GeorgiiChertkoev/wb-test-task => ../../
