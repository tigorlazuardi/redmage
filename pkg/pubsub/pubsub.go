package pubsub

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-bolt/pkg/bolt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"go.etcd.io/bbolt"
)

func NewDB(cfg *config.Config) (*bbolt.DB, error) {
	db, err := bbolt.Open(cfg.String("pubsub.db.name"), 0o600, &bbolt.Options{
		Timeout: cfg.Duration("pubsub.db.timeout"),
	})
	if err != nil {
		return db, errs.Wrapw(err, "failed to open db")
	}
	return db, err
}

func NewPublisher(db *bbolt.DB) (message.Publisher, error) {
	return bolt.NewPublisher(db, bolt.PublisherConfig{
		Common: bolt.CommonConfig{
			Bucket: []bolt.BucketName{bolt.BucketName("watermill")},
			Logger: watermill.NopLogger{},
		},
	})
}

func NewSubscriber(db *bbolt.DB) (message.Subscriber, error) {
	return bolt.NewSubscriber(db, bolt.SubscriberConfig{
		Common: bolt.CommonConfig{
			Bucket:    []bolt.BucketName{bolt.BucketName("watermill")},
			Marshaler: nil,
			Logger:    watermill.NopLogger{},
		},
	})
}
