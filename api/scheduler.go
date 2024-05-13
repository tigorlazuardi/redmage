package api

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

func (api *API) startScheduler() func() {
	now := time.Now()

	stop := make(chan struct{})

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	var ticker *time.Ticker

	nextMinute := now.Truncate(time.Minute).Add(time.Minute)
	log.New(context.Background()).Infof("starting scheduler at %s", nextMinute)
	timer := time.AfterFunc(nextMinute.Sub(now), func() {
		api.scheduleRun(time.Now().Truncate(time.Second).Truncate(0), parser)
		ticker = time.NewTicker(time.Minute)
		go func() {
			for {
				select {
				case <-stop:
					return
				case now := <-ticker.C:
					api.scheduleRun(now.Truncate(time.Second).Truncate(0), parser)
				}
			}
		}()
	})
	return func() {
		log.New(context.Background()).Info("scheduler: stop called")
		timer.Stop()
		if ticker != nil {
			ticker.Stop()
		}
		stop <- struct{}{}
	}
}

func (api *API) scheduleRun(now time.Time, parser cron.Parser) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx, span := tracer.Start(ctx, "scheduler:tick")
	defer span.End()

	previous := now.Add(-time.Minute).Truncate(time.Minute)

	log.New(ctx).Info("scheduler: run")
	subreddits, err := models.Subreddits.Query(ctx, api.db, models.SelectWhere.Subreddits.EnableSchedule.EQ(1)).All()
	if err != nil {
		log.New(ctx).Err(err).Error("scheduler: failed to query subreddits")
		return
	}
	for _, subreddit := range subreddits {
		schedule, err := parser.Parse(subreddit.Schedule)
		if err != nil {
			log.New(ctx).Err(err).Error("scheduler: failed to parse schedule")
			continue
		}
		next := schedule.Next(previous)
		log.New(ctx).Info("scheduler: check time", "subreddit", subreddit.Name, "trigger_time", next, "now", now, "should_run", now.After(next))
		if now.After(next) {
			err := api.PubsubStartDownloadSubreddit(ctx, PubsubStartDownloadSubredditParams{Subreddit: subreddit.Name})
			if err != nil {
				log.New(ctx).Err(err).Error("scheduler: failed to start download", "subreddit", subreddit.Name)
				continue
			}
		}
	}
}
