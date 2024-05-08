package routes

import (
	"net/http"

	"github.com/gosimple/slug"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views/devicesview/put"
)

func (routes *Routes) DevicesValidateSlugHTMX(rw http.ResponseWriter, req *http.Request) {
	ctx, span := tracer.Start(req.Context(), "*Routes.ValidateSlugHTMX")
	defer span.End()

	var data put.SlugInputData
	data.Value = slug.Make(req.FormValue("slug"))
	if data.Value == "" {
		if err := put.SlugInput(data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render slug input")
		}
		return
	}

	exist, err := routes.API.DevicesValidateSlug(ctx, data.Value)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to validate slug")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		data.Error = message
		if err := put.SlugInput(data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render slug input")
		}
		return
	}

	if exist {
		data.Error = "Device with this identifier already exist"
		rw.WriteHeader(http.StatusConflict)
		if err := put.SlugInput(data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render slug input")
		}
		return
	}
	data.Valid = "Identifier is available"

	if err := put.SlugInput(data).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render slug input")
	}
}
