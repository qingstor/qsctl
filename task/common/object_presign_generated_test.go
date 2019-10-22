// Code generated by go generate; DO NOT EDIT.
package common

import (
	"errors"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/utils"
)

var _ navvy.Pool
var _ types.Pool
var _ = utils.SubmitNextTask

func TestNewObjectPresignTask(t *testing.T) {
	m := &mockObjectPresignTask{}
	task := NewObjectPresignTask(m)
	assert.NotNil(t, task)
}

func TestObjectPresignTask_TriggerFault(t *testing.T) {
	m := &mockObjectPresignTask{}
	task := &ObjectPresignTask{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.objectPresignTaskRequirement.ValidateFault())
}
