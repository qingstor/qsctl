package task

import (
	"io/ioutil"

	"github.com/Xuanwo/storage/types/pairs"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *SegmentInitTask) new() {}
func (t *SegmentInitTask) run() {
	id, err := t.GetSegmenter().InitSegment(t.GetPath(),
		pairs.WithPartSize(t.GetPartSize()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	t.SetSegmentID(id)
}

func (t *SegmentFileCopyTask) new() {}
func (t *SegmentFileCopyTask) run() {
	r, err := t.GetSourceStorage().Read(t.GetSourcePath(), pairs.WithSize(t.GetSize()), pairs.WithOffset(t.GetOffset()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	defer r.Close()

	// TODO: Add checksum support.
	err = t.GetDestinationSegmenter().WriteSegment(t.GetSegmentID(), t.GetOffset(), t.GetSize(), r)
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}

func (t *SegmentStreamCopyTask) new() {}
func (t *SegmentStreamCopyTask) run() {
	// TODO: Add checksum support
	err := t.GetSegmenter().WriteSegment(t.GetSegmentID(), t.GetOffset(), t.GetSize(), ioutil.NopCloser(t.GetContent()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}

func (t *SegmentCompleteTask) new() {}
func (t *SegmentCompleteTask) run() {
	err := t.GetSegmenter().CompleteSegment(t.GetSegmentID())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}
