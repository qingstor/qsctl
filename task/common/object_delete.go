package common

import (
	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *RemoveDirTask) new() {
	oc := make(chan *typ.Object)
	t.SetObjectChannel(oc)
	// done to notify get object from channel has done
	done := false
	t.SetDone(&done)
	// set recursive for list async task to list recursively
	t.SetRecursive(true)

	t.SetScheduler(types.NewScheduler(NewObjectDeleteScheduledTask))

	t.AddTODOs(
		NewObjectListAsyncTask,
		NewObjectDeleteIterateTask,
		NewWaitTask,
	)
}

func (t *ObjectDeleteTask) run() {
	log.Debugf("Task <%s> for key <%s> started",
		"ObjectDeleteTask", t.GetDestinationPath())

	if err := t.GetDestinationStorage().Delete(t.GetDestinationPath()); err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for key <%s> finished",
		"ObjectDeleteTask", t.GetDestinationPath())
}

func (t *ObjectDeleteIterateTask) run() {
	log.Debugf("Task <%s> started", "ObjectDeleteIterateTask")
	for {
		if *t.GetDone() {
			break
		}
		t.GetScheduler().New(t.objectDeleteIterateTaskRequirement)
	}
	log.Debugf("Task <%s> finished", "ObjectDeleteIterateTask")
}

func (t *ObjectDeleteScheduledTask) new() {
	obj, ok := <-t.GetObjectChannel()
	if !ok {
		*t.GetDone() = true
		t.AddTODOs(
			NewDoneSchedulerTask,
		)
		return
	}

	t.SetDestinationPath(obj.Name)
	t.AddTODOs(
		NewObjectDeleteTask,
		NewDoneSchedulerTask,
	)
}
