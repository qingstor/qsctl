package common

import (
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *BucketCreateTask) run() {
	err := t.GetStorage().PutBucket()
	if err != nil {
		t.TriggerError(fault.NewUnhandled(err))
		return
	}
	log.Debugf("Task <%s> for Bucket <%s> finished.", "BucketCreateTask", t.GetBucketName())
}
