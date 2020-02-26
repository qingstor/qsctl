package utils

import (
	"fmt"
)

func ExampleParseWd() {
	var wd, file string
	wd, file, _ = ParseWd("/path/to/file")
	fmt.Println(fmt.Sprintf("wd: <%s>, file: <%s>", wd, file))

	wd, file, _ = ParseWd("/path/to/")
	fmt.Println(fmt.Sprintf("wd: <%s>, file: <%s>", wd, file))

	// wd, file, _ = ParseWd(".")
	// fmt.Println(fmt.Sprintf("wd: <%s>, file: <%s>", wd, file))
	// Will get --> wd: <{pwd}>, file: <>

	// Output:
	// wd: </path/to/>, file: <file>
	// wd: </path/to/>, file: <>
}
