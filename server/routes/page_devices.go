package routes

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/devicesview"
)

func (routes *Routes) PageDevices(rw http.ResponseWriter, req *http.Request) {
	ctx, start := tracer.Start(req.Context(), "*Routes.PageDevices")
	defer start.End()

	vc := views.NewContext(routes.Config, req)
	var data devicesview.Data
	data.Params.FillFromQuery(req.URL.Query())

	result, err := routes.API.DevicesList(ctx, data.Params)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to query devices")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		data.Error = message
		if err := devicesview.Devices(vc, data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render devices error view")
		}
	}
	data.Devices = result.Devices
	data.Total = result.Total

	if err := devicesview.Devices(vc, data).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render devices view")
	}
}
