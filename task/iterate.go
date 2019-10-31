package task

import (
	typ "github.com/Xuanwo/storage/types"
)

func (t *IterateFileTask) new() {
	oc := make(chan *typ.Object)
	t.SetObjectChannel(oc)
}

func (t *IterateFileTask) run() {
	t.GetScheduler().Async(t, NewListFileTask)

	for o := range t.GetObjectChannel() {
		x := NewFileShim(t)
		x.SetPath(o.Name)

		t.GetScheduler().Async(x, t.GetScheduleFunc())
	}
}
