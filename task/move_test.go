package task

import (
	"fmt"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/pkg/mock"
)

func TestMoveDirTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("normal", func(t *testing.T) {
		sche := mock.NewMockScheduler(ctrl)
		srcStore := mock.NewMockStorager(ctrl)
		dstStore := mock.NewMockStorager(ctrl)

		task := MoveDirTask{}
		task.SetFault(fault.New())
		task.SetPool(navvy.NewPool(10))
		task.SetSourcePath("source")
		task.SetSourceStorage(srcStore)
		task.SetDestinationPath("destination")
		task.SetDestinationStorage(dstStore)
		task.SetScheduler(sche)

		sche.EXPECT().Sync(gomock.Any()).Do(func(task navvy.Task) {
			switch v := task.(type) {
			case *ListDirTask:
				v.validateInput()
			default:
				panic(fmt.Errorf("unexpected task %v", v))
			}
		})

		task.run()
		assert.Empty(t, task.GetFault().Error())
	})
}

func TestMoveFileTask_run(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sche := mock.NewMockScheduler(ctrl)
		srcStore := mock.NewMockStorager(ctrl)
		dstStore := mock.NewMockStorager(ctrl)

		task := MoveFileTask{}
		task.SetFault(fault.New())
		task.SetPool(navvy.NewPool(10))
		task.SetSourcePath("source")
		task.SetSourceStorage(srcStore)
		task.SetDestinationPath("destination")
		task.SetDestinationStorage(dstStore)
		task.SetScheduler(sche)

		sche.EXPECT().Sync(gomock.Any()).Do(func(task navvy.Task) {
			switch v := task.(type) {
			case *CopyFileTask:
				v.validateInput()
			case *DeleteFileTask:
				v.validateInput()
			default:
				panic(fmt.Errorf("unexpected task %v", v))
			}
		}).AnyTimes()

		task.run()
		assert.Empty(t, task.GetFault().Error())
	})
}
