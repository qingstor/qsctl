package common

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestObjectPresignTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// store := NewMockStorager(ctrl)
	//
	// pool, err := navvy.NewPool(10)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// cases := []struct {
	// 	bucket    string
	// 	objectKey string
	// 	expire    int
	// 	url       string
	// }{
	// 	// private bucket url
	// 	{storage.MockZoneAlpha, storage.MockMBObject, 300, fmt.Sprintf("%s.%s/%s?expire=%d",
	// 		storage.MockZoneAlpha,
	// 		storage.MockZoneAlpha,
	// 		storage.MockMBObject,
	// 		300)},
	// }

	// for _, v := range cases {
	// 	bucket := store.Buckets[v.bucket]
	// 	err = store.SetupBucket(v.bucket, bucket.Location)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	//
	// 	x := &mockObjectPresignTask{}
	// 	x.SetPool(pool)
	// 	x.SetDestinationStorage(store)
	//
	// 	x.SetBucketName(v.bucket)
	// 	x.SetKey(v.objectKey)
	// 	x.SetExpire(v.expire)
	//
	// 	task := NewObjectPresignTask(x)
	// 	task.Run()
	// 	pool.Wait()
	//
	// 	assert.Equal(t, v.url, x.GetURL())
	// }
}
