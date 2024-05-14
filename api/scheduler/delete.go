package scheduler

// Delete removes a job from the scheduler.
//
// If job does not exist, it will be a no-op.
func (scheduler *Scheduler) Delete(subreddit string) {
	scheduler.mu.Lock()
	defer scheduler.mu.Unlock()

	job := scheduler.Get(subreddit)
	if job != nil {
		scheduler.scheduler.Remove(job.ID)
	}
}
