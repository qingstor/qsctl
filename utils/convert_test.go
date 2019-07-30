package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/constants"
)

func TestParseByteSize(t *testing.T) {
	cases := []struct {
		input    string
		expected int64
	}{
		{"1B", 1},
		{"1GB", 1024 * 1024 * 1024},
		{"1 GB", 1024 * 1024 * 1024},
		{"1 G", 1024 * 1024 * 1024},
	}

	for _, v := range cases {
		x, err := ParseByteSize(v.input)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, v.expected, x)
	}
}

func TestUnixReadableSize(t *testing.T) {
	cases := []struct {
		input    string
		expected string
		err      error
	}{
		{"1 EB", "1E", nil},
		{"1 PB", "1P", nil},
		{"1 TB", "1T", nil},
		{"1 GB", "1G", nil},
		{"1 MB", "1M", nil},
		{"1 KB", "1K", nil},
		{"1 B", "1B", nil},
		{"1PB", "", constants.ErrorReadableSizeFormat},
		{"1 ", "", constants.ErrorReadableSizeFormat},
		{" 1", "", constants.ErrorReadableSizeFormat},
	}
	for _, v := range cases {
		x, err := UnixReadableSize(v.input)
		assert.Equal(t, v.expected, x)
		assert.Equal(t, v.err, err)
	}
}
