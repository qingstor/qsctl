package task

import (
	"github.com/Xuanwo/navvy"
	log "github.com/sirupsen/logrus"

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
	t.AddTODOs(NewStatObjectTask)
	return t
}

func (t *StatObjectTask) run() {
	om, err := t.GetStorage().HeadObject(t.GetKey())
	if err != nil {
		panic(err)
	}
	oriOm := t.GetObjectMeta()
	// replace the original om
	*oriOm = *om
	log.Debugf("Task <%s> for Key <%s> finished.", "StatObjectTask", t.GetKey())
}
