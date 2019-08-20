package utils

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/constants"
)

func TestParseFilePathForRead(t *testing.T) {
	x, err := ParseFilePathForRead("-")
	assert.NoError(t, err)
	assert.Equal(t, os.Stdin, x)

	x, err = ParseFilePathForRead("/xxxxxxxxxxxx")
	assert.Equal(t, constants.ErrorFileNotExist, err)
	assert.Nil(t, x)

	file, err := ioutil.TempFile(os.TempDir(), "example")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())

	x, err = ParseFilePathForRead(file.Name())
	assert.NoError(t, err)
	assert.NotNil(t, x)
}

func TestParseFilePathForWrite(t *testing.T) {
	x, err := ParseFilePathForWrite("-")
	assert.NoError(t, err)
	assert.Equal(t, os.Stdout, x)

	file, err := ioutil.TempFile(os.TempDir(), "example")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())

	x, err = ParseFilePathForWrite(file.Name())
	assert.NoError(t, err)
	assert.NotNil(t, x)
}

func TestParseQsPath(t *testing.T) {
	cases := []struct {
		msg       string
		input     string
		expected1 string
		expected2 string
		err       error
	}{
		{"case1", "qs://abcdef/def", "abcdef", "def", nil},
		{"case2", "qs://abcdef", "abcdef", "", nil},
		{"case3", "qs://abcdef/", "abcdef", "", nil},
		{"case4", "qs://", "", "", constants.ErrorQsPathInvalid},
		{"case5", "qs://abcdef/def/ghi", "abcdef", "def/ghi", nil},
		{"case6", "abcdef", "abcdef", "", nil},
		{"case7", "abcdef/def/ghi", "abcdef", "def/ghi", nil},
		{"case8", "xx://abcdef", "", "", constants.ErrorQsPathInvalid},
		{"case9", "qs://a-bcdef", "a-bcdef", "", nil},
		{"case10", "abcdef/👾 🙇 💁 🙅 🙆 🙋 🙎 🙍 ", "abcdef", "👾 🙇 💁 🙅 🙆 🙋 🙎 🙍 ", nil},
		{"case11", "ABCDEF", "", "", constants.ErrorQsPathInvalid},
		{"case12", "-abced", "", "", constants.ErrorQsPathInvalid},
		{"case13", "qs://abcdef/def/", "abcdef", "def/", nil},
	}

	for _, v := range cases {
		x1, x2, err := ParseQsPath(v.input)
		assert.Equal(t, v.err, err, v.msg)
		assert.Equal(t, v.expected1, x1, v.msg)
		assert.Equal(t, v.expected2, x2, v.msg)
	}
}

func TestCalculatePartSize(t *testing.T) {
	cases := []struct {
		msg      string
		input    int64
		expected int64
		err      error
	}{
		{"1B", 1, constants.DefaultPartSize, nil},
		{"1G", 1024 * 1024 * 1024, constants.DefaultPartSize, nil},
		{"10G", 10 * 1024 * 1024 * 1024, constants.DefaultPartSize, nil},
		{"2TB", 2 * 1024 * 1024 * 1024 * 1024, constants.DefaultPartSize << 1, nil},
		{"10TB", 10 * 1024 * 1024 * 1024 * 1024, 1099511628, nil},
		{"100TB", 101 * 1024 * 1024 * 1024 * 1024, 0, constants.ErrorFileTooLarge},
	}

	for _, v := range cases {
		x, err := CalculatePartSize(v.input)
		assert.Equal(t, v.err, err, v.msg)
		assert.Equal(t, v.expected, x, v.msg)
	}
}
