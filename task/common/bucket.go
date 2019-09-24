package common

import (
	log "github.com/sirupsen/logrus"
)

func (t *BucketCreateTask) run() {
	err := t.GetStorage().PutBucket()
	if err != nil {
		panic(err)
	}
	log.Debugf("Task <%s> for Bucket <%s> finished.", "BucketCreateTask", t.GetBucketName())
}
