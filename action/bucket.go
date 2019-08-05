package action

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
)

// MakeBucket will make a bucket with specific name.
func MakeBucket(ctx context.Context) (err error) {
	// Get params from context
	zone := contexts.FromContext(ctx, constants.ZoneFlag).(string)
	remote := contexts.FromContext(ctx, "remote").(string)
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
func RemoveBucket(ctx context.Context) (err error) {
	// Get params from context
	zone := contexts.FromContext(ctx, constants.ZoneFlag).(string)
	remote := contexts.FromContext(ctx, "remote").(string)
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
func ListBuckets(ctx context.Context) (err error) {
	// Get params from context
	zone := contexts.FromContext(ctx, constants.ZoneFlag).(string)

	buckets, err := contexts.Storage.ListBuckets(zone)
	if err != nil {
		return
	}
	for _, b := range buckets {
		fmt.Println(b)
	}
	return nil
}
