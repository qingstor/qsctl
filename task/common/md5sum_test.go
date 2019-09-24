package common

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/utils"
)

func TestFileMD5SumTask_Run(t *testing.T) {
	x := &mockFileMD5SumTask{}

	name, size, md5sum := utils.GenerateTestFile()

	x.SetPath(name)
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
