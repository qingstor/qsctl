package task

import (
	"errors"
	"io/ioutil"

	"github.com/Xuanwo/storage/pkg/iterator"
	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *SegmentInitTask) new() {}
func (t *SegmentInitTask) run() {
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

		// TODO: rethink the logic here.
		t.GetScheduler().Async(t.segmentInitTaskRequirement, t.GetScheduleFunc())
	}

	log.Debugf("Task <%s> for Object <%s> finished.", "MultipartInitTask", t.GetDestinationPath())
}
func (t *SegmentFileCopyTask) new() {}
func (t *SegmentFileCopyTask) run() {
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
func (t *SegmentStreamCopyTask) new() {}
func (t *SegmentStreamCopyTask) run() {
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
func (t *SegmentCompleteTask) new() {}
func (t *SegmentCompleteTask) run() {
	log.Debugf("Task <%s> for Object <%s> started.", "MultipartCompleteTask", t.GetPath())

	err := t.GetStorage().CompleteSegment(t.GetSegmentID())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for Object <%s> finished.", "MultipartCompleteTask", t.GetPath())
}
func (t *SegmentAbortAllTask) new() {}
func (t *SegmentAbortAllTask) run() {
	log.Debugf("Task <%s> for Bucket <%s> started.", "AbortMultipartTask")
	it := t.GetStorage().ListSegments("")

	for {
		o, err := it.Next()
		if err != nil && errors.Is(err, iterator.ErrDone) {
			break
		}
		if err != nil {
			t.TriggerFault(fault.NewUnhandled(err))
			return
		}
		if err := t.GetStorage().AbortSegment(o.ID); err != nil {
			t.TriggerFault(fault.NewUnhandled(err))
			return
		}
	}

	log.Debugf("Task <%s> for Bucket <%s> finished.", "AbortMultipartTask")
}
