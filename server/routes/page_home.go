package routes

import (
	"encoding/json"
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

	vc := views.NewContext(routes.Config, r)

	listSubredditParams := parseSubredditListQuery(r)

	list, err := routes.API.ListSubreddits(ctx, listSubredditParams)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to list subreddits")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": message})
		return
	}

	imageListParams := api.ImageListParams{}
	imageListParams.FillFromQuery(r.URL.Query())
	imageListParams.CreatedAt = time.Now().Add(-time.Hour * 24) // images in the last 24 hours
	imageListParams.Limit = 0

	imageList, err := routes.API.ImagesListWithDevicesAndSubreddits(ctx, imageListParams)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to list subreddits")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": message})
		return
	}

	data := homeview.Data{
		SubredditsList:      list,
		RecentlyAddedImages: homeview.NewRecentlyAddedImages(imageList.Images),
		Error:               err,
	}

	if err := homeview.Home(vc, data).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render home view")
	}
}
