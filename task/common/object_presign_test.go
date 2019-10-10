package common

import (
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
)

func TestObjectPresignTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key, bucketName := uuid.New().String(), uuid.New().String()
	store := mock.NewMockStorager(ctrl)

	pool := navvy.NewPool(10)

	x := &mockObjectPresignTask{}
	x.SetPool(pool)
	x.SetDestinationStorage(store)
	x.SetBucketName(bucketName)
	x.SetKey(key)

	store.EXPECT().Reach(gomock.Any()).Do(func(inputPath string) {
		assert.Equal(t, key, inputPath)
	})

	task := NewObjectPresignTask(x)
	task.Run()
	pool.Wait()
}
