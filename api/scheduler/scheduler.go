package scheduler

import (
	"sync"

	"github.com/robfig/cron/v3"
)

type Runner = func(subreddit string)

type Scheduler struct {
	scheduler *cron.Cron
	mu        sync.RWMutex
	entries   map[string]*Job
	run       Runner
}

type Job struct {
	ID    cron.EntryID
	Entry cron.Entry
}

func (job *Job) clone() *Job {
	return &Job{
		ID:    job.ID,
		Entry: job.Entry,
	}
}

func New(runner Runner) *Scheduler {
	return &Scheduler{
		scheduler: cron.New(),
		entries:   make(map[string]*Job),
		run:       runner,
	}
}

// Start starts the scheduler in the background.
func (s *Scheduler) Start() {
	s.scheduler.Start()
}

// Stop stops the scheduler.
func (s *Scheduler) Stop() {
	s.scheduler.Stop()
}
