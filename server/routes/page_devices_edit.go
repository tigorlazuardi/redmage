package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/components"
	"github.com/tigorlazuardi/redmage/views/devices/put"
)

func (routes *Routes) PageDevicesEdit(rw http.ResponseWriter, req *http.Request) {
	ctx, span := tracer.Start(req.Context(), "*Routes.PageDevicesEdit")
	defer span.End()

	c := views.NewContext(routes.Config, req)

	slug := chi.URLParam(req, "slug")

	device, err := routes.API.DeviceBySlug(ctx, slug)
	if err != nil {
		code, message := errs.HTTPMessage(err)
		if code >= 500 {
			log.New(ctx).Err(err).Error("failed to get device by slug")
		}
		rw.WriteHeader(code)
		msg := fmt.Sprintf("%d: %s", code, message)
		if err := components.PageError(c, msg).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render device edit page")
		}
		return
	}

	data := put.Data{
		PageTitle:  fmt.Sprintf("Edit Device %q", device.Name),
		PostAction: fmt.Sprintf("/devices/edit/%s", device.Slug),
		EditMode:   true,
		NameInput: put.NameInputData{
			Value:    device.Name,
			EditMode: true,
		},
		SlugInput: put.SlugInputData{
			Value: device.Slug,
		},
		ResolutionX: put.ResolutionData{
			Value: int(device.ResolutionX),
		},
		ResolutionY: put.ResolutionData{
			Value: int(device.ResolutionY),
		},
		AspectRatioTolerance: put.AspectRatioToleranceData{
			Value: device.AspectRatioTolerance,
		},
		NSFWCheckbox: put.NSFWCheckboxData{
			Checked:  device.NSFW == 1,
			EditMode: true,
		},
		WindowsWallpaperCheckbox: put.WindowsWallpaperCheckboxData{
			Checked: device.WindowsWallpaperMode == 1,
		},
		MinImageResolutionXInput: put.ResolutionData{
			Value: int(device.MinX),
		},
		MinImageResolutionYInput: put.ResolutionData{
			Value: int(device.MinY),
		},
		MaxImageResolutionXInput: put.ResolutionData{
			Value: int(device.MaxX),
		},
		MaxImageResolutionYInput: put.ResolutionData{
			Value: int(device.MaxY),
		},
	}

	if err := put.View(c, data).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render device edit page")
	}
}
