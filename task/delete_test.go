package task

import (
	"fmt"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/qingstor/qsctl/v2/pkg/fault"
	"github.com/qingstor/qsctl/v2/pkg/mock"
)

func TestDeleteFileTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStorager(ctrl)
	testPath := uuid.New().String()

	task := DeleteFileTask{}
	task.SetFault(fault.New())
	task.SetStorage(store)
	task.SetPath(testPath)

	store.EXPECT().Delete(gomock.Any()).Do(func(name string) error {
		assert.Equal(t, testPath, name)
		return nil
	})

	task.run()
	assert.Empty(t, task.GetFault().Error())
}

func TestDeleteDirTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStorager(ctrl)
	sche := mock.NewMockScheduler(ctrl)
	path := uuid.New().String()

	task := DeleteDirTask{}
	task.SetFault(fault.New())
	task.SetPool(navvy.NewPool(10))
	task.SetScheduler(sche)
	task.SetPath(path)
	task.SetStorage(store)

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
}

func TestDeleteSegmentTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	segmenter := mock.NewMockSegmenter(ctrl)
	segmentID := uuid.New().String()

	task := DeleteSegmentTask{}
	task.SetFault(fault.New())
	task.SetSegmenter(segmenter)
	task.SetSegmentID(segmentID)

	segmenter.EXPECT().AbortSegment(gomock.Any()).Do(func(id string) error {
		assert.Equal(t, segmentID, id)
		return nil
	})

	task.run()
	assert.Empty(t, task.GetFault().Error())
}

func TestNewDeleteStorageTask(t *testing.T) {
	t.Run("delete without force", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		srv := mock.NewMockServicer(ctrl)
		storageName := uuid.New().String()

		task := DeleteStorageTask{}
		task.SetFault(fault.New())
		task.SetService(srv)
		task.SetStorageName(storageName)
		task.SetForce(false)

		srv.EXPECT().Delete(gomock.Any()).Do(func(name string) error {
			assert.Equal(t, storageName, name)
			return nil
		})

		task.run()
		assert.Empty(t, task.GetFault().Error())
	})

	t.Run("delete with force", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sche := mock.NewMockScheduler(ctrl)
		srv := mock.NewMockServicer(ctrl)
		store := mock.NewMockStorager(ctrl)
		storageName := uuid.New().String()

		task := DeleteStorageTask{}
		task.SetFault(fault.New())
		task.SetService(srv)
		task.SetStorageName(storageName)
		task.SetForce(true)
		task.SetPool(navvy.NewPool(10))
		task.SetScheduler(sche)

		srv.EXPECT().Delete(gomock.Any()).Do(func(name string) error {
			assert.Equal(t, storageName, name)
			return nil
		})
		srv.EXPECT().Get(gomock.Any()).DoAndReturn(func(name string, pairs ...*types.Pair) (storage.Storager, error) {
			assert.Equal(t, storageName, name)
			return store, nil
		})
		sche.EXPECT().Async(gomock.Any()).Do(func(task navvy.Task) {
			switch v := task.(type) {
			case *DeleteDirTask:
				v.validateInput()
			default:
				panic(fmt.Errorf("unexpected task %v", v))
			}
		}).AnyTimes()
		sche.EXPECT().Wait()

		task.run()
		assert.Empty(t, task.GetFault().Error())
	})

	t.Run("delete with segmenter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sche := mock.NewMockScheduler(ctrl)
		srv := mock.NewMockServicer(ctrl)
		store := mock.NewMockStorager(ctrl)
		segmenter := mock.NewMockSegmenter(ctrl)
		storageName := uuid.New().String()

		task := DeleteStorageTask{}
		task.SetFault(fault.New())
		task.SetService(srv)
		task.SetStorageName(storageName)
		task.SetForce(true)
		task.SetPool(navvy.NewPool(10))
		task.SetScheduler(sche)

		srv.EXPECT().Delete(gomock.Any()).Do(func(name string) error {
			assert.Equal(t, storageName, name)
			return nil
		})
		srv.EXPECT().Get(gomock.Any()).DoAndReturn(func(name string, pairs ...*types.Pair) (storage.Storager, error) {
			assert.Equal(t, storageName, name)
			return struct {
				storage.Storager
				storage.Segmenter
			}{
				store,
				segmenter,
			}, nil
		})
		sche.EXPECT().Async(gomock.Any()).Do(func(task navvy.Task) {
			switch v := task.(type) {
			case *DeleteDirTask:
				v.validateInput()
			case *ListSegmentTask:
				v.validateInput()
			default:
				panic(fmt.Errorf("unexpected task %v", v))
			}
		}).AnyTimes()
		sche.EXPECT().Wait()

		task.run()
		assert.Empty(t, task.GetFault().Error())
	})
}
