package task

func (t *IterateFileTask) new() {}

func (t *IterateFileTask) run() {
	listTask := NewListFile(t)
	t.GetScheduler().Async(listTask)

	for o := range listTask.GetObjectChannel() {
		x := t.GetPathScheduleFunc()(t)
		x.SetPath(o.Name)

		t.GetScheduler().Async(x)
	}
}
