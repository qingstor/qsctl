package task

import (
	"errors"

	typ "github.com/Xuanwo/storage/types"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/utils"
)

func (t *SyncTask) new() {}

func (t *SyncTask) run() {
	x := NewIterateFile(t)
	utils.ChooseSourceStorage(x, t)
	x.SetPathFunc(func(key string) {
		sf := NewSyncFile(t)
		sf.SetSourcePath(key)
		sf.SetDestinationPath(key)
		t.GetScheduler().Async(sf)
	})
	x.SetRecursive(true)
	t.GetScheduler().Sync(x)

	// if delete flag not set, return now
	if !t.GetDelete() {
		return
	}
	// otherwise, iterate in destination storage and delete files not exist in source storage
	t.GetScheduler().Wait()
	df := NewIterateFile(t)
	utils.ChooseDestinationStorage(df, t)
	df.SetPathFunc(func(key string) {
		sf := NewSyncFileDelete(t)
		sf.SetDestinationPath(key)
		t.GetScheduler().Async(sf)
	})
	df.SetRecursive(true)
	t.GetScheduler().Sync(df)
}

func (t *SyncFileTask) new() {}

func (t *SyncFileTask) run() {
	checkTask := NewCopyCheck(t)
	t.GetScheduler().Sync(checkTask)

	if !checkTask.ValidatePassed() {
		return
	}

	if !checkTask.GetPassed() {
		return
	}

	sf := NewCopyFile(t)
	t.GetScheduler().Async(sf)
}

func (t *SyncFileDeleteTask) new() {}

func (t *SyncFileDeleteTask) run() {
	_, err := t.GetSourceStorage().Stat(t.GetDestinationPath())
	if err != nil && !errors.Is(err, typ.ErrObjectNotExist) {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	if err != nil && errors.Is(err, typ.ErrObjectNotExist) {
		sf := NewDeleteFile(t)
		utils.ChooseDestinationStorage(sf, t)
		t.GetScheduler().Sync(sf)
	}
}
