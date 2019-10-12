package common

import (
	"github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *BucketCreateTask) run() {
	_, err := t.GetDestinationService().Create(t.GetBucketName(), types.WithLocation(t.GetZone()))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	log.Debugf("Task <%s> for Bucket <%s> finished.", "BucketCreateTask", t.GetBucketName())
}

func (t *BucketDeleteTask) run() {
	err := t.GetDestinationService().Delete(t.GetBucketName())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	log.Debugf("Task <%s> for Bucket <%s> finished.", "BucketDeleteTask", t.GetBucketName())
}

func (t *BucketListTask) run() {
	resp, err := t.GetDestinationService().List(types.WithLocation(t.GetZone()))
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
