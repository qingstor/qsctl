package utils

import (
	"fmt"
	"testing"
)

func ExampleParseQsWorkDir() {
	var wd, file string
	wd, file = ParseQsWorkDir("/path/to/file")
	fmt.Println(fmt.Sprintf("wd: <%s>, file: <%s>", wd, file))

	wd, file = ParseQsWorkDir("/path/to/")
	fmt.Println(fmt.Sprintf("wd: <%s>, file: <%s>", wd, file))

	wd, file = ParseQsWorkDir("/")
	fmt.Println(fmt.Sprintf("wd: <%s>, file: <%s>", wd, file))

	// Output:
	// wd: </path/to/>, file: <file>
	// wd: </path/to/>, file: <>
	// wd: </>, file: <>
}

func TestParseQsWorkDir(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantWd   string
		wantFile string
	}{
		{
			name:     "file without prefix",
			path:     "path/to/file",
			wantWd:   "/path/to/",
			wantFile: "file",
		},
		{
			name:     "file with prefix",
			path:     "/path/to/file",
			wantWd:   "/path/to/",
			wantFile: "file",
		},
		{
			name:     "dir without prefix",
			path:     "path/to/dir/",
			wantWd:   "/path/to/dir/",
			wantFile: "",
		},
		{
			name:     "dir with prefix",
			path:     "/path/to/dir/",
			wantWd:   "/path/to/dir/",
			wantFile: "",
		},
		{
			name:     "dir with more than one prefix",
			path:     "path/to/dir/",
			wantWd:   "/path/to/dir/",
			wantFile: "",
		},
		{
			name:     "dir with redundant separator",
			path:     "path///to///dir/",
			wantWd:   "/path/to/dir/",
			wantFile: "",
		},
		{
			name:     "root dir",
			path:     "/",
			wantWd:   "/",
			wantFile: "",
		},
		{
			name:     "blank dir",
			path:     "",
			wantWd:   "/",
			wantFile: "",
		},
		{
			name:     "redundant dir",
			path:     "////",
			wantWd:   "/",
			wantFile: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWd, gotFile := ParseQsWorkDir(tt.path)
			if gotWd != tt.wantWd {
				t.Errorf("ParseQsWorkDir() gotWd = %v, want %v", gotWd, tt.wantWd)
			}
			if gotFile != tt.wantFile {
				t.Errorf("ParseQsWorkDir() gotFile = %v, want %v", gotFile, tt.wantFile)
			}
		})
	}
}
