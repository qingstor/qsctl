package task

import (
	typ "github.com/Xuanwo/storage/types"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/utils"
)

func (t *MoveDirTask) new() {}

func (t *MoveDirTask) run() {
	x := NewIterateFile(t)
	utils.ChooseSourceStorage(x, t)
	x.SetPathFunc(func(key string) {
		sf := NewMoveFile(t)
		sf.SetSourcePath(key)
		sf.SetDestinationPath(key)
		t.GetScheduler().Async(sf)
	})
	x.SetRecursive(true)
	t.GetScheduler().Sync(x)

	t.GetScheduler().Wait()
	// Use storage.Delete() directly, because DeleteDirTask will skip dir while iterating.
	if err := t.GetSourceStorage().Delete(t.GetSourcePath(), typ.WithRecursive(true)); err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}

func (t *MoveFileTask) new() {}

func (t *MoveFileTask) run() {
	ct := NewCopyFile(t)
	t.GetScheduler().Sync(ct)

	dt := NewDeleteFile(t)
	utils.ChooseSourceStorage(dt, t)
	t.GetScheduler().Sync(dt)
}
