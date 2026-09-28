package tests

import (
	"context"
	"log"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go"
	kafkatest "github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/wait"
)

var Ctx = context.Background()
var Dsn string
var Brokers []string

func TestMain(m *testing.M) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASS": "test",
			"POSTGRES_USER": "test",
			"POSTGRES_DB":   "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}
	pgContainer, err := testcontainers.GenericContainer(Ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	defer pgContainer.Terminate(Ctx)
	if err != nil {
		log.Fatalf("failed to start container, %v", err)
	}
	host, err := pgContainer.Host(Ctx)
	if err != nil {
		log.Fatalf("failed to get host, %v", err)
	}
	port, err := pgContainer.MappedPort(Ctx, "5432")
	if err != nil {
		log.Fatalf("failed to get port, %v", err)
	}
	Dsn = "postgres://test:test@" + host + ":" + port.Port() + "/testdb?sslmode=disable"
	kafkaContainer, err := kafkatest.Run(Ctx, "confluentinc/confluent-local:7.5.0", kafkatest.WithClusterID("test-cluster"))
	if err != nil {
		log.Fatalf("failed to run kafka container, %v", err)
	}
	defer kafkaContainer.Terminate(Ctx)
	Brokers, err := kafkaContainer.Brokers(Ctx)
	if err != nil {
		log.Fatalf("failed to get brokers from container, %v", err)
	}
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
	m.Run()
}
