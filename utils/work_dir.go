package utils

import (
	"path/filepath"
	"strings"
)

var separator = "/"

// ParseWd get a path as input, split the work dir and file by following rules
// 1. if the path is like /path/to/target, parse the wd as /path/to/ and file as target;
// 2. if the path is like /path/to/, parse the wd as /path/to/ and file as "";
// 3. if the path is like . , parse the wd as {pwd}/ and file as "".
func ParseWd(path string) (wd, file string, err error) {
	var absPath string
	absPath, err = filepath.Abs(path)
	if err != nil {
		return
	}
	// because filepath.Abs will clean the last '/', so we need add it back for dir to Split
	if strings.HasSuffix(path, separator) || strings.HasSuffix(path, ".") {
		absPath = absPath + separator
	}
	wd, file = filepath.Split(absPath)
	return
}
