package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseByteSize(t *testing.T) {
	cases := []struct {
		input    string
		expected int64
		wantErr  bool
	}{
		{"1B", 1, false},
		{"1GB", 1024 * 1024 * 1024, false},
		{"1 GB", 1024 * 1024 * 1024, false},
		{"1 G", 1024 * 1024 * 1024, false},
		{"1 G", 1024 * 1024 * 1024, false},
		{"Gb", 1024 * 1024 * 1024, true},
		{"G", 1024 * 1024 * 1024, true},
	}

	for _, v := range cases {
		x, err := ParseByteSize(v.input)
		if v.wantErr {
			assert.Error(t, err)
			continue
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
		{"1PB", "", ErrReadableSizeFormat},
		{"1 ", "", ErrReadableSizeFormat},
		{" 1", "", ErrReadableSizeFormat},
	}
	for _, v := range cases {
		x, err := UnixReadableSize(v.input)
		assert.Equal(t, v.expected, x, v.input)
		if v.err != nil {
			assert.Equal(t, true, errors.Is(err, v.err), v.input)
		} else {
			assert.Nil(t, err, v.err)
		}
	}
}
