package action

import (
	"context"
	"io"
	"io/ioutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/utils"
)

type CopyTestSuite struct {
	suite.Suite
}

func (suite CopyTestSuite) SetupTest() {
	contexts.Storage = storage.NewMockObjectStorage()
}

func (suite CopyTestSuite) TestCopy() {
	expectSize := int64(1024 * 1024)
	cases := []struct {
		msg       string
		inputSrc  string
		inputDest string
		err       error
	}{
		{"stdin to remote", "-", "qs://bucket/object", nil},
		{"remote to stdout", "qs://bucket/" + storage.Mock0BObject, "-", nil},
	}

	for _, v := range cases {
		// Package context
		var ctx context.Context
		ctx = contexts.NewMockCmdContext()
		ctx = contexts.SetContext(ctx, constants.BenchFlag, true)
		ctx = contexts.SetContext(ctx, constants.ExpectSizeFlag, expectSize)
		ctx = contexts.SetContext(ctx, constants.MaximumMemoryContentFlag, int64(0))
		ctx = contexts.SetContext(ctx, constants.ZoneFlag, "")
		ctx = contexts.SetContext(ctx, "src", v.inputSrc)
		ctx = contexts.SetContext(ctx, "dest", v.inputDest)

		err := Copy(ctx)
		assert.Equal(suite.T(), v.err, err, v.msg)
	}
}

func (suite CopyTestSuite) TestCopyNotSeekableFileToRemote() {
	size := int64(1024 * 1024 * 1024) // 1G
	r := io.LimitReader(utils.NewRand(), size)
	objectKey := uuid.New().String()

	// Package context
	var ctx context.Context
	ctx = contexts.NewMockCmdContext()
	ctx = contexts.SetContext(ctx, constants.BenchFlag, true)
	ctx = contexts.SetContext(ctx, constants.ExpectSizeFlag, size)
	ctx = contexts.SetContext(ctx, constants.MaximumMemoryContentFlag, int64(0))
	ctx = contexts.SetContext(ctx, "objectKey", objectKey)
	ctx = contexts.SetContext(ctx, "reader", r)

	total, err := CopyNotSeekableFileToRemote(ctx)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), size, total)
}

func (suite CopyTestSuite) TestCopyObjectToNotSeekableFile() {
	size := int64(1024 * 1024 * 1024) // 1G
	w := ioutil.Discard

	// Package context
	var ctx context.Context
	ctx = contexts.NewMockCmdContext()
	ctx = contexts.SetContext(ctx, constants.BenchFlag, true)
	ctx = contexts.SetContext(ctx, "objectKey", storage.MockGBObject)
	ctx = contexts.SetContext(ctx, "writer", w)

	total, err := CopyObjectToNotSeekableFile(ctx)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), size, total)
}

func TestCopyTestSuite(t *testing.T) {
	suite.Run(t, new(CopyTestSuite))
}
