package api

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"

	"github.com/robfig/cron/v3"
	"github.com/stephenafamo/bob"
	"github.com/teivah/broadcast"
	"github.com/tigorlazuardi/redmage/api/bmessage"
	"github.com/tigorlazuardi/redmage/api/reddit"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"

	watermillSql "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	watermillSqlite "github.com/walterwanderley/watermill-sqlite"
)

type API struct {
	db bob.Executor

	scheduler   *cron.Cron
	scheduleMap map[cron.EntryID]*models.Subreddit

	downloadBroadcast *broadcast.Relay[bmessage.ImageDownloadMessage]

	config *config.Config

	imageSemaphore     chan struct{}
	subredditSemaphore chan struct{}

	reddit *reddit.Reddit

	subscriber message.Subscriber
	publisher  message.Publisher
}

type Dependencies struct {
	DB       *sql.DB
	PubsubDB *sql.DB
	Config   *config.Config
	Reddit   *reddit.Reddit
}

const downloadTopic = "subreddit_download"

var watermillLogger = &log.WatermillLogger{}

func New(deps Dependencies) *API {
	ackDeadline := deps.Config.Duration("download.pubsub.ack.deadline")
	subscriber, err := watermillSql.NewSubscriber(deps.PubsubDB, watermillSql.SubscriberConfig{
		AckDeadline:      &ackDeadline,
		SchemaAdapter:    watermillSqlite.DefaultSQLiteSchema{},
		OffsetsAdapter:   watermillSqlite.DefaultSQLiteOffsetsAdapter{},
		InitializeSchema: true,
	}, watermillLogger)
	if err != nil {
		panic(err)
	}
	publisher, err := watermillSql.NewPublisher(deps.PubsubDB, watermillSql.PublisherConfig{
		SchemaAdapter:        watermillSqlite.DefaultSQLiteSchema{},
		AutoInitializeSchema: true,
	}, watermillLogger)
	if err != nil {
		panic(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	ch, err := subscriber.Subscribe(ctx, downloadTopic)
	if err != nil {
		panic(err)
	}
	api := &API{
		db:                 bob.New(deps.DB),
		scheduler:          cron.New(),
		scheduleMap:        make(map[cron.EntryID]*models.Subreddit, 8),
		downloadBroadcast:  broadcast.NewRelay[bmessage.ImageDownloadMessage](),
		config:             deps.Config,
		imageSemaphore:     make(chan struct{}, deps.Config.Int("download.concurrency.images")),
		subredditSemaphore: make(chan struct{}, deps.Config.Int("download.concurrency.subreddits")),
		reddit:             deps.Reddit,
		subscriber:         subscriber,
		publisher:          publisher,
	}

	api.startSubredditDownloadPubsub(ch)
	return api
}

func (api *API) StartScheduler(ctx context.Context) error {
	subreddits, err := models.Subreddits.Query(ctx, api.db, nil).All()
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
	})
	if err != nil {
		return errs.Wrap(err)
	}
	api.scheduleMap[id] = subreddit

	return nil
}
