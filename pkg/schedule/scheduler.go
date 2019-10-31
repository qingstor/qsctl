package schedule

import (
	"sync"

	"github.com/Xuanwo/navvy"
	"github.com/yunify/qsctl/v2/pkg/fault"
)

// TaskFunc will be used create a new task.
type TaskFunc func(navvy.Task) navvy.Task

// Scheduler will schedule tasks.
//go:generate mockgen -package mock -destination ../mock/scheduler.go github.com/yunify/qsctl/v2/pkg/schedule Scheduler
type Scheduler interface {
	Sync(task navvy.Task, fn TaskFunc)
	Async(task navvy.Task, fn TaskFunc)

	Wait()
}

// Schedulable is the task that can be used in RealScheduler.
type Schedulable interface {
	navvy.Task

	GetID() string
	GetFault() *fault.Fault
}

type task struct {
	s *RealScheduler
	t Schedulable
}

func newTask(s *RealScheduler, t Schedulable) *task {
	return &task{
		s: s,
		t: t,
	}
}

func (t *task) Run() {
	defer func() {
		t.s.wg.Done()
	}()

	t.t.Run()
}

// RealScheduler will hold the task's sub tasks.
type RealScheduler struct {
	wg   *sync.WaitGroup
	pool *navvy.Pool
}

// NewScheduler will create a new RealScheduler.
func NewScheduler(pool *navvy.Pool) *RealScheduler {
	return &RealScheduler{
		wg:   &sync.WaitGroup{},
		pool: pool,
	}
}

// Sync will create a new task after wait.
func (s *RealScheduler) Sync(task navvy.Task, fn TaskFunc) {
	s.Wait()
	s.Async(task, fn)
}

// Async will create a new task immediately.
func (s *RealScheduler) Async(task navvy.Task, fn TaskFunc) {
	// Don't submit to pool if task has an error.
	t := fn(task).(Schedulable)
	if t.GetFault().HasError() {
		return
	}

	s.wg.Add(1)
	s.pool.Submit(newTask(s, t))
}

// Wait will wait until a task finished.
func (s *RealScheduler) Wait() {
	s.wg.Wait()
}
