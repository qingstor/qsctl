package task

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Xuanwo/storage/pkg/iterator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func TestNewListTask(t *testing.T) {
	cases := []struct {
		listType         constants.ListType
		expectedTodoFunc types.TaskFunc
		wantPanic        bool
	}{
		{constants.ListTypeBucket, NewBucketListTask, false},
		{constants.ListTypeInvalid, nil, true},
	}

	for _, v := range cases {
		pt := new(ListTask)
		panicFunc := func() {
			pt = NewListTask(func(task *ListTask) {
				task.SetListType(v.listType)
			})
		}
		if v.wantPanic {
			assert.Panics(t, panicFunc)
		} else {
			panicFunc()
		}
		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))
	}
}

func TestObjectListTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStorager(ctrl)
	key, listErr := uuid.New().String(), errors.New("list-object-err")

	cases := []struct {
		name      string
		key       string
		recursive bool
		fault     bool
		err       error
	}{
		{"non-recursive ok", key, false, false, iterator.ErrDone},
		{"recursive ok", key, true, false, iterator.ErrDone},
		{"non-recursive error not done", key, false, true, listErr},
		{"recursive error not done", key, true, true, listErr},
	}

	for _, ca := range cases {
		x := &mockObjectListTask{}
		x.SetDestinationStorage(store)
		x.SetDestinationPath(ca.key)
		x.SetRecursive(ca.recursive)
		x.SetObjectChannel(make(chan *types.Object))

		go func() {
			// make channel not blocked when set
			for {
				<-x.GetObjectChannel()
			}
		}()

		store.EXPECT().ListDir(gomock.Any(), gomock.Any()).DoAndReturn(func(inputPath string, paris ...*types.Pair) iterator.ObjectIterator {
			assert.Equal(t, inputPath, key)
			if ca.recursive {
				assert.Equal(t, 0, len(paris), ca.name)
			} else {
				assert.Equal(t, "/", paris[0].Value.(string), ca.name)
			}
			count := 3
			return iterator.NewObjectIterator(func(object *[]*types.Object) error {
				*object = make([]*types.Object, 1)
				count--
				if count > 0 {
					return nil
				}
				return ca.err
			})
		})

		task := NewObjectListTask(x)
		task.Run()

		assert.Equal(t, ca.fault, x.ValidateFault(), ca.name)
		if ca.fault {
			assert.Error(t, x.GetFault(), ca.name)
			assert.Equal(t, true, errors.Is(x.GetFault(), ca.err), ca.name)
		}
	}
}

func TestObjectListAsyncTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStorager(ctrl)
	key, listErr := uuid.New().String(), errors.New("list-object-async-err")

	cases := []struct {
		name      string
		key       string
		recursive bool
		fault     bool
		err       error
	}{
		{"non-recursive ok", key, false, false, iterator.ErrDone},
		{"recursive ok", key, true, false, iterator.ErrDone},
		{"non-recursive error not done", key, false, true, listErr},
		{"recursive error not done", key, true, true, listErr},
	}

	for _, ca := range cases {
		x := &mockObjectListAsyncTask{}
		x.SetDestinationStorage(store)
		x.SetDestinationPath(ca.key)
		x.SetRecursive(ca.recursive)
		x.SetObjectChannel(make(chan *types.Object))

		// set fault manually to trigger fault
		// because async list would not trigger fault instantly
		if ca.fault {
			x.SetFault(ca.err)
		}

		store.EXPECT().ListDir(gomock.Any(), gomock.Any()).DoAndReturn(func(inputPath string, paris ...*types.Pair) iterator.ObjectIterator {
			assert.Equal(t, inputPath, key)
			if ca.recursive {
				assert.Equal(t, 0, len(paris), ca.name)
			} else {
				assert.Equal(t, "/", paris[0].Value.(string), ca.name)
			}
			count := 3
			return iterator.NewObjectIterator(func(object *[]*types.Object) error {
				*object = make([]*types.Object, 1)
				count--
				if count > 0 {
					return nil
				}
				return ca.err
			})
		})

		task := NewObjectListAsyncTask(x)
		task.Run()

		// make channel blocked before all get
		for range x.GetObjectChannel() {
		}

		assert.Equal(t, ca.fault, x.ValidateFault(), ca.name)
		if ca.fault {
			assert.Error(t, x.GetFault(), ca.name)
			assert.Equal(t, true, errors.Is(x.GetFault(), ca.err), ca.name)
		}
	}
}
