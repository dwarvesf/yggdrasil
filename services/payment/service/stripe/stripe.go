package stripe

import (
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"

	"github.com/dwarvesf/yggdrasil/services/payment/model"
)

// Client represents stripe client
type Client struct{}

// New return an stripe client
func New(apiKey string) *Client {
	stripe.Key = apiKey
	return &Client{}
}

// Pay implement Pay function in Payer interface
func (sc Client) Pay(p model.Payload) error {
	customerParams := &stripe.CustomerParams{}
	customerParams.SetSource(p.Token)
	newCustomer, err := customer.New(customerParams)
	if err != nil {
		return err
	}

	chargeParams := &stripe.ChargeParams{
		Amount:      stripe.Int64(p.Amount),
		Currency:    stripe.String(string(p.Currency)),
		Description: stripe.String(p.Desc),
		Customer:    stripe.String(newCustomer.ID),
	}

	_, err = charge.New(chargeParams)
	return err
}
