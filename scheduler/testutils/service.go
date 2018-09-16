package testutils

import (
	"encoding/json"

	"github.com/dwarvesf/yggdrasil/scheduler/model"
	"github.com/dwarvesf/yggdrasil/scheduler/service/stream"
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

// MockStreamService to mock StreamService
type MockStreamService struct {
	ReadData  chan []byte
	WriteData chan Output
}

// NewWriter is mock implementation
func (s *MockStreamService) NewWriter(topic string) stream.Writer {
	return &mockWriter{
		Service: s,
		Topic:   topic,
	}
}

// NewReader is mock implementation
func (s *MockStreamService) NewReader() stream.Reader {
	return &mockReader{
		Service: s,
	}
}

type mockWriter struct {
	Service *MockStreamService
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
	Service *MockStreamService
}

func (r *mockReader) ReadMessage() ([]byte, error) {
	return <-r.Service.ReadData, nil
}

func (*mockReader) Close() error {
	return nil
}
