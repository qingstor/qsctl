package action

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/contexts"
)

// MakeBucket will make a bucket with specific name.
func MakeBucket(name string) (err error) {
	bucket, err := contexts.Service.Bucket(name, contexts.Zone)
	if err != nil {
		return
	}
	putBucketOutput, err := bucket.Put()
	if err != nil {
		return
	}
	log.Debugf("put code: %d\n", *putBucketOutput.StatusCode)
	fmt.Printf("bucket: %s created successfully", name)
	return nil
}
