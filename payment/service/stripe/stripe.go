package stripe

import (
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"

	"github.com/dwarvesf/yggdrasil/payment/model"
)

// Client represents stripe client
type Client struct{}

// New return an stripe client
func New(apiKey string) *Client {
	stripe.Key = apiKey
	return &Client{}
}

// Pay implement Pay function in Payer interface
func (sc Client) Pay(request model.Request) error {
	customerParams := &stripe.CustomerParams{}
	customerParams.SetSource(request.Token)
	newCustomer, err := customer.New(customerParams)
	if err != nil {
		return err
	}

	chargeParams := &stripe.ChargeParams{
		Amount:      stripe.Int64(request.Amount),
		Currency:    stripe.String(string(request.Currency)),
		Description: stripe.String(request.Desc),
		Customer:    stripe.String(newCustomer.ID),
	}

	_, err = charge.New(chargeParams)
	return err
}
