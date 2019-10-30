package fault

import (
	"strings"
	"sync"
)

type Fault struct {
	errs []error
	lock sync.RWMutex
}

func New() *Fault {
	return &Fault{}
}

func (f *Fault) Append(err ...error) {
	f.lock.Lock()
	f.errs = append(f.errs, err...)
	f.lock.Unlock()
}

func (f *Fault) HasError() bool {
	f.lock.RLock()
	defer f.lock.RUnlock()

	return len(f.errs) != 0
}

func (f *Fault) Error() string {
	f.lock.RLock()
	defer f.lock.RUnlock()

	x := make([]string, 0)
	for _, v := range f.errs {
		x = append(x, v.Error())
	}
	return strings.Join(x, "\n")
}

func (f *Fault) Unwrap() error {
	f.lock.RLock()
	defer f.lock.RUnlock()

	if len(f.errs) == 0 {
		return nil
	}
	return f.errs[0]
}
