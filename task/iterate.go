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

func (t *IterateSegmentTask) new() {}

func (t *IterateSegmentTask) run() {
	listTask := NewListSegment(t)
	t.GetScheduler().Async(listTask)

	for o := range listTask.GetSegmentChannel() {
		x := t.GetSegmentIDScheduleFunc()(t)
		x.SetSegmentID(o.ID)

		t.GetScheduler().Async(x)
	}
}
