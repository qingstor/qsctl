package task

import (
	"fmt"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/qingstor/qsctl/v2/pkg/fault"
	"github.com/qingstor/qsctl/v2/pkg/mock"
)

func TestSyncTask_run(t *testing.T) {
	t.Run("without ignore existing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sche := mock.NewMockScheduler(ctrl)
		srcStore := mock.NewMockStorager(ctrl)
		dstStore := mock.NewMockStorager(ctrl)
		sourcePath := uuid.New().String()
		dstPath := uuid.New().String()

		task := SyncTask{}
		task.SetPool(navvy.NewPool(10))
		task.SetScheduler(sche)
		task.SetFault(fault.New())
		task.SetSourcePath(sourcePath)
		task.SetSourceStorage(srcStore)
		task.SetDestinationStorage(dstStore)
		task.SetDestinationPath(dstPath)
		task.SetIgnoreExisting(false)

		sche.EXPECT().Sync(gomock.Any()).Do(func(task navvy.Task) {
			switch v := task.(type) {
			case *ListDirTask:
				v.validateInput()
			default:
				panic(fmt.Errorf("unexpected task %v", v))
			}
		}).AnyTimes()

		task.run()
		assert.Empty(t, task.GetFault().Error())
	})

	t.Run("with ignore existing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sche := mock.NewMockScheduler(ctrl)
		srcStore := mock.NewMockStorager(ctrl)
		dstStore := mock.NewMockStorager(ctrl)
		sourcePath := uuid.New().String()
		dstPath := uuid.New().String()

		task := SyncTask{}
		task.SetPool(navvy.NewPool(10))
		task.SetScheduler(sche)
		task.SetFault(fault.New())
		task.SetSourcePath(sourcePath)
		task.SetSourceStorage(srcStore)
		task.SetDestinationStorage(dstStore)
		task.SetDestinationPath(dstPath)
		task.SetIgnoreExisting(true)

		sche.EXPECT().Sync(gomock.Any()).Do(func(task navvy.Task) {
			switch v := task.(type) {
			case *ListDirTask:
				v.validateInput()
			default:
				panic(fmt.Errorf("unexpected task %v", v))
			}
		}).AnyTimes()

		task.run()
		assert.Empty(t, task.GetFault().Error())
	})
}
