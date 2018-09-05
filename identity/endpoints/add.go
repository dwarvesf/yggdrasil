package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/yggdrasil/identity/service"
	"github.com/dwarvesf/yggdrasil/identity/service/add"
)

type AddRequest struct {
	Add *add.Add `json:"add"`
}
type AddResponse struct {
	Result int `json:"result"`
}

// MakeAddEndpoint ...
func MakeAddEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRequest)
		result, err := s.AddService.Add(ctx, req.Add)
		if err != nil {
			return nil, err
		}
		return AddResponse{Result: result}, nil
	}
}
