package task

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Xuanwo/navvy"
	typ "github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yunify/qsctl/v2/pkg/schedule"
	"github.com/yunify/qsctl/v2/pkg/types"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/mock"
)

func TestCopyFileTask_new(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	paths := make([]string, 100)
	for k := range paths {
		paths[k] = uuid.New().String()
	}
	tests := []struct {
		name string
		size int64
	}{
		{
			"normal",
			constants.MaximumAutoMultipartSize - 1,
		},
	}

	for k, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			srcStore := mock.NewMockStorager(ctrl)
			srcStore.EXPECT().Stat(gomock.Any()).DoAndReturn(func(inputPath string) (o *typ.Object, err error) {
				assert.Equal(t, paths[k], inputPath)
				return &typ.Object{
					Name: inputPath,
					Type: typ.ObjectTypeFile,
					Metadata: typ.Metadata{
						typ.Size: v.size,
					},
				}, nil
			})

			m := &types.MockCopyFileTask{}
			m.SetSourceStorage(srcStore)
			m.SetSourcePath(paths[k])
			task := &CopyFileTask{CopyFileRequirement: m}
			task.new()

			assert.Equal(t, v.size, task.GetTotalSize())
		})
	}
}

func TestCopyFileTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		size int64
		fn   schedule.TaskFunc
	}{
		{
			"large file",
			constants.MaximumAutoMultipartSize + 1,
			NewCopyLargeFileTask,
		},
		{
			"small file",
			constants.MaximumAutoMultipartSize - 1,
			NewCopySmallFileTask,
		},
	}
	for _, v := range tests {
		m := &types.MockCopyFileTask{}
		m.SetPool(navvy.NewPool(10))
		task := &CopyFileTask{CopyFileRequirement: m}
		task.SetTotalSize(v.size)

		sch := mock.NewMockScheduler(ctrl)
		sch.EXPECT().Sync(gomock.Any()).Do(func(inputTask navvy.Task) {
			assert.Equal(t, reflect.TypeOf(v.fn(task)), reflect.TypeOf(inputTask))
		})
		sch.EXPECT().Wait().Do(func() {})
		task.SetScheduler(sch)

		task.Run()
	}
}

func TestCopyLargeFileTask_new(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := &types.MockCopyLargeFileTask{}
	m.SetPool(navvy.NewPool(10))
	m.SetTotalSize(1024)

	task := NewCopyLargeFile(m)

	assert.True(t, task.ValidateScheduleFunc())
	assert.Equal(t,
		fmt.Sprint(schedule.TaskFunc(NewCopyPartialFileTask)),
		fmt.Sprint(task.GetScheduleFunc()))
	assert.True(t, task.ValidatePartSize())
}

func TestCopyLargeFileTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	name := uuid.New().String()
	key := uuid.New().String()
	size := int64(1234)

	srcStore := mock.NewMockStorager(ctrl)
	dstStore := mock.NewMockStorager(ctrl)

	pool := navvy.NewPool(10)

	x := &types.MockCopyLargeFileTask{}
	x.SetPool(pool)
	x.SetSourcePath(name)
	x.SetSourceStorage(srcStore)
	x.SetDestinationPath(key)
	x.SetDestinationStorage(dstStore)
	x.SetTotalSize(size)

	task := NewCopyLargeFileTask(x)
	tt := task.(*CopyLargeFileTask)
	assert.Equal(t, int64(constants.DefaultPartSize), tt.GetPartSize())
	assert.NotNil(t, tt.GetScheduler())
}

func TestCopyStreamTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStorager(ctrl)
	key := uuid.New().String()

	pool := navvy.NewPool(10)

	x := &types.MockCopyStreamTask{}
	x.SetDestinationStorage(store)
	x.SetDestinationPath(key)
	x.SetSourcePath("-")
	x.SetPool(pool)

	task := NewCopyStreamTask(x)

	tt := task.(*CopyStreamTask)
	assert.NotNil(t, tt.GetBytesPool())
	assert.Equal(t, int64(constants.DefaultPartSize), tt.GetPartSize())
	assert.NotNil(t, tt.GetScheduler())
	assert.Equal(t, int64(0), *tt.GetCurrentOffset())
	assert.Equal(t, int64(-1), tt.GetTotalSize())
}
