package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stephenafamo/bob"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (api *API) StartSubredditDownloadPubsub(messages <-chan *message.Message) {
	for msg := range messages {
		var subreddit *models.Subreddit
		if err := json.Unmarshal(msg.Payload, &subreddit); err != nil {
			log.New(context.Background()).Err(err).Error("failed to unmarshal json for download pubsub", "topic", downloadTopic)
			return
		}
		log.New(context.Background()).Debug("received pubsub message",
			"message", string(msg.Payload),
		)
		func(msg *message.Message, subreddit *models.Subreddit) {
			defer msg.Ack()
			var err error
			ctx, span := tracer.Start(context.Background(), "Download Subreddit Pubsub")
			defer func() {
				if err != nil {
					if _, err := api.ScheduleSet(ctx, ScheduleSetParams{
						Subreddit:    subreddit.Name,
						Status:       ScheduleStatusError,
						ErrorMessage: err.Error(),
					}); err != nil {
						log.New(ctx).Err(err).Error("failed to set schedule status", "subreddit", subreddit.Name, "status", ScheduleStatusError.String())
					}
				} else {
					if _, err := api.ScheduleSet(ctx, ScheduleSetParams{
						Subreddit: subreddit.Name,
						Status:    ScheduleStatusStandby,
					}); err != nil {
						log.New(ctx).Err(err).Error("failed to set schedule status", "subreddit", subreddit.Name, "status", ScheduleStatusStandby.String())
					}
				}
				telemetry.EndWithStatus(span, err)
			}()
			span.AddEvent("pubsub." + downloadTopic)
			_, err = api.ScheduleSet(ctx, ScheduleSetParams{
				Subreddit: subreddit.Name,
				Status:    ScheduleStatusDownloading,
			})
			if err != nil {
				log.New(ctx).Err(err).Error("failed to set schedule status", "subreddit", subreddit.Name, "status", ScheduleStatusDownloading.String())
			}

			devices, err := models.Devices.Query(ctx, api.db).All()
			if err != nil {
				log.New(ctx).Err(err).Error("failed to query devices")
				return
			}

			err = api.DownloadSubredditImages(ctx, subreddit, devices)
			if err != nil {
				log.New(ctx).Err(err).Error("failed to download subreddit images", "subreddit", subreddit)
				return
			}
		}(msg, subreddit)
	}
}

type PubsubStartDownloadSubredditParams struct {
	Subreddit string `json:"subreddit"`
}

func (api *API) PubsubStartDownloadSubreddit(ctx context.Context, params PubsubStartDownloadSubredditParams) (err error) {
	ctx, span := tracer.Start(ctx, "*API.PubsubStartDownloadSubreddit", trace.WithAttributes(attribute.String("subreddit", params.Subreddit)))
	defer span.End()

	subreddit, err := models.Subreddits.Query(ctx, api.db, models.SelectWhere.Subreddits.Name.EQ(params.Subreddit)).One()
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errs.Wrapw(err, "subreddit not registered", "params", params).Code(http.StatusNotFound)
		}
		return errs.Wrapw(err, "failed to verify subreddit existence", "params", params)
	}

	err = api.withTransaction(ctx, func(exec bob.Executor) error {
		_, err := api.scheduleSet(ctx, exec, ScheduleSetParams{
			Subreddit: subreddit.Name,
			Status:    ScheduleStatusEnqueued,
		})
		if err != nil {
			return err
		}

		payload, err := json.Marshal(subreddit)
		if err != nil {
			return errs.Wrapw(err, "failed to marshal subreddit")
		}

		err = api.publisher.Publish(downloadTopic, message.NewMessage(watermill.NewUUID(), payload))
		if err != nil {
			return errs.Wrapw(err, "failed to enqueue reddit download", "params", params)
		}
		return nil
	})

	return nil
}
