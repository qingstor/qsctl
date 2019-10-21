package common

import (
	log "github.com/sirupsen/logrus"

	"github.com/Xuanwo/storage/types"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *ObjectDeleteTask) run() {
	if err := t.GetDestinationStorage().Delete(t.GetKey()); err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for key <%s> finished",
		"ObjectDeleteTask", t.GetKey())
}

func (t *RemoveDirTask) new() {
	log.Debugf("new remove dir task")
	t.SetPrefix(t.GetDeleteKey())
	oc := make(chan *types.Object)
	t.SetObjectChannel(oc)

	t.GetPool().Submit(NewObjectListTask(t))

	for obj := range t.GetObjectChannel() {
		task := *t
		task.SetKey(obj.Name)
		t.GetPool().Submit(NewObjectDeleteTask(&task))
	}
}
