package kafka

import "github.com/dwarvesf/yggdrasil/email/queue"

type Kafka struct {
	// Protocol
}

func New() queue.Queuer {
	return &Kafka{}
}
