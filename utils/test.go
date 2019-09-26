package utils

import (
	"bytes"
	"crypto/md5"
	"io"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/c2h5oh/datasize"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/pkg/types"
)

const maxTestSize = 64 * int64(datasize.MB)

// EmptyTask is used for test.
type EmptyTask struct {
}

// Run implement navvy.Task interface.
func (t *EmptyTask) Run() {
}

// EmptyTasker is a valid Tasker for test.
type EmptyTasker struct {
	EmptyTask

	types.Todo
	types.Fault
	types.Pool
}

// GenerateTestFile will generate a test file.
func GenerateTestFile() (name string, size int64, md5sum []byte) {
	f, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		log.Fatal(err)
	}

	name = f.Name()
	size = rand.Int63n(maxTestSize)

	r := NewRand()
	h := md5.New()
	w := io.MultiWriter(f, h)

	_, err = io.CopyN(w, r, size)
	if err != nil {
		log.Fatal(err)
	}
	md5sum = h.Sum(nil)

	err = f.Sync()
	if err != nil {
		log.Fatal(err)
	}
	return
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
