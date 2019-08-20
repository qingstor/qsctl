package action

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/task/utils"
)

var stor storage.ObjectStorage

func init() {
	var err error
	stor, err = storage.NewQingStorObjectStorage()
	if err != nil {
		panic(err)
	}
}

// BucketHandler is all params for Bucket func
type BucketHandler struct {
	// Remote is the remote qs path
	Remote string `json:"remote"`
	// Zone specifies the zone for bucket action
	Zone string `json:"zone"`
}

// WithRemote sets the Remote field with given remote path
func (bh *BucketHandler) WithRemote(path string) *BucketHandler {
	bh.Remote = path
	return bh
}

// WithZone sets the Zone field with given zone
func (bh *BucketHandler) WithZone(z string) *BucketHandler {
	bh.Zone = z
	return bh
}

// MakeBucket will make a bucket with specific name.
func (bh *BucketHandler) MakeBucket() (err error) {
	// Get bucket name from path
	bucketName, _, err := utils.ParseQsPath(bh.Remote)
	if err != nil {
		return
	}
	// Init bucket
	err = stor.SetupBucket(bucketName, bh.Zone)
	if err != nil {
		return err
	}
	if err = stor.PutBucket(); err != nil {
		return
	}
	log.Infof("Bucket <%s> created.", bucketName)
	return nil
}

// RemoveBucket remove a bucket with specific remote qs path.
func (bh *BucketHandler) RemoveBucket() (err error) {
	// Get bucket name from path
	bucketName, _, err := utils.ParseQsPath(bh.Remote)
	if err != nil {
		return
	}
	// Init bucket
	err = stor.SetupBucket(bucketName, bh.Zone)
	if err != nil {
		return err
	}
	if err = stor.DeleteBucket(); err != nil {
		return
	}
	log.Infof("Bucket <%s> removed.", bucketName)
	return nil
}

// ListBuckets list all buckets.
func (bh *BucketHandler) ListBuckets() (err error) {
	buckets, err := stor.ListBuckets(bh.Zone)
	if err != nil {
		return
	}
	for _, b := range buckets {
		fmt.Println(b)
	}
	return nil
}
