package task

import (
	"bytes"
	"io"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/Xuanwo/navvy"
	typ "github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/mock"
	"github.com/yunify/qsctl/v2/pkg/types"
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
		fn   types.TaskFunc
	}{
		{
			"small file",
			constants.MaximumAutoMultipartSize - 1,
			NewCopySmallFileTask,
		},
		{
			"large file",
			constants.MaximumAutoMultipartSize + 1,
			NewCopyLargeFileTask,
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

			m := &mockCopyFileTask{}
			m.SetSourceStorage(srcStore)
			m.SetSourcePath(paths[k])
			task := &CopyFileTask{copyFileTaskRequirement: m}
			task.new()

			assert.Equal(t, v.size, task.GetTotalSize())
		})
	}
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

	x := &mockCopyLargeFileTask{}
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
	assert.Equal(t, int64(0), *tt.GetCurrentOffset())
}

func TestCopyPartialFileTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	name := uuid.New().String()
	key := uuid.New().String()
	segmentID := uuid.New().String()
	size := int64(1234)
	buf := bytes.NewReader([]byte("Hello, World"))

	srcStore := mock.NewMockStorager(ctrl)
	srcStore.EXPECT().Read(gomock.Any(), gomock.Any()).DoAndReturn(func(inputPath string, pairs ...*typ.Pair) (r io.ReadCloser, err error) {
		assert.Equal(t, name, inputPath)
		return ioutil.NopCloser(buf), nil
	}).AnyTimes()
	dstStore := mock.NewMockStorager(ctrl)
	dstStore.EXPECT().WriteSegment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(inputPath string, inputOffset, inputSize int64, _ io.ReadCloser) {
		assert.Equal(t, segmentID, inputPath)
		assert.Equal(t, int64(0), inputOffset)
		assert.Equal(t, size, inputSize)
	})

	pool := navvy.NewPool(10)

	x := &mockCopyPartialFileTask{}
	x.SetPool(pool)
	x.SetSourcePath(name)
	x.SetSourceStorage(srcStore)
	x.SetDestinationPath(key)
	x.SetDestinationStorage(dstStore)
	x.SetPartSize(64 * 1024 * 1024)
	x.SetTotalSize(size)
	x.SetSegmentID(segmentID)

	currentOffset := int64(0)
	x.SetCurrentOffset(&currentOffset)

	task := NewCopyPartialFileTask(x)
	task.Run()
	pool.Wait()
}

func TestCopySmallFileTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	name := uuid.New().String()
	key := uuid.New().String()
	size := int64(1234)
	buf := bytes.NewReader([]byte("Hello, World"))

	srcStore := mock.NewMockStorager(ctrl)
	srcStore.EXPECT().Read(gomock.Any()).DoAndReturn(func(inputPath string) (r io.ReadCloser, err error) {
		assert.Equal(t, name, inputPath)
		return ioutil.NopCloser(buf), nil
	})
	srcStore.EXPECT().Read(gomock.Any(), gomock.Any()).DoAndReturn(func(inputPath string, pairs ...*typ.Pair) (r io.ReadCloser, err error) {
		assert.Equal(t, name, inputPath)
		return ioutil.NopCloser(buf), nil
	})
	dstStore := mock.NewMockStorager(ctrl)
	dstStore.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(inputPath string, _ io.ReadCloser, option ...*typ.Pair) {
		assert.Equal(t, key, inputPath)
		assert.Equal(t, size, option[0].Value.(int64))
	})

	pool := navvy.NewPool(10)

	x := &mockCopySmallFileTask{}
	x.SetPool(pool)
	x.SetSourcePath(name)
	x.SetSourceStorage(srcStore)
	x.SetDestinationPath(key)
	x.SetDestinationStorage(dstStore)
	x.SetTotalSize(size)

	task := NewCopySmallFileTask(x)
	task.Run()
	pool.Wait()
}

func TestCopyStreamTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStorager(ctrl)
	key := uuid.New().String()

	pool := navvy.NewPool(10)

	x := &mockCopyStreamTask{}
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

func TestCopyPartialStreamTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := uuid.New().String()
	localPath := "-"
	segmentID := uuid.New().String()
	buf := bytes.NewReader([]byte("Hello, world!"))

	srcStore := mock.NewMockStorager(ctrl)
	srcStore.EXPECT().Read(gomock.Any(), gomock.Any()).DoAndReturn(func(inputPath string, pairs ...*typ.Pair) (r io.ReadCloser, err error) {
		assert.Equal(t, localPath, inputPath)
		assert.Equal(t, int64(constants.DefaultPartSize), pairs[0].Value.(int64))
		return ioutil.NopCloser(buf), nil
	})

	dstStore := mock.NewMockStorager(ctrl)
	dstStore.EXPECT().WriteSegment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(inputPath string, inputOffset, inputSize int64, _ io.ReadCloser) {
		assert.Equal(t, segmentID, inputPath)
		assert.Equal(t, int64(0), inputOffset)
		assert.Equal(t, int64(13), inputSize)
	})

	pool := navvy.NewPool(10)

	x := &mockCopyPartialStreamTask{}
	x.SetPool(pool)
	x.SetSourcePath(localPath)
	x.SetSourceStorage(srcStore)
	x.SetDestinationPath(key)
	x.SetDestinationStorage(dstStore)
	x.SetPartSize(constants.DefaultPartSize)
	x.SetSegmentID(segmentID)
	x.SetBytesPool(&sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, x.GetPartSize()))
		},
	})

	currentOffset := int64(0)
	x.SetCurrentOffset(&currentOffset)

	task := NewCopyPartialStreamTask(x)
	task.Run()
	pool.Wait()
}
