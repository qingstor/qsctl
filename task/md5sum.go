package task

import (
	"crypto/md5"
	"io"

	"github.com/Xuanwo/storage/types/pairs"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *MD5SumFileTask) new() {}
func (t *MD5SumFileTask) run() {
	r, err := t.GetStorage().Read(t.GetPath(), pairs.WithSize(t.GetSize()), pairs.WithOffset(t.GetOffset()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	defer r.Close()

	h := md5.New()
	_, err = io.Copy(h, r)
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}

	t.SetMD5Sum(h.Sum(nil)[:])
}

func (t *MD5SumStreamTask) new() {}
func (t *MD5SumStreamTask) run() {
	md5Sum := md5.Sum(t.GetContent().Bytes())
	t.SetMD5Sum(md5Sum[:])
}
