package producerkafka

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

func LoadKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		Brokers: []string{"localhost:9092"},
	}
}
