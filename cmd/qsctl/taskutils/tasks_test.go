package taskutils

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func randIntP() int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for {
		tmp := r.Int31()
		if tmp > 0 {
			return int(tmp)
		}
	}
}

func TestNewAtServiceTask(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "normal",
			size: randIntP(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAtServiceTask(tt.size)
			assert.Equal(t, tt.size, got.GetPool().Cap(), tt.name)
			assert.False(t, got.GetFault().HasError(), tt.name)
		})
	}
}

func TestNewAtStorageTask(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "normal",
			size: randIntP(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAtStorageTask(tt.size)
			assert.Equal(t, tt.size, got.GetPool().Cap(), tt.name)
			assert.False(t, got.GetFault().HasError(), tt.name)
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
			name: "normal",
			size: randIntP(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBetweenStorageTask(tt.size)
			assert.Equal(t, tt.size, got.GetPool().Cap(), tt.name)
			assert.False(t, got.GetFault().HasError(), tt.name)
		})
	}
}
