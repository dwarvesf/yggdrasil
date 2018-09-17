package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dwarvesf/yggdrasil/toolkit"
	"github.com/dwarvesf/yggdrasil/toolkit/queue"
	consul "github.com/hashicorp/consul/api"
	gokafka "github.com/segmentio/kafka-go"
)

// Kafka implementation
type Kafka struct {
	Consul *consul.Client
	Writer *gokafka.Writer
	Reader *gokafka.Reader
}

// New return a new kafka queue
func New(consulClient *consul.Client) queue.Queue {
	return &Kafka{
		Consul: consulClient,
	}
}

// Write return a new writer for kafka Queue
func (k *Kafka) Write(topic string, value []byte) error {
	w := gokafka.NewWriter(gokafka.WriterConfig{
		Brokers: k.getBrokers(),
		Topic:   topic,
	})
	defer w.Close()
	return w.WriteMessages(
		context.Background(),
		gokafka.Message{
			Key:   []byte(topic),
			Value: value,
		},
	)
}

// Read return a new reader for kafka Queue
func (k *Kafka) Read(topic string) []byte {
	r := gokafka.NewReader(gokafka.ReaderConfig{
		Brokers: k.getBrokers(),
		Topic:   topic,
	})
	defer r.Close()

	var b []byte
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println(fmt.Sprintf("err: %v", err))
			break
		}

		if string(m.Value) == "" {
			continue
		}

		b = m.Value
		break
	}
	return b
}

func (k *Kafka) getBrokers() []string {
	var kafkaInfo struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	}

	v, _ := toolkit.GetConsulValueFromKey(k.Consul, "kafka")
	if err := json.Unmarshal([]byte(v), &kafkaInfo); err != nil {
		return []string{}
	}

	return []string{fmt.Sprintf("%v:%v", kafkaInfo.Address, kafkaInfo.Port)}
}
