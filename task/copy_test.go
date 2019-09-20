package task

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/yunify/qsctl/v2/storage"
)

func TestCopyFile(t *testing.T) {
	s := storage.NewMockObjectStorage()
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("Hello, world")
	f.Sync()
	defer os.Remove(f.Name())

	task := NewCopyFileTask(f.Name(), storage.MockGBObject, s)
	task.Run()
	task.GetPool().Wait()
}
