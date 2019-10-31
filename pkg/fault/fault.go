package fault

import (
	"strings"
	"sync"
)

// Fault will handle multi error in tasks.
type Fault struct {
	errs []error
	lock sync.RWMutex
}

// New will create a new Fault.
func New() *Fault {
	return &Fault{}
}

// Append will append errors in fault.
func (f *Fault) Append(err ...error) {
	f.lock.Lock()
	f.errs = append(f.errs, err...)
	f.lock.Unlock()
}

// HasError checks whether this fault has error or not.
func (f *Fault) HasError() bool {
	f.lock.RLock()
	defer f.lock.RUnlock()

	return len(f.errs) != 0
}

// Error will print all errors in fault.
func (f *Fault) Error() string {
	f.lock.RLock()
	defer f.lock.RUnlock()

	x := make([]string, 0)
	for _, v := range f.errs {
		x = append(x, v.Error())
	}
	return strings.Join(x, "\n")
}

// Unwrap implements unwarp interface.
func (f *Fault) Unwrap() error {
	f.lock.RLock()
	defer f.lock.RUnlock()

	if len(f.errs) == 0 {
		return nil
	}
	return f.errs[0]
}
