package common

import (
	"bufio"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *MultipartInitTask) run() {
	log.Debugf("Task <%s> for Object <%s> started.", "MultipartInitTask", t.GetKey())

	uploadID, err := t.GetStorage().InitiateMultipartUpload(t.GetKey())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	t.SetUploadID(uploadID)

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

	err = t.GetStorage().UploadMultipart(t.GetKey(), t.GetUploadID(), t.GetSize(), t.GetPartNumber(), t.GetMD5Sum(), r)
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for File <%s> at Offset <%d> finished.", "MultipartFileUploadTask", t.GetPath(), t.GetOffset())
}

func (t *MultipartStreamUploadTask) run() {
	log.Debugf("Task <%s> for Stream at PartNumber <%d> started.", "MultipartStreamUploadTask", t.GetPartNumber())

	defer t.GetScheduler().Done(t.GetID())

	err := t.GetStorage().UploadMultipart(
		t.GetKey(), t.GetUploadID(), t.GetSize(),
		t.GetPartNumber(), t.GetMD5Sum(), t.GetContent())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for Stream at PartNumber <%d> finished.", "MultipartStreamUploadTask", t.GetPartNumber())
}

func (t *MultipartCompleteTask) run() {
	log.Debugf("Task <%s> for Object <%s> UploadID <%s> started.", "MultipartCompleteTask", t.GetKey(), t.GetUploadID())

	err := t.GetStorage().CompleteMultipartUpload(t.GetKey(), t.GetUploadID(), int(*t.GetCurrentPartNumber()))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for Object <%s> UploadID <%s> finished.", "MultipartCompleteTask", t.GetKey(), t.GetUploadID())
}
