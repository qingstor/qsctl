package task

import (
	"github.com/Xuanwo/storage/types"
	"github.com/yunify/qsctl/v2/utils"
)

func (t *SyncTask) new() {}

func (t *SyncTask) run() {
	x := NewListDir(t)
	utils.ChooseSourceStorage(x, t)
	x.SetFileFunc(func(o *types.Object) {
		sf := NewCopyFile(t)
		sf.SetSourcePath(o.Name)
		sf.SetDestinationPath(o.Name)
		t.GetScheduler().Async(sf)
	})
	x.SetRecursive(true)
	t.GetScheduler().Sync(x)
}
