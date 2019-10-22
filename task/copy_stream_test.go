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

func TestCopyStreamTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStorager(ctrl)
	key := uuid.New().String()

	pool := navvy.NewPool(10)

	x := NewCopyTask(func(task *CopyTask) {
		task.SetSourceType(typ.ObjectTypeStream)
		task.SetDestinationStorage(store)
		task.SetDestinationPath(key)
		task.SetSourcePath("-")
		task.SetPool(pool)
	})

	task := NewCopyStreamTask(x)

	tt := task.(*CopyStreamTask)
	assert.NotNil(t, tt.GetBytesPool())
	assert.Equal(t, int64(constants.DefaultPartSize), tt.GetPartSize())
	assert.NotNil(t, tt.GetScheduler())
	assert.Equal(t, int64(0), *tt.GetCurrentOffset())
	assert.Equal(t, int64(-1), tt.GetTotalSize())
	assert.NotNil(t, tt.NextTODO())
}

func TestCopyPartialStreamTask_Run(t *testing.T) {
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

	sche := types.NewMockScheduler(nil)
	sche.New(nil)
	x.SetScheduler(sche)

	currentOffset := int64(0)
	x.SetCurrentOffset(&currentOffset)

	task := NewCopyPartialStreamTask(x)
	task.Run()
	pool.Wait()
}
