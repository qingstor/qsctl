//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
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
		{
			name: "monkey err",
			args: args{
				path: "whatever",
			},
			wantWd:   "",
			wantFile: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.NewPatches()
			if tt.wantErr {
				patches.ApplyFunc(filepath.Abs, func(string) (string, error) {
					return "", filepath.ErrBadPattern
				})
			}
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
			patches.Reset()
		})
	}
}
