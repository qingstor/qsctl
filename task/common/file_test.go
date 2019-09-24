package common

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/utils"
)

func TestFileUploadTask_Run(t *testing.T) {
	x := &mockFileUploadTask{}

	store := storage.NewMockObjectStorage()
	x.SetStorage(store)

	key := uuid.New().String()
	x.SetKey(key)

	name, size, md5sum := utils.GenerateTestFile()
	defer os.Remove(name)

	x.SetPath(name)
	x.SetMD5Sum(md5sum)

	task := NewFileUploadTask(x)
	task.Run()

	om, err := store.HeadObject(key)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, om.ContentLength, size)
}
