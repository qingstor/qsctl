package common

import (
	"fmt"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

// Run implement navvy.Task.
func (t *WaitTask) run() {
	t.GetScheduler().Wait()
	errs := t.GetScheduler().Errors()
	if len(errs) != 0 {
		t.SetFault(fault.NewUnhandled(fmt.Errorf("%v", errs)))
	}
}
