package task

import "github.com/Xuanwo/navvy"

// NewPresignTask will create a presign task
func NewPresignTask(fn func(*PresignTask)) *PresignTask {
	t := &PresignTask{}

	pool, err := navvy.NewPool(10)
	if err != nil {
		panic(err)
	}
	t.SetPool(pool)

	fn(t)

	t.AddTODOs(NewObjectPresignTask)
	return t
}
