// +build !solution

package schedule

import (
	"time"
)

// strideScheduler may be defined as an alias for rrScheduler; it has the same fields.
type strideScheduler struct {
	baseScheduler
	//TODO(student) add missing fields, if necessary
}

// newStrideScheduler returns a stride scheduler.
// With this scheduler, jobs are executed similar to round robin,
// but with exact proportions determined by how many tickets each job is assigned.
func newStrideScheduler(quantum time.Duration) *strideScheduler {
	//TODO(student) construct new Stride scheduler
	return nil
}

// schedule schedules the provided jobs according to a stride scheduler's order.
// The task with the lowest pass is scheduled to run first.
func (s *strideScheduler) schedule(inJobs jobs) {
	//TODO(student) implement scheduler
	close(s.runQueue)
}

func minPass(theJobs jobs) int {
	//TODO(student) implement minPass and use it from schedule()
	lowest := 0
	return lowest
}
