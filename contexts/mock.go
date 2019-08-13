package contexts

import (
	"sync"
	"time"
)

// MockContext mock the CmdContext interface
type MockContext struct {
	sync.Mutex
	Ctx map[interface{}]interface{}
}

// Deadline implements CmdContext interface
func (m *MockContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

// Done implements CmdContext interface
func (m *MockContext) Done() <-chan struct{} {
	return nil
}

// Err implements CmdContext interface
func (m *MockContext) Err() error {
	return nil
}

// Value implements CmdContext interface
func (m *MockContext) Value(key interface{}) interface{} {
	m.Lock()
	defer m.Unlock()
	return m.Ctx[key]
}

// NewMockCmdContext will return a pointer of MockContext with a blank map
func NewMockCmdContext() *MockContext {
	return &MockContext{
		Ctx: make(map[interface{}]interface{}),
	}
}
