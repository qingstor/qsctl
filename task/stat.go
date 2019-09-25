package task

import (
	"github.com/Xuanwo/navvy"
)

// NewStatTask will create a stat task.
func NewStatTask(fn func(*StatTask)) *StatTask {
	t := &StatTask{}

	pool, err := navvy.NewPool(10)
	if err != nil {
		panic(err)
	}
	t.SetPool(pool)

	fn(t)
	t.AddTODOs(NewObjectStatTask)
	return t
}
