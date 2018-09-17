package message

import (
	"context"
	"fmt"

	"github.com/dwarvesf/yggdrasil/toolkit"
	consul "github.com/hashicorp/consul/api"
	kafka "github.com/segmentio/kafka-go"
)

// Message writer implementation
type kafkaWriter struct {
	Writer *kafka.Writer
}

func (w *kafkaWriter) Close() error {
	return w.Writer.Close()
}

func (w *kafkaWriter) WriteMessage(value []byte) error {
	return w.Writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte("scheduler"),
			Value: value,
		},
	)
}

// Message reader implementation
type kafkaReader struct {
	Reader *kafka.Reader
}

func (r *kafkaReader) Close() error {
	return r.Reader.Close()
}

func (r *kafkaReader) ReadMessage() ([]byte, error) {
	message, err := r.Reader.ReadMessage(context.Background())
	if err != nil {
		return nil, err
	}

	return message.Value, nil
}

// Service implementation
type kafkaService struct {
	ConsulClient *consul.Client
}

// NewKafkaService return a new kafka service
func NewKafkaService(consulClient *consul.Client) Service {
	return &kafkaService{
		ConsulClient: consulClient,
	}
}

// NewReader return a new writer for kafka service
func (s *kafkaService) NewWriter(topic string) Writer {
	return &kafkaWriter{
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: s.getBrokers(),
			Topic:   topic,
		}),
	}
}

// NewReader return a new reader for kafka service
func (s *kafkaService) NewReader() Reader {
	return &kafkaReader{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: s.getBrokers(),
			Topic:   "scheduler",
		}),
	}
}

func (s *kafkaService) getBrokers() []string {
	kafkaAddr, kafkaPort, err := toolkit.GetServiceAddress(s.ConsulClient, "kafka")
	if err != nil {
		panic(err)
	}

	return []string{fmt.Sprintf("%v:%v", kafkaAddr, kafkaPort)}
}
