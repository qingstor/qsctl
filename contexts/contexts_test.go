package contexts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCmdCtx(t *testing.T) {
	NewCmdCtx()
}

func TestCmdCtx_Set(t *testing.T) {
	c := NewCmdCtx()
	cases := []struct {
		k    interface{}
		v    interface{}
		want interface{}
		msg  string
	}{
		{"firstKey", "firstValue", "firstValue", "first"},
		{"secondKey", "secondValue", "secondValue", "second"},
		{"secondKey", "thirdValue", "thirdValue", "third"},
	}
	for _, ca := range cases {
		c = SetContext(c, ca.k, ca.v)
		assert.Equal(t, ca.want, FromContext(c, ca.k), ca.msg)
	}
}
