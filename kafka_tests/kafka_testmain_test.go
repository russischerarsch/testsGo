package kafkatests

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	kafkatest "github.com/testcontainers/testcontainers-go/modules/kafka"
)

type EventUser struct {
	EventID   int64     `json:"event_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

var Brokers []string
var Ctx = context.Background()

func TestMain(m *testing.M) {

	container, err := kafkatest.Run(Ctx, "confluentinc/confluent-local:7.5.0", kafkatest.WithClusterID("test-cluster"))
	if err != nil {
		log.Fatalf("failed to run test kafka cluster, %v", err)
	}
	defer container.Terminate(Ctx)
	brokers, err := container.Brokers(Ctx)
	Brokers = brokers
	if err != nil {
		log.Fatalf("failed to get brokers from container, %v", err)
	}
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		log.Fatalf("failed to create a conn, %v", err)
	}
	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             "user-event",
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		log.Fatalf("failed to create topic, %v", err)
	}
	m.Run()
}
