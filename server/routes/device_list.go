package routes

import (
	"encoding/json"
	"net/http"

	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

func (routes *Routes) APIDeviceList(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.APIDeviceList")
	defer span.End()

	var params api.DevicesListParams
	params.FillFromQuery(r.URL.Query())

	result, err := routes.API.DevicesList(ctx, params)
	if err != nil {
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": message})
		return
	}

	if err := json.NewEncoder(rw).Encode(result); err != nil {
		log.New(ctx).Err(err).Error("failed to marshal json api devices")
	}
}
