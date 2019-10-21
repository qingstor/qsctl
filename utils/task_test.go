package utils

import (
	"testing"

	"github.com/Xuanwo/storage/types"
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

func TestParseKey(t *testing.T) {
	cases := []struct {
		input              string
		expectedKeyType    types.ObjectType
		expectedBucketName string
		expectedKey        string
	}{
		{"qs://xxxx-bucket/abc", types.ObjectTypeFile, "xxxx-bucket", "abc"},
		{"qs://abcdef", types.ObjectTypeDir, "abcdef", ""},
		{"qs://abcdef/", types.ObjectTypeDir, "abcdef", ""},
		{"qs://abcdef/def/ghi", types.ObjectTypeFile, "abcdef", "def/ghi"},
		{"qs://abcdef/def/ghi/", types.ObjectTypeDir, "abcdef", "def/ghi/"},
		{"abcdef", types.ObjectTypeDir, "abcdef", ""},
		{"abcdef/", types.ObjectTypeDir, "abcdef", ""},
		{"abcdef/def/ghi", types.ObjectTypeFile, "abcdef", "def/ghi"},
		{"abcdef/👾 🙇 💁 🙅 🙆 🙋 🙎 🙍", types.ObjectTypeFile, "abcdef", "👾 🙇 💁 🙅 🙆 🙋 🙎 🙍"},
	}

	for k, v := range cases {
		actualKeyType, actualBucketName, actualKey, err := ParseQsPath(v.input)
		assert.Equal(t, v.expectedKeyType, actualKeyType, k)
		assert.Equal(t, v.expectedBucketName, actualBucketName, k)
		assert.Equal(t, v.expectedKey, actualKey, k)
		assert.NoError(t, err, k)
	}
}
