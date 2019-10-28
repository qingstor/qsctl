package common

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Xuanwo/navvy"
	typ "github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/utils"
)

func TestRemoveDirTask_New(t *testing.T) {
	cases := []struct {
		name     string
		nextFunc types.TaskFunc
		err      error
	}{
		{
			name:     "ok",
			nextFunc: NewObjectListAsyncTask,
			err:      nil,
		},
	}

	for _, tt := range cases {
		m := &mockRemoveDirTask{}
		task := NewRemoveDirTask(m).(*RemoveDirTask)

		assert.Equal(t,
			fmt.Sprintf("%v", tt.nextFunc),
			fmt.Sprintf("%v", task.NextTODO()))

	}
}

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
		x.SetDestinationPath(ca.key)

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

func TestObjectDeleteIterateTask_run(t *testing.T) {
	pool := navvy.NewPool(10)
	err := errors.New("test error")

	x := &mockObjectDeleteIterateTask{}

	id := uuid.New().String()

	fn := func(task types.Todoist) navvy.Task {
		assert.Equal(t, false, *x.GetDone())
		*x.GetDone() = true

		t := &utils.EmptyTask{}
		t.SetID(id)
		t.SetPool(pool)
		return t
	}
	x.SetScheduler(types.NewScheduler(fn))

	{
		done := false
		x.SetDone(&done)
		task := NewObjectDeleteIterateTask(x)
		task.Run()
		assert.Equal(t, true, *x.GetDone())
	}
	{
		done := false
		x.SetDone(&done)
		x.SetFault(err)
		task := NewObjectDeleteIterateTask(x)
		task.Run()
		assert.Equal(t, true, *x.GetDone())
		assert.Equal(t, true, x.ValidateFault())
		assert.Equal(t, true, errors.Is(x.GetFault(), err))
	}
}

func TestObjectDeleteScheduledTask_New(t *testing.T) {
	path := uuid.New().String()

	cases := []struct {
		name            string
		nextFunc        types.TaskFunc
		done            bool
		destinationPath string
	}{
		{
			name:            "ok",
			nextFunc:        NewObjectDeleteTask,
			done:            false,
			destinationPath: path,
		},
		{
			name:            "done",
			nextFunc:        NewDoneSchedulerTask,
			done:            true,
			destinationPath: "",
		},
	}

	for _, tt := range cases {
		oc := make(chan *typ.Object)
		m := &mockObjectDeleteScheduledTask{}
		m.SetObjectChannel(oc)
		done := false
		m.SetDone(&done)

		if tt.done {
			close(oc)
		} else {
			go func() {
				oc <- &typ.Object{Name: tt.destinationPath}
			}()
		}

		task := NewObjectDeleteScheduledTask(m).(*ObjectDeleteScheduledTask)

		assert.Equal(t,
			fmt.Sprintf("%v", tt.nextFunc),
			fmt.Sprintf("%v", task.NextTODO()))

		if tt.done {
			assert.Equal(t, *m.GetDone(), tt.done)
		} else {
			assert.Equal(t, task.GetDestinationPath(), tt.destinationPath)
		}
	}
}
