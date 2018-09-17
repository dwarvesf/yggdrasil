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
	Reader  *gokafka.Reader
	Writter *gokafka.Writer
}

// New return a new kafka queue
func New(consulClient *consul.Client) queue.Queue {
	k := &Kafka{
		Consul: consulClient,
	}
	k.address = k.getBrokers()
	return k
}

// NewWriter return kafka reader
func (k *Kafka) NewWriter(topic string) {
	k.Writter = gokafka.NewWriter(gokafka.WriterConfig{
		Brokers: k.address,
		Topic:   topic,
	})
}

// Write msg to kafka queue
func (k *Kafka) Write(topic string, values [][]byte) error {
	if k.Writter == nil {
		k.NewWriter(topic)
	}

	var msgs []gokafka.Message
	for _, v := range values {
		msgs = append(msgs, gokafka.Message{
			Key:   []byte(topic),
			Value: v,
		})
	}
	return k.Writter.WriteMessages(
		context.Background(),
		msgs...,
	)
}

// NewReader return kafka reader
func (k *Kafka) NewReader(topic string) {
	k.Reader = gokafka.NewReader(gokafka.ReaderConfig{
		Brokers: k.address,
		Topic:   topic,
	})
}

// Read msg from kafka queue§
func (k *Kafka) Read(topic string) []byte {
	if k.Reader == nil {
		k.NewReader(topic)
	}

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
func (k *Kafka) Close() error {
	if k.Reader != nil {
		return k.Reader.Close()
	}
	if k.Writter != nil {
		return k.Writter.Close()
	}
	return nil
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
