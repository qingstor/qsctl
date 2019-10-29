package types

import (
	"sync"

	"github.com/Xuanwo/navvy"
)

// TaskFunc will be used create a new task.
type TaskFunc func(navvy.Task) navvy.Task

// schedulable is the task that can be used in RealScheduler.
type schedulable interface {
	IDGetter
	FaultValidator
	FaultGetter
}

type scheduler interface {
	Sync(task navvy.Task, fn TaskFunc)
	Async(task navvy.Task, fn TaskFunc)

	Done(string)
	Wait()
	Errors() []error
}

// RealScheduler will hold the task's sub tasks.
// TODO: we need a better way to handle type name conflict.
type RealScheduler struct {
	meta map[string]schedulable
	wg   *sync.WaitGroup
	errs []error
	pool *navvy.Pool

	lock sync.Mutex
}

// NewScheduler will create a new RealScheduler.
func NewScheduler(pool *navvy.Pool) *RealScheduler {
	return &RealScheduler{
		meta: make(map[string]schedulable),
		wg:   &sync.WaitGroup{},
		pool: pool,
	}
}

// New will create a new task.
func (s *RealScheduler) Sync(task navvy.Task, fn TaskFunc) {
	s.Wait()
	s.Async(task, fn)
}

// New will create a new task.
func (s *RealScheduler) Async(task navvy.Task, fn TaskFunc) {
	v := fn(task)
	s.pool.Submit(v)

	sch := v.(schedulable)
	s.lock.Lock()
	s.meta[sch.GetID()] = sch
	s.lock.Unlock()
}

// Done will mark a task as done.
func (s *RealScheduler) Done(taskID string) {
	s.lock.Lock()
	t := s.meta[taskID]
	if t.ValidateFault() {
		s.errs = append(s.errs, t.GetFault())
	}
	delete(s.meta, taskID)
	s.lock.Unlock()

	s.wg.Done()
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

// MockScheduler is a mock of scheduler to be test.
type MockScheduler struct {
	wg *sync.WaitGroup
}

// NewMockScheduler will create a new mock RealScheduler.
func NewMockScheduler(fn TaskFunc) *MockScheduler {
	return &MockScheduler{
		wg: &sync.WaitGroup{},
	}
}

// New implements scheduler.New
func (m MockScheduler) New(taskFunc TaskFunc) {
	m.wg.Add(1)
}

// New will create a new task.
func (m *MockScheduler) Sync(task navvy.Task, fn TaskFunc) {
	m.wg.Add(1)
}

// New will create a new task.
func (m *MockScheduler) Async(task navvy.Task, fn TaskFunc) {
	m.wg.Add(1)
}

// Done implements scheduler.Done
func (m MockScheduler) Done(string) {
	m.wg.Done()
}

// Wait implements scheduler.Wait
func (m MockScheduler) Wait() {
	m.wg.Wait()
}

// Errors implements scheduler.Errors
func (m MockScheduler) Errors() []error {
	return nil
}
