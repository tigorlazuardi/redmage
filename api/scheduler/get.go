package scheduler

// List returns a copy of the list of jobs.
func (scheduler *Scheduler) List() map[string]*Job {
	scheduler.mu.RLock()
	defer scheduler.mu.RUnlock()

	m := make(map[string]*Job, len(scheduler.list))
	for k, v := range scheduler.list {
		m[k] = v.clone()
	}

	return m
}

// Get returns the job for the given subreddit.
//
// Returns nil if the subreddit is not found or active.
func (scheduler *Scheduler) Get(subreddit string) *Job {
	scheduler.mu.RLock()
	defer scheduler.mu.RUnlock()

	schedule := scheduler.list[subreddit]
	if schedule != nil {
		return schedule.clone()
	}
	return nil
}

func (scheduler *Scheduler) Iter(f func(string, *Job) bool) {
	scheduler.mu.RLock()
	defer scheduler.mu.RUnlock()

	for k, v := range scheduler.list {
		if !f(k, v.clone()) {
			break
		}
	}
}
