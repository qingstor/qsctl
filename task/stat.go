package task

import (
	"github.com/Xuanwo/navvy"

	"github.com/yunify/qsctl/v2/storage"
)

// NewStatTask will create a stat task.
func NewStatTask(fn func(t *StatTask)) *StatTask {
	t := &StatTask{}

	pool, err := navvy.NewPool(10)
	if err != nil {
		panic(err)
	}
	t.SetPool(pool)
	t.SetObjectMeta(&storage.ObjectMeta{})

	fn(t)
	t.AddTODOs(NewObjectStatTask)
	return t
}
