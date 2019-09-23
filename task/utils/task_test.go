package utils

import (
	"testing"

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
		{"xxxx", "", constants.FlowAtLocal},
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
		expectedKeyType    constants.KeyType
		expectedBucketName string
		expectedKey        string
	}{
		{"qs://xxxx/abc", constants.KeyTypeObject, "xxxx", "abc"},
	}

	for _, v := range cases {
		actualKeyType, actualBucketName, actualKey, err := ParseKey(v.input)
		assert.Equal(t, v.expectedKeyType, actualKeyType)
		assert.Equal(t, v.expectedBucketName, actualBucketName)
		assert.Equal(t, v.expectedKey, actualKey)
		assert.NoError(t, err)
	}
}
