package task

import (
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/storage"
)

func TestPutBucketTask_Run(t *testing.T) {
	bucketName, zone := uuid.New().String(), "t1"
	store := storage.NewMockObjectStorage()
	err := store.SetupBucket(bucketName, zone)
	if err != nil {
		t.Fatal(err)
	}

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}

	x := &mockPutBucketTask{}
	x.SetBucketName(bucketName)
	x.SetPool(pool)
	x.SetStorage(store)

	task := NewPutBucketTask(x)
	task.Run()
	pool.Wait()

	buckets, err := store.ListBuckets(zone)
	assert.NoError(t, err)
	assert.Contains(t, buckets, bucketName)
}
