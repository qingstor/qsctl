package utils

import (
	"fmt"
	"os"
	"testing"
)

func ExampleParseFsWorkDir() {
	var wd, file string
	wd, file, _ = ParseFsWorkDir("/path/to/file")
	fmt.Println(fmt.Sprintf("wd: <%s>, file: <%s>", wd, file))

	wd, file, _ = ParseFsWorkDir("/path/to/")
	fmt.Println(fmt.Sprintf("wd: <%s>, file: <%s>", wd, file))

	wd, file, _ = ParseFsWorkDir("/")
	fmt.Println(fmt.Sprintf("wd: <%s>, file: <%s>", wd, file))

	// wd, file, _ = ParseWorkDir(".")
	// fmt.Println(fmt.Sprintf("wd: <%s>, file: <%s>", wd, file))
	// Will get --> wd: <{pwd}/>, file: <>

	// wd, file, _ = ParseWorkDir("-")
	// fmt.Println(fmt.Sprintf("wd: <%s>, file: <%s>", wd, file))
	// Will get --> wd: <{pwd}/>, file: <->

	// Output:
	// wd: </path/to/>, file: <file>
	// wd: </path/to/>, file: <>
	// wd: </>, file: <>
}

func TestParseWd(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		path string
	}
	tests := []struct {
		name     string
		args     args
		wantWd   string
		wantFile string
		wantErr  bool
	}{
		{
			name: "normal abs path",
			args: args{
				path: "/path/to/file",
			},
			wantWd:   "/path/to/",
			wantFile: "file",
			wantErr:  false,
		},
		{
			name: "normal rel path",
			args: args{
				path: "path/to/file",
			},
			wantWd:   fmt.Sprintf("%s/path/to/", pwd),
			wantFile: "file",
			wantErr:  false,
		},
		{
			name: "root path",
			args: args{
				path: "/",
			},
			wantWd:   "/",
			wantFile: "",
			wantErr:  false,
		},
		{
			name: "blank path",
			args: args{
				path: "",
			},
			wantWd:   fmt.Sprintf("%s/", pwd),
			wantFile: "",
			wantErr:  false,
		},
		{
			name: "dot path",
			args: args{
				path: ".",
			},
			wantWd:   fmt.Sprintf("%s/", pwd),
			wantFile: "",
			wantErr:  false,
		},
		{
			name: "stdin",
			args: args{
				path: "-",
			},
			wantWd:   fmt.Sprintf("%s/", pwd),
			wantFile: "-",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWd, gotFile, err := ParseFsWorkDir(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWorkDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWd != tt.wantWd {
				t.Errorf("ParseWorkDir() gotWd = %v, want %v", gotWd, tt.wantWd)
			}
			if gotFile != tt.wantFile {
				t.Errorf("ParseWorkDir() gotFile = %v, want %v", gotFile, tt.wantFile)
			}
		})
	}
}

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
