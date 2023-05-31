package scheduler

import (
	"fmt"
	"sync"
	"time"
)

type Scheduler interface {
	Schedule(fn func(), executeAt time.Time) string
	ScheduleAtFixedInterval(fn func(), interval int) string
	Start()
	Stop()
	GetTaskStatus(id string) status
	GetTaskExecutions(id string) int
}

type inMemoryScheduler struct {
	tasks       []Task
	taskHistory map[string][]Task
	stop        chan struct{}
	mutex       sync.Mutex
}

func New() Scheduler {
	return &inMemoryScheduler{
		tasks:       make([]Task, 0),
		taskHistory: make(map[string][]Task),
		stop:        make(chan struct{}),
	}
}

func (sc *inMemoryScheduler) Schedule(fn func(), executeAt time.Time) string {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	return sc.scheduledTask(fn, executeAt, 0, "")
}

func (sc *inMemoryScheduler) ScheduleAtFixedInterval(fn func(), interval int) string {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	return sc.scheduledTask(fn, time.Now(), interval, "")
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
					go sc.execute(i)
					if task.interval > 0 {
						var originalId string
						if task.originalId == "" {
							originalId = task.id
						} else {
							originalId = task.originalId
						}

						sc.scheduledTask(task.fn,
							time.Now().Add(time.Duration(task.interval)*time.Second),
							task.interval, originalId)
					}
				} else if task.status == Finished || task.status == Error {
					sc.tasks = append(sc.tasks[:i], sc.tasks[i+1:]...)
					i--
				}
			}
			sc.mutex.Unlock()

		case <-sc.stop:
			ticker.Stop()
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

func (sc *inMemoryScheduler) execute(index int) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.tasks[index].status = Running
	fmt.Printf("Executing task %s\n", sc.tasks[index].id)

	sc.tasks[index].fn()

	sc.tasks[index].status = Finished
	fmt.Printf("Finished task %s\n", sc.tasks[index].id)

	sc.taskHistory[sc.tasks[index].originalId] = append(sc.taskHistory[sc.tasks[index].originalId], sc.tasks[index])
}

func (sc *inMemoryScheduler) scheduledTask(fn func(), executeAt time.Time, interval int, originalId string) string {
	id := fmt.Sprintf("task%v-%d", executeAt, len(sc.tasks))
	if originalId == "" {
		originalId = id
	}
	task := Task{
		id:         id,
		fn:         fn,
		executeAt:  executeAt,
		interval:   interval,
		status:     NotStarted,
		originalId: originalId,
	}

	sc.tasks = append(sc.tasks, task)

	return task.id
}
