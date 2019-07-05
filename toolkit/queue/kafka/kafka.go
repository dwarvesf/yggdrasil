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

// Queue implementation
type kafka struct {
	ConsulClient *consul.Client
}

// New return a new kafka
func New(consulClient *consul.Client) queue.Queue {
	return &kafka{
		ConsulClient: consulClient,
	}
}

// Message writer implementation
type kafkaWriter struct {
	Writer *gokafka.Writer
}

// NewWriter return a new writer for kafka
func (k *kafka) NewWriter(topic string) queue.Writer {
	return &kafkaWriter{
		Writer: gokafka.NewWriter(gokafka.WriterConfig{
			Brokers: k.getBrokers(),
			Topic:   topic,
		}),
	}
}

// Close writer
func (w *kafkaWriter) Close() error {
	return w.Writer.Close()
}

func (w *kafkaWriter) Write(key string, value []byte) error {
	return w.Writer.WriteMessages(
		context.Background(),
		gokafka.Message{
			Key:   []byte(key),
			Value: value,
		},
	)
}

// Message reader implementation
type kafkaReader struct {
	Reader *gokafka.Reader
}

// NewReader return a new reader for kafka service
func (k *kafka) NewReader(topic string) queue.Reader {
	return &kafkaReader{
		Reader: gokafka.NewReader(gokafka.ReaderConfig{
			Brokers: k.getBrokers(),
			Topic:   topic,
			GroupID: topic,
		}),
	}
}

// Close reader
func (r *kafkaReader) Close() error {
	return r.Reader.Close()
}

func (r *kafkaReader) Read() ([]byte, error) {
	message, err := r.Reader.ReadMessage(context.Background())
	if err != nil {
		return nil, err
	}
	return message.Value, nil
}

func (k *kafka) getBrokers() []string {
	var kafkaInfo struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	}

	v, _ := toolkit.GetConsulValueFromKey(k.ConsulClient, "kafka")
	if err := json.Unmarshal([]byte(v), &kafkaInfo); err != nil {
		return []string{}
	}

	return []string{fmt.Sprintf("%v:%v", kafkaInfo.Address, kafkaInfo.Port)}
}
