package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/api/reddit"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views/subreddits/put"
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

func (routes *Routes) SubredditCheckHTMX(rw http.ResponseWriter, r *http.Request) {
	var (
		data    put.NameInputData
		subtype reddit.SubredditType
	)
	data.Value = r.FormValue("name")
	_ = subtype.Parse(r.FormValue("type"))

	if data.Value == "" {
		if err := put.NameInput(data).Render(r.Context(), rw); err != nil {
			log.New(r.Context()).Err(err).Error("failed to render subreddit input form")
		}
		return
	}

	ctx, span := tracer.Start(r.Context(), "*Routes.SubredditCheckHTMX")
	defer span.End()

	params := api.SubredditCheckParam{
		Subreddit:     data.Value,
		SubredditType: subtype,
	}

	actual, err := routes.API.SubredditCheck(ctx, params)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to check subreddit")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		data.Error = message
		if err := put.NameInput(data).Render(r.Context(), rw); err != nil {
			log.New(r.Context()).Err(err).Error("failed to render subreddit input form")
		}
		return
	}
	params.Subreddit = actual
	data.Value = actual

	exist, err := routes.API.SubredditRegistered(ctx, params)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to check subreddit")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		data.Error = message
		if err := put.NameInput(data).Render(r.Context(), rw); err != nil {
			log.New(r.Context()).Err(err).Error("failed to render subreddit input form")
		}
		return
	}

	if exist {
		rw.WriteHeader(http.StatusConflict)
		data.Error = "subreddit already registered"
		if err := put.NameInput(data).Render(r.Context(), rw); err != nil {
			log.New(r.Context()).Err(err).Error("failed to render subreddit input form")
		}
		return
	}

	data.Valid = fmt.Sprintf("%s is valid", subtype)

	if err := put.NameInput(data).Render(r.Context(), rw); err != nil {
		log.New(r.Context()).Err(err).Error("failed to render subreddit input form")
	}
}

func validateSubredditCheckParam(body api.SubredditCheckParam) error {
	if body.Subreddit == "" {
		return errors.New("subreddit name is required")
	}

	return nil
}
