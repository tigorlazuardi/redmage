package routes

import (
	"encoding/json"
	"net/http"

	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

func (r *Routes) SubredditStartDownloadAPI(rw http.ResponseWriter, req *http.Request) {
	enc := json.NewEncoder(rw)
	if r.Config.String("download.directory") == "" {
		rw.WriteHeader(http.StatusBadRequest)
		_ = enc.Encode(map[string]string{"error": "cannot download subreddits when download directory is not configured"})
		return
	}

	ctx := req.Context()

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
