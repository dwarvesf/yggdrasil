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
	Consul  *consul.Client
	address []string
	topic   string
	Reader  *gokafka.Reader
	Writter *gokafka.Writer
}

// New return a new kafka queue
func New(consulClient *consul.Client, topic string) queue.Queue {
	k := &Kafka{Consul: consulClient}
	k.address = k.getBrokers()
	k.topic = topic
	k.Writter = gokafka.NewWriter(gokafka.WriterConfig{
		Brokers: k.address,
		Topic:   k.topic,
	})
	k.Reader = gokafka.NewReader(gokafka.ReaderConfig{
		Brokers: k.address,
		Topic:   k.topic,
	})
	return k
}

// Write msg to kafka queue
func (k *Kafka) Write(values [][]byte) error {
	var msgs []gokafka.Message
	for _, v := range values {
		msgs = append(msgs, gokafka.Message{
			Key:   []byte(k.topic),
			Value: v,
		})
	}
	return k.Writter.WriteMessages(
		context.Background(),
		msgs...,
	)
}

// Read msg from kafka queue§
func (k *Kafka) Read() []byte {
	var b []byte
	for {
		m, err := k.Reader.ReadMessage(context.Background())
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

// Close kafka reader/writer
func (k *Kafka) Close() (error, error) {
	return k.Reader.Close(), k.Writter.Close()
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
