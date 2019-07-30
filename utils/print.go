package utils

import (
	"strings"

	"github.com/jedib0t/go-pretty/text"
)

// AlignPrintWithColon will align print with ":"
func AlignPrintWithColon(s ...string) string {
	a := make([]string, 0)

	width := 0
	for _, v := range s {
		x := strings.Split(v, ":")
		if len(x) < 2 {
			panic("input string must have : for split")
		}
		if width < len(x[0]) {
			width = len(x[0])
		}
	}

	for _, v := range s {
		x := strings.Split(v, ":")
		a = append(a, text.AlignRight.Apply(x[0], width)+text.AlignLeft.Apply(v[len(x[0]):], 1))
	}
	return strings.Join(a, "\n")
}

// AlignLinux will align print by lines (element in slice)
func AlignLinux(input ...[]string) [][]string {
	if len(input) == 0 {
		return nil
	}
	// init colNum equal with the first line
	lineNum, colNum := len(input), len(input[0])
	res := make([][]string, 0, lineNum)
	widths := make([]int, colNum)
	for _, line := range input {
		for i, str := range line {
			// if current line column greater than colNum,
			// directly append str-len to width, and colNum++
			if i >= colNum {
				widths = append(widths, len(str))
				colNum++
				continue
			}
			if widths[i] < len(str) {
				widths[i] = len(str)
			}
		}
	}
	for _, line := range input {
		newLine := make([]string, colNum)
		for i, str := range line {
			// align left if last column
			if i == colNum-1 {
				newLine[i] = text.AlignLeft.Apply(str, 1)
			} else {
				newLine[i] = text.AlignRight.Apply(str, widths[i])
			}
		}
		res = append(res, newLine)
	}
	return res
}
