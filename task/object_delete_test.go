package task

import (
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/storage"
)

func TestObjectDeleteTask_run(t *testing.T) {
	store := storage.NewMockObjectStorage()
	pool, _ := navvy.NewPool(10)

	cases := []struct {
		objectKey string
		getPanic  bool
	}{
		{storage.MockMBObject, false},
		{"not exist object", true},
	}

	for _, v := range cases {
		x := &mockObjectDeleteTask{}
		x.SetPool(pool)
		x.SetStorage(store)
		x.SetKey(v.objectKey)

		task := NewObjectDeleteTask(x)

		if v.getPanic {
			assert.Panics(t, func() {
				task.Run()
			})
		} else {
			assert.NotPanics(t, func() {
				task.Run()
			})
		}

		pool.Wait()
	}
}
