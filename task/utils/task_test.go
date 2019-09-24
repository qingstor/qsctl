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
		{"qs://xxxx-bucket/abc", constants.KeyTypeObject, "xxxx-bucket", "abc"},
	}

	for _, v := range cases {
		actualKeyType, actualBucketName, actualKey, err := ParseKey(v.input)
		assert.Equal(t, v.expectedKeyType, actualKeyType)
		assert.Equal(t, v.expectedBucketName, actualBucketName)
		assert.Equal(t, v.expectedKey, actualKey)
		assert.NoError(t, err)
	}
}

func TestIsValidBucketName(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"start with letter", args{"a-bucket-test"}, true},
		{"start with digit", args{"0-bucket-test"}, true},
		{"start with strike", args{"-bucket-test"}, false},
		{"end with strike", args{"bucket-test-"}, false},
		{"too short", args{"abcd"}, false},
		{"too long (64)", args{"abcdefghijklmnopqrstuvwxyz123456abcdefghijklmnopqrstuvwxyz123456"}, false},
		{"contains illegal char", args{"abcdefg_1234"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidBucketName(tt.args.s); got != tt.want {
				t.Errorf("IsValidBucketName() = %v, want %v", got, tt.want)
			}
		})
	}
}
