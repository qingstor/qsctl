package task

import (
	"io"
	"testing"

	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
)

func TestFileUploadTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	x := &mockFileUploadTask{}

	srcStore := mock.NewMockStorager(ctrl)
	x.SetSourceStorage(srcStore)

	dstStore := mock.NewMockStorager(ctrl)
	x.SetDestinationStorage(dstStore)

	mockReader := mock.NewMockReadCloser(ctrl)

	key := uuid.New().String()
	x.SetDestinationPath(key)

	name := uuid.New().String()
	size := int64(10)

	x.SetSourcePath(name)
	x.SetSize(size)

	mockReader.EXPECT().Close().Do(func() {})

	srcStore.EXPECT().Read(gomock.Any()).DoAndReturn(func(inputPath string) (r io.ReadCloser, err error) {
		assert.Equal(t, name, inputPath)
		return mockReader, nil
	})

	dstStore.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(inputPath string, r io.ReadCloser, option ...*types.Pair) {
			assert.Equal(t, key, inputPath)
			assert.Equal(t, size, option[0].Value.(int64))
		})

	task := NewFileUploadTask(x)
	task.Run()
}
