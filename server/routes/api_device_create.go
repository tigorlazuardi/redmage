package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/pkg/telemetry"
)

func (routes *Routes) APIDeviceCreate(rw http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := tracer.Start(r.Context(), "*Routes.APIDeviceCreate")
	defer func() { telemetry.EndWithStatus(span, err) }()

	var body api.DeviceCreateParams

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.New(ctx).Err(err).Error("failed to decode json body")
		rw.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": fmt.Sprintf("cannot decode json body: %s", err)})
		return
	}

	if err = validateDeviceCreateParams(body); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": err.Error()})
		return
	}

	device, err := routes.API.DevicesCreate(ctx, body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := json.NewEncoder(rw).Encode(device); err != nil {
		log.New(ctx).Err(err).Error("failed to marshal json api devices")
	}
}

func validateDeviceCreateParams(params api.DeviceCreateParams) error {
	if params.Name == "" {
		return errors.New("name is required")
	}
	if params.Slug == "" {
		return errors.New("slug is required")
	}
	if params.ResolutionX < 1 {
		return errors.New("device width resolution is required")
	}
	if params.ResolutionY < 1 {
		return errors.New("device height resolution is required")
	}
	if params.MaxX < 0 {
		params.MaxX = 0
	}
	if params.MaxY < 0 {
		params.MaxY = 0
	}
	if params.MinX < 0 {
		params.MinX = 0
	}
	if params.MinY < 0 {
		params.MinY = 0
	}
	return nil
}
