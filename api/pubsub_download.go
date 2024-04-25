package api

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/tigorlazuardi/redmage/api/reddit"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
)

func (api *API) StartSubredditDownloadPubsub(messages <-chan *message.Message) {
	for msg := range messages {
		api.subredditSemaphore <- struct{}{}
		go func(msg *message.Message) {
			defer func() {
				msg.Ack()
				<-api.subredditSemaphore
			}()
			var err error
			ctx, span := tracer.Start(context.Background(), "Download Subreddit Pubsub")
			defer func() { telemetry.EndWithStatus(span, err) }()
			span.AddEvent("pubsub." + downloadTopic)
			subredditName := string(msg.Payload)
			span.SetAttributes(attribute.String("subreddit", subredditName))

			subreddit, err := models.Subreddits.Query(ctx, api.db, models.SelectWhere.Subreddits.Name.EQ(subredditName)).One()
			if err != nil {
				log.New(ctx).Err(err).Error("failed to find subreddit", "subreddit", subredditName)
				return
			}

			devices, err := models.Devices.Query(ctx, api.db).All()
			if err != nil {
				log.New(ctx).Err(err).Error("failed to query devices")
				return
			}

			err = api.DownloadSubredditImages(ctx, subredditName, DownloadSubredditParams{
				Countback:     int(subreddit.Countback),
				Devices:       devices,
				SubredditType: reddit.SubredditType(subreddit.Subtype),
			})
			if err != nil {
				log.New(ctx).Err(err).Error("failed to download subreddit images", "subreddit", subredditName)
				return
			}
		}(msg)
	}
}
