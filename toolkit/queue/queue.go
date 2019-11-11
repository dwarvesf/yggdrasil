package queue

import (
	"io"

	"github.com/nats-io/stan.go"
)

// Writer implement methods related to sending msg to queue
type Writer interface {
	Write(channel string, msg []byte) error
	io.Closer
}

// Reader implement methods related to reading msg from queue
type Reader interface {
	Read(handleMsg func(msg *stan.Msg)) error
	io.Closer
}

// Queue represent message queue functions
type Queue interface {
	// NewWriter return a new writer for Queue
	NewWriter() (Writer, error)
	// NewWriter return a new reader for Queue
	NewReader(topic string) (Reader, error)
}
