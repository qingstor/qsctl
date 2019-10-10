package common

import (
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
)

func TestObjectStatTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	objectKey := uuid.New().String()
	store := mock.NewMockStorager(ctrl)

	store.EXPECT().Stat(gomock.Any(), gomock.Any()).Do(func(inputPath string, option ...*types.Pair) {
		assert.Equal(t, objectKey, inputPath)
	})

	pool := navvy.NewPool(10)

	x := &mockObjectStatTask{}
	x.SetKey(objectKey)
	x.SetPool(pool)
	x.SetDestinationStorage(store)

	task := NewObjectStatTask(x)
	task.Run()
	pool.Wait()
}
