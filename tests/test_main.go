package tests

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	kafkatest "github.com/testcontainers/testcontainers-go/modules/kafka"
)

func SetupContainers() (func(), []string, error) {
	Ctx := context.Background()
	kafkaContainer, err := kafkatest.Run(Ctx, "confluentinc/confluent-local:7.5.0", kafkatest.WithClusterID("test-cluster"))
	if err != nil {
		log.Fatalf("failed to run kafka container, %v", err)
	}

	Brokers, err := kafkaContainer.Brokers(Ctx)
	if err != nil {
		log.Fatalf("failed to get brokers from container, %v", err)
	}
	log.Println("brokers:", Brokers)

	conn, err := kafka.Dial("tcp", Brokers[0])
	if err != nil {
		log.Fatalf("failed to set connection to kafka, %v", err)
	}

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             "user-event",
		ReplicationFactor: 1,
		NumPartitions:     1,
	})
	if err != nil {
		log.Fatalf("failed to create topic, %v", err)
	}

	return func() {
		conn.Close()
		kafkaContainer.Terminate(Ctx)
	}, Brokers, nil
}
