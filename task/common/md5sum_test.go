package common

import (
	"bytes"
	"crypto/md5"
	"io"
	"io/ioutil"
	"testing"

	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yunify/qsctl/v2/pkg/mock"

	"github.com/yunify/qsctl/v2/utils"
)

func TestFileMD5SumTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	name := uuid.New().String()
	size := int64(1234)
	md5sum := md5.Sum([]byte("abc"))
	buf := bytes.NewBufferString("abc")

	srcStore := mock.NewMockStorager(ctrl)
	srcStore.EXPECT().Read(gomock.Any(), gomock.Any()).DoAndReturn(func(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
		assert.Equal(t, name, path)
		assert.Equal(t, size, pairs[0].Value.(int64))
		assert.Equal(t, int64(0), pairs[1].Value.(int64))
		return ioutil.NopCloser(buf), nil
	})

	x := &mockFileMD5SumTask{}
	x.SetSourcePath(name)
	x.SetSourceStorage(srcStore)
	x.SetOffset(0)
	x.SetSize(size)

	task := NewFileMD5SumTask(x)
	task.Run()

	assert.Equal(t, x.GetMD5Sum(), md5sum[:])
}

func TestStreamMD5SumTask_Run(t *testing.T) {
	x := &mockStreamMD5SumTask{}

	buf, _, md5sum := utils.GenerateTestStream()

	x.SetContent(buf)

	task := NewStreamMD5SumTask(x)
	task.Run()

	assert.Equal(t, x.GetMD5Sum(), md5sum[:])
}
