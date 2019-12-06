package task

import (
	"errors"

	typ "github.com/Xuanwo/storage/types"

	"github.com/qingstor/qsctl/v2/pkg/types"
)

func (t *BetweenStorageCheckTask) new() {}
func (t *BetweenStorageCheckTask) run() {
	// Source Object must be exist.
	src, err := t.GetSourceStorage().Stat(t.GetSourcePath())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	t.SetSourceObject(src)

	// If Destination Object not exist, we will set DestinationObject to nil.
	// So we can check its existences later.
	dst, err := t.GetDestinationStorage().Stat(t.GetDestinationPath())
	if err != nil && !errors.Is(err, typ.ErrObjectNotExist) {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	t.SetDestinationObject(dst)
}

func (t *IsDestinationObjectExistTask) new() {}
func (t *IsDestinationObjectExistTask) run() {
	t.SetResult(t.GetDestinationObject() != nil)
}

func (t *IsSizeEqualTask) new() {}
func (t *IsSizeEqualTask) run() {
	t.SetResult(t.GetSourceObject().Size == t.GetDestinationObject().Size)
}

func (t *IsUpdateAtGreaterTask) new() {}
func (t *IsUpdateAtGreaterTask) run() {
	t.SetResult(
		t.GetSourceObject().UpdatedAt.After(t.GetDestinationObject().UpdatedAt),
	)
}
