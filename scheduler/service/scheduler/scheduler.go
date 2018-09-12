package scheduler

import (
	"github.com/dwarvesf/yggdrasil/scheduler/model"
)

type Service interface {
	SaveRequest(r model.RequestEntity) error
}
