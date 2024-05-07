package routes

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/devicesview/adddevice"
)

func (routes *Routes) PageDevicesAdd(rw http.ResponseWriter, req *http.Request) {
	ctx, span := tracer.Start(req.Context(), "*Routes.PageDevicesAdd")
	defer span.End()

	c := views.NewContext(routes.Config, req)

	if err := adddevice.View(c).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render add device page")
	}
}
