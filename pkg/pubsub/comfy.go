package pubsub

import (
	"github.com/ThreeDotsLabs/watermill"
	watermillsql "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/davidroman0O/comfylite3"
	comfymill "github.com/davidroman0O/watermill-comfymill"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

func New(cfg *config.Config) (*comfylite3.ComfyDB, error) {
	target := cfg.String("pubsub.db.name")
	db, err := comfymill.NewDatabase(comfylite3.WithPath(target))
	if err != nil {
		return db, errs.Wrapf(err, "pubsub: failed to create/open comfy db at %q", target)
	}
	return db, nil
}

func NewPublisher(db *comfylite3.ComfyDB) (message.Publisher, error) {
	return watermillsql.NewPublisher(
		db,
		watermillsql.PublisherConfig{
			SchemaAdapter:        comfymill.DefaultSQLite3Schema{},
			AutoInitializeSchema: true,
		},
		watermill.NopLogger{})
}

func NewSubscriber(cfg *config.Config, db *comfylite3.ComfyDB) (message.Subscriber, error) {
	deadline := cfg.Duration("pubsub.ack.deadline")
	pollInterval := cfg.Duration("pubsub.poll.interval")
	retryInterval := cfg.Duration("pubsub.retry.interval")

	return watermillsql.NewSubscriber(
		db,
		watermillsql.SubscriberConfig{
			ConsumerGroup:    cfg.String("pubsub.consumer.group"),
			AckDeadline:      &deadline,
			ResendInterval:   retryInterval,
			BackoffManager:   watermillsql.NewDefaultBackoffManager(pollInterval, retryInterval),
			SchemaAdapter:    comfymill.DefaultSQLite3Schema{},
			OffsetsAdapter:   comfymill.DefaultSQLite3OffsetsAdapter{},
			InitializeSchema: true,
		},
		watermill.NopLogger{})
}
