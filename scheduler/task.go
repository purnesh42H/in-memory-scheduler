package scheduler

import "time"

const (
	NotStarted = "NotStarted"
	Running    = "Running"
	Error      = "Error"
	Finished   = "Finished"
)

type status string

type Task struct {
	id         string
	name       string
	fn         func()
	interval   int
	executeAt  time.Time
	status     status
	failure    string
	originalId string
}
