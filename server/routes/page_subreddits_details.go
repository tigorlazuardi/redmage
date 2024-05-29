package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/subreddits/details"
)

func (routes *Routes) PageSubredditsDetails(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.PageSubredditsDetails")
	defer span.End()

	name := chi.URLParam(r, "name")

	var params api.SubredditGetByNameImageParams
	params.FillFromQuery(r.URL.Query())

	var data details.Data
	data.FlashMessageSuccess = r.Header.Get("X-Flash-Message-Success")
	var err error
	data.Params = params

	c := views.NewContext(routes.Config, r)

	result, err := routes.API.SubredditGetByNameWithImages(ctx, name, params)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to get subreddit by name")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		data.Error = message
		if err := details.View(c, data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render subreddit details page")
		}
		return
	}
	data.Subreddit = result.Subreddit
	data.Images = result.Images
	data.TotalImages = result.Total
	data.Devices, err = routes.API.GetDevices(ctx, api.DevicesListParams{Status: -1})
	if err != nil {
		log.New(ctx).Err(err).Error("failed to get devices")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		data.Error = message
		if err := details.View(c, data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render subreddit details page")
		}
	}

	if err := details.View(c, data).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render subreddit details page")
	}
}
