package common

import (
	"github.com/yunify/qsctl/v2/task/utils"
)

// Run implement navvy.Task.
func (t *WaitTask) Run() {
	t.GetWaitGroup().Wait()

	utils.SubmitNextTask(t.WaitTaskRequirement)
}
