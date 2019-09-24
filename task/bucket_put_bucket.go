package task

import (
	log "github.com/sirupsen/logrus"
)

func (t *PutBucketTask) run() {
	err := t.GetStorage().PutBucket()
	if err != nil {
		panic(err)
	}
	log.Debugf("Task <%s> for Bucket <%s> finished.", "PutBucketTask", t.GetBucketName())
}
