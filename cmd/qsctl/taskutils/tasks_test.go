package taskutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAtServiceTask(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "normal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewAtServiceTask()
		})
	}
}

func TestNewAtStorageTask(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "normal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAtStorageTask()
			assert.False(t, got.ValidatePath(), tt.name)
			assert.False(t, got.ValidateStorage(), tt.name)
			assert.False(t, got.ValidateType(), tt.name)
		})
	}
}

func TestNewBetweenStorageTask(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "normal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBetweenStorageTask()
			assert.False(t, got.ValidateDestinationPath(), tt.name)
			assert.False(t, got.ValidateDestinationStorage(), tt.name)
			assert.False(t, got.ValidateDestinationType(), tt.name)
			assert.False(t, got.ValidateSourcePath(), tt.name)
			assert.False(t, got.ValidateSourceStorage(), tt.name)
			assert.False(t, got.ValidateSourceType(), tt.name)
		})
	}
}
