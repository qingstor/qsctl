package common

import (
	"bytes"
	"crypto/md5"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/c2h5oh/datasize"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/task/utils"
)

func TestFileMD5SumTask_Run(t *testing.T) {
	x := &mockFileMD5SumTask{}

	f, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	x.SetPath(f.Name())

	size := rand.Int63n(5 * int64(datasize.MB))
	r := utils.NewRand()

	var buf bytes.Buffer

	_, err = io.CopyN(&buf, r, size)
	if err != nil {
		t.Fatal(err)
	}
	md5sum := md5.Sum(buf.Bytes())
	_, err = io.Copy(f, &buf)
	if err != nil {
		t.Fatal(err)
	}

	x.SetOffset(0)
	x.SetSize(size)

	task := NewFileMD5SumTask(x)
	task.Run()

	assert.Equal(t, x.GetMD5Sum(), md5sum[:])
}

func TestStreamMD5SumTask_Run(t *testing.T) {
	x := &mockStreamMD5SumTask{}

	size := rand.Int63n(5 * int64(datasize.MB))
	r := utils.NewRand()

	var buf bytes.Buffer

	_, err := io.CopyN(&buf, r, size)
	if err != nil {
		t.Fatal(err)
	}
	md5sum := md5.Sum(buf.Bytes())

	x.SetPath("")
	x.SetContent(&buf)

	task := NewStreamMD5SumTask(x)
	task.Run()

	assert.Equal(t, x.GetMD5Sum(), md5sum[:])
}
