package queue

import "io"

// Writer implement methods related to sending msg to queue
type Writer interface {
	Write(key string, value []byte) error
	io.Closer
}

// Reader implement methods related to reading msg from queue
type Reader interface {
	Read() ([]byte, error)
	io.Closer
}

// Queue represent message queue functions
type Queue interface {
	NewWriter(topic string) Writer
	NewReader(topic string) Reader
}
