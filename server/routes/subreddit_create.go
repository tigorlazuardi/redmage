package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/robfig/cron/v3"
	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/api/reddit"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views/components"
	"github.com/tigorlazuardi/redmage/views/subreddits/put"
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

	actual, err := routes.API.SubredditCheck(ctx, api.SubredditCheckParam{
		Subreddit:     body.Name,
		SubredditType: reddit.SubredditType(body.Subtype),
	})
	if err != nil {
		log.New(ctx).Err(err).Error("subreddit check returns error")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		_ = enc.Encode(map[string]string{"error": message})
		return
	}

	body.Name = actual

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

func (routes *Routes) SubredditsCreateHTMX(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.SubredditsCreateHTMX")
	defer span.End()

	sub, err := subredditsDataFromRequest(r)
	if err != nil {
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		if err := components.ErrorNotication(message).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render error notification")
		}
		return
	}
	actual, err := routes.API.SubredditCheck(ctx, api.SubredditCheckParam{
		Subreddit:     sub.Name,
		SubredditType: reddit.SubredditType(sub.Subtype),
	})
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.New(ctx).Err(err).Error("subreddit check returns error")
		renderer := put.NameInput(put.NameInputData{
			Value: sub.Name,
			Error: err.Error(),
		})
		if err := renderer.Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render error")
		}
		return
	}
	sub.Name = actual

	_, err = routes.API.SubredditsCreate(ctx, sub)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to create subreddit")
		code, message := errs.HTTPMessage(err)
		rw.Header().Set("HX-Retarget", components.NotificationContainerID)
		rw.WriteHeader(code)
		if err := components.ErrorNotication(message).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render error")
		}
		return
	}
	if fetch, _ := strconv.ParseBool(r.FormValue("fetch")); fetch {
		_ = routes.API.PubsubStartDownloadSubreddit(ctx, api.PubsubStartDownloadSubredditParams{
			Subreddit: sub.Name,
		})
	}
	rw.Header().Set("HX-Redirect", "/subreddits")
	rw.WriteHeader(http.StatusCreated)
	_, _ = rw.Write([]byte("Subreddit created"))
}

func subredditsDataFromRequest(r *http.Request) (sub *models.Subreddit, err error) {
	sub = &models.Subreddit{}

	var t reddit.SubredditType
	err = t.Parse(r.FormValue("type"))
	if err != nil {
		return nil, errs.
			Wrapw(err, "invalid subreddit type", "type", r.FormValue("type")).
			Code(http.StatusBadRequest)
	}
	sub.Subtype = int32(t)

	sub.Name = r.FormValue("name")
	if sub.Name == "" {
		return nil, errs.Fail("name is required").Code(http.StatusBadRequest)
	}

	enableSchedule, _ := strconv.Atoi(r.FormValue("enable_schedule"))
	sub.EnableSchedule = int32(enableSchedule)
	if sub.EnableSchedule > 1 {
		sub.EnableSchedule = 1
	} else if sub.EnableSchedule < 0 {
		sub.EnableSchedule = 0
	}
	if sub.EnableSchedule == 0 {
		sub.Schedule = "@daily"
	}

	if sub.EnableSchedule == 1 {
		schedule := r.FormValue("schedule")
		_, err = cron.ParseStandard(schedule)
		if err != nil {
			return nil, errs.Wrapf(err, "invalid cron schedule: %s", err).Code(http.StatusBadRequest)
		}
		sub.Schedule = schedule
	}

	countback, _ := strconv.Atoi(r.FormValue("countback"))
	sub.Countback = int32(countback)
	if sub.Countback < 1 {
		return nil, errs.Fail("countback must be 1 or higher").Code(http.StatusBadRequest)
	}

	return sub, nil
}

var cronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

func validateSubredditsCreate(body *models.Subreddit) error {
	if body.Name == "" {
		return errors.New("name is required")
	}
	if body.EnableSchedule > 1 {
		body.EnableSchedule = 1
	} else if body.EnableSchedule < 0 {
		body.EnableSchedule = 0
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
