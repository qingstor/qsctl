package task

import (
	"errors"

	typ "github.com/Xuanwo/storage/types"
)

func (t *IsDestinationObjectExistTask) new() {}
func (t *IsDestinationObjectExistTask) run() {
	_, err := t.GetDestinationStorage().Stat(t.GetDestinationPath())
	if err == nil {
		t.SetResult(true)
		return
	}
	if errors.Is(err, typ.ErrObjectNotExist) {
		t.SetResult(false)
		return
	}
	t.TriggerFault(err)
}

func (t *IsSizeEqualTask) new() {}
func (t *IsSizeEqualTask) run() {
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
	t.SetResult(src.Size == dst.Size)
}

func (t *IsUpdateAtGreaterTask) new() {}
func (t *IsUpdateAtGreaterTask) run() {
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
	t.SetResult(src.UpdatedAt.After(dst.UpdatedAt))
}
