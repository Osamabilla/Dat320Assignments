// +build !solution

package schedule

type sjfScheduler struct {
	baseScheduler
}

// newSJFScheduler returns a shortest job first scheduler.
// With this scheduler, jobs are executed in the order of shortest job first.
func newSJFScheduler() *sjfScheduler {
	//TODO(student) construct new SJF scheduler
	return nil
}

// schedule schedules the provided jobs according to SJF order.
// The tasks with the lowest estimate is scheduled to run first.
func (s *sjfScheduler) schedule(inJobs jobs) {
	//TODO(student) implement scheduler
	close(s.runQueue)
}
