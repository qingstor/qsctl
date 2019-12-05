package task

import (
	"testing"

	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/pkg/mock"
)

func TestStatFileTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStorager(ctrl)
	expectedPath := uuid.New().String()

	task := StatFileTask{}
	task.SetFault(fault.New())
	task.SetStorage(store)
	task.SetPath(expectedPath)

	store.EXPECT().Stat(gomock.Any()).DoAndReturn(func(path string) (o *types.Object, err error) {
		assert.Equal(t, expectedPath, path)
		return &types.Object{}, nil
	})

	task.run()
	assert.Empty(t, task.GetFault().Error())
	assert.NotNil(t, task.GetObject())
}
