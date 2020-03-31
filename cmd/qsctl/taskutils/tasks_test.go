package taskutils

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// random an int more than 0 and less than 20 to restrict the memory use.
func randIntP() int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return int(r.Int31()%100 + 1)
}

func TestNewAtServiceTask(t *testing.T) {
	tests := []struct {
		name      string
		size      int
		wantPanic bool
	}{
		{
			name:      "normal",
			size:      randIntP(),
			wantPanic: false,
		},
		{
			name:      "zero size",
			size:      0,
			wantPanic: true,
		},
		{
			name:      "negative size",
			size:      -randIntP(),
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() {
					NewAtServiceTask(tt.size)
				}, tt.name)
			} else {
				got := NewAtServiceTask(tt.size)
				assert.Equal(t, tt.size, got.GetPool().Cap(), tt.name)
				assert.False(t, got.GetFault().HasError(), tt.name)
			}
		})
	}
}

func TestNewAtStorageTask(t *testing.T) {
	tests := []struct {
		name      string
		size      int
		wantPanic bool
	}{
		{
			name:      "normal",
			size:      randIntP(),
			wantPanic: false,
		},
		{
			name:      "zero size",
			size:      0,
			wantPanic: true,
		},
		{
			name:      "negative size",
			size:      -randIntP(),
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() {
					NewAtStorageTask(tt.size)
				}, tt.name)
			} else {
				got := NewAtStorageTask(tt.size)
				assert.Equal(t, tt.size, got.GetPool().Cap(), tt.name)
				assert.False(t, got.GetFault().HasError(), tt.name)
			}
		})
	}
}

func TestNewBetweenStorageTask(t *testing.T) {
	tests := []struct {
		name      string
		size      int
		wantPanic bool
	}{
		{
			name:      "normal",
			size:      randIntP(),
			wantPanic: false,
		},
		{
			name:      "zero size",
			size:      0,
			wantPanic: true,
		},
		{
			name:      "negative size",
			size:      -randIntP(),
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() {
					NewBetweenStorageTask(tt.size)
				}, tt.name)
			} else {
				got := NewBetweenStorageTask(tt.size)
				assert.Equal(t, tt.size, got.GetPool().Cap(), tt.name)
				assert.False(t, got.GetFault().HasError(), tt.name)
			}
		})
	}
}
