package task

import (
	"errors"

	typ "github.com/Xuanwo/storage/types"

	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *CopyCheckTask) new() {}

func (t *CopyCheckTask) run() {
	if t.GetWholeFile() {
		t.SetPassed(true)
		return
	}

	dstObj, err := t.GetDestinationStorage().Stat(t.GetDestinationPath())
	// if got error, and error not not-exist
	if err != nil && !errors.Is(err, typ.ErrObjectNotExist) {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	// if obj does not exist
	if err != nil && errors.Is(err, typ.ErrObjectNotExist) {
		// if existing was set, don't copy, otherwise, copy
		t.SetPassed(!t.GetExisting())
		return
	}
	// if obj exists, and set update flag, don't copy
	if t.GetUpdate() {
		t.SetPassed(false)
		return
	}

	srcObj, err := t.GetSourceStorage().Stat(t.GetSourcePath())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}

	dstUpdate, dstOk := dstObj.GetUpdatedAt()
	srcUpdate, srcOk := srcObj.GetUpdatedAt()
	// both update get and src is newer than dst, then copy
	if dstOk && srcOk && dstUpdate.Unix() < srcUpdate.Unix() {
		t.SetPassed(true)
		return
	}
	t.SetPassed(false)
}
