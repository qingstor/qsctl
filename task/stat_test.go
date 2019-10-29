package task

import (
	"fmt"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yunify/qsctl/v2/pkg/mock"

	"github.com/yunify/qsctl/v2/pkg/types"

	"github.com/yunify/qsctl/v2/utils"
)

func TestObjectStatTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	objectKey := uuid.New().String()
	store := mock.NewMockStorager(ctrl)

	store.EXPECT().Stat(gomock.Any(), gomock.Any()).Do(func(inputPath string, option ...*types.Pair) {
		assert.Equal(t, objectKey, inputPath)
	})

	pool := navvy.NewPool(10)

	x := &mockObjectStatTask{}
	x.SetDestinationPath(objectKey)
	x.SetPool(pool)
	x.SetDestinationStorage(store)

	task := NewObjectStatTask(x)
	task.Run()
	pool.Wait()
}
