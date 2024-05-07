package routes

import (
	"net/http"

	"github.com/gosimple/slug"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views/devicesview/adddevice"
)

func (routes *Routes) DevicesValidateNameHTMX(rw http.ResponseWriter, req *http.Request) {
	ctx, span := tracer.Start(req.Context(), "*Routes.ValidateName")
	defer span.End()

	var nameData adddevice.NameInputData
	nameData.Value = req.FormValue("name")
	nameComponent := adddevice.NameInput(nameData)
	s := req.FormValue("slug")
	if s != "" || nameData.Value == "" {
		if err := nameComponent.Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render name input")
		}
		return
	}
	s = slug.Make(nameData.Value)

	slugData := adddevice.SlugInputData{
		Value:     s,
		HXSwapOOB: true,
	}
	exist, err := routes.API.DevicesValidateSlug(ctx, s)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to validate slug")
		_, message := errs.HTTPMessage(err)
		slugData.Error = message
		_ = nameComponent.Render(ctx, rw)
		if err := adddevice.SlugInput(slugData).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render name input")
		}
	}

	if exist {
		slugData.Error = "Device with this identifier already exist. Please change the value manually."
		// rw.WriteHeader(http.StatusConflict)
		_ = nameComponent.Render(ctx, rw)
		if err := adddevice.SlugInput(slugData).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render name input")
		}
		return
	}

	slugData.Valid = "Identifier is available"

	_ = nameComponent.Render(ctx, rw)
	if err := adddevice.SlugInput(slugData).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render name input")
	}
}
