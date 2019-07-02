package testutil

import (
	"encoding/json"

	"github.com/dwarvesf/yggdrasil/toolkit/queue"
)

// MockPayload to mock payload data
func MockPayload(content string) string {
	data := make(map[string]interface{})
	data["content"] = content

	res, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return string(res)
}

// Output to mock output of writer
type Output struct {
	Topic string
	Data  []byte
}

// MockQueueService to mock queue Service
type MockQueueService struct {
	ReadData  chan []byte
	WriteData chan Output
}

// NewWriter is mock implementation
func (s *MockQueueService) NewWriter(topic string) queue.Writer {
	return &mockWriter{
		Service: s,
		Topic:   topic,
	}
}

// NewReader is mock implementation
func (s *MockQueueService) NewReader(topic string) queue.Reader {
	return &mockReader{
		Service: s,
		Topic:   topic,
	}
}

type mockWriter struct {
	Service *MockQueueService
	Topic   string
}

func (w *mockWriter) Close() error {
	return nil
}

func (w *mockWriter) Write(key string, data []byte) error {
	w.Service.WriteData <- Output{Data: data, Topic: w.Topic}

	return nil
}

type mockReader struct {
	Service *MockQueueService
	Topic   string
}

func (r *mockReader) Read() ([]byte, error) {
	return <-r.Service.ReadData, nil
}

func (*mockReader) Close() error {
	return nil
}
