package testutils

import (
	"encoding/json"

	"github.com/dwarvesf/yggdrasil/scheduler/model"
	"github.com/dwarvesf/yggdrasil/scheduler/service/message"
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

// MockRequestService to mock RequestService
type MockRequestService struct {
	Requests   []model.RequestEntity
	DeletedIds []uint
}

// SaveRequest is mock implementation
func (s *MockRequestService) SaveRequest(r model.RequestEntity) error {
	s.Requests = append(s.Requests, r)

	return nil
}

// GetRequests is mock implementation
func (s *MockRequestService) GetRequests() ([]model.RequestEntity, error) {
	return s.Requests, nil
}

// DeleteRequests is mock implementation
func (s *MockRequestService) DeleteRequests(ids []uint) error {
	s.DeletedIds = ids

	return nil
}

// Output to mock output of writer
type Output struct {
	Topic string
	Data  []byte
}

// MockMessageService to mock message Service
type MockMessageService struct {
	ReadData  chan []byte
	WriteData chan Output
}

// NewWriter is mock implementation
func (s *MockMessageService) NewWriter(topic string) message.Writer {
	return &mockWriter{
		Service: s,
		Topic:   topic,
	}
}

// NewReader is mock implementation
func (s *MockMessageService) NewReader() message.Reader {
	return &mockReader{
		Service: s,
	}
}

type mockWriter struct {
	Service *MockMessageService
	Topic   string
}

func (w *mockWriter) Close() error {
	return nil
}

func (w *mockWriter) WriteMessage(data []byte) error {
	w.Service.WriteData <- Output{Data: data, Topic: w.Topic}

	return nil
}

type mockReader struct {
	Service *MockMessageService
}

func (r *mockReader) ReadMessage() ([]byte, error) {
	return <-r.Service.ReadData, nil
}

func (*mockReader) Close() error {
	return nil
}
