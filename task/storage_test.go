package task

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/Xuanwo/storage"
	typ "github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yunify/qsctl/v2/utils"

	"github.com/yunify/qsctl/v2/pkg/mock"
	"github.com/yunify/qsctl/v2/pkg/types"
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
			store.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(inputPath string, option *typ.Pair) (_ storage.Storager, err error) {
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
		store.EXPECT().Delete(gomock.Any(), gomock.Any()).DoAndReturn(func(inputPath string, option ...*typ.Pair) error {
			assert.Equal(t, ca.bucketName, inputPath)
			return ca.err
		}).Times(1)

		x := &mockBucketDeleteTask{}
		x.SetBucketName(bucketName)
		x.SetPool(pool)
		x.SetDestinationService(store)

		task := NewDeleteStorageTask(x)
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
		mStorage.EXPECT().Metadata().DoAndReturn(func() (typ.Metadata, error) {
			return ca.meta, ca.metadataErr
		}).MaxTimes(1)

		store.EXPECT().List(gomock.Any()).DoAndReturn(func(option ...*typ.Pair) ([]storage.Storager, error) {
			assert.Equal(t, typ.WithLocation(ca.zone), option[0])
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

func TestRemoveBucketForceTask_new(t *testing.T) {
	cases := []struct {
		name     string
		nextFunc types.TaskFunc
	}{
		{
			name:     "ok",
			nextFunc: NewObjectListAsyncTask,
		},
	}

	for _, tt := range cases {
		m := &mockRemoveBucketForceTask{}

		task := NewRemoveBucketForceTask(m).(*RemoveBucketForceTask)

		assert.Equal(t,
			fmt.Sprintf("%v", tt.nextFunc),
			fmt.Sprintf("%v", task.NextTODO()))
	}
}

func TestNewMakeBucketTask(t *testing.T) {
	cases := []struct {
		input            string
		expectedTodoFunc types.TaskFunc
		expectErr        error
	}{
		{"qs://test-bucket", NewBucketCreateTask, nil},
		{"test-bucket", NewBucketCreateTask, nil},
	}

	for _, v := range cases {
		pt := NewMakeBucketTask(func(task *MakeBucketTask) {
			_, bucketName, _, err := utils.ParseQsPath(v.input)
			if err != nil {
				t.Fatal(err)
			}
			task.SetBucketName(bucketName)
		})

		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))
	}
}
func TestNewRemoveBucketTask(t *testing.T) {
	removeBucketErr := errors.New("remove bucket error")
	cases := []struct {
		input            string
		force            bool
		expectedTodoFunc types.TaskFunc
		expectErr        error
	}{
		{"qs://test-bucket/obj", false, NewDeleteStorageTask, nil},
		{"qs://test-bucket/obj", true, NewRemoveBucketForceTask, nil},
		{"error", false, nil, removeBucketErr},
	}

	for _, v := range cases {
		pt := NewRemoveBucketTask(func(task *RemoveBucketTask) {
			task.SetForce(v.force)
			if v.expectErr != nil {
				task.TriggerFault(v.expectErr)
			}
		})

		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))

		if v.expectErr != nil {
			assert.Equal(t, true, pt.ValidateFault())
			assert.Equal(t, true, errors.Is(pt.GetFault(), v.expectErr))
		} else {
			assert.Equal(t, false, pt.ValidateFault())
		}
	}
}
