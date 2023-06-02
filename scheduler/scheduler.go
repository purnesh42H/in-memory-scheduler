package scheduler

import (
	"fmt"
	"sync"
	"time"
)

type Scheduler interface {
	Schedule(taskName string, fn func(), executeAt time.Time) string
	ScheduleAtFixedInterval(taskName string, fn func(), interval int) string
	Start()
	Stop()
	GetTaskStatus(id string) status
	GetTaskExecutions(id string) int
}

type inMemoryScheduler struct {
	threads     int
	tasks       []Task
	taskHistory map[string][]Task
	stop        chan struct{}
	semaphore   chan struct{}
	mutex       sync.Mutex
	wg          sync.WaitGroup
}

func New(threads int) Scheduler {
	return &inMemoryScheduler{
		threads:     threads,
		tasks:       make([]Task, 0),
		taskHistory: make(map[string][]Task),
		semaphore:   make(chan struct{}, threads),
		stop:        make(chan struct{}),
	}
}

func (sc *inMemoryScheduler) Schedule(taskName string, fn func(), executeAt time.Time) string {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	return sc.scheduledTask(taskName, fn, executeAt, 0, "")
}

func (sc *inMemoryScheduler) ScheduleAtFixedInterval(taskName string, fn func(), interval int) string {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	return sc.scheduledTask(taskName, fn, time.Now(), interval, "")
}

func (sc *inMemoryScheduler) Start() {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			now := time.Now()

			sc.mutex.Lock()
			for i := 0; i < len(sc.tasks); i++ {
				task := sc.tasks[i]

				if task.executeAt.Before(now) && task.status == NotStarted {

					select {
					case sc.semaphore <- struct{}{}:
						sc.wg.Add(1)
						go sc.handle(i)

					default:
						fmt.Printf("All threads are occupied, can't run %s\n", task.id)
					}
				} else if task.status == Finished || task.status == Error {
					sc.tasks = append(sc.tasks[:i], sc.tasks[i+1:]...)
					i--
				}
			}
			sc.mutex.Unlock()

		case <-sc.stop:
			ticker.Stop()
			sc.wg.Wait()
			return
		}
	}
}

func (sc *inMemoryScheduler) Stop() {
	close(sc.stop)
}

func (sc *inMemoryScheduler) GetTaskStatus(id string) status {
	return sc.taskHistory[id][0].status
}

func (sc *inMemoryScheduler) GetTaskExecutions(id string) int {
	return len(sc.taskHistory[id])
}

func (sc *inMemoryScheduler) handle(index int) {
	defer func() {
		<-sc.semaphore
		sc.wg.Done()
	}()

	task := sc.tasks[index]

	sc.execute(index)

	if task.interval > 0 {
		var originalId string
		if task.originalId == "" {
			originalId = task.id
		} else {
			originalId = task.originalId
		}

		sc.scheduledTask(task.name, task.fn,
			time.Now().Add(time.Duration(task.interval)*time.Second),
			task.interval, originalId)
	}
}

func (sc *inMemoryScheduler) execute(index int) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.tasks[index].status = Running
	fmt.Printf("Executing task %s at timestamp %v\n", sc.tasks[index].id, time.Now())

	sc.tasks[index].fn()

	sc.tasks[index].status = Finished
	fmt.Printf("Finished task %s at timestamp %v\n\n", sc.tasks[index].id, time.Now())

	sc.taskHistory[sc.tasks[index].originalId] = append(sc.taskHistory[sc.tasks[index].originalId], sc.tasks[index])
}

func (sc *inMemoryScheduler) scheduledTask(taskName string, fn func(), executeAt time.Time, interval int, originalId string) string {
	id := fmt.Sprintf("%s-%d", taskName, len(sc.taskHistory[originalId]))
	if originalId == "" {
		originalId = id
	}
	task := Task{
		id:         id,
		name:       taskName,
		fn:         fn,
		executeAt:  executeAt,
		interval:   interval,
		status:     NotStarted,
		originalId: originalId,
	}

	sc.tasks = append(sc.tasks, task)

	return task.id
}
