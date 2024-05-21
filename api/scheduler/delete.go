package scheduler

// Delete removes a job from the scheduler.
//
// If job does not exist, it will be a no-op.
func (scheduler *Scheduler) Delete(subreddit string) {
	scheduler.mu.Lock()
	defer scheduler.mu.Unlock()
	scheduler.delete(subreddit)
}

func (scheduler *Scheduler) delete(subreddit string) {
	job := scheduler.get(subreddit, false)
	if job != nil {
		scheduler.scheduler.Remove(job.ID)
	}
}
