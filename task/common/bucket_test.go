package common

import (
	"errors"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
)

func TestBucketCreateTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketName, zone, errBucket := uuid.New().String(), uuid.New().String(), "error-bucket"
	store := mock.NewMockServicer(ctrl)
	bucketErr := errors.New(errBucket)
	pool := navvy.NewPool(10)

	cases := []struct {
		name       string
		bucketName string
		zone       string
		err        error
	}{
		{"ok", bucketName, zone, nil},
		{"error", errBucket, zone, bucketErr},
	}

	for _, ca := range cases {
		t.Run(ca.name, func(t *testing.T) {
			store.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(inputPath string, option *types.Pair) (_ storage.Storager, err error) {
				assert.Equal(t, ca.bucketName, inputPath)
				assert.Equal(t, ca.zone, option.Value.(string))
				return nil, ca.err
			}).Times(1)

			x := &mockBucketCreateTask{}
			x.SetBucketName(ca.bucketName)
			x.SetZone(zone)
			x.SetPool(pool)
			x.SetDestinationService(store)

			task := NewBucketCreateTask(x)
			task.Run()
			pool.Wait()

			if ca.err == nil {
				assert.Equal(t, false, x.ValidateFault())
				return
			}
			assert.Equal(t, x.ValidateFault(), true)
			assert.Error(t, x.GetFault())
			assert.Equal(t, true, errors.Is(x.GetFault(), ca.err))
		})
	}
}

func TestBucketDeleteTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketName, errBucket := uuid.New().String(), "err-bucket"
	bucketErr := errors.New(errBucket)
	pool := navvy.NewPool(10)
	store := mock.NewMockServicer(ctrl)

	cases := []struct {
		name       string
		bucketName string
		err        error
	}{
		{"ok", bucketName, nil},
		{"error", bucketName, bucketErr},
	}

	for _, ca := range cases {
		// different case different behaviour
		store.EXPECT().Delete(gomock.Any(), gomock.Any()).DoAndReturn(func(inputPath string, option ...*types.Pair) error {
			assert.Equal(t, ca.bucketName, inputPath)
			return ca.err
		}).Times(1)

		x := &mockBucketDeleteTask{}
		x.SetBucketName(bucketName)
		x.SetPool(pool)
		x.SetDestinationService(store)

		task := NewBucketDeleteTask(x)
		task.Run()
		pool.Wait()

		if ca.err == nil {
			assert.Equal(t, false, x.ValidateFault())
			continue
		}
		assert.Equal(t, x.ValidateFault(), true)
		assert.Error(t, x.GetFault())
		assert.Equal(t, true, errors.Is(x.GetFault(), ca.err))
	}
}
