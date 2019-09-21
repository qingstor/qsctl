package common

import (
	"bytes"
	"crypto/md5"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/c2h5oh/datasize"
	"github.com/magiconair/properties/assert"
	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/task/utils"
)

func TestFileUploadTask_Run(t *testing.T) {
	x := &mockFileUploadTask{}
	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}
	x.SetPool(pool)

	store := storage.NewMockObjectStorage()
	x.SetStorage(store)
	x.SetKey(storage.MockGBObject)

	f, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
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
	x.SetMD5Sum(md5sum[:])

	task := NewFileUploadTask(x)
	task.Run()
	pool.Wait()

	om, err := store.HeadObject(storage.MockGBObject)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, om.ContentLength, size)
}
