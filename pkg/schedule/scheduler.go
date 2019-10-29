package schedule

import (
	"sync"

	"github.com/Xuanwo/navvy"
)

// TaskFunc will be used create a new task.
type TaskFunc func(navvy.Task) navvy.Task

// Scheduler will schedule tasks.
type Scheduler interface {
	Sync(task navvy.Task, fn TaskFunc)
	Async(task navvy.Task, fn TaskFunc)

	Wait()
	Errors() []error
}

// Schedulable is the task that can be used in RealScheduler.
type Schedulable interface {
	navvy.Task

	GetID() string
	ValidateFault() bool
	GetFault() error
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
		if t.t.ValidateFault() {
			t.s.errs = append(t.s.errs, t.t.GetFault())
		}
		t.s.wg.Done()
	}()

	t.t.Run()
}

// RealScheduler will hold the task's sub tasks.
type RealScheduler struct {
	wg   *sync.WaitGroup
	errs []error
	pool *navvy.Pool

	lock sync.Mutex
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
	s.pool.Submit(newTask(s, fn(task).(Schedulable)))
}

// Wait will wait until a task finished.
func (s *RealScheduler) Wait() {
	s.wg.Wait()
}

// Errors will return all errors.
func (s *RealScheduler) Errors() []error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.errs
}
