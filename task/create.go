package task

import (
	typ "github.com/Xuanwo/storage/types"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *CreateStorageTask) new() {}

func (t *CreateStorageTask) run() {
	_, err := t.GetService().Create(t.GetStorageName(), typ.WithLocation(t.GetZone()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}
