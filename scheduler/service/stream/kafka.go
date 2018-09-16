package stream

import (
	"context"
	"fmt"

	"github.com/dwarvesf/yggdrasil/toolkit"
	consul "github.com/hashicorp/consul/api"
	kafka "github.com/segmentio/kafka-go"
)

// Stream writer implementation
type kafkaStreamWriter struct {
	Writer *kafka.Writer
}

func (w *kafkaStreamWriter) Close() error {
	return w.Writer.Close()
}

func (w *kafkaStreamWriter) WriteMessage(value []byte) error {
	return w.Writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte("scheduler"),
			Value: value,
		},
	)
}

// Stream reader implementation
type kafkaStreamReader struct {
	Reader *kafka.Reader
}

func (r *kafkaStreamReader) Close() error {
	return r.Reader.Close()
}

func (r *kafkaStreamReader) ReadMessage() ([]byte, error) {
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
	return &kafkaStreamWriter{
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: s.getBrokers(),
			Topic:   topic,
		}),
	}
}

// NewReader return a new reader for kafka service
func (s *kafkaService) NewReader() Reader {
	return &kafkaStreamReader{
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
