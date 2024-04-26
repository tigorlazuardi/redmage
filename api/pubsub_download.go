package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/tigorlazuardi/redmage/api/reddit"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (api *API) StartSubredditDownloadPubsub(messages <-chan *message.Message) {
	for msg := range messages {
		api.subredditSemaphore <- struct{}{}
		go func(msg *message.Message) {
			defer func() {
				msg.Ack()
				<-api.subredditSemaphore
			}()
			var (
				err       error
				subreddit *models.Subreddit
			)
			ctx, span := tracer.Start(context.Background(), "Download Subreddit Pubsub")
			defer func() { telemetry.EndWithStatus(span, err) }()
			span.AddEvent("pubsub." + downloadTopic)

			err = json.Unmarshal(msg.Payload, &subreddit)
			if err != nil {
				log.New(ctx).Err(err).Error("failed to unmarshal json for download pubsub", "topic", downloadTopic)
				return
			}

			devices, err := models.Devices.Query(ctx, api.db).All()
			if err != nil {
				log.New(ctx).Err(err).Error("failed to query devices")
				return
			}

			err = api.DownloadSubredditImages(ctx, subreddit.Name, DownloadSubredditParams{
				Countback:     int(subreddit.Countback),
				Devices:       devices,
				SubredditType: reddit.SubredditType(subreddit.Subtype),
			})
			if err != nil {
				log.New(ctx).Err(err).Error("failed to download subreddit images", "subreddit", subreddit)
				return
			}
		}(msg)
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

	payload, _ := json.Marshal(subreddit)

	err = api.publisher.Publish(downloadTopic, message.NewMessage(watermill.NewUUID(), payload))
	if err != nil {
		return errs.Wrapw(err, "failed to enqueue reddit download", "params", params)
	}

	return nil
}
