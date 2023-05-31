package scheduler

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testScheduler = New()
	fn1           = func() {
		fmt.Printf("2 + 3 = 5\n")
	}
	fn2 = func() {
		fmt.Printf("No Op\n")
	}
	fn3 = func() {
		fmt.Printf("Hello World!\n")
	}
)

func TestSchedule(t *testing.T) {
	id1 := testScheduler.Schedule(fn1, time.Now().Add(1*time.Second))
	id2 := testScheduler.Schedule(fn2, time.Now().Add(1*time.Second))

	go testScheduler.Start()
	time.Sleep(2 * time.Second)
	testScheduler.Stop()

	assert.Equal(t, Finished, string(testScheduler.GetTaskStatus(id1)))
	assert.Equal(t, Finished, string(testScheduler.GetTaskStatus(id2)))
}

func TestScheduleAtFixedInterval(t *testing.T) {
	id1 := testScheduler.ScheduleAtFixedInterval(fn2, 2)
	id2 := testScheduler.ScheduleAtFixedInterval(fn3, 3)

	go testScheduler.Start()
	time.Sleep(10 * time.Second)
	testScheduler.Stop()

	assert.Equal(t, Finished, string(testScheduler.GetTaskStatus(id1)))
	assert.Equal(t, Finished, string(testScheduler.GetTaskStatus(id2)))
	assert.Less(t, 1, testScheduler.GetTaskExecutions(id1))
	assert.Less(t, 1, testScheduler.GetTaskExecutions(id2))
}
