package message

import (
	"io"
)

// Writer abstracts message writer functions
type Writer interface {
	io.Closer
	WriteMessage([]byte) error
}

// Reader abstracts message reader functions
type Reader interface {
	io.Closer
	ReadMessage() ([]byte, error)
}

// Service represent message service functions
type Service interface {
	NewWriter(topic string) Writer
	NewReader() Reader
}
