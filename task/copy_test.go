package task

import (
	"fmt"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/Xuanwo/storage"
	typ "github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/fault"
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

func TestCopyFileTask_run(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cases := []struct {
			name string
			size int64
		}{
			{
				"large file",
				constants.MaximumAutoMultipartSize + 1,
			},
			{
				"small file",
				constants.MaximumAutoMultipartSize - 1,
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				sche := mock.NewMockScheduler(ctrl)
				srcStore := mock.NewMockStorager(ctrl)
				srcPath := uuid.New().String()
				dstStore := mock.NewMockStorager(ctrl)
				dstPath := uuid.New().String()

				task := &CopyFileTask{}
				task.SetPool(navvy.NewPool(10))
				task.SetScheduler(sche)
				task.SetCheckTasks(nil)
				task.SetSourcePath(srcPath)
				task.SetSourceStorage(srcStore)
				task.SetDestinationPath(dstPath)
				task.SetDestinationStorage(dstStore)

				sche.EXPECT().Sync(gomock.Any()).Do(func(task navvy.Task) {
					switch v := task.(type) {
					case *BetweenStorageCheckTask:
						v.SetSourceObject(&typ.Object{Name: srcPath, Size: tt.size})
						v.SetDestinationObject(&typ.Object{Name: dstPath})
					case *CopyLargeFileTask:
						assert.True(t, tt.size >= constants.MaximumAutoMultipartSize)
					case *CopySmallFileTask:
						assert.True(t, tt.size < constants.MaximumAutoMultipartSize)
					}
				}).AnyTimes()

				task.run()
			})
		}
	})
}

func TestCopySmallFileTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sche := mock.NewMockScheduler(ctrl)
	srcStore := mock.NewMockStorager(ctrl)
	srcPath := uuid.New().String()
	dstStore := mock.NewMockStorager(ctrl)
	dstPath := uuid.New().String()

	task := &CopySmallFileTask{}
	task.SetPool(navvy.NewPool(10))
	task.SetSourcePath(srcPath)
	task.SetSourceStorage(srcStore)
	task.SetDestinationPath(dstPath)
	task.SetDestinationStorage(dstStore)
	task.SetScheduler(sche)
	task.SetSize(1024)

	sche.EXPECT().Sync(gomock.Any()).Do(func(task navvy.Task) {
		switch v := task.(type) {
		case *MD5SumFileTask:
			assert.Equal(t, srcPath, v.GetPath())
			assert.Equal(t, int64(0), v.GetOffset())
			v.SetMD5Sum([]byte("string"))
		case *CopySingleFileTask:
			assert.Equal(t, []byte("string"), v.GetMD5Sum())
		default:
			panic(fmt.Errorf("unexpected task %v", v))
		}
	}).AnyTimes()

	task.run()
}

func TestCopyLargeFileTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sche := mock.NewMockScheduler(ctrl)
	srcStore := mock.NewMockStorager(ctrl)
	srcPath := uuid.New().String()
	dstStore := mock.NewMockStorager(ctrl)
	dstSegmenter := mock.NewMockSegmenter(ctrl)
	dstPath := uuid.New().String()
	segmentID := uuid.New().String()

	task := &CopyLargeFileTask{}
	task.SetPool(navvy.NewPool(10))
	task.SetSourcePath(srcPath)
	task.SetSourceStorage(srcStore)
	task.SetDestinationPath(dstPath)
	task.SetDestinationStorage(struct {
		storage.Storager
		storage.Segmenter
	}{
		dstStore,
		dstSegmenter,
	})
	task.SetScheduler(sche)
	task.SetFault(fault.New())
	// 50G
	task.SetTotalSize(10 * constants.MaximumPartSize)

	sche.EXPECT().Sync(gomock.Any()).Do(func(task navvy.Task) {
		switch v := task.(type) {
		case *SegmentInitTask:
			assert.Equal(t, dstPath, v.GetPath())
			v.SetSegmentID(segmentID)
		case *SegmentCompleteTask:
			assert.Equal(t, dstPath, v.GetPath())
			assert.Equal(t, segmentID, v.GetSegmentID())
		default:
			panic(fmt.Errorf("invalid task %v", v))
		}
	}).AnyTimes()
	sche.EXPECT().Async(gomock.Any()).Do(func(task navvy.Task) {
		switch v := task.(type) {
		case *CopyPartialFileTask:
			assert.Equal(t, srcPath, v.GetSourcePath())
			assert.Equal(t, dstPath, v.GetDestinationPath())
			assert.Equal(t, segmentID, v.GetSegmentID())
			v.SetDone(true)
		default:
			panic(fmt.Errorf("unexpected task %v", v))
		}
	}).AnyTimes()
	sche.EXPECT().Wait().Do(func() {})

	task.run()
	assert.Empty(t, task.GetFault().Error())
}

func TestCopyPartialFileTask_new(t *testing.T) {
	cases := []struct {
		name       string
		totalsize  int64
		offset     int64
		partsize   int64
		expectSize int64
		expectDone bool
	}{
		{
			"middle part",
			1024,
			128,
			128,
			128,
			false,
		},
		{
			"last fulfilled part",
			1024,
			512,
			512,
			512,
			true,
		},
		{
			"last not fulfilled part",
			1024,
			768,
			512,
			256,
			true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			task := &CopyPartialFileTask{}
			task.SetTotalSize(tt.totalsize)
			task.SetOffset(tt.offset)
			task.SetPartSize(tt.partsize)

			task.new()

			assert.Equal(t, tt.expectSize, task.GetSize())
			assert.Equal(t, tt.expectDone, task.GetDone())
		})
	}
}

func TestCopyPartialFileTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sche := mock.NewMockScheduler(ctrl)
	srcStore := mock.NewMockStorager(ctrl)
	srcPath := uuid.New().String()
	dstStore := mock.NewMockStorager(ctrl)
	dstSegmenter := mock.NewMockSegmenter(ctrl)
	dstPath := uuid.New().String()
	segmentID := uuid.New().String()

	task := &CopyPartialFileTask{}
	task.SetFault(fault.New())
	task.SetPool(navvy.NewPool(10))
	task.SetSourcePath(srcPath)
	task.SetSourceStorage(srcStore)
	task.SetDestinationPath(dstPath)
	task.SetDestinationStorage(struct {
		storage.Storager
		storage.Segmenter
	}{
		dstStore,
		dstSegmenter,
	})
	task.SetScheduler(sche)
	task.SetSize(1024)
	task.SetOffset(512)
	task.SetSegmentID(segmentID)

	srcStore.EXPECT().String().Do(func() {}).AnyTimes()
	dstStore.EXPECT().String().Do(func() {}).AnyTimes()

	sche.EXPECT().Sync(gomock.Any()).Do(func(task navvy.Task) {
		t.Logf("Got task %v", task)

		switch v := task.(type) {
		case *MD5SumFileTask:
			assert.Equal(t, srcPath, v.GetPath())
			assert.Equal(t, int64(512), v.GetOffset())
			v.SetMD5Sum([]byte("string"))
		case *SegmentFileCopyTask:
			assert.Equal(t, []byte("string"), v.GetMD5Sum())
		default:
			panic(fmt.Errorf("unexpected task %v", v))
		}
	}).AnyTimes()

	task.run()
	assert.Empty(t, task.GetFault().Error())
}
