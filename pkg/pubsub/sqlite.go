package pubsub

import (
	"database/sql"

	"github.com/ThreeDotsLabs/watermill"
	wsql "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/pkg/errs"

	comfymill "github.com/davidroman0O/watermill-comfymill"
)

func NewSqlite(cfg *config.Config) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", cfg.String("pubsub.db.string"))
	if err != nil {
		return nil, errs.Wrapw(err, "failed to open sqlite3 database", "db.string", cfg.String("pubsub.db.string"))
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func NewSQLPublisher(db *sql.DB) (message.Publisher, error) {
	pub, err := wsql.NewPublisher(db, wsql.PublisherConfig{
		SchemaAdapter:        comfymill.DefaultSQLite3Schema{},
		AutoInitializeSchema: true,
	}, watermill.NopLogger{})
	if err != nil {
		return nil, errs.Wrapw(err, "failed to create sqlite3 subscriber")
	}
	return pub, nil
}

func NewSQLSubscriber(cfg *config.Config, db *sql.DB) (message.Subscriber, error) {
	deadline := cfg.Duration("pubsub.ack.deadline")
	pollInterval := cfg.Duration("pubsub.poll.interval")
	retryInterval := cfg.Duration("pubsub.retry.interval")
	sub, err := wsql.NewSubscriber(db, wsql.SubscriberConfig{
		ConsumerGroup:    "redmage",
		AckDeadline:      &deadline,
		ResendInterval:   0,
		BackoffManager:   wsql.NewDefaultBackoffManager(pollInterval, retryInterval),
		SchemaAdapter:    comfymill.DefaultSQLite3Schema{},
		OffsetsAdapter:   comfymill.DefaultSQLite3OffsetsAdapter{},
		InitializeSchema: true,
	}, watermill.NopLogger{})
	if err != nil {
		return nil, errs.Wrapw(err, "failed to create sqlite3 subscriber")
	}

	return sub, nil
}
