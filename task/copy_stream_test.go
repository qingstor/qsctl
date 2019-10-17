package task

import (
	"bytes"
	"io"
	"sync"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/mock"
	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/utils"
)

func TestCopyStreamTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStorager(ctrl)
	key := uuid.New().String()
	segmentID := uuid.New().String()

	buf, size, _ := utils.GenerateTestStream()

	pool := navvy.NewPool(10)

	x := NewCopyTask(func(task *CopyTask) {
		task.SetDestinationStorage(store)
		task.SetKey(key)
		task.SetPath("-")
		task.SetPool(pool)
		task.SetKeyType(constants.KeyTypeObject)
		task.SetPathType(constants.PathTypeStream)
		task.SetFlowType(constants.FlowToRemote)
		task.SetStream(buf)
	})

	store.EXPECT().InitSegment(gomock.Any()).DoAndReturn(func(inputPath string) (id string, err error) {
		assert.Equal(t, key, inputPath)
		return segmentID, nil
	})
	store.EXPECT().WriteSegment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(inputPath string, inputOffset, inputSize int64, _ io.ReadCloser) {
		assert.Equal(t, segmentID, inputPath)
		assert.Equal(t, int64(0), inputOffset)
		assert.Equal(t, size, inputSize)
	})
	store.EXPECT().CompleteSegment(gomock.Any()).Do(func(inputPath string) {
		assert.Equal(t, segmentID, inputPath)
	})

	task := NewCopyStreamTask(x)
	task.Run()
	pool.Wait()
}

func TestCopyPartialStreamTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := uuid.New().String()
	segmentID := uuid.New().String()
	buf, size, _ := utils.GenerateTestStream()

	store := mock.NewMockStorager(ctrl)
	store.EXPECT().WriteSegment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(inputPath string, inputOffset, inputSize int64, _ io.ReadCloser) {
		assert.Equal(t, segmentID, inputPath)
		assert.Equal(t, int64(0), inputOffset)
		assert.Equal(t, size, inputSize)
	})

	pool := navvy.NewPool(10)

	x := &mockCopyPartialStreamTask{}
	x.SetPool(pool)
	x.SetStream(buf)
	x.SetKey(key)
	x.SetDestinationStorage(store)
	x.SetPartSize(64 * 1024 * 1024)
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
