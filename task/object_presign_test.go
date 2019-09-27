package task

import (
	"fmt"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/storage"
)

func TestObjectPresignTask_Run(t *testing.T) {
	store := storage.NewMockObjectStorage()
	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		bucket    string
		objectKey string
		expire    int
		url       string
	}{
		// private bucket url
		{storage.MockZoneAlpha, storage.MockMBObject, 300, fmt.Sprintf("%s.%s/%s?expire=%d",
			storage.MockZoneAlpha,
			storage.MockZoneAlpha,
			storage.MockMBObject,
			300)},
	}

	for _, v := range cases {
		bucket := store.Buckets[v.bucket]
		err = store.SetupBucket(v.bucket, bucket.Location)
		if err != nil {
			t.Fatal(err)
		}

		x := &mockObjectPresignTask{}
		x.SetPool(pool)
		x.SetStorage(store)

		x.SetBucketName(v.bucket)
		x.SetKey(v.objectKey)
		x.SetExpire(v.expire)
		x.SetURL("")

		task := NewObjectPresignTask(x)
		task.Run()
		pool.Wait()

		assert.Equal(t,
			fmt.Sprintf("%v", v.url),
			fmt.Sprintf("%v", x.GetURL()),
		)
	}
}
