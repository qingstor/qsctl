package common

import (
	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *BucketCreateTask) run() {
	_, err := t.GetDestinationService().Create(t.GetBucketName(), typ.WithLocation(t.GetZone()))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	log.Debugf("Task <%s> for Bucket <%s> finished.", "BucketCreateTask", t.GetBucketName())
}

func (t *BucketDeleteTask) run() {
	log.Debugf("Task <%s> for Bucket <%s> started.", "BucketDeleteTask", t.GetBucketName())
	err := t.GetDestinationService().Delete(t.GetBucketName())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	log.Debugf("Task <%s> for Bucket <%s> finished.", "BucketDeleteTask", t.GetBucketName())
}

func (t *BucketListTask) run() {
	resp, err := t.GetDestinationService().List(typ.WithLocation(t.GetZone()))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	buckets := make([]string, 0, len(resp))
	for _, v := range resp {
		b, err := v.Metadata()
		if err != nil {
			t.TriggerFault(fault.NewUnhandled(err))
			return
		}
		if name, ok := b.GetName(); ok {
			buckets = append(buckets, name)
		}
	}
	t.SetBucketList(buckets)
	log.Debugf("Task <%s> in zone <%s> finished.", "BucketListTask", t.GetZone())
}

func (t *RemoveBucketForceTask) new() {
	oc := make(chan *typ.Object)
	t.SetObjectChannel(oc)
	// done to notify get object from channel has done
	done := false
	t.SetDone(&done)
	// set recursive for list async task to list recursively
	t.SetRecursive(true)

	t.SetScheduler(types.NewScheduler(NewObjectDeleteScheduledTask))

	t.GetScheduler().Sync(NewObjectListAsyncTask, t)
	t.GetScheduler().Sync(NewObjectDeleteIterateTask, t)
	t.GetScheduler().Sync(NewAbortMultipartTask, t)
	t.GetScheduler().Sync(NewBucketDeleteTask, t)
}
