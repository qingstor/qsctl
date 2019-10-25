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

func TestNewObjectStatTask(t *testing.T) {
	m := &mockObjectStatTask{}
	task := NewObjectStatTask(m)
	assert.NotNil(t, task)
}

func TestObjectStatTask_TriggerFault(t *testing.T) {
	m := &mockObjectStatTask{}
	task := &ObjectStatTask{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.objectStatTaskRequirement.ValidateFault())
}

func TestMockObjectStatTask_Run(t *testing.T) {
	task := &mockObjectStatTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}