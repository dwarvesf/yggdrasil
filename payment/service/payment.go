package payment

import "github.com/dwarvesf/yggdrasil/payment/model"

// Payer is an interface for creating payment
type Payer interface {
	Pay(request model.Request) error
}
