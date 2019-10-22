package utils

import (
	"bytes"
	"crypto/md5"
	"io"
	"math/rand"

	"github.com/c2h5oh/datasize"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/pkg/types"
)

const maxTestSize = 64 * int64(datasize.MB)

var _ = uuid.New()

// EmptyTask is used for test.
type EmptyTask struct {
	types.ID
	types.Fault
	types.Pool
	types.Todo
}

// Run implement navvy.Task interface.
func (t *EmptyTask) Run() {
}

// NewCallbackTask will create a new callback test.
func NewCallbackTask(fn func()) *CallbackTask {
	return &CallbackTask{
		fn: fn,
	}
}

// CallbackTask is the callback task.
type CallbackTask struct {
	types.ID
	types.Fault
	types.Pool
	types.Todo

	fn func()
}

// Run implement navvy.Task interface.
func (t *CallbackTask) Run() {
	t.fn()
}

// GenerateTestStream will generate a test stream.
func GenerateTestStream() (buf *bytes.Buffer, size int64, md5sum []byte) {
	buf = bytes.NewBuffer(nil)
	size = rand.Int63n(maxTestSize)

	r := NewRand()
	h := md5.New()
	w := io.MultiWriter(buf, h)

	_, err := io.CopyN(w, r, size)
	if err != nil {
		log.Fatal(err)
	}
	md5sum = h.Sum(nil)
	return
}
