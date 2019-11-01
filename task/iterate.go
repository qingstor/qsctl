package task

import (
	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"
)

func (t *IterateFileTask) new() {
	oc := make(chan *typ.Object)
	t.SetObjectChannel(oc)
}

func (t *IterateFileTask) run() {
	log.Debugf("Task <%s> for path <%s> started",
		"IterateFileTask", t.GetPath())
	t.GetScheduler().Async(NewListFileTask(t))

	for o := range t.GetObjectChannel() {
		x := t.GetPathScheduleFunc()(t)
		x.SetPath(o.Name)

		t.GetScheduler().Async(x)
	}
	log.Debugf("Task <%s> for path <%s> finished",
		"IterateFileTask", t.GetPath())
}
