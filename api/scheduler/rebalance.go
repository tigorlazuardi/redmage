package scheduler

import (
	"context"
	"errors"

	"github.com/stephenafamo/bob"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

// Sync empties the scheduler and re-adds all enabled jobs from the database.
//
// This is costly but ensures that the scheduler is in sync with the database.
//
// For simpler operation consider using Put and Delete instead.
func (scheduler *Scheduler) Sync(ctx context.Context, db bob.Executor) (err error) {
	ctx, span := tracer.Start(ctx, "*Scheduler.Rebalance")
	defer span.End()

	subs, err := models.Subreddits.Query(ctx, db, models.SelectWhere.Subreddits.EnableSchedule.EQ(1)).All()
	if err != nil {
		return errs.Wrapw(err, "scheduler: rebalance: failed to query subreddits")
	}

	scheduler.mu.Lock()
	defer scheduler.mu.Unlock()

	for _, entry := range scheduler.scheduler.Entries() {
		scheduler.scheduler.Remove(entry.ID)
	}

	errcoll := make([]error, 0, len(subs))

	for _, sub := range subs {
		_, err := scheduler.put(sub.Name, sub.Schedule)
		errcoll = append(errcoll, err)
	}

	if err := errors.Join(errcoll...); err != nil {
		return errs.Wrapw(err, "scheduler: encountered errors while rebalancing jobs")
	}

	return nil
}
