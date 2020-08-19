package taskutils

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_truncateBefore(t *testing.T) {
	type args struct {
		s string
		l int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not longer",
			args: args{
				s: "abcd",
				l: 10,
			},
			want: "abcd",
		},
		{
			name: "equal",
			args: args{
				s: "abcd",
				l: 4,
			},
			want: "abcd",
		},
		{
			name: "longer",
			args: args{
				s: "abcdefg",
				l: 5,
			},
			want: "...cdefg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncateBefore(tt.args.s, tt.args.l); got != tt.want {
				t.Errorf("truncateBefore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calBarSize(t *testing.T) {
	tests := []struct {
		name        string
		fullWid     int
		wantNameWid int
		wantBarWid  int
	}{
		{
			name:        "normal",
			fullWid:     200,
			wantNameWid: 100,
			wantBarWid:  40,
		},
		{
			name:        "half",
			fullWid:     100,
			wantNameWid: 20,
			wantBarWid:  20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNameWid, gotBarWid := calBarSize(tt.fullWid)
			if gotNameWid != tt.wantNameWid {
				t.Errorf("calBarSize() gotNameWid = %v, want %v", gotNameWid, tt.wantNameWid)
			}
			if gotBarWid != tt.wantBarWid {
				t.Errorf("calBarSize() gotBarWid = %v, want %v", gotBarWid, tt.wantBarWid)
			}
		})
	}
}

func Test_barGroup_GetPBarByID(t *testing.T) {
	handler, clearFunc := NewHandler(context.Background())
	defer clearFunc()
	bar := handler.pbPool.AddBar(time.Now().Unix())
	id := uuid.New().String()
	handler.bGroup.SetBarByID(id, bar)
	tests := []struct {
		name string
		id   string
		want bool
	}{
		{
			name: "blank",
			id:   uuid.New().String(),
			want: false,
		},
		{
			name: "not blank",
			id:   id,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handler.bGroup.GetBarByID(tt.id)
			assert.Equal(t, got == nil, !tt.want, tt.name)
		})
	}
}
