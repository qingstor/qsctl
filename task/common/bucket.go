package common

import (
	"github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *BucketCreateTask) run() {
	err := t.GetDestinationStorage().CreateDir(t.GetBucketName(), types.WithLocation(t.GetZone()))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	log.Debugf("Task <%s> for Bucket <%s> finished.", "BucketCreateTask", t.GetBucketName())
}

func (t *BucketDeleteTask) run() {
	err := t.GetDestinationStorage().Delete(t.GetBucketName())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	log.Debugf("Task <%s> for Bucket <%s> finished.", "BucketDeleteTask", t.GetBucketName())
}
