// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package utils

import (
	"errors"
	"os"
	"testing"

	"bou.ke/monkey"
	fs "github.com/aos-dev/go-service-fs"
	qingstor "github.com/aos-dev/go-service-qingstor"
	"github.com/aos-dev/go-storage/v2"
	typ "github.com/aos-dev/go-storage/v2/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
)

func TestParseLocalPath(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		wantPathType typ.ObjectType
		wantErr      error
	}{
		{
			name:         "not exist file",
			path:         "/" + uuid.New().String(),
			wantPathType: typ.ObjectTypeFile,
			wantErr:      nil,
		},
		{
			name:         "not exist dir",
			path:         "/" + uuid.New().String() + "/",
			wantPathType: typ.ObjectTypeDir,
			wantErr:      nil,
		},
		{
			name:         "stream",
			path:         "-",
			wantPathType: typ.ObjectTypeStream,
			wantErr:      nil,
		},
		{
			name:         "path err",
			path:         uuid.New().String(),
			wantPathType: typ.ObjectTypeInvalid,
			wantErr:      errTmp,
		},
		{
			name:         "normal file",
			path:         "/etc/profile",
			wantPathType: typ.ObjectTypeFile,
			wantErr:      nil,
		},
		{
			name:         "normal dir",
			path:         "/etc/",
			wantPathType: typ.ObjectTypeDir,
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr != nil {
				monkey.Patch(os.Stat, func(path string) (os.FileInfo, error) {
					assert.Equal(t, tt.path, path, tt.name)
					return nil, tt.wantErr
				})
				defer monkey.UnpatchAll()
			}
			gotPathType, err := ParseLocalPath(tt.path)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ParseLocalPath() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if gotPathType != tt.wantPathType {
				t.Errorf("ParseLocalPath() gotPathType = %v, want %v", gotPathType, tt.wantPathType)
			}
		})
	}
}

func TestParseStorageInputFs(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		workDir  string
		path     string
		pathErr  error
		wdErr    error
		fsNewErr error
	}{
		{
			name:     "valid local path",
			input:    "/etc",
			workDir:  "/",
			path:     "etc",
			pathErr:  nil,
			wdErr:    nil,
			fsNewErr: nil,
		},
		{
			name:     "invalid path err",
			input:    "/etc",
			workDir:  "",
			path:     "",
			pathErr:  errTmp,
			wdErr:    nil,
			fsNewErr: nil,
		},
		{
			name:     "wd error",
			input:    "/etc",
			workDir:  "",
			path:     "",
			pathErr:  nil,
			wdErr:    errTmp,
			fsNewErr: nil,
		},
		{
			name:     "new fs storage error",
			input:    "/etc",
			workDir:  "/",
			path:     "etc",
			pathErr:  nil,
			wdErr:    nil,
			fsNewErr: errTmp,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			if v.pathErr != nil {
				monkey.Patch(ParseLocalPath, func(p string) (_ typ.ObjectType, err error) {
					assert.Equal(t, v.input, p, v.name)
					err = v.pathErr
					return
				})
			}
			if v.wdErr != nil {
				monkey.Patch(ParseFsWorkDir, func(p string) (_, _ string, err error) {
					assert.Equal(t, v.input, p, v.name)
					err = v.wdErr
					return
				})
			}
			if v.fsNewErr != nil {
				monkey.Patch(fs.NewStorager, func(pairs ...*typ.Pair) (_ storage.Storager, err error) {
					err = v.fsNewErr
					return
				})
			}
			gotWorkDir, gotPath, gotObjectType, gotStore, gotErr := ParseStorageInput(v.input, fs.Type)
			assert.Equal(t, v.pathErr == nil && v.wdErr == nil && v.fsNewErr == nil, gotErr == nil)
			if gotErr == nil {
				assert.NotZero(t, gotObjectType, v.name)
				assert.NotNil(t, gotStore, v.name)
			} else {
				assert.Nil(t, gotStore, v.name)
				assert.True(t, errors.Is(gotErr, errTmp), v.name)
			}
			assert.Equal(t, v.workDir, gotWorkDir, v.name)
			assert.Equal(t, v.path, gotPath, v.name)
		})
	}
}

func TestParseBetweenStorageInput(t *testing.T) {
	tests := []struct {
		name           string
		src            string
		dst            string
		wantSrcWorkDir string
		wantDstWorkDir string
		failType       StoragerType
		err            error
	}{
		{
			name:           "normal local to remote",
			src:            "/etc/host",
			dst:            "qs://bucket/path/to/dir/",
			wantSrcWorkDir: "/etc/",
			wantDstWorkDir: "/path/to/dir/",
			err:            nil,
		},
		{
			name:           "normal remote to local",
			src:            "qs://bucket/path/to/file",
			dst:            "/etc/host",
			wantSrcWorkDir: "/path/to/",
			wantDstWorkDir: "/etc/",
			err:            nil,
		},
		{
			name:           "invalid flow",
			src:            "qs://etc/host",
			dst:            "qs://bucket/path/to/file",
			wantSrcWorkDir: "",
			wantDstWorkDir: "",
			err:            ErrInvalidFlow,
		},
		{
			name:           "parse local to remote src failed",
			src:            "/etc/host",
			dst:            "qs://bucket/path/to/dir/",
			wantSrcWorkDir: "/etc/",
			wantDstWorkDir: "",
			failType:       fs.Type,
			err:            errTmp,
		},
		{
			name:           "parse local to remote dst failed",
			src:            "/etc/host",
			dst:            "qs://bucket/path/to/dir/",
			wantSrcWorkDir: "/etc/",
			wantDstWorkDir: "/path/to/dir/",
			failType:       qingstor.Type,
			err:            errTmp,
		},
		{
			name:           "parse remote to local src failed",
			src:            "qs://bucket/path/to/dir/",
			dst:            "/etc/host",
			wantSrcWorkDir: "/path/to/dir/",
			wantDstWorkDir: "",
			failType:       qingstor.Type,
			err:            errTmp,
		},
		{
			name:           "parse remote to local dst failed",
			src:            "qs://bucket/path/to/dir/",
			dst:            "/etc/host",
			wantSrcWorkDir: "/path/to/dir/",
			wantDstWorkDir: "/etc/",
			failType:       fs.Type,
			err:            errTmp,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			monkey.Patch(ParseStorageInput, func(input string, storageType StoragerType) (
				workDir, path string, objectType typ.ObjectType, store storage.Storager, err error) {
				switch storageType {
				case fs.Type:
					workDir, path, _ = ParseFsWorkDir(input)
				case qingstor.Type:
					_, _, key, _ := ParseQsPath(input)
					workDir, path = ParseQsWorkDir(key)
				}
				if tt.failType == storageType && tt.err != nil {
					err = errTmp
				}
				return
			})

			task := &taskutils.BetweenStorageTask{}
			srcWorkDir, dstWorkDir, err := ParseBetweenStorageInput(task, tt.src, tt.dst)
			assert.Equal(t, tt.err != nil, err != nil, tt.name)
			assert.Equal(t, tt.wantSrcWorkDir, srcWorkDir, tt.name)
			assert.Equal(t, tt.wantDstWorkDir, dstWorkDir, tt.name)
			if tt.err != nil {
				assert.True(t, errors.Is(err, tt.err), tt.name)
			}
		})
	}
}
