package utils

import (
	"errors"
	"testing"

	"bou.ke/monkey"
	"github.com/aos-dev/go-service-fs"
	"github.com/aos-dev/go-service-qingstor"
	"github.com/aos-dev/go-storage/v2"
	"github.com/aos-dev/go-storage/v2/pkg/endpoint"
	typ "github.com/aos-dev/go-storage/v2/types"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
)

var errTmp = errors.New("temp error")

func TestParseFlow(t *testing.T) {
	cases := []struct {
		input1   string
		input2   string
		expected constants.FlowType
	}{
		{"xxxx", "qs://xxxx", constants.FlowToRemote},
		{"qs://xxxx", "xxxx", constants.FlowToLocal},
		{"xxxx", "xxxx", constants.FlowInvalid},
		{"qs://xxxx", "qs://xxxx", constants.FlowInvalid},
		{"xxxx", "", constants.FlowAtRemote},
		{"qs://xxxx", "", constants.FlowAtRemote},
	}

	for _, v := range cases {
		x := ParseFlow(v.input1, v.input2)
		assert.Equal(t, v.expected, x)
	}
}

func TestParseQsPath(t *testing.T) {
	cases := []struct {
		input              string
		expectedKeyType    typ.ObjectType
		expectedBucketName string
		expectedKey        string
	}{
		{"qs://xxxx-bucket/abc", typ.ObjectTypeFile, "xxxx-bucket", "abc"},
		{"qs://abcdef", typ.ObjectTypeDir, "abcdef", ""},
		{"qs://abcdef/", typ.ObjectTypeDir, "abcdef", ""},
		{"qs://abcdef/def/ghi", typ.ObjectTypeFile, "abcdef", "def/ghi"},
		{"qs://abcdef/def/ghi/", typ.ObjectTypeDir, "abcdef", "def/ghi/"},
		{"abcdef", typ.ObjectTypeDir, "abcdef", ""},
		{"abcdef/", typ.ObjectTypeDir, "abcdef", ""},
		{"abcdef/def/ghi", typ.ObjectTypeFile, "abcdef", "def/ghi"},
		{"abcdef/👾 🙇 💁 🙅 🙆 🙋 🙎 🙍", typ.ObjectTypeFile, "abcdef", "👾 🙇 💁 🙅 🙆 🙋 🙎 🙍"},
	}

	for k, v := range cases {
		actualKeyType, actualBucketName, actualKey, err := ParseQsPath(v.input)
		assert.Equal(t, v.expectedKeyType, actualKeyType, k)
		assert.Equal(t, v.expectedBucketName, actualBucketName, k)
		assert.Equal(t, v.expectedKey, actualKey, k)
		assert.NoError(t, err, k)
	}
}

func TestParseStorageInputInvalidType(t *testing.T) {
	workDir, path, objectType, store, err := ParseStorageInput(uuid.New().String(), "invalid type")
	assert.Empty(t, workDir)
	assert.Empty(t, path)
	assert.Empty(t, objectType)
	assert.Empty(t, store)
	assert.True(t, errors.Is(err, ErrStoragerTypeInvalid))
}

func TestParseStorageInputQingstor(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		workDir string
		path    string
		pathErr error
		srvErr  error
		getErr  error
	}{
		{
			name:    "normal",
			input:   "qs://testaaa",
			workDir: "/",
			path:    "",
			pathErr: nil,
			srvErr:  nil,
			getErr:  nil,
		},
		{
			name:    "path error",
			input:   "qs://testaaa",
			workDir: "",
			path:    "",
			pathErr: errTmp,
			srvErr:  nil,
			getErr:  nil,
		},
		{
			name:    "new qingstor storage error",
			input:   "qs://testaaa",
			workDir: "",
			path:    "",
			pathErr: nil,
			srvErr:  errTmp,
			getErr:  nil,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			if v.pathErr != nil {
				monkey.Patch(ParseQsPath, func(p string) (_ typ.ObjectType, _, _ string, err error) {
					err = v.pathErr
					return
				})
			}

			monkey.Patch(NewQingStorStorage, func(...*typ.Pair) (stor storage.Storager, err error) {
				if v.srvErr != nil {
					err = v.srvErr
				} else {
					stor = &qingstor.Storage{}
				}
				return
			})

			gotWorkDir, gotPath, gotObjectType, gotStore, gotErr := ParseStorageInput(v.input, qingstor.Type)
			assert.Equal(t, v.pathErr == nil && v.srvErr == nil && v.getErr == nil, gotErr == nil)
			if gotErr == nil {
				assert.Equal(t, v.workDir, gotWorkDir, v.name)
				assert.Equal(t, v.path, gotPath, v.name)
				assert.NotZero(t, gotObjectType)
				assert.NotNil(t, gotStore)
			} else {
				assert.True(t, errors.Is(gotErr, errTmp))
			}
		})
	}
}

