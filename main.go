package main

import (
	"fmt"
	"in-memory-task-scheduler/scheduler"
	"time"
)

func main() {
	inMemoryScheduler := scheduler.New()

	fn1 := func() {
		fmt.Printf("2 + 3 = 5\n")
	}
	fn2 := func() {
		fmt.Printf("No Op\n")
	}
	fn3 := func() {
		fmt.Printf("Hello World!\n")
	}

	inMemoryScheduler.Schedule(fn1, time.Now().Add(1*time.Second))
	inMemoryScheduler.Schedule(fn2, time.Now().Add(1*time.Second))
	inMemoryScheduler.ScheduleAtFixedInterval(fn2, 2)
	inMemoryScheduler.ScheduleAtFixedInterval(fn3, 5)

	go inMemoryScheduler.Start()
	time.Sleep(100 * time.Second)
	inMemoryScheduler.Stop()

	fmt.Printf("Done\n")
}
