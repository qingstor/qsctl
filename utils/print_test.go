package utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlignPrintWithColon(t *testing.T) {
	cases := []struct {
		name      string
		input     []string
		expected  string
		wantPanic bool
	}{
		{
			name: "normal",
			input: []string{
				"e1: test1",
				"e12: example2",
				"long1: s3",
			},
			expected: "   e1: test1\n" +
				"  e12: example2\n" +
				"long1: s3",
			wantPanic: false,
		},
		{
			name: "panic",
			input: []string{
				"string without colon",
			},
			expected:  "",
			wantPanic: true,
		},
	}
	for _, v := range cases {
		if v.wantPanic {
			assert.Panics(t, func() {
				AlignPrintWithColon(v.input...)
			}, v.name)
		} else {
			x := AlignPrintWithColon(v.input...)
			assert.Equal(t, v.expected, x, v.name)
		}
	}
}

func TestAlignLinux(t *testing.T) {
	ori := make([][]string, 0)
	ori = append(ori, []string{"a00-t", "a01-text", "a02-longtext"})
	ori = append(ori, []string{"a10", "a11", "a12", "a13"})
	ori = append(ori, []string{"a20", "a21-longtext", "a22-long", "a23-long"})
	ori = AlignLinux(ori...)

	assert.Equal(t, len(ori), 3)

	res := make([]string, len(ori))
	for i, line := range ori {
		res[i] = strings.Join(line, " ")
		fmt.Println(res[i])
	}

	expected :=
		"a00-t     a01-text a02-longtext /n" +
			"  a10          a11          a12 a13/n" +
			"  a20 a21-longtext     a22-long a23-long"

	assert.Equal(t, strings.Join(res, "/n"), expected, "normal")

	var input [][]string
	input = AlignLinux(input...)
	assert.Equal(t, 0, len(input), "nil condition")
}
