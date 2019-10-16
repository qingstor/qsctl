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

	bucketName, errBucket := uuid.New().String(), "listErr-bucket"
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

func TestBucketListTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listErr, zone, name := errors.New("listErr-list-bucket"), uuid.New().String(), uuid.New().String()
	pool := navvy.NewPool(10)
	store := mock.NewMockServicer(ctrl)

	cases := []struct {
		name        string
		zone        string
		meta        map[string]interface{}
		metadataErr error
		listErr     error
	}{
		{"ok", zone, map[string]interface{}{"name": name}, nil, nil},
		{"blank_zone", "", map[string]interface{}{"name": name}, nil, nil},
		{"meta-data-error", zone, nil, listErr, nil},
		{"list-error", zone, nil, nil, listErr},
	}

	for _, ca := range cases {
		mStorage := mock.NewMockStorager(ctrl)
		mStorage.EXPECT().Metadata().DoAndReturn(func() (types.Metadata, error) {
			return ca.meta, ca.metadataErr
		}).MaxTimes(1)

		store.EXPECT().List(gomock.Any()).DoAndReturn(func(option ...*types.Pair) ([]storage.Storager, error) {
			assert.Equal(t, types.WithLocation(ca.zone), option[0])
			return []storage.Storager{mStorage}, ca.listErr
		}).Times(1)

		x := &mockBucketListTask{}
		x.SetZone(ca.zone)
		x.SetPool(pool)
		x.SetDestinationService(store)

		task := NewBucketListTask(x)
		task.Run()
		pool.Wait()

		if ca.listErr == nil && ca.metadataErr == nil {
			assert.Equal(t, len(ca.meta), len(x.GetBucketList()), ca.name)
			assert.Equal(t, false, x.ValidateFault(), ca.name)
			continue
		}
		assert.Equal(t, x.ValidateFault(), true, ca.name)
		assert.Error(t, x.GetFault(), ca.name)
		assert.Equal(t, true, errors.Is(x.GetFault(), listErr), ca.name)
	}
}