func TestParseServiceInput(t *testing.T) {
	cases := []struct {
		name         string
		servicerType StoragerType
		err          error
	}{
		{
			name:         "invalid",
			servicerType: "invalid",
			err:          ErrStoragerTypeInvalid,
		},
		{
			name:         "valid",
			servicerType: qingstor.Type,
			err:          nil,
		},
		{
			name:         "new service failed",
			servicerType: qingstor.Type,
			err:          errTmp,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			if v.err == errTmp {
				monkey.Patch(NewQingStorService, func() (_ storage.Servicer, err error) {
					err = errTmp
					return
				})
			}
			gotStore, gotErr := ParseServiceInput(v.servicerType)
			assert.Equal(t, v.err == nil, gotErr == nil)
			if v.err == nil {
				assert.NotNil(t, gotStore)
			} else {
				assert.True(t, errors.Is(gotErr, v.err))
			}
		})
	}
}

func TestParseAtServiceInput(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "normal",
			wantErr: false,
		},
		{
			name:    "error with parse",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			monkey.Patch(ParseServiceInput, func(serviceType StoragerType) (service storage.Servicer, err error) {
				assert.Equal(t, qingstor.Type, serviceType, tt.name)
				if tt.wantErr {
					err = errTmp
				}
				return
			})
			task := &taskutils.AtServiceTask{}
			if err := ParseAtServiceInput(task); (err != nil) != tt.wantErr {
				t.Errorf("ParseAtServiceInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseAtStorageInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantWorkDir string
		err         error
	}{
		{
			name:        "normal",
			input:       "qs://bucket/path/to/file",
			wantWorkDir: "/path/to/",
			err:         nil,
		},
		{
			name:        "error with parse flow",
			input:       "/bucket/path/to/file",
			wantWorkDir: "",
			err:         ErrInvalidFlow,
		},
		{
			name:        "error with parse storage",
			input:       "qs://bucket/path/to/file",
			wantWorkDir: "/path/to/",
			err:         errTmp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			monkey.Patch(ParseStorageInput, func(input string, storageType StoragerType) (
				workDir, path string, objectType typ.ObjectType, store storage.Storager, err error) {
				assert.Equal(t, qingstor.Type, storageType, tt.name)
				assert.Equal(t, tt.input, input, tt.name)
				_, _, key, _ := ParseQsPath(input)
				workDir, path = ParseQsWorkDir(key)
				if tt.err != nil {
					err = errTmp
				}
				return
			})

			if tt.err == ErrInvalidFlow {
				monkey.Patch(ParseFlow, func(_, _ string) constants.FlowType {
					return constants.FlowInvalid
				})
			}

			task := &taskutils.AtStorageTask{}
			workDir, err := ParseAtStorageInput(task, tt.input)
			assert.Equal(t, tt.err != nil, err != nil, tt.name)
			assert.Equal(t, tt.wantWorkDir, workDir, tt.name)
			if tt.err != nil {
				assert.True(t, errors.Is(err, tt.err), tt.name)
			}
		})
	}
}

func TestSetupDestinationStorage(t *testing.T) {
	cases := []struct {
		name      string
		path      string
		storeType typ.ObjectType
		stor      storage.Storager
	}{
		{
			name:      "fs",
			path:      uuid.New().String(),
			storeType: typ.ObjectType(uuid.New().String()),
			stor:      &fs.Storage{},
		},
		{
			name:      "qingstor",
			path:      uuid.New().String(),
			storeType: typ.ObjectType(uuid.New().String()),
			stor:      &qingstor.Storage{},
		},
	}

	for _, tt := range cases {
		task := &taskutils.BetweenStorageTask{}
		t.Run(tt.name, func(t *testing.T) {
			setupDestinationStorage(task, tt.path, tt.storeType, tt.stor)
			assert.Equal(t, tt.path, task.GetDestinationPath(), tt.name)
			assert.Equal(t, tt.storeType, task.GetDestinationType(), tt.name)
			assert.Equal(t, tt.stor, task.GetDestinationStorage(), tt.name)
		})
	}
}

