package task

import (
	"github.com/Xuanwo/storage/types/pairs"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *CreateStorageTask) new() {}

func (t *CreateStorageTask) run() {
	_, err := t.GetService().Create(t.GetStorageName(), pairs.WithLocation(t.GetZone()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}
