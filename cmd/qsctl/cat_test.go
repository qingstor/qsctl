package main

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/google/uuid"
	"github.com/qingstor/noah/pkg/schedule"
	"github.com/qingstor/noah/pkg/types"
	"github.com/qingstor/noah/task"
	"github.com/stretchr/testify/assert"

	"github.com/qingstor/qsctl/v2/utils"
)

func TestCatRun(t *testing.T) {
	tmpErr := errors.New("temp err")
	cases := []struct {
		name     string
		input    string
		parseErr error
		runErr   error
	}{
		{
			name:     "normal",
			input:    uuid.New().String(),
			parseErr: nil,
			runErr:   nil,
		},
		{
			name:     "parse err",
			input:    uuid.New().String(),
			parseErr: tmpErr,
			runErr:   tmpErr,
		},
		{
			name:     "run err",
			input:    uuid.New().String(),
			parseErr: nil,
			runErr:   tmpErr,
		},
	}

	for _, tt := range cases {
		monkey.Patch(utils.ParseBetweenStorageInput, func(_ interface {
			types.SourcePathSetter
			types.SourceStorageSetter
			types.SourceTypeSetter
			types.DestinationPathSetter
			types.DestinationStorageSetter
			types.DestinationTypeSetter
		}, src, dst string) (_, _ string, err error) {
			assert.Equal(t, tt.input, src, tt.name)
			assert.Equal(t, "-", dst, tt.name)
			if tt.parseErr != nil {
				err = tmpErr
			}
			return
		})
		var tk *task.CopyFileTask
		monkey.PatchInstanceMethod(reflect.TypeOf(tk), "Run", func(task *task.CopyFileTask, ctx context.Context) {
			if tt.runErr != nil {
				task.TriggerFault(tmpErr)
			} else {
				task.SetScheduler(schedule.New())
			}
		})
		gotErr := catRun(CatCommand, []string{tt.input})
		if tt.runErr == nil {
			assert.Nil(t, gotErr, tt.name)
		} else {
			assert.True(t, errors.Is(gotErr, tt.runErr), tt.name)
		}
		monkey.UnpatchAll()
	}
}
