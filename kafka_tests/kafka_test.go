package kafkatests

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestKafkaReadWrite(t *testing.T) {
	ctx, cancel := context.WithTimeout(Ctx, 10*time.Second)
	defer cancel()
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      Brokers,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	})
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: Brokers,
		Topic:   "user-event",
	})
	var req = EventUser{
		UserID:    20,
		CreatedAt: time.Now(),
	}
	data, err := json.Marshal(req)
	assert.NoError(t, err)
	err = writer.WriteMessages(ctx, kafka.Message{
		Topic: "user-event",
		Value: data,
	})
	assert.NoError(t, err)
	message, err := reader.ReadMessage(ctx)
	assert.NoError(t, err)
	var response EventUser
	err = json.Unmarshal(message.Value, &response)
	assert.NoError(t, err)
	diff := cmp.Diff(req, response)
	assert.Equal(t, "", diff)
	writer.Close()
	reader.Close()
}
