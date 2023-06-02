package main

import (
	"fmt"
	"in-memory-task-scheduler/scheduler"
	"time"
)

func main() {
	inMemoryScheduler := scheduler.New(3)

	fn1Name := "addition"
	fn1 := func() {
		fmt.Printf("2 + 3 = 5\n")
		time.Sleep(time.Second)
	}

	fn2Name := "noOp"
	fn2 := func() {
		fmt.Printf("No Op\n")
		time.Sleep(5 * time.Second)
	}

	fn3Name := "hellowWorld"
	fn3 := func() {
		fmt.Printf("Hello World!\n")
		time.Sleep(time.Second)
	}

	fn4Name := "fooBar"
	fn4 := func() {
		fmt.Printf("Foo Bar\n")
		time.Sleep(time.Second)
	}

	fmt.Printf("Running tasks within threadpool\n\n")

	go inMemoryScheduler.Schedule(fn1Name, fn1, time.Now().Add(2*time.Second))
	go inMemoryScheduler.Schedule(fn1Name, fn2, time.Now().Add(2*time.Second))
	go inMemoryScheduler.ScheduleAtFixedInterval(fn3Name, fn3, 1)
	go inMemoryScheduler.ScheduleAtFixedInterval(fn4Name, fn4, 5)

	go inMemoryScheduler.Start()

	time.Sleep(20 * time.Second)

	fmt.Printf("\nRunning tasks more than threadpool\n\n")

	go inMemoryScheduler.ScheduleAtFixedInterval(fn1Name, fn1, 1)
	go inMemoryScheduler.ScheduleAtFixedInterval(fn2Name, fn2, 5)
	go inMemoryScheduler.ScheduleAtFixedInterval(fn4Name, fn4, 1)

	time.Sleep(20 * time.Second)

	inMemoryScheduler.Stop()

	fmt.Printf("Done\n")
}
