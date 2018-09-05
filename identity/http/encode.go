package http

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
)

// encodeResponse is the common method to encode all response types to the client.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return kithttp.EncodeJSONResponse(ctx, w, response)
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	// maybe we can be smart here by returning text/json error based on request's
	// content-type header
	encodeJSONError(ctx, err, w)
}

func encodeJSONError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// we can have custom response headers by implementing kithttp.Headerer in
	// our response struct
	if headerer, ok := err.(kithttp.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := http.StatusInternalServerError
	// and custom status code
	if sc, ok := err.(kithttp.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	// enforce json err response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
