package common

func (t *DoneSchedulerTask) run() {
	t.GetScheduler().Done(t.GetID())
}
