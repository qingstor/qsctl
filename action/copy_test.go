package action

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

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
		// Package input params
		input := CopyHandler{
			FlagHandler: (&FlagHandler{}).WithExpectSize(expectSize),
		}

		err := input.WithSrc(v.inputSrc).WithDest(v.inputDest).Copy()
		assert.Equal(suite.T(), v.err, err, v.msg)
	}
}

func (suite CopyTestSuite) TestCopyNotSeekableFileToRemote() {
	size := int64(1024 * 1024 * 1024) // 1G
	r := io.LimitReader(utils.NewRand(), size)
	objectKey := uuid.New().String()

	// Package input params
	input := CopyHandler{
		FlagHandler: (&FlagHandler{}).WithBench(true).WithExpectSize(size),
	}

	total, err := input.WithObjectKey(objectKey).WithReader(r).CopyNotSeekableFileToRemote()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), size, total)
}

func (suite CopyTestSuite) TestCopyObjectToNotSeekableFile() {
	size := int64(1024 * 1024 * 1024) // 1G
	w := ioutil.Discard

	// Package input params
	input := CopyHandler{}
	total, err := input.WithObjectKey(storage.MockGBObject).WithWriter(w).CopyObjectToNotSeekableFile()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), size, total)
}

func TestCopyTestSuite(t *testing.T) {
	suite.Run(t, new(CopyTestSuite))
}
