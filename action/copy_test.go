package action

import (
	"io"
	"io/ioutil"
	"os"
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
	contexts.Bench = true
}

func (suite CopyTestSuite) TestCopy() {
	contexts.ExpectSize = int64(1024 * 1024)

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
		err := Copy(v.inputSrc, v.inputDest)
		assert.Equal(suite.T(), v.err, err, v.msg)
	}
}

func (suite CopyTestSuite) TestCopyNotSeekableFileToRemote() {
	size := int64(1024 * 1024 * 1024) // 1G

	contexts.ExpectSize = size

	r := io.LimitReader(utils.NewRand(), size)
	objectKey := uuid.New().String()

	total, err := CopyNotSeekableFileToRemote(r, objectKey)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), size, total)

	contexts.ExpectSize = 0
	_, err = CopyNotSeekableFileToRemote(r, objectKey)
	assert.Equal(suite.T(), err, constants.ErrorExpectSizeRequired)
}

func (suite CopyTestSuite) TestCopyObjectToNotSeekableFile() {
	size := int64(1024 * 1024 * 1024) // 1G

	w := ioutil.Discard

	total, err := CopyObjectToNotSeekableFile(w, storage.MockGBObject)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), size, total)
}

func (suite CopyTestSuite) TestCopySeekableFileToRemote() {
	size := int64(1024 * 1024) // 1M
	f, err := os.Create("seekable_test_file")
	defer os.Remove("seekable_test_file")
	defer f.Close()
	if err != nil {
		suite.T().Fatal(err)
	}
	total, err := CopyObjectToSeekableFile(f, storage.MockMBObject)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), size, total)
}

func TestCopyTestSuite(t *testing.T) {
	suite.Run(t, new(CopyTestSuite))
}
