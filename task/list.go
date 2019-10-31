package task

import (
	"errors"

	"github.com/Xuanwo/storage/pkg/iterator"
	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *ListFileTask) new() {}

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
			t.TriggerFault(types.NewErrUnhandled(err))
			return
		}
		t.GetObjectChannel() <- o
	}

	log.Debugf("Task <%s> for key <%s> finished", "ObjectListTask", t.GetPath())
}

func (t *ListStorageTask) new() {}
func (t *ListStorageTask) run() {
	resp, err := t.GetService().List(typ.WithLocation(t.GetZone()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	buckets := make([]string, 0, len(resp))
	for _, v := range resp {
		b, err := v.Metadata()
		if err != nil {
			t.TriggerFault(types.NewErrUnhandled(err))
			return
		}
		if name, ok := b.GetName(); ok {
			buckets = append(buckets, name)
		}
	}
	t.SetBucketList(buckets)
	log.Debugf("Task <%s> in zone <%s> finished.", "BucketListTask", t.GetZone())
}
