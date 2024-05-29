package routes

import (
	"net/http"
	"time"

	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/homeview"
)

func (routes *Routes) PageHome(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "routes.PageHome")
	defer span.End()

	var data homeview.Data

	vc := views.NewContext(routes.Config, r)

	data.ListSubredditParams.FillFromQuery(r.URL.Query())
	data.ListSubredditParams.Limit = 0
	data.ListSubredditParams.Offset = 0
	list, err := routes.API.ListSubreddits(ctx, data.ListSubredditParams)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to list subreddits")
		code, message := errs.HTTPMessage(err)
		data.Error = message
		rw.WriteHeader(code)
		if err := homeview.Home(vc, data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render home view")
		}
		return
	}

	data.ImageListParams.FillFromQuery(r.URL.Query())
	if data.ImageListParams.CreatedAt.IsZero() {
		data.ImageListParams.CreatedAt = time.Now().Add(-time.Hour * 24) // images in the last 24 hours
	}
	if r.URL.Query().Get("limit") == "" {
		data.ImageListParams.Limit = 100
	}

	imageList, err := routes.API.ImagesListWithDevicesAndSubreddits(ctx, data.ImageListParams)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to list subreddits")
		code, message := errs.HTTPMessage(err)
		data := homeview.Data{Error: message}
		rw.WriteHeader(code)
		if err := homeview.Home(vc, data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render home view")
		}
		return
	}

	data.Devices, err = routes.API.GetDevices(ctx, api.DevicesListParams{Status: -1})
	if err != nil {
		log.New(ctx).Err(err).Error("failed to list subreddits")
		code, message := errs.HTTPMessage(err)
		data := homeview.Data{Error: message}
		rw.WriteHeader(code)
		if err := homeview.Home(vc, data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render home view")
		}
		return
	}
	data.SubredditsList = list
	data.RecentlyAddedImages = homeview.NewRecentlyAddedImages(imageList.Images)
	data.Now = time.Now()
	data.TotalImages = imageList.Total

	if err := homeview.Home(vc, data).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render home view")
	}
}
