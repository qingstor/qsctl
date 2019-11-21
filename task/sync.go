package task

import (
	"github.com/Xuanwo/navvy"
	typ "github.com/Xuanwo/storage/types"
	"github.com/yunify/qsctl/v2/utils"
)

func (t *SyncTask) new() {}
func (t *SyncTask) run() {
	x := NewListDir(t)
	utils.ChooseSourceStorage(x, t)
	x.SetDirFunc(func(o *typ.Object) {
		sf := NewSync(t)
		sf.SetSourcePath(o.Name)
		sf.SetDestinationPath(o.Name)
		t.GetScheduler().Sync(sf)
	})

	var fn []func(task navvy.Task) navvy.Task
	if t.GetIgnoreExisting() {
		fn = append(fn,
			NewIsDestinationObjectExistTask,
			NewIsSizeEqualTask,
			NewIsUpdateAtGreaterTask,
		)
	}
	x.SetFileFunc(func(o *typ.Object) {
		sf := NewCopyFile(t)
		sf.SetSourcePath(o.Name)
		sf.SetDestinationPath(o.Name)
		sf.SetCheckTasks(fn)

		t.GetScheduler().Async(sf)
	})
	t.GetScheduler().Sync(x)
}
