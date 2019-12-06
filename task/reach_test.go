package task

import (
	"testing"

	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/qingstor/qsctl/v2/pkg/fault"
	"github.com/qingstor/qsctl/v2/pkg/mock"
)

func TestReachFileTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reacher := mock.NewMockReacher(ctrl)
	reachPath := uuid.New().String()
	reachExpire := 1024
	reachedURL := uuid.New().String()

	task := ReachFileTask{}
	task.SetFault(fault.New())
	task.SetReacher(reacher)
	task.SetPath(reachPath)
	task.SetExpire(reachExpire)

	reacher.EXPECT().Reach(gomock.Any(), gomock.Any()).DoAndReturn(func(path string, pairs ...*types.Pair) (url string, err error) {
		assert.Equal(t, reachPath, path)
		assert.Equal(t, reachExpire, pairs[0].Value.(int))
		return reachedURL, nil
	})

	task.run()
	assert.Empty(t, task.GetFault().Error())
	assert.Equal(t, reachedURL, task.GetURL())
}
