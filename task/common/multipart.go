package common

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *MultipartInitTask) run() {
	log.Debugf("Task <%s> for Object <%s> started.", "MultipartInitTask", t.GetKey())

	err := t.GetDestinationStorage().InitSegment(t.GetKey())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	for {
		if *t.GetCurrentOffset() == t.GetTotalSize() {
			break
		}

		t.GetScheduler().New(t.multipartInitTaskRequirement)
	}

	log.Debugf("Task <%s> for Object <%s> finished.", "MultipartInitTask", t.GetKey())
}

func (t *MultipartFileUploadTask) run() {
	log.Debugf("Task <%s> for File <%s> at Offset <%d> started.", "MultipartFileUploadTask", t.GetPath(), t.GetOffset())

	defer t.GetScheduler().Done(t.GetID())

	f, err := os.Open(t.GetPath())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	defer f.Close()

	r := bufio.NewReader(io.NewSectionReader(f, t.GetOffset(), t.GetSize()))

	// TODO: Add checksum support.
	// TODO: storage should not handle file close?
	err = t.GetDestinationStorage().WriteSegment(t.GetKey(), t.GetOffset(), t.GetSize(), ioutil.NopCloser(r))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for File <%s> at Offset <%d> finished.", "MultipartFileUploadTask", t.GetPath(), t.GetOffset())
}

func (t *MultipartStreamUploadTask) run() {
	log.Debugf("Task <%s> for Stream at Offset <%d> started.", "MultipartStreamUploadTask", t.GetOffset())

	defer t.GetScheduler().Done(t.GetID())

	// TODO: Add checksum support
	err := t.GetDestinationStorage().WriteSegment(t.GetKey(), t.GetOffset(), t.GetSize(), ioutil.NopCloser(t.GetContent()))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for Stream at Offset <%d> finished.", "MultipartStreamUploadTask", t.GetOffset())
}

func (t *MultipartCompleteTask) run() {
	log.Debugf("Task <%s> for Object <%s> started.", "MultipartCompleteTask", t.GetKey())

	err := t.GetDestinationStorage().CompleteSegment(t.GetKey())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for Object <%s> finished.", "MultipartCompleteTask", t.GetKey())
}
