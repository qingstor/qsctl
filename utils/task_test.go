package utils

import (
	"testing"

	"github.com/Xuanwo/storage/services/posixfs"
	"github.com/Xuanwo/storage/services/qingstor"
	typ "github.com/Xuanwo/storage/types"
	"github.com/stretchr/testify/assert"
	"github.com/yunify/qsctl/v2/constants"
)

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

func TestParseStorageInput(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		storageType typ.StoragerType
		hasPanic    bool
		err         error
	}{
		{
			"invalid storager type",
			"qs://testaaa",
			"test",
			true,
			nil,
		},
		{
			"valid local path",
			"/etc",
			posixfs.StoragerType,
			false,
			nil,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			if v.hasPanic {
				assert.Panics(t, func() {
					_, _, _, _ = ParseStorageInput(v.input, v.storageType)
				})
				return
			}

			gotPath, gotObjectType, gotStore, gotErr := ParseStorageInput(v.input, v.storageType)
			assert.Equal(t, v.err == nil, gotErr == nil)
			if v.err == nil {
				assert.Zero(t, gotPath)
				assert.NotZero(t, gotObjectType)
				assert.NotNil(t, gotStore)
			}
		})
	}
}

func TestParseServiceInput(t *testing.T) {
	cases := []struct {
		name         string
		servicerType typ.ServicerType
		hasPanic     bool
		err          error
	}{
		{
			"invalid",
			"invalid",
			true,
			nil,
		},
		{
			"valid",
			qingstor.ServicerType,
			false,
			nil,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			if v.hasPanic {
				assert.Panics(t, func() {
					_, _ = ParseServiceInput(v.servicerType)
				})
				return
			}

			gotStore, gotErr := ParseServiceInput(v.servicerType)
			assert.Equal(t, v.err == nil, gotErr == nil)
			if v.err == nil {
				assert.NotNil(t, gotStore)
			}
		})
	}
}
