package task

import (
	"errors"
	"io/ioutil"

	"github.com/Xuanwo/storage/pkg/iterator"
	typ "github.com/Xuanwo/storage/types"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *SegmentInitTask) new() {}
func (t *SegmentInitTask) run() {
	id, err := t.GetStorage().InitSegment(t.GetPath(),
		typ.WithPartSize(t.GetPartSize()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	t.SetSegmentID(id)
}

func (t *SegmentFileCopyTask) new() {}
func (t *SegmentFileCopyTask) run() {
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
}
func (t *SegmentStreamCopyTask) new() {}
func (t *SegmentStreamCopyTask) run() {
	// TODO: Add checksum support
	err := t.GetDestinationStorage().WriteSegment(t.GetSegmentID(), t.GetOffset(), t.GetSize(), ioutil.NopCloser(t.GetContent()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}
func (t *SegmentCompleteTask) new() {}
func (t *SegmentCompleteTask) run() {
	err := t.GetStorage().CompleteSegment(t.GetSegmentID())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}
func (t *SegmentAbortAllTask) new() {}
func (t *SegmentAbortAllTask) run() {
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
}
