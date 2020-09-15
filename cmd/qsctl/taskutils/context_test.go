package taskutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vbauerster/mpb/v5"
)

func TestContextWithHandler(t *testing.T) {
	h := &PbHandler{
		pbPool: mpb.New(),
	}
	tests := []struct {
		name    string
		ctx     context.Context
		handler *PbHandler
		res     *PbHandler
	}{
		{
			name:    "nil ctx with nil handler",
			ctx:     nil,
			handler: nil,
			res:     nil,
		},
		{
			name:    "nil ctx with non-nil handler",
			ctx:     nil,
			handler: h,
			res:     h,
		},
		{
			name:    "bg ctx with nil handler",
			ctx:     context.Background(),
			handler: nil,
			res:     nil,
		},
		{
			name:    "bg ctx with non-nil handler",
			ctx:     context.Background(),
			handler: h,
			res:     h,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextWithHandler(tt.ctx, tt.handler)
			assert.Equal(t, tt.res, HandlerFromContext(ctx))
		})
	}
}
