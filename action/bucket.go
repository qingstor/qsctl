package action

import (
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/contexts"
	"github.com/yunify/qsctl/helper"
)

// MakeBucket will make a bucket with specific name.
func MakeBucket(bucketPath string) (err error) {
	// Get bucket name from path
	bucketName, _, err := ParseQsPath(bucketPath)
	if err != nil {
		return
	}
	// Init bucket
	if err = helper.PutBucket(bucketName, contexts.Zone); err != nil {
		return
	}
	log.Infof("Bucket <%s> created.", bucketName)
	return nil
}
