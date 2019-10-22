package common

import (
	"io/ioutil"

	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *MultipartInitTask) run() {
	log.Debugf("Task <%s> for Object <%s> started.", "MultipartInitTask", t.GetDestinationPath())

	id, err := t.GetDestinationStorage().InitSegment(t.GetDestinationPath())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	t.SetSegmentID(id)

	for {
		if *t.GetCurrentOffset() == t.GetTotalSize() {
			break
		}

		t.GetScheduler().New(t.multipartInitTaskRequirement)
	}

	log.Debugf("Task <%s> for Object <%s> finished.", "MultipartInitTask", t.GetDestinationPath())
}

func (t *MultipartFileUploadTask) run() {
	log.Debugf("Task <%s> for File <%s> at Offset <%d> started.", "MultipartFileUploadTask", t.GetSourcePath(), t.GetOffset())

	defer t.GetScheduler().Done(t.GetID())

	r, err := t.GetSourceStorage().Read(t.GetSourcePath(), typ.WithSize(t.GetSize()), typ.WithOffset(t.GetOffset()))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	defer r.Close()

	// TODO: Add checksum support.
	err = t.GetDestinationStorage().WriteSegment(t.GetSegmentID(), t.GetOffset(), t.GetSize(), r)
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for File <%s> at Offset <%d> finished.", "MultipartFileUploadTask", t.GetSourcePath(), t.GetOffset())
}

func (t *MultipartStreamUploadTask) run() {
	log.Debugf("Task <%s> for Stream at Offset <%d> started.", "MultipartStreamUploadTask", t.GetOffset())

	defer t.GetScheduler().Done(t.GetID())

	// TODO: Add checksum support
	err := t.GetDestinationStorage().WriteSegment(t.GetSegmentID(), t.GetOffset(), t.GetSize(), ioutil.NopCloser(t.GetContent()))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for Stream at Offset <%d> finished.", "MultipartStreamUploadTask", t.GetOffset())
}

func (t *MultipartCompleteTask) run() {
	log.Debugf("Task <%s> for Object <%s> started.", "MultipartCompleteTask", t.GetDestinationPath())

	err := t.GetDestinationStorage().CompleteSegment(t.GetSegmentID())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for Object <%s> finished.", "MultipartCompleteTask", t.GetDestinationPath())
}
