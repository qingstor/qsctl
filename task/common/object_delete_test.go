package common

import (
	"errors"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
)

const (
	MockMbObj = "mock-mb-object"
	MockGbObj = "mock-gb-object"
)

var objList = map[string]string{
	MockGbObj: MockGbObj,
	MockMbObj: MockMbObj,
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
		{"ok", MockGbObj, nil},
		{"error", key, removeErr},
	}

	for _, ca := range cases {
		objCount := len(objList)
		store.EXPECT().Delete(gomock.Any()).DoAndReturn(func(key string) error {
			if key, ok := objList[key]; ok {
				delete(objList, key)
				return nil
			}
			return ca.err
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
			assert.Equal(t, len(objList), objCount)
			continue
		}

		assert.Equal(t, x.ValidateFault(), false)
		assert.Equal(t, len(objList), objCount-1)
	}
}
