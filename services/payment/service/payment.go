package payment

import "github.com/dwarvesf/yggdrasil/services/payment/model"

// Payer is an interface for creating payment
type Payer interface {
	Pay(p model.Payload) error
}
