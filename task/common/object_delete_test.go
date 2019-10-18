package common

import (
	"errors"
	"testing"
	"time"

	"github.com/Xuanwo/navvy"
	storTypes "github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
	"github.com/yunify/qsctl/v2/pkg/types"
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

func TestObjectDeleteWithSchedulerTask_run(t *testing.T) {
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

		x := &mockObjectDeleteWithSchedulerTask{}
		x.SetPool(pool)
		x.SetDestinationStorage(store)
		x.SetKey(ca.key)
		x.SetID(uuid.New().String())

		sche := types.NewMockScheduler(nil)
		sche.New(nil)
		x.SetScheduler(sche)
		task := NewObjectDeleteWithSchedulerTask(x)
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

func TestDirDeleteInitTask_run(t *testing.T) {
	pool := navvy.NewPool(10)
	objKey := uuid.New().String()
	objList := []storTypes.Object{
		{Name: objKey},
		{Name: objKey},
		{Name: objKey},
	}

	oc := make(chan *storTypes.Object, len(objList))
	for _, obj := range objList {
		oc <- &obj
	}
	close(oc)

	x := &mockDirDeleteInitTask{}
	x.SetPool(pool)
	x.SetObjectChannel(oc)
	x.SetPrefix(uuid.New().String())

	x.SetID(uuid.New().String())
	sche := types.NewMockScheduler(nil)
	sche.New(nil)
	x.SetScheduler(sche)

	task := NewDirDeleteInitTask(x)
	task.Run()
	pool.Wait()

	time.Sleep(time.Second) // wait goroutine to set delete key
	assert.Equal(t, x.GetDeleteKey(), objKey)
	assert.Equal(t, x.ValidateFault(), false)
}
