package common

// Run implement navvy.Task.
func (t *WaitTask) run() {
	t.GetWaitGroup().Wait()
}
