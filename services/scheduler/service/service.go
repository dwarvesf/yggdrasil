package service

import (
	"github.com/dwarvesf/yggdrasil/services/scheduler/service/request"
	"github.com/dwarvesf/yggdrasil/toolkit/queue"
)

type Service struct {
	RequestService request.Service
	QueueService   queue.Queue
}
