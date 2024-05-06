package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/devicesview/devicedetails"
)

func (routes *Routes) PageDeviceDetails(rw http.ResponseWriter, req *http.Request) {
	ctx, span := tracer.Start(req.Context(), "*Routees.PageDeviceDetails")
	defer span.End()

	c := views.NewContext(routes.Config, req)

	slug := chi.URLParam(req, "slug")

	var data devicedetails.Data
	data.Params.FillFromQuery(req.URL.Query())

	var err error

	data.Device, err = routes.API.DeviceBySlug(ctx, slug)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to get device by slug")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		data.Error = message
		if err := devicedetails.View(c, data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render device details page")
		}
		return
	}

	data.Params.Device = data.Device.Slug

	result, err := routes.API.ImagesList(ctx, data.Params)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to get image by device")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		data.Error = message
		if err := devicedetails.View(c, data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render device details page")
		}
		return
	}

	data.Images = result.Images
	data.TotalImages = result.Total

	if err := devicedetails.View(c, data).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render device details page")
	}
}
