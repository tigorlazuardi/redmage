package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/robfig/cron/v3"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

func (routes *Routes) SubredditsCreateAPI(rw http.ResponseWriter, req *http.Request) {
	ctx, span := tracer.Start(req.Context(), "*Routes.SubredditsCreate")
	defer span.End()

	var (
		body *models.Subreddit
		enc  = json.NewEncoder(rw)
	)

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_ = enc.Encode(map[string]string{"error": fmt.Sprintf("failed to decode json body: %s", err)})
		return
	}

	if err := validateSubredditsCreate(body); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_ = enc.Encode(map[string]string{"error": err.Error()})
		return
	}

	//  TODO: check if the subreddit actually exists on reddit

	sub, err := routes.API.SubredditsCreate(ctx, body)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to create subreddit")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		_ = enc.Encode(map[string]string{"error": message})
		return
	}

	rw.WriteHeader(http.StatusCreated)
	if err := enc.Encode(sub); err != nil {
		log.New(ctx).Err(err).Error("failed to encode subreddit into json")
	}
}

var cronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

func validateSubredditsCreate(body *models.Subreddit) error {
	if body.Name == "" {
		return errors.New("name is required")
	}
	if body.Enable > 1 {
		body.Enable = 1
	} else if body.Enable < 0 {
		body.Enable = 0
	}
	if body.Subtype > 1 {
		body.Subtype = 1
	} else if body.Subtype < 0 {
		body.Subtype = 0
	}
	if body.Schedule == "" {
		return errors.New("schedule is required")
	}
	_, err := cronParser.Parse(body.Schedule)
	if err != nil {
		return fmt.Errorf("bad cron schedule: %w", err)
	}
	if body.Countback < 1 {
		return errors.New("countback must be 1 or higher")
	}
	return nil
}
