package task

import (
	typ "github.com/Xuanwo/storage/types"
)

func (t *CopyTask) run() {
	switch t.GetSourceType() {
	case typ.ObjectTypeStream:
		t.GetScheduler().Sync(NewCopyStreamTask, t)
	case typ.ObjectTypeFile:
		t.GetScheduler().Sync(NewCopyFileTask, t)
	default:
		panic("not supported object type")
	}
}

func (t *CopyTask) new() {}
