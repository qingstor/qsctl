package common

import (
	log "github.com/sirupsen/logrus"

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

func (t *ObjectInitDirDeleteTask) run() {
	log.Debugf("Task <%s> for prefix <%s> started.", "ObjectInitDirDeleteTask", t.GetPrefix())
	go func() {
		for obj := range t.GetObjectChannel() {
			t.SetDeleteKey(obj.Name)
			t.GetScheduler().New(t.objectInitDirDeleteTaskRequirement)
		}
	}()
	log.Debugf("Task <%s> for prefix <%s> finished.", "ObjectInitDirDeleteTask", t.GetPrefix())
}

func (t *ObjectDeleteRecursivelyTask) new() {
	t.SetKey(t.GetDeleteKey())
	t.AddTODOs(
		NewObjectDeleteWithSchedulerTask,
	)
}
