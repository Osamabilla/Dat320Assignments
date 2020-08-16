package schedule

import (
	"time"
)

type rrScheduler struct {
	baseScheduler
	//TODO(student) add missing fields, if necessary
}

// newRRScheduler returns a Round Robin scheduler with the time slice, quantum.
func newRRScheduler(quantum time.Duration) *rrScheduler {
	//TODO(student) construct new RR scheduler
	return nil
}

// schedule schedules the provided jobs in round robin order.
func (s *rrScheduler) schedule(jobs jobs) {
	//TODO(student) implement scheduler
	close(s.runQueue)
}
