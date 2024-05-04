package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/stephenafamo/bob"
	"github.com/teivah/broadcast"
	"github.com/tigorlazuardi/redmage/api/bmessage"
	"github.com/tigorlazuardi/redmage/api/reddit"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type API struct {
	db    bob.Executor
	sqldb *sql.DB

	scheduler   *cron.Cron
	scheduleMap map[cron.EntryID]*models.Subreddit

	downloadBroadcast *broadcast.Relay[bmessage.ImageDownloadMessage]

	config *config.Config

	imageSemaphore chan struct{}

	reddit *reddit.Reddit

	subscriber message.Subscriber
	publisher  message.Publisher
}

type Dependencies struct {
	DB         *sql.DB
	Config     *config.Config
	Reddit     *reddit.Reddit
	Publisher  message.Publisher
	Subscriber message.Subscriber
}

const downloadTopic = "subreddit_download"

func New(deps Dependencies) *API {
	ch, err := deps.Subscriber.Subscribe(context.Background(), downloadTopic)
	if err != nil {
		panic(err)
	}
	api := &API{
		db:                bob.New(deps.DB),
		sqldb:             deps.DB,
		scheduler:         cron.New(),
		scheduleMap:       make(map[cron.EntryID]*models.Subreddit, 8),
		downloadBroadcast: broadcast.NewRelay[bmessage.ImageDownloadMessage](),
		config:            deps.Config,
		imageSemaphore:    make(chan struct{}, deps.Config.Int("download.concurrency.images")),
		reddit:            deps.Reddit,
		subscriber:        deps.Subscriber,
		publisher:         deps.Publisher,
	}

	if err := api.StartScheduler(context.Background()); err != nil {
		panic(err)
	}
	go api.StartSubredditDownloadPubsub(ch)
	return api
}

func (api *API) StartScheduler(ctx context.Context) error {
	subreddits, err := models.Subreddits.Query(ctx, api.db, models.SelectWhere.Subreddits.EnableSchedule.EQ(1)).All()
	if err != nil {
		return errs.Wrapw(err, "failed to get all subreddits")
	}

	for _, subreddit := range subreddits {
		err := api.scheduleSubreddit(subreddit)
		if err != nil {
			log.New(ctx).Err(err).Error(
				fmt.Sprintf("failed to start scheduler for subreddit '%s'", subreddit.Name),
				"subreddit", subreddit,
			)
			continue
		}
	}

	return nil
}

func (api *API) scheduleSubreddit(subreddit *models.Subreddit) error {
	id, err := api.scheduler.AddFunc(subreddit.Schedule, func() {
		payload, _ := json.Marshal(subreddit)
		_ = api.publisher.Publish(downloadTopic, message.NewMessage(watermill.NewUUID(), payload))
	})
	if err != nil {
		return errs.Wrap(err)
	}
	api.scheduleMap[id] = subreddit

	return nil
}
