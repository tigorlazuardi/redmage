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

	listSubredditParams := parseSubredditListQuery(r)

	list, err := routes.API.ListSubreddits(ctx, listSubredditParams)
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

	imageListParams := api.ImageListParams{}
	imageListParams.FillFromQuery(r.URL.Query())
	if imageListParams.CreatedAt.IsZero() {
		imageListParams.CreatedAt = time.Now().Add(-time.Hour * 24) // images in the last 24 hours
	}
	imageListParams.Limit = 0

	imageList, err := routes.API.ImagesListWithDevicesAndSubreddits(ctx, imageListParams)
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
