package task

import (
	typ "github.com/Xuanwo/storage/types"

	"github.com/yunify/qsctl/v2/utils"
)

func (t *MoveDirTask) new() {}
func (t *MoveDirTask) run() {
	x := NewListDir(t)
	utils.ChooseSourceStorage(x, t)
	// TODO: we should handle dir here.
	x.SetFileFunc(func(o *typ.Object) {
		sf := NewMoveFile(t)
		sf.SetSourcePath(o.Name)
		sf.SetDestinationPath(o.Name)
		t.GetScheduler().Async(sf)
	})
	t.GetScheduler().Sync(x)
}

func (t *MoveFileTask) new() {}
func (t *MoveFileTask) run() {
	ct := NewCopyFile(t)
	t.GetScheduler().Sync(ct)

	dt := NewDeleteFile(t)
	utils.ChooseSourceStorage(dt, t)
	t.GetScheduler().Sync(dt)
}
