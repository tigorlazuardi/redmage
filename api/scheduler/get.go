package scheduler

// List returns a copy of the list of jobs.
func (scheduler *Scheduler) List() map[string]*Job {
	return scheduler.list(true)
}

func (scheduler *Scheduler) list(lock bool) map[string]*Job {
	if lock {
		scheduler.mu.RLock()
		defer scheduler.mu.RUnlock()
	}

	m := make(map[string]*Job, len(scheduler.entries))
	for k, v := range scheduler.entries {
		m[k] = v.clone()
	}

	return m
}

// Get returns the job for the given subreddit.
//
// Returns nil if the subreddit is not found or active.
func (scheduler *Scheduler) Get(subreddit string) *Job {
	return scheduler.get(subreddit, true)
}

func (scheduler *Scheduler) get(subreddit string, lock bool) *Job {
	if lock {
		scheduler.mu.RLock()
		defer scheduler.mu.RUnlock()
	}

	schedule := scheduler.entries[subreddit]
	if schedule != nil {
		return schedule.clone()
	}
	return nil
}

func (scheduler *Scheduler) Iter(f func(string, *Job) bool) {
	scheduler.iter(f, true)
}

func (scheduler *Scheduler) iter(f func(string, *Job) bool, lock bool) {
	if lock {
		scheduler.mu.RLock()
		defer scheduler.mu.RUnlock()
	}

	for k, v := range scheduler.entries {
		if !f(k, v.clone()) {
			break
		}
	}
}
