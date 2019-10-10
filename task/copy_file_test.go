package task

import (
	"io"
	"os"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/utils"
)

func TestCopyLargeFileTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := uuid.New().String()
	name, size, _ := utils.GenerateTestFile()
	defer os.Remove(name)

	store := mock.NewMockStorager(ctrl)

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}

	x := &mockCopyLargeFileTask{}
	x.SetPool(pool)
	x.SetPath(name)
	x.SetKey(key)
	x.SetDestinationStorage(store)
	x.SetTotalSize(size)

	store.EXPECT().InitSegment(gomock.Any()).Do(func(inputPath string) {
		assert.Equal(t, key, inputPath)
	})
	store.EXPECT().WriteSegment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(inputPath string, inputOffset, inputSize int64, _ io.ReadCloser) {
		assert.Equal(t, key, inputPath)
		assert.Equal(t, int64(0), inputOffset)
		assert.Equal(t, size, inputSize)
	})
	store.EXPECT().CompleteSegment(gomock.Any()).Do(func(inputPath string) {
		assert.Equal(t, key, inputPath)
	})

	task := NewCopyLargeFileTask(x)
	task.Run()
	pool.Wait()
}

func TestCopyPartialFileTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := uuid.New().String()

	name, size, _ := utils.GenerateTestFile()
	defer os.Remove(name)

	store := mock.NewMockStorager(ctrl)
	store.EXPECT().WriteSegment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(inputPath string, inputOffset, inputSize int64, _ io.ReadCloser) {
		assert.Equal(t, key, inputPath)
		assert.Equal(t, int64(0), inputOffset)
		assert.Equal(t, size, inputSize)
	})

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}

	x := &mockCopyPartialFileTask{}
	x.SetPool(pool)
	x.SetPath(name)
	x.SetKey(key)
	x.SetDestinationStorage(store)
	x.SetPartSize(64 * 1024 * 1024)
	x.SetTotalSize(size)

	currentOffset := int64(0)
	x.SetCurrentOffset(&currentOffset)

	sche := types.NewMockScheduler(nil)
	sche.New(nil)
	x.SetScheduler(sche)

	task := NewCopyPartialFileTask(x)
	task.Run()
	pool.Wait()
}

func TestCopySmallFileTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := uuid.New().String()
	name, size, _ := utils.GenerateTestFile()
	defer os.Remove(name)

	store := mock.NewMockStorager(ctrl)

	store.EXPECT().WriteFile(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(inputPath string, inputSize int64, _ io.ReadCloser) {
		assert.Equal(t, key, inputPath)
		assert.Equal(t, size, inputSize)
	})

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}

	x := &mockCopySmallFileTask{}
	x.SetPool(pool)
	x.SetPath(name)
	x.SetKey(key)
	x.SetDestinationStorage(store)
	x.SetTotalSize(size)

	task := NewCopySmallFileTask(x)
	task.Run()
	pool.Wait()
}
