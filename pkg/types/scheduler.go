package types

import (
	"sync"

	"github.com/Xuanwo/navvy"
)

// schedulable is the task that can be used in RealScheduler.
type schedulable interface {
	navvy.Task

	IDGetter
	FaultValidator
	FaultGetter
	PoolGetter
}

type scheduler interface {
	New(Todoist)
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

	fn TodoFunc
}

// NewScheduler will create a new RealScheduler.
func NewScheduler(fn TodoFunc) *RealScheduler {
	return &RealScheduler{
		meta: make(map[string]schedulable),
		wg:   &sync.WaitGroup{},
		fn:   fn,
	}
}

// New will create a new task.
func (s *RealScheduler) New(t Todoist) {
	x := s.fn(t)
	v := x.(schedulable)

	s.meta[v.GetID()] = v
	s.wg.Add(1)

	go v.GetPool().Submit(x)
}

// Done will mark a task as done.
func (s *RealScheduler) Done(taskID string) {
	t := s.meta[taskID]
	if t.ValidateFault() {
		s.errs = append(s.errs, t.GetFault())
	}

	delete(s.meta, taskID)
	s.wg.Done()
}

// Wait will wait until a task finished.
func (s *RealScheduler) Wait() {
	s.wg.Wait()
}

// Errors will return all errors.
func (s *RealScheduler) Errors() []error {
	return s.errs
}

// MockScheduler is a mock of scheduler to be test.
type MockScheduler struct {
	wg *sync.WaitGroup
}

// NewMockScheduler will create a new mock RealScheduler.
func NewMockScheduler(fn TodoFunc) *MockScheduler {
	return &MockScheduler{
		wg: &sync.WaitGroup{},
	}
}

// New implements scheduler.New
func (m MockScheduler) New(Todoist) {
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
