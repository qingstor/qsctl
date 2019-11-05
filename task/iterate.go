package task

func (t *IterateFileTask) new() {}

func (t *IterateFileTask) run() {
	listTask := NewListFile(t)
	t.GetScheduler().Async(listTask)

	for o := range listTask.GetObjectChannel() {
		t.GetPathFunc()(o.Name)
	}
}

func (t *IterateSegmentTask) new() {}

func (t *IterateSegmentTask) run() {
	listTask := NewListSegment(t)
	t.GetScheduler().Async(listTask)

	for o := range listTask.GetSegmentChannel() {
		t.GetSegmentIDFunc()(o.ID)
	}
}
