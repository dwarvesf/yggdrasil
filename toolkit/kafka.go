package toolkit

import (
	"encoding/json"
	"fmt"

	consul "github.com/hashicorp/consul/api"
	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	Consul *consul.Client
	Writer *kafka.Writer
	Reader *kafka.Reader
}

func (k *Kafka) New(topic string) error {
	var kafkaInfo struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	}
	v, _ := GetConsulValueFromKey(k.Consul, "kafka")
	if err := json.Unmarshal([]byte(v), &kafkaInfo); err != nil {
		return err
	}

	address := fmt.Sprintf("%v:%v", kafkaInfo.Address, kafkaInfo.Port)
	k.Writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{address},
		Topic:   topic,
	})
	k.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{address},
		Topic:   topic,
	})

	return nil
}
