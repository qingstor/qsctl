package common

import (
	"crypto/md5"
	"io"

	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *FileMD5SumTask) run() {
	log.Debugf("Task <%s> for File <%s> at Offset <%d> started.", "FileMD5SumTask", t.GetSourcePath(), t.GetOffset())

	r, err := t.GetSourceStorage().Read(t.GetSourcePath(), typ.WithSize(t.GetSize()), typ.WithOffset(t.GetOffset()))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	defer r.Close()

	h := md5.New()
	_, err = io.Copy(h, r)
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	t.SetMD5Sum(h.Sum(nil)[:])
	log.Debugf("Task <%s> for File <%s> at Offset <%d> finished.", "FileMD5SumTask", t.GetSourcePath(), t.GetOffset())
}

func (t *StreamMD5SumTask) run() {
	log.Debugf("Task <%s> for Stream started.", "StreamMD5SumTask")

	md5Sum := md5.Sum(t.GetContent().Bytes())
	t.SetMD5Sum(md5Sum[:])

	log.Debugf("Task <%s> for Stream finished.", "StreamMD5SumTask")
}
