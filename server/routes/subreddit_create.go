package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/robfig/cron/v3"
	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/api/reddit"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views/components"
	"github.com/tigorlazuardi/redmage/views/subredditsview/addview"
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

	sub, errComponents := subredditsDataFromRequest(r)
	if len(errComponents) > 0 {
		rw.WriteHeader(http.StatusBadRequest)
		for _, err := range errComponents {
			if e := err.Render(ctx, rw); e != nil {
				log.New(ctx).Err(e).Error("failed to render error")
			}
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
		renderer := addview.SubredditInputForm(addview.SubredditInputData{
			Value:     sub.Name,
			Error:     err.Error(),
			Type:      reddit.SubredditType(sub.Subtype),
			HXSwapOOB: "true",
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
	rw.Header().Set("HX-Redirect", "/subreddits")
	rw.WriteHeader(http.StatusCreated)
	_, _ = rw.Write([]byte("Subreddit created"))
}

func subredditsDataFromRequest(r *http.Request) (sub *models.Subreddit, errs []templ.Component) {
	sub = &models.Subreddit{}

	var t reddit.SubredditType
	err := t.Parse(r.FormValue("type"))
	if err != nil {
		errs = append(errs, addview.SubredditTypeInput(addview.SubredditTypeData{
			Value:     strconv.Itoa(int(t)),
			Error:     err.Error(),
			HXSwapOOB: "true",
		}))

		return nil, errs
	}
	sub.Subtype = int32(t)

	sub.Name = r.FormValue("name")
	if sub.Name == "" {
		errs = append(errs, addview.SubredditInputForm(addview.SubredditInputData{
			Value:     sub.Name,
			Error:     "name is required",
			Type:      t,
			HXSwapOOB: "true",
		}))
		return nil, errs
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
		sub.Schedule = r.FormValue("schedule")
		_, err = cronParser.Parse(sub.Schedule)
		if err != nil {
			errs = append(errs, addview.ScheduleInput(addview.ScheduleInputData{
				Value:     sub.Schedule,
				Error:     fmt.Sprintf("invalid cron schedule: %s", err),
				HXSwapOOB: "true",
			}))
		}
	}

	countback, _ := strconv.Atoi(r.FormValue("countback"))
	sub.Countback = int32(countback)
	if sub.Countback < 1 {
		errs = append(errs, addview.CountbackInput(addview.CountbackInputData{
			Value: int64(sub.Countback),
			Error: "countback must be 1 or higher",
		}))
	}

	return sub, errs
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
