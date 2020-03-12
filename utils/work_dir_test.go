package utils

import (
	"fmt"
	"os"
	"testing"
)

func ExampleParseWorkDir() {
	var wd, file string
	wd, file, _ = ParseWorkDir("/path/to/file", "/")
	fmt.Println(fmt.Sprintf("wd: <%s>, file: <%s>", wd, file))

	wd, file, _ = ParseWorkDir("/path/to/", "/")
	fmt.Println(fmt.Sprintf("wd: <%s>, file: <%s>", wd, file))

	wd, file, _ = ParseWorkDir("/", "/")
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
		path      string
		separator string
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
				path:      "/path/to/file",
				separator: "/",
			},
			wantWd:   "/path/to/",
			wantFile: "file",
			wantErr:  false,
		},
		{
			name: "normal rel path",
			args: args{
				path:      "path/to/file",
				separator: "/",
			},
			wantWd:   fmt.Sprintf("%s/path/to/", pwd),
			wantFile: "file",
			wantErr:  false,
		},
		{
			name: "root path",
			args: args{
				path:      "/",
				separator: "/",
			},
			wantWd:   "/",
			wantFile: "",
			wantErr:  false,
		},
		{
			name: "blank path",
			args: args{
				path:      "",
				separator: "/",
			},
			wantWd:   fmt.Sprintf("%s/", pwd),
			wantFile: "",
			wantErr:  false,
		},
		{
			name: "dot path",
			args: args{
				path:      ".",
				separator: "/",
			},
			wantWd:   fmt.Sprintf("%s/", pwd),
			wantFile: "",
			wantErr:  false,
		},
		{
			name: "stdin",
			args: args{
				path:      "-",
				separator: "/",
			},
			wantWd:   fmt.Sprintf("%s/", pwd),
			wantFile: "-",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWd, gotFile, err := ParseWorkDir(tt.args.path, tt.args.separator)
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
