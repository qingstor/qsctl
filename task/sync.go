package task

import "github.com/yunify/qsctl/v2/utils"

func (t *SyncTask) new() {}

func (t *SyncTask) run() {
	x := NewIterateFile(t)
	utils.ChooseSourceStorage(x, t)
	x.SetPathFunc(func(key string) {
		sf := NewCopyFile(t)
		sf.SetSourcePath(key)
		sf.SetDestinationPath(key)
		t.GetScheduler().Async(sf)
	})
	x.SetRecursive(true)
	t.GetScheduler().Sync(x)
}
