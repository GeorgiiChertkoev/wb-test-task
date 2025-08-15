package main

// test producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/GeorgiiChertkoev/wb-test-task/pkg/models"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
)

type TestProducer struct {
	writer *kafka.Writer
}

func MakeTestProducer() *TestProducer {
	return &TestProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP("kafka:9092"),
			Topic:    "orders",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (t *TestProducer) ProduceRandomOrder() error {

	order := models.MakeRandomOrder()

	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Printf("failed to Marshal order: %v", err)
		return err
	}
	err = t.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: orderJSON,
		},
	)

	if err != nil {
		log.Printf("Failed to send order: %v", err)
		return err
	}
	return nil
}

func (t *TestProducer) Close() {
	t.writer.Close()
}

func main() {
	testKafka()

	fmt.Println("Producer is up")

	producer := MakeTestProducer()
	defer producer.Close()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "producer is fine")
	})
	r.HandleFunc("/produce", func(w http.ResponseWriter, r *http.Request) {
		err := producer.ProduceRandomOrder()
		if err != nil {
			fmt.Fprintf(w, "Failed to produce: %v", err)
		} else {
			fmt.Fprintf(w, "Successfully produced")
		}
	})
	http.ListenAndServe(":8082", r)

}

func testKafka() {
	conn, err := kafka.Dial("tcp", "kafka:9092")
	if err != nil {
		log.Printf("Failed to connect: %v", err)
	}
	defer conn.Close()

	topics, _ := conn.ReadPartitions()
	found := false
	for _, p := range topics {
		if p.Topic == "orders" {
			found = true
			break
		}
	}

	if !found {
		err := conn.CreateTopics(kafka.TopicConfig{
			Topic:             "orders",
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
		if err != nil {
			log.Fatalf("Failed to create topic: %v", err)
		}
	}
}
