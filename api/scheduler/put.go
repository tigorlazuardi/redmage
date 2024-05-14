package scheduler

import (
	"net/http"

	"github.com/robfig/cron/v3"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

// Put adds a job to the scheduler.
//
// If job already exists, it will be replaced.
func (scheduler *Scheduler) Put(subreddit string, schedule string) (job *Job, err error) {
	return scheduler.put(subreddit, schedule, true)
}

func (scheduler *Scheduler) put(subreddit string, schedule string, lock bool) (job *Job, err error) {
	sched, err := cron.ParseStandard(schedule)
	if err != nil {
		return nil, errs.
			Wrapw(err, "scheduler: failed to parse schedule", "subreddit", subreddit, "schedule", schedule).
			Code(http.StatusBadRequest)
	}

	scheduler.delete(subreddit, false)

	id := scheduler.scheduler.Schedule(sched, cron.FuncJob(func() { scheduler.run(subreddit) }))

	e := scheduler.scheduler.Entry(id)
	job = &Job{ID: id, Entry: e}

	if lock {
		scheduler.mu.Lock()
		defer scheduler.mu.Unlock()
	}

	scheduler.entries[subreddit] = job

	return job, nil
}
