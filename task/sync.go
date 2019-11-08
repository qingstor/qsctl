package task

import (
	"errors"
	"time"

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
		_, err := t.GetSourceStorage().Stat(key)
		if err != nil && !errors.Is(err, typ.ErrObjectNotExist) {
			t.TriggerFault(types.NewErrUnhandled(err))
			return
		}
		if err != nil && errors.Is(err, typ.ErrObjectNotExist) {
			sf := NewDeleteFile(t)
			utils.ChooseDestinationStorage(sf, t)
			sf.SetPath(key)
			t.GetScheduler().Async(sf)
		}
	})
	df.SetRecursive(true)
	t.GetScheduler().Sync(df)
}

func (t *SyncFileTask) new() {}

func (t *SyncFileTask) run() {
	needCopy, err := t.needCopy()
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	if !needCopy {
		return
	}

	sf := NewCopyFile(t)
	t.GetScheduler().Async(sf)
}

// needCopy checks flags and time and return whether an object should be copied or not.
func (t *SyncFileTask) needCopy() (bool, error) {
	var srcUpdate time.Time
	var dstUpdate time.Time

	if t.GetWholeFile() {
		return true, nil
	}

	dstObj, err := t.GetDestinationStorage().Stat(t.GetDestinationPath())
	// if got error, and error not not-exist
	if err != nil && !errors.Is(err, typ.ErrObjectNotExist) {
		return false, err
	}
	// if obj does not exist
	if err != nil && errors.Is(err, typ.ErrObjectNotExist) {
		// if existing was set, don't copy, otherwise, copy
		return !t.GetExisting(), nil
	}
	// if obj exists, and set update flag, don't copy
	if t.GetUpdate() {
		return false, nil
	}

	srcObj, err := t.GetSourceStorage().Stat(t.GetSourcePath())
	if err != nil {
		return false, err
	}
	dstUpdate, dstOk := dstObj.GetUpdatedAt()
	srcUpdate, srcOk := srcObj.GetUpdatedAt()
	// both update get and src is newer than dst, then copy
	if dstOk && srcOk && dstUpdate.Unix() < srcUpdate.Unix() {
		return true, nil
	}
	return false, nil
}
