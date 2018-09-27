package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/device/endpoints"
)

//DecodeCreateDeviceRequest decode create device request
func DecodeCreateDeviceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

//DecodeGetDeviceRequest decode request for get device by device id
func DecodeGetDeviceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	deviceID, err := uuid.FromString(chi.URLParam(r, "device_id"))
	if err != nil {
		return nil, err
	}
	var req endpoints.GetRequest
	req.DeviceID = deviceID
	return req, err
}

//DecodeGetListDeviceRequest decode request for get list device
func DecodeGetListDeviceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID, err := uuid.FromString(chi.URLParam(r, "user_id"))
	if err != nil {
		return nil, err
	}
	var req endpoints.GetListDeviceRequest
	req.UserID = userID
	return req, err
}

//DecodeUpdateDeviceRequest decode request for update device
func DecodeUpdateDeviceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	deviceID, err := uuid.FromString(chi.URLParam(r, "device_id"))
	if err != nil {
		return nil, err
	}
	var req endpoints.UpdateRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	req.ID = deviceID
	return req, err
}
