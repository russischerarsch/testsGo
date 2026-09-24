package producerkafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

func CreateWriter(brokers []string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		BatchTimeout: 10 * time.Millisecond,
	}
}
