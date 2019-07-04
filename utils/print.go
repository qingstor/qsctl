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
		a = append(a, text.AlignLeft.Apply(x[0], width)+text.AlignLeft.Apply(v[len(x[0]):], 1))
	}
	return strings.Join(a, "\n")
}
