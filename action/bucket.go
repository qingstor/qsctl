package action

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/contexts"
)

// BucketHandler is all params for Bucket func
type BucketHandler struct {
	*FlagHandler
	// Remote is the remote qs path
	Remote string `json:"remote"`
}

// WithZone rewrite the WithZone method
func (sh *BucketHandler) WithZone(z string) *BucketHandler {
	sh.FlagHandler = sh.FlagHandler.WithZone(z)
	return sh
}

// WithRemote sets the Remote field with given remote path
func (sh *BucketHandler) WithRemote(path string) *BucketHandler {
	sh.Remote = path
	return sh
}

// MakeBucket will make a bucket with specific name.
func (sh *BucketHandler) MakeBucket() (err error) {
	// Get params from handler
	zone := sh.GetZone()
	remote := sh.Remote
	// Get bucket name from path
	bucketName, _, err := ParseQsPath(remote)
	if err != nil {
		return
	}
	// Init bucket
	err = contexts.Storage.SetupBucket(bucketName, zone)
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
func (sh *BucketHandler) RemoveBucket() (err error) {
	// Get params from handler
	zone := sh.GetZone()
	remote := sh.Remote
	// Get bucket name from path
	bucketName, _, err := ParseQsPath(remote)
	if err != nil {
		return
	}
	// Init bucket
	err = contexts.Storage.SetupBucket(bucketName, zone)
	if err != nil {
		return err
	}
	if err = contexts.Storage.DeleteBucket(); err != nil {
		return
	}
	log.Infof("Bucket <%s> removed.", bucketName)
	return nil
}

// ListBuckets list all buckets.
func (sh *BucketHandler) ListBuckets() (err error) {
	// Get params from handler
	zone := sh.GetZone()

	buckets, err := contexts.Storage.ListBuckets(zone)
	if err != nil {
		return
	}
	for _, b := range buckets {
		fmt.Println(b)
	}
	return nil
}
