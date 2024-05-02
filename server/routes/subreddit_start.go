package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views/components"
)

func (r *Routes) SubredditStartDownloadAPI(rw http.ResponseWriter, req *http.Request) {
	ctx, span := tracer.Start(req.Context(), "Routes.SubredditStartDownloadAPI")
	defer span.End()

	enc := json.NewEncoder(rw)
	if r.Config.String("download.directory") == "" {
		rw.WriteHeader(http.StatusBadRequest)
		_ = enc.Encode(map[string]string{"error": "cannot download subreddits when download directory is not configured"})
		return
	}

	var body api.PubsubStartDownloadSubredditParams

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_ = enc.Encode(map[string]string{"error": err.Error()})
		return
	}

	err = r.API.PubsubStartDownloadSubreddit(ctx, body)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to start subreddit download", "subreddit", body.Subreddit)
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		_ = enc.Encode(map[string]string{"error": message})
		return
	}

	_ = enc.Encode(map[string]string{"message": "subreddit enqueued"})
}

func (routes *Routes) SubredditStartDownloadHTMX(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "Routes.SubredditStartDownloadHTMX")
	defer span.End()

	if routes.Config.String("download.directory") == "" {
		rw.WriteHeader(http.StatusBadRequest)
		if err := components.ErrorNotication("Cannot download subreddits when download directory is not configured").Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render error notification")
		}
		return
	}

	var body api.PubsubStartDownloadSubredditParams
	body.Subreddit = r.FormValue("subreddit")

	err := routes.API.PubsubStartDownloadSubreddit(ctx, body)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to start subreddit download", "subreddit", body.Subreddit)
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		if err := components.ErrorNotication(message).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render error notification")
		}
		return
	}

	msg := fmt.Sprintf("Subreddit %s enqueued", body.Subreddit)

	if err := components.SuccessNotification(msg).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render success notification")
	}
}
