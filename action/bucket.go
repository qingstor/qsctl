package action

import (
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/contexts"
)

// MakeBucket will make a bucket with specific name.
func MakeBucket(remote string) (err error) {
	// Get bucket name from path
	bucketName, _, err := ParseQsPath(remote)
	if err != nil {
		return
	}
	// Init bucket
	if err = contexts.Storage.SetupBucket(bucketName, contexts.Zone); err != nil {
		return
	}
	log.Infof("Bucket <%s> created.", bucketName)
	return nil
}
