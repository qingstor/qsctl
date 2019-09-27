package utils

import (
	"bytes"
	"crypto/md5"
	"io"
	"io/ioutil"
	"math/rand"
	"os"

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
