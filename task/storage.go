package task

import (
	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *CreateStorageTask) new() {}

func (t *CreateStorageTask) run() {
	_, err := t.GetService().Create(t.GetStorageName(), typ.WithLocation(t.GetZone()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	log.Debugf("Task <%s> for Bucket <%s> finished.", "BucketCreateTask", t.GetStorageName())
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
