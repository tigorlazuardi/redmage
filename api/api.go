package api

import (
	"context"
	"database/sql"
	"sync"

	"github.com/stephenafamo/bob"
	"github.com/teivah/broadcast"
	"github.com/tigorlazuardi/redmage/api/bmessage"
	"github.com/tigorlazuardi/redmage/api/reddit"
	"github.com/tigorlazuardi/redmage/api/scheduler"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"go.opentelemetry.io/otel/attribute"

	"github.com/ThreeDotsLabs/watermill/message"
)

type API struct {
	db    bob.Executor
	sqldb *sql.DB

	scheduler *scheduler.Scheduler

	downloadBroadcast *broadcast.Relay[bmessage.ImageDownloadMessage]

	config *config.Config

	imageSemaphore chan struct{}

	reddit *reddit.Reddit

	subscriber message.Subscriber
	publisher  message.Publisher

	mu *sync.Mutex
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
		downloadBroadcast: broadcast.NewRelay[bmessage.ImageDownloadMessage](),
		config:            deps.Config,
		imageSemaphore:    make(chan struct{}, deps.Config.Int("download.concurrency.images")),
		reddit:            deps.Reddit,
		subscriber:        deps.Subscriber,
		publisher:         deps.Publisher,
		mu:                &sync.Mutex{},
	}

	api.scheduler = scheduler.New(api.scheduleRun)

	if err := api.scheduler.Sync(context.Background(), api.db); err != nil {
		panic(err)
	}
	go api.scheduler.Start()

	go api.StartSubredditDownloadPubsub(ch)
	return api
}

func (api *API) scheduleRun(subreddit string) {
	ctx, cancel := context.WithTimeout(context.Background(), api.config.Duration("scheduler.timeout"))
	defer cancel()

	ctx, span := tracer.Start(ctx, "*API.scheduleRun")
	defer span.End()
	span.SetAttributes(attribute.String("subreddit", subreddit))

	log.New(ctx).Info("api: schedule run", "subreddit", subreddit)

	err := api.PubsubStartDownloadSubreddit(ctx, PubsubStartDownloadSubredditParams{Subreddit: subreddit})
	if err != nil {
		log.New(ctx).Err(err).Error("api: failed to start download subreddit", "subreddit", subreddit)
	}
}
