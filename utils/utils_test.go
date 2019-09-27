package utils

import (
	"fmt"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/stretchr/testify/assert"
	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func TestCalculatePartSize(t *testing.T) {
	cases := []struct {
		msg      string
		input    int64
		expected int64
		err      error
	}{
		{"1B", 1, constants.DefaultPartSize, nil},
		{"1G", 1024 * 1024 * 1024, constants.DefaultPartSize, nil},
		{"10G", 10 * 1024 * 1024 * 1024, constants.DefaultPartSize, nil},
		{"2TB", 2 * 1024 * 1024 * 1024 * 1024, constants.DefaultPartSize << 1, nil},
		{"10TB", 10 * 1024 * 1024 * 1024 * 1024, 1099511628, nil},
		{"100TB", 101 * 1024 * 1024 * 1024 * 1024, 0, fmt.Errorf("calculate part size failed: {%w}", fault.NewLocalFileTooLarge(nil, 101*1024*1024*1024*1024))},
	}

	for _, v := range cases {
		x, err := CalculatePartSize(v.input)
		assert.Equal(t, v.err, err, v.msg)
		assert.Equal(t, v.expected, x, v.msg)
	}
}

func TestSubmitNextTask(t *testing.T) {
	task := &EmptyTask{
		Todo:  types.Todo{},
		Fault: types.Fault{},
	}
	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}
	task.SetPool(pool)

	assert.NotPanics(t, func() {
		SubmitNextTask(task)
	})
}
