package main

// test producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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

func (t *TestProducer) ProduceRandomOrder() (*models.Order, error) {

	order := models.MakeRandomOrder()

	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Printf("failed to Marshal order: %v", err)
		return nil, err
	}
	err = t.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: orderJSON,
		},
	)

	if err != nil {
		log.Printf("Failed to send order: %v", err)
		return nil, err
	}
	return order, nil
}

func (t *TestProducer) Close() {
	t.writer.Close()
}

func main() {
	// пытаемя подключиться к кафке и создать топик
	for i := 0; i < 5; i++ {
		err := testKafka()
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	fmt.Println("Producer is up")

	producer := MakeTestProducer()
	defer producer.Close()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "producer is fine")
	})
	r.HandleFunc("/produce", func(w http.ResponseWriter, r *http.Request) {
		order, err := producer.ProduceRandomOrder()
		if err != nil {
			fmt.Fprintf(w, "Failed to produce: %v", err)
		} else {
			fmt.Fprintf(w, "Successfully produced order with id: %v", order.OrderUID)
		}
	})
	http.ListenAndServe(":8082", r)

}

func testKafka() error {
	conn, err := kafka.Dial("tcp", "kafka:9092")
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return err
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
			log.Printf("Failed to create topic: %v", err)
			return err
		}
	}
	return nil
}
