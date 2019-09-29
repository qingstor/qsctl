package common

import (
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/storage"
)

func TestBucketCreateTask_Run(t *testing.T) {
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

	x := &mockBucketCreateTask{}
	x.SetBucketName(bucketName)
	x.SetPool(pool)
	x.SetStorage(store)

	task := NewBucketCreateTask(x)
	task.Run()
	pool.Wait()

	buckets, err := store.ListBuckets(zone)
	assert.NoError(t, err)
	assert.Contains(t, buckets, bucketName)
}

func TestBucketDeleteTask_Run(t *testing.T) {
	bucketName, zone := storage.MockZoneAlpha, storage.MockZoneAlpha
	store := storage.NewMockObjectStorage()
	// get ori bucket counts
	oriBucketCount := len(store.Buckets)

	err := store.SetupBucket(bucketName, zone)
	if err != nil {
		t.Fatal(err)
	}

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}

	x := &mockBucketDeleteTask{}
	x.SetBucketName(bucketName)
	x.SetPool(pool)
	x.SetStorage(store)

	task := NewBucketDeleteTask(x)
	task.Run()
	pool.Wait()
	// list all buckets from different zones
	buckets, err := store.ListBuckets("")
	assert.NoError(t, err)
	// buckets not contains removed bucket's name
	assert.NotContains(t, buckets, bucketName)
	// bucket count reduce 1
	assert.Equal(t, oriBucketCount-1, len(buckets))
}
