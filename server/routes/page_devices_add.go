package routes

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/devicesview/put"
)

func (routes *Routes) PageDevicesAdd(rw http.ResponseWriter, req *http.Request) {
	ctx, span := tracer.Start(req.Context(), "*Routes.PageDevicesAdd")
	defer span.End()

	c := views.NewContext(routes.Config, req)

	data := put.Data{
		PageTitle:  "Add Device",
		PostAction: "/devices/add",
		AspectRatioTolerance: put.AspectRatioToleranceData{
			Value: 0.2,
		},
	}

	if err := put.View(c, data).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render add device page")
	}
}
