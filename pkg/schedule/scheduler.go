package schedule

import (
	"sync"

	"github.com/Xuanwo/navvy"
)

// TaskFunc will be used create a new task.
type TaskFunc func(navvy.Task) navvy.Task

// Scheduler will schedule tasks.
//go:generate mockgen -package mock -destination ../mock/scheduler.go github.com/yunify/qsctl/v2/pkg/schedule Scheduler
type Scheduler interface {
	Sync(task navvy.Task)
	Async(task navvy.Task)

	Wait()
}

type task struct {
	s *RealScheduler
	t navvy.Task
	c *sync.Cond
}

func newTask(s *RealScheduler, t navvy.Task) *task {
	return &task{
		s: s,
		t: t,
	}
}

func newSyncTask(s *RealScheduler, t navvy.Task) *task {
	lock := &sync.Mutex{}
	lock.Lock()

	return &task{
		s: s,
		t: t,
		c: sync.NewCond(lock),
	}
}

func (t *task) Run() {
	defer func() {
		t.s.wg.Done()
		if t.c != nil {
			t.c.Broadcast()
		}
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

// Sync will return after this task finished.
func (s *RealScheduler) Sync(task navvy.Task) {
	s.wg.Add(1)
	t := newSyncTask(s, task)
	s.pool.Submit(t)
	t.c.Wait()
}

// Async will create a new task immediately.
func (s *RealScheduler) Async(task navvy.Task) {
	s.wg.Add(1)
	s.pool.Submit(newTask(s, task))
}

// Wait will wait until a task finished.
func (s *RealScheduler) Wait() {
	s.wg.Wait()
}
