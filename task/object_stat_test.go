package task

import (
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/storage"

	storageType "github.com/yunify/qsctl/v2/pkg/types/storage"
)

func TestObjectStatTask_Run(t *testing.T) {
	bucketName, objectKey, zone := uuid.New().String(), storage.MockMBObject, "t1"
	store := storage.NewMockObjectStorage()
	err := store.SetupBucket(bucketName, zone)
	if err != nil {
		t.Fatal(err)
	}

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}

	x := &mockObjectStatTask{}
	x.SetKey(objectKey)
	x.SetPool(pool)
	x.SetStorage(store)
	x.SetObjectMeta(&storageType.ObjectMeta{})

	task := NewObjectStatTask(x)
	task.Run()
	pool.Wait()

	cases := []struct {
		input          string
		expectedLength int64
		expectErr      error
	}{
		{storage.Mock0BObject, int64(0), nil},
		{storage.MockMBObject, int64(1024 * 1024), nil},
		{storage.MockGBObject, int64(1024 * 1024 * 1024), nil},
		{storage.MockTBObject, int64(1024 * 1024 * 1024 * 1024), nil},
	}

	for _, v := range cases {
		om, err := store.HeadObject(v.input)
		assert.Equal(t, v.expectErr, err)
		assert.Equal(t, v.expectedLength, om.ContentLength)
	}

}
