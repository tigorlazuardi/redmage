package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

func (routes *Routes) SubredditsCheckAPI(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.SubredditsCheck")
	defer span.End()

	var (
		enc = json.NewEncoder(rw)
		dec = json.NewDecoder(r.Body)
	)

	var body api.SubredditCheckParam
	if err := dec.Decode(&body); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_ = enc.Encode(map[string]string{"error": fmt.Sprintf("failed to decode json body: %s", err)})
		return
	}

	if err := validateSubredditCheckParam(body); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_ = enc.Encode(map[string]string{"error": err.Error()})
		return
	}

	actual, err := routes.API.SubredditCheck(ctx, body)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to check subreddit")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		_ = enc.Encode(map[string]string{"error": message})
		return
	}

	_ = enc.Encode(map[string]string{"subreddit": actual})
}

func validateSubredditCheckParam(body api.SubredditCheckParam) error {
	if body.Subreddit == "" {
		return errors.New("subreddit name is required")
	}

	return nil
}
