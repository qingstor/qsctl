package common

import (
	"errors"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/Xuanwo/storage/pkg/iterator"
	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
)

func TestObjectDeleteTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key, errKey := uuid.New().String(), "remove-object-error"
	removeErr := errors.New(errKey)

	store := mock.NewMockStorager(ctrl)

	pool := navvy.NewPool(10)

	cases := []struct {
		name string
		key  string
		err  error
	}{
		{"ok", key, nil},
		{"error", errKey, removeErr},
	}

	for _, ca := range cases {
		store.EXPECT().Delete(gomock.Any()).DoAndReturn(func(key string) error {
			assert.Equal(t, ca.key, key)
			if ca.err != nil {
				return ca.err
			}
			return nil
		}).Times(1)

		x := &mockObjectDeleteTask{}
		x.SetPool(pool)
		x.SetDestinationStorage(store)
		x.SetKey(ca.key)

		task := NewObjectDeleteTask(x)
		task.Run()
		pool.Wait()

		if ca.err != nil {
			assert.Equal(t, x.ValidateFault(), true)
			assert.Equal(t, errors.Is(x.GetFault(), ca.err), true)
			continue
		}

		assert.Equal(t, x.ValidateFault(), false)
	}
}

func TestRemoveDirTask_new(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key, errKey := uuid.New().String(), "remove-dir-error"
	removeDirErr := errors.New(errKey)

	store := mock.NewMockStorager(ctrl)

	pool := navvy.NewPool(10)

	cases := []struct {
		name string
		key  string
		err  error
	}{
		{"ok", key, nil},
		{"error", errKey, removeDirErr},
	}

	for _, tt := range cases {
		store.EXPECT().Delete(gomock.Any()).DoAndReturn(func(key string) error {
			assert.Equal(t, tt.key, key)
			if tt.err != nil {
				return tt.err
			}
			return nil
		}).AnyTimes()

		store.EXPECT().ListDir(gomock.Any(), gomock.Any()).DoAndReturn(func(inputPath string, paris ...*types.Pair) iterator.ObjectIterator {
			assert.Equal(t, inputPath, tt.key)
			count := 3
			return iterator.NewObjectIterator(func(object *[]*types.Object) error {
				*object = append(*object, &types.Object{Name: tt.key})
				count--
				if count > 0 {
					return nil
				}
				return iterator.ErrDone
			})
		})

		x := &mockRemoveDirTask{}
		x.SetPool(pool)
		x.SetDestinationStorage(store)
		x.SetRecursive(true)
		x.SetDeleteKey(tt.key)

		task := NewRemoveDirTask(x)
		if tt.err != nil {
			x.SetFault(tt.err)
			task.(*RemoveDirTask).TriggerFault(tt.err)
		}
		task.Run()
		pool.Wait()

		if tt.err != nil {
			assert.Equal(t, x.ValidateFault(), true)
			assert.Equal(t, errors.Is(x.GetFault(), tt.err), true)
			continue
		}

		assert.Equal(t, x.ValidateFault(), false)
	}
}
