package task

import (
	"errors"

	"github.com/Xuanwo/storage/pkg/iterator"
	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *ListFileTask) new() {

}

func (t *ListFileTask) run() {
	log.Debugf("Task <%s> for key <%s> started", "ObjectListTask", t.GetPath())

	pairs := make([]*typ.Pair, 0)

	if !t.GetRecursive() {
		pairs = append(pairs, typ.WithDelimiter("/"))
	}

	it := t.GetStorage().ListDir(t.GetPath(), pairs...)

	// Always close the object channel.
	defer close(t.GetObjectChannel())

	for {
		o, err := it.Next()
		if err != nil && errors.Is(err, iterator.ErrDone) {
			break
		}
		if err != nil {
			t.TriggerFault(fault.NewUnhandled(err))
			return
		}
		t.GetObjectChannel() <- o
	}

	log.Debugf("Task <%s> for key <%s> finished", "ObjectListTask", t.GetPath())
}

func (t *IterateFileTask) new() {
	oc := make(chan *typ.Object)
	t.SetObjectChannel(oc)
}

func (t *IterateFileTask) run() {
	t.GetScheduler().Async(t, NewListFileTask)

	for o := range t.GetObjectChannel() {
		x := NewIterateFile(t)
		x.SetPath(o.Name)

		t.GetScheduler().Async(x, t.GetScheduleFunc())
	}
}
