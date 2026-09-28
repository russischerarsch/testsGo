package clients

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	Writer *kafka.Writer
	Reader *kafka.Reader
}

func CreateKafkaClient(brokers []string, topic string) *KafkaClient {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	})
	reader := kafka.NewReader(kafka.ReaderConfig{
		Topic:   topic,
		Brokers: brokers,
	})
	return &KafkaClient{Writer: writer, Reader: reader}
}
func (c *KafkaClient) Close() {
	c.Reader.Close()
	c.Writer.Close()
}
func (c *KafkaClient) Producer(ctx context.Context, value []byte) error {
	return c.Writer.WriteMessages(ctx, kafka.Message{Value: value})
}
func (c *KafkaClient) Consume(ctx context.Context) (kafka.Message, error) {
	return c.Reader.ReadMessage(ctx)
}
