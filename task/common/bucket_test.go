package common

import (
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
)

func TestBucketCreateTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketName, zone := uuid.New().String(), uuid.New().String()
	store := mock.NewMockStorager(ctrl)

	store.EXPECT().CreateDir(gomock.Any(), gomock.Any()).Do(func(inputPath string, option ...*types.Pair) {
		assert.Equal(t, bucketName, inputPath)
		assert.Equal(t, zone, option[0].Value.(string))
	})

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}

	x := &mockBucketCreateTask{}
	x.SetBucketName(bucketName)
	x.SetZone(zone)
	x.SetPool(pool)
	x.SetDestinationStorage(store)

	task := NewBucketCreateTask(x)
	task.Run()
	pool.Wait()
}

func TestBucketDeleteTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketName := uuid.New().String()
	store := mock.NewMockStorager(ctrl)

	store.EXPECT().Delete(gomock.Any(), gomock.Any()).Do(func(inputPath string, option ...*types.Pair) {
		assert.Equal(t, bucketName, inputPath)
	})

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}

	x := &mockBucketDeleteTask{}
	x.SetBucketName(bucketName)
	x.SetPool(pool)
	x.SetDestinationStorage(store)

	task := NewBucketDeleteTask(x)
	task.Run()
	pool.Wait()
}
