package common

import (
	"testing"

	"github.com/Xuanwo/storage/pkg/iterator"
	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
)

func TestObjectListTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStorager(ctrl)

	{
		// Test list without recursive
		key := uuid.New().String()

		x := &mockObjectListTask{}
		x.SetDestinationStorage(store)
		x.SetKey(key)
		x.SetRecursive(false)
		x.SetObjectChannel(make(chan *types.Object))

		store.EXPECT().ListDir(gomock.Any(), gomock.Any()).DoAndReturn(func(inputPath string, paris ...*types.Pair) iterator.Iterator {
			assert.Equal(t, inputPath, key)
			assert.Equal(t, "/", paris[0].Value.(string))
			return iterator.NewPrefixBasedIterator(func(object *[]*types.Object) error {
				return iterator.ErrDone
			})
		})

		task := NewObjectListTask(x)
		task.Run()
	}

	{
		// Test list with recursive
		key := uuid.New().String()

		x := &mockObjectListTask{}
		x.SetDestinationStorage(store)
		x.SetKey(key)
		x.SetRecursive(true)
		x.SetObjectChannel(make(chan *types.Object))

		store.EXPECT().ListDir(gomock.Any(), gomock.Any()).DoAndReturn(func(inputPath string, paris ...*types.Pair) iterator.Iterator {
			assert.Equal(t, inputPath, key)
			assert.Equal(t, 0, len(paris))
			return iterator.NewPrefixBasedIterator(func(object *[]*types.Object) error {
				return iterator.ErrDone
			})
		})

		task := NewObjectListTask(x)
		task.Run()
	}
}
