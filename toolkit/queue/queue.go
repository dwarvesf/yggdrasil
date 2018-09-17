package queue

// Queue represent message queue functions
type Queue interface {
	Write(topic string, value [][]byte) error
	Read(topic string) []byte
	Close() error
}
