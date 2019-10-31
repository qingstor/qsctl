package task

import (
	typ "github.com/Xuanwo/storage/types"
)

func (t *IterateFileTask) new() {
	oc := make(chan *typ.Object)
	t.SetObjectChannel(oc)
}

func (t *IterateFileTask) run() {
	t.GetScheduler().Async(NewListFileTask(t))

	for o := range t.GetObjectChannel() {
		x := NewFileShim(t)
		x.SetPath(o.Name)

		t.GetScheduler().Async(t.GetScheduleFunc()(x))
	}
}
