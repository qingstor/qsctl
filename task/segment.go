package task

import (
	"errors"
	"io/ioutil"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/types/pairs"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *SegmentInitTask) new() {}
func (t *SegmentInitTask) run() {
	segmenter, ok := t.GetStorage().(storage.Segmenter)
	if !ok {
		t.TriggerFault(types.NewErrUnhandled(errors.New("no supported")))
		return
	}

	id, err := segmenter.InitSegment(t.GetPath(),
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

	segmenter, ok := t.GetDestinationStorage().(storage.Segmenter)
	if !ok {
		t.TriggerFault(types.NewErrUnhandled(errors.New("no supported")))
		return
	}

	// TODO: Add checksum support.
	err = segmenter.WriteSegment(t.GetSegmentID(), t.GetOffset(), t.GetSize(), r)
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}
func (t *SegmentStreamCopyTask) new() {}
func (t *SegmentStreamCopyTask) run() {
	segmenter, ok := t.GetDestinationStorage().(storage.Segmenter)
	if !ok {
		t.TriggerFault(types.NewErrUnhandled(errors.New("no supported")))
		return
	}

	// TODO: Add checksum support
	err := segmenter.WriteSegment(t.GetSegmentID(), t.GetOffset(), t.GetSize(), ioutil.NopCloser(t.GetContent()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}
func (t *SegmentCompleteTask) new() {}
func (t *SegmentCompleteTask) run() {
	segmenter, ok := t.GetStorage().(storage.Segmenter)
	if !ok {
		t.TriggerFault(types.NewErrUnhandled(errors.New("no supported")))
		return
	}

	err := segmenter.CompleteSegment(t.GetSegmentID())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}
func (t *SegmentAbortAllTask) new() {}
func (t *SegmentAbortAllTask) run() {
	segmenter, ok := t.GetStorage().(storage.Segmenter)
	if !ok {
		t.TriggerFault(types.NewErrUnhandled(errors.New("no supported")))
		return
	}

	err := segmenter.ListSegments("", pairs.WithSegmentFunc(func(segment *segment.Segment) {
		if err := segmenter.AbortSegment(segment.ID); err != nil {
			t.TriggerFault(types.NewErrUnhandled(err))
			return
		}
	}))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}