func TestSetupStorage(t *testing.T) {
	cases := []struct {
		name      string
		path      string
		storeType typ.ObjectType
		stor      storage.Storager
	}{
		{
			name:      "fs",
			path:      uuid.New().String(),
			storeType: typ.ObjectType(uuid.New().String()),
			stor:      &fs.Storage{},
		},
		{
			name:      "qingstor",
			path:      uuid.New().String(),
			storeType: typ.ObjectType(uuid.New().String()),
			stor:      &qingstor.Storage{},
		},
	}

	for _, tt := range cases {
		task := &taskutils.AtStorageTask{}
		t.Run(tt.name, func(t *testing.T) {
			setupStorage(task, tt.path, tt.storeType, tt.stor)
			assert.Equal(t, tt.path, task.GetPath(), tt.name)
			assert.Equal(t, tt.storeType, task.GetType(), tt.name)
			assert.Equal(t, tt.stor, task.GetStorage(), tt.name)
		})
	}
}

func TestSetupService(t *testing.T) {
	cases := []struct {
		name string
		stor storage.Servicer
	}{
		{
			name: "qingstor",
			stor: &qingstor.Service{},
		},
	}

	for _, tt := range cases {
		task := &taskutils.AtServiceTask{}
		t.Run(tt.name, func(t *testing.T) {
			setupService(task, tt.stor)
			assert.Equal(t, tt.stor, task.GetService(), tt.name)
		})
	}
}

func TestSetupSourceStorage(t *testing.T) {
	cases := []struct {
		name      string
		path      string
		storeType typ.ObjectType
		stor      storage.Storager
	}{
		{
			name:      "fs",
			path:      uuid.New().String(),
			storeType: typ.ObjectType(uuid.New().String()),
			stor:      &fs.Storage{},
		},
		{
			name:      "qingstor",
			path:      uuid.New().String(),
			storeType: typ.ObjectType(uuid.New().String()),
			stor:      &qingstor.Storage{},
		},
	}

	for _, tt := range cases {
		task := &taskutils.BetweenStorageTask{}
		t.Run(tt.name, func(t *testing.T) {
			setupSourceStorage(task, tt.path, tt.storeType, tt.stor)
			assert.Equal(t, tt.path, task.GetSourcePath(), tt.name)
			assert.Equal(t, tt.storeType, task.GetSourceType(), tt.name)
			assert.Equal(t, tt.stor, task.GetSourceStorage(), tt.name)
		})
	}
}

func TestNewQingStorService(t *testing.T) {
	cases := []struct {
		name     string
		protocol string
		wantErr  bool
	}{
		{
			"https",
			endpoint.ProtocolHTTPS,
			false,
		},
		{
			"http",
			endpoint.ProtocolHTTP,
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(constants.ConfigProtocol, tt.protocol)
			srv, err := NewQingStorService()
			assert.Nil(t, err, tt.name)
			_, ok := srv.(*qingstor.Service)
			assert.True(t, ok, tt.name)
		})
	}
}

func TestNewQingStorStorage(t *testing.T) {
	cases := []struct {
		name     string
		protocol string
		wantErr  bool
	}{
		{
			"https",
			endpoint.ProtocolHTTPS,
			false,
		},
		{
			"http",
			endpoint.ProtocolHTTP,
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			monkey.Patch(qingstor.NewStorager, func(pairs ...*typ.Pair) (storage.Storager, error) {
				assert.Equal(t, 2, len(pairs), tt.name)
				return &qingstor.Storage{}, nil
			})

			viper.Set(constants.ConfigProtocol, tt.protocol)
			stor, err := NewQingStorStorage()
			assert.Nil(t, err, tt.name)
			_, ok := stor.(*qingstor.Storage)
			assert.True(t, ok, tt.name)
		})
	}
}
