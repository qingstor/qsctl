package task

import (
	"errors"
	"io/ioutil"

	"github.com/Xuanwo/storage/pkg/iterator"
	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *SegmentInitTask) new() {}
func (t *SegmentInitTask) run() {
	log.Debugf("Task <%s> for Object <%s> started.", "SegmentInitTask", t.GetPath())

	id, err := t.GetStorage().InitSegment(t.GetPath())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	t.SetSegmentID(id)

	log.Debugf("Task <%s> for Object <%s> finished.", "SegmentInitTask", t.GetPath())
}

func (t *SegmentFileCopyTask) new() {}
func (t *SegmentFileCopyTask) run() {
	log.Debugf("Task <%s> for File <%s> at Offset <%d> started.", "SegmentFileUploadTask", t.GetSourcePath(), t.GetOffset())

	r, err := t.GetSourceStorage().Read(t.GetSourcePath(), typ.WithSize(t.GetSize()), typ.WithOffset(t.GetOffset()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	defer r.Close()

	// TODO: Add checksum support.
	err = t.GetDestinationStorage().WriteSegment(t.GetSegmentID(), t.GetOffset(), t.GetSize(), r)
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for File <%s> at Offset <%d> finished.", "SegmentFileUploadTask", t.GetSourcePath(), t.GetOffset())
}
func (t *SegmentStreamCopyTask) new() {}
func (t *SegmentStreamCopyTask) run() {
	log.Debugf("Task <%s> for Stream at Offset <%d> started.", "SegmentStreamUploadTask", t.GetOffset())

	// TODO: Add checksum support
	err := t.GetDestinationStorage().WriteSegment(t.GetSegmentID(), t.GetOffset(), t.GetSize(), ioutil.NopCloser(t.GetContent()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for Stream at Offset <%d> finished.", "SegmentStreamUploadTask", t.GetOffset())
}
func (t *SegmentCompleteTask) new() {}
func (t *SegmentCompleteTask) run() {
	log.Debugf("Task <%s> for Object <%s> started.", "SegmentCompleteTask", t.GetPath())

	err := t.GetStorage().CompleteSegment(t.GetSegmentID())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for Object <%s> finished.", "SegmentCompleteTask", t.GetPath())
}
func (t *SegmentAbortAllTask) new() {}
func (t *SegmentAbortAllTask) run() {
	log.Debugf("Task <%s> for Bucket started.", "AbortSegmentTask")
	it := t.GetStorage().ListSegments("")

	for {
		o, err := it.Next()
		if err != nil && errors.Is(err, iterator.ErrDone) {
			break
		}
		if err != nil {
			t.TriggerFault(types.NewErrUnhandled(err))
			return
		}
		if err := t.GetStorage().AbortSegment(o.ID); err != nil {
			t.TriggerFault(types.NewErrUnhandled(err))
			return
		}
	}

	log.Debugf("Task <%s> for Bucket finished.", "AbortSegmentTask")
}
