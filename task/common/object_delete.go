package common

import (
	log "github.com/sirupsen/logrus"

	storType "github.com/Xuanwo/storage/types"

	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *ObjectDeleteTask) run() {
	if err := t.GetDestinationStorage().Delete(t.GetKey()); err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for key <%s> finished",
		"ObjectDeleteTask", t.GetKey())
}

func (t *ObjectDeleteWithSchedulerTask) run() {
	log.Debugf("Task <%s> for key <%s> started",
		"ObjectDeleteWithSchedulerTask", t.GetKey())
	defer t.GetScheduler().Done(t.GetID())

	if err := t.GetDestinationStorage().Delete(t.GetKey()); err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for key <%s> finished",
		"ObjectDeleteWithSchedulerTask", t.GetKey())
}

func (t *DirDeleteInitTask) run() {
	log.Debugf("Task <%s> for prefix <%s> started.", "DirDeleteInitTask", t.GetPrefix())
	go func() {
		for obj := range t.GetObjectChannel() {
			t.SetDeleteKey(obj.Name)
			t.GetScheduler().New(t.dirDeleteInitTaskRequirement)
		}
	}()
	log.Debugf("Task <%s> for prefix <%s> finished.", "DirDeleteInitTask", t.GetPrefix())
}

func (t *DirDeleteTask) new() {
	t.SetKey(t.GetDeleteKey())
	t.AddTODOs(
		NewObjectDeleteWithSchedulerTask,
	)
}

func (t *RemoveDirTask) new() {
	oc := make(chan *storType.Object)
	t.SetObjectChannel(oc)

	t.SetScheduler(types.NewScheduler(NewDirDeleteTask))

	t.AddTODOs(
		NewDirDeleteInitTask,
		NewObjectListTask,
		NewWaitTask,
	)
}
