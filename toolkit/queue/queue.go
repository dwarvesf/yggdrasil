package queue

// Queue represent message queue functions
type Queue interface {
	Write(value [][]byte) error
	Read() []byte
	Close() (error, error)
}
