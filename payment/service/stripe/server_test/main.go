package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/stripe/stripe-go"

	"github.com/dwarvesf/yggdrasil/payment/model"
	microStripe "github.com/dwarvesf/yggdrasil/payment/service/stripe"
)

func main() {
	publishableKey := os.Getenv("STRIPE_API_PUBLIC")
	secretKey := os.Getenv("STRIPE_API_SECRET")

	tmpls, err := template.ParseFiles(filepath.Join("views", "payment.html"))
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := tmpls.Lookup("payment.html")
		tmpl.Execute(w, map[string]string{"Key": publishableKey})
	})

	http.HandleFunc("/charge", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		request := model.Request{
			Currency: string(stripe.CurrencyUSD),
			Amount:   500,
			Desc:     "Sample Charge",
			Token:    r.Form.Get("stripeToken"),
		}

		payer := microStripe.New(secretKey)

		if err := payer.Pay(request); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Charge completed successfully!")
	})

	http.ListenAndServe(":4567", nil)
}
