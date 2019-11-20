package task

import (
	"errors"

	typ "github.com/Xuanwo/storage/types"
)

func (t *CheckExistenceTask) new() {}
func (t *CheckExistenceTask) run() {
	_, err := t.GetStorage().Stat(t.GetPath())
	if err == nil {
		t.SetBoolResult(true)
		return
	}
	if errors.Is(err, typ.ErrObjectNotExist) {
		t.SetBoolResult(false)
		return
	}
	t.TriggerFault(err)
}

func (t *CheckSizeTask) new() {}
func (t *CheckSizeTask) run() {
	src, err := t.GetSourceStorage().Stat(t.GetSourcePath())
	if err != nil {
		t.TriggerFault(err)
		return
	}
	dst, err := t.GetDestinationStorage().Stat(t.GetDestinationPath())
	if err != nil {
		t.TriggerFault(err)
		return
	}

	delta := src.Size - dst.Size
	if delta > 0 {
		t.SetCompareResult(1)
	} else if delta < 0 {
		t.SetCompareResult(-1)
	} else {
		t.SetCompareResult(0)
	}
}

func (t *CheckUpdateAtTask) new() {}
func (t *CheckUpdateAtTask) run() {
	src, err := t.GetSourceStorage().Stat(t.GetSourcePath())
	if err != nil {
		t.TriggerFault(err)
		return
	}
	dst, err := t.GetDestinationStorage().Stat(t.GetDestinationPath())
	if err != nil {
		t.TriggerFault(err)
		return
	}

	delta := src.UpdatedAt.Sub(dst.UpdatedAt)
	if delta > 0 {
		t.SetCompareResult(1)
	} else if delta < 0 {
		t.SetCompareResult(-1)
	} else {
		t.SetCompareResult(0)
	}
}
