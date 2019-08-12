package action

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/contexts"
)

// BucketHandler is all params for Bucket func
type BucketHandler struct {
	BaseHandler
	// Remote is the remote qs path
	Remote string `json:"remote"`
}

// WithZone rewrite the WithZone method
func (bh *BucketHandler) WithZone(z string) *BucketHandler {
	bh.Zone = z
	return bh
}

// WithRemote sets the Remote field with given remote path
func (bh *BucketHandler) WithRemote(path string) *BucketHandler {
	bh.Remote = path
	return bh
}

// MakeBucket will make a bucket with specific name.
func (bh *BucketHandler) MakeBucket() (err error) {
	// Get bucket name from path
	bucketName, _, err := ParseQsPath(bh.Remote)
	if err != nil {
		return
	}
	// Init bucket
	err = contexts.Storage.SetupBucket(bucketName, bh.Zone)
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
func (bh *BucketHandler) RemoveBucket() (err error) {
	// Get bucket name from path
	bucketName, _, err := ParseQsPath(bh.Remote)
	if err != nil {
		return
	}
	// Init bucket
	err = contexts.Storage.SetupBucket(bucketName, bh.Zone)
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
func (bh *BucketHandler) ListBuckets() (err error) {
	buckets, err := contexts.Storage.ListBuckets(bh.Zone)
	if err != nil {
		return
	}
	for _, b := range buckets {
		fmt.Println(b)
	}
	return nil
}
