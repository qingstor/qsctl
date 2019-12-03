package task

import (
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
)

func TestCopyDirTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("normal", func(t *testing.T) {
		sche := mock.NewMockScheduler(ctrl)
		srcStore := mock.NewMockStorager(ctrl)
		dstStore := mock.NewMockStorager(ctrl)

		task := CopyDirTask{}
		task.SetPool(navvy.NewPool(10))
		task.SetSourcePath("source")
		task.SetSourceStorage(srcStore)
		task.SetDestinationPath("destination")
		task.SetDestinationStorage(dstStore)
		task.SetScheduler(sche)

		sche.EXPECT().Sync(gomock.Any()).Do(func(task navvy.Task) {
			_, ok := task.(*ListDirTask)
			assert.True(t, ok)
		})
		task.run()
	})
}
