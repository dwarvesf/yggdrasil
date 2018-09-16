package stream

import (
	"io"
)

// Writer abstracts stream writer functions
type Writer interface {
	io.Closer
	WriteMessage([]byte) error
}

// Reader abstracts stream reader functions
type Reader interface {
	io.Closer
	ReadMessage() ([]byte, error)
}

// Service represent stream service functions
type Service interface {
	NewWriter(topic string) Writer
	NewReader() Reader
}
