package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
