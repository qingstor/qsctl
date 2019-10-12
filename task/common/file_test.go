package common

import (
	"io"
	"os"
	"testing"

	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
	"github.com/yunify/qsctl/v2/utils"
)

func TestFileUploadTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	x := &mockFileUploadTask{}

	store := mock.NewMockStorager(ctrl)
	x.SetDestinationStorage(store)

	key := uuid.New().String()
	x.SetKey(key)

	name, size, md5sum := utils.GenerateTestFile()
	defer os.Remove(name)

	x.SetPath(name)
	x.SetSize(size)
	x.SetMD5Sum(md5sum)

	store.EXPECT().WriteFile(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(inputPath string, inputSize int64, r io.ReadCloser, option ...*types.Pair) {
			assert.Equal(t, key, inputPath)
			assert.Equal(t, size, inputSize)
		})

	task := NewFileUploadTask(x)
	task.Run()
}
