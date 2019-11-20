package task

import (
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
	x.SetFileFunc(func(o *typ.Object) {
		sf := NewSyncFile(t)
		sf.SetSourcePath(o.Name)
		sf.SetDestinationPath(o.Name)
		if t.GetIgnoreExisting() {
			sf.SetCheckFunc(func() bool {
				existence := NewCheckExistence(t)
				utils.ChooseDestinationStorage(existence, t)
				existence.SetPath(o.Name)
				t.GetScheduler().Sync(existence)
				if existence.ValidateBoolResult() && !existence.GetBoolResult() {
					return false
				}

				sizeTask := NewCheckSize(t)
				sizeTask.SetSourcePath(o.Name)
				sizeTask.SetDestinationPath(o.Name)
				t.GetScheduler().Sync(sizeTask)
				if sizeTask.ValidateCompareResult() && sizeTask.GetCompareResult() != 0 {
					return false
				}

				updateAtTask := NewCheckUpdateAt(t)
				updateAtTask.SetSourcePath(o.Name)
				updateAtTask.SetDestinationPath(o.Name)
				t.GetScheduler().Sync(updateAtTask)
				if updateAtTask.ValidateCompareResult() && updateAtTask.GetCompareResult() > 0 {
					return false
				}

				return true
			})
		}
		t.GetScheduler().Async(sf)
	})
	t.GetScheduler().Sync(x)
}

func (t *SyncFileTask) new() {}
func (t *SyncFileTask) run() {
	if t.GetCheckFunc()() {
		return
	}

	sf := NewCopyFile(t)
	t.GetScheduler().Sync(sf)
}
