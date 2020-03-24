package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ParseFsWorkDir get a path as input, split the work dir and file by following rules
// 1. if the path is like /path/to/target, parse the wd as /path/to/ and file as target;
// 2. if the path is like /path/to/, parse the wd as /path/to/ and file as "";
// 3. if the path is like . , parse the wd as {pwd}/ and file as "".
// ParseFsWorkDir use os.PathSeparator as the separator, for capability of windows.
func ParseFsWorkDir(path string) (wd, file string, err error) {
	separator := string(os.PathSeparator)
	// parse blank wd as dot, aka pwd
	if path == "" {
		path = "."
	}
	var absPath string
	absPath, err = filepath.Abs(path)
	if err != nil {
		return
	}
	// because filepath.Abs will clean the last '/', so we need add it back for dir to Split
	if strings.HasSuffix(path, separator) || strings.HasSuffix(path, ".") {
		// add TrimRight to remove more than one separator
		absPath = strings.TrimRight(absPath, separator) + separator
	}
	wd, file = filepath.Split(absPath)
	return
}

// ParseQsWorkDir get a path as input, split the work dir and file by like ParseFsWorkDir
// What's difference is ParseQsWorkDir use '/' as separator, which is defined by qingstor service.
func ParseQsWorkDir(path string) (wd, file string) {
	// always treat qs path as abs path, so add "/" before
	separator := "/"
	path = separator + path
	parts := strings.SplitAfter(path, separator)
	// because path always start with "/", so parts' len will always longer than 2
	// trim redundant "/" in the middle of path (not the first one)
	for i := 1; i < len(parts); i++ {
		if parts[i] == separator {
			parts[i] = ""
		}
	}
	wd = strings.Join(parts[0:len(parts)-1], "")
	file = parts[len(parts)-1]
	return
}
