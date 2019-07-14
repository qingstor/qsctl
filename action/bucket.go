package action

import (
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/contexts"
)

// MakeBucket will make a bucket with specific name.
func MakeBucket(remote string) (err error) {
	// Get bucket name from path
	bucketName, _, err := ParseQsPath(remote)
	if err != nil {
		return
	}
	// Init bucket
	err = contexts.Storage.SetupBucket(bucketName, contexts.Zone)
	if err != nil {
		return err
	}
	if err = contexts.Storage.PutBucket(); err != nil {
		return
	}
	log.Infof("Bucket <%s> created.", bucketName)
	return nil
}

// RemoveBucket remove a bucket with specific remote qs path.
func RemoveBucket(remote string) (err error) {
	// Get bucket name from path
	bucketName, _, err := ParseQsPath(remote)
	if err != nil {
		return
	}
	// Init bucket
	err = contexts.Storage.SetupBucket(bucketName, contexts.Zone)
	if err != nil {
		return err
	}
	if err = contexts.Storage.DeleteBucket(); err != nil {
		return
	}
	log.Infof("Bucket <%s> removed.", bucketName)
	return nil
}
