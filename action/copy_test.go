package action

import (
	"errors"
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
	f, err := NewMockReadAtCloseSeeker(size)
	defer f.Close()
	if err != nil {
		suite.T().Fatal(err)
	}
	total, err := CopySeekableFileToRemote(f, storage.MockMBObject)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), size, total)
}

func TestCopyTestSuite(t *testing.T) {
	suite.Run(t, new(CopyTestSuite))
}

// MockReadAtCloseSeeker implements the io.Reader, io.ReaderAt, io.Seeker
// and io.Closer interfaces by reading from a byte slice.
// a MockReadAtCloseSeeker is read-only and supports seeking.
// The zero value for Reader operates like a Reader of an empty slice.
type MockReadAtCloseSeeker struct {
	s        []byte
	i        int64 // current reading index
	prevRune int   // index of previous rune; or < 0
}

// Len returns the number of bytes of the unread portion of the
// slice.
func (r *MockReadAtCloseSeeker) Len() int {
	if r.i >= int64(len(r.s)) {
		return 0
	}
	return int(int64(len(r.s)) - r.i)
}

// Read implements the io.Reader interface.
func (r *MockReadAtCloseSeeker) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	r.prevRune = -1
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}

// ReadAt implements the io.ReaderAt interface.
func (r *MockReadAtCloseSeeker) ReadAt(b []byte, off int64) (n int, err error) {
	// cannot modify state - see io.ReaderAt
	if off < 0 {
		return 0, errors.New("bytes.Reader.ReadAt: negative offset")
	}
	if off >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(b, r.s[off:])
	if n < len(b) {
		err = io.EOF
	}
	return
}

// Seek implements the io.Seeker interface.
func (r *MockReadAtCloseSeeker) Seek(offset int64, whence int) (int64, error) {
	r.prevRune = -1
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.i + offset
	case io.SeekEnd:
		abs = int64(len(r.s)) + offset
	default:
		return 0, errors.New("bytes.Reader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("bytes.Reader.Seek: negative position")
	}
	r.i = abs
	return abs, nil
}

func (r *MockReadAtCloseSeeker) Close() error {
	return nil
}

// NewMockReadAtCloseSeeker return a MockReadAtCloseSeeker pointer
func NewMockReadAtCloseSeeker(size int64) (*MockReadAtCloseSeeker, error) {
	s := make([]byte, size)
	lr := io.LimitReader(utils.NewRand(), size)
	_, err := lr.Read(s)
	if err != nil {
		return nil, err
	}
	r := &MockReadAtCloseSeeker{s: s}
	return r, nil
}

