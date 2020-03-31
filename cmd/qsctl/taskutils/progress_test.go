package taskutils

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/google/uuid"
	"github.com/qingstor/noah/pkg/progress"
	"github.com/stretchr/testify/assert"
	"github.com/vbauerster/mpb/v4"
)

func Test_pBar_GetBar(t *testing.T) {
	bar := pbPool.AddBar(time.Now().Unix())
	type fields struct {
		status pBarStatus
		bar    *mpb.Bar
	}
	tests := []struct {
		name   string
		fields fields
		want   *mpb.Bar
	}{
		{
			name: "normal",
			fields: fields{
				bar: bar,
			},
			want: bar,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer bar.Abort(true)
			b := pBar{
				status: tt.fields.status,
				bar:    tt.fields.bar,
			}
			if got := b.GetBar(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBar() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			wantNameWid: 120,
			wantBarWid:  40,
		},
		{
			name:        "half",
			fullWid:     100,
			wantNameWid: 30,
			wantBarWid:  30,
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

func Test_pBar_MarkFinished(t *testing.T) {
	tests := []struct {
		name   string
		status pBarStatus
	}{
		{
			name: "true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &pBar{
				status: tt.status,
			}
			b.MarkFinished()
			assert.True(t, b.Finished(), tt.name)
		})
	}
}

func Test_pBar_Finished(t *testing.T) {
	tests := []struct {
		name   string
		status pBarStatus
		want   bool
	}{
		{
			name:   "true",
			status: pbFinished,
			want:   true,
		},
		{
			name:   "false",
			status: pbNotExist,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := pBar{
				status: tt.status,
			}
			if got := b.Finished(); got != tt.want {
				t.Errorf("Finished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pBar_Shown(t *testing.T) {
	tests := []struct {
		name   string
		status pBarStatus
		want   bool
	}{
		{
			name:   "true",
			status: pbShown,
			want:   true,
		},
		{
			name:   "false",
			status: pbNotExist,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := pBar{
				status: tt.status,
			}
			if got := b.Shown(); got != tt.want {
				t.Errorf("Shown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pBar_NotExist(t *testing.T) {
	tests := []struct {
		name   string
		status pBarStatus
		want   bool
	}{
		{
			name:   "true",
			status: pbNotExist,
			want:   true,
		},
		{
			name:   "false",
			status: pbFinished,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := pBar{
				status: tt.status,
			}
			if got := b.NotExist(); got != tt.want {
				t.Errorf("NotExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pBarGroup_GetActiveCount(t *testing.T) {
	tests := []struct {
		name     string
		incTimes int
		decTimes int
		want     int
	}{
		{
			name:     "normal",
			incTimes: 10,
			decTimes: 5,
			want:     5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.incTimes; i++ {
				pbGroup.IncActive()
			}

			for i := 0; i < tt.decTimes; i++ {
				pbGroup.DecActive()
			}

			if got := pbGroup.GetActiveCount(); got != tt.want {
				t.Errorf("GetActiveCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pBarGroup_GetPBarByID(t *testing.T) {
	id, status := uuid.New().String(), pbFinished
	pbGroup.SetPBarByID(id, &pBar{status: status})
	tests := []struct {
		name   string
		id     string
		status pBarStatus
	}{
		{
			name:   "blank",
			id:     uuid.New().String(),
			status: pbNotExist,
		},
		{
			name:   "not blank",
			id:     id,
			status: status,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pbGroup.GetPBarByID(tt.id); !reflect.DeepEqual(got.status, tt.status) {
				t.Errorf("GetPBarByID() = %v, want %v", got.status, tt.status)
			}
		})
	}
}

func TestStartProgress(t *testing.T) {
	id := uuid.New().String()
	type args struct {
		d           time.Duration
		maxBarCount int
	}
	tests := []struct {
		name           string
		args           args
		seconds        int
		activeBarCount int
		barCount       int
	}{
		{
			name:           "3s",
			args:           args{d: time.Second, maxBarCount: 3},
			seconds:        3,
			activeBarCount: 2,
			barCount:       4,
		},
	}
	for _, tt := range tests {
		pbGroup = &pBarGroup{
			bars: make(map[string]*pBar),
		}
		wg = new(sync.WaitGroup)
		pbPool = mpb.New(mpb.WithWaitGroup(wg), mpb.WithOutput(nil))
		monkey.Patch(progress.GetStates, func() map[string]progress.State {
			return map[string]progress.State{
				id + "1": {Name: "list spinner", Status: "", Type: 0, Done: 0, Total: 1},
				id + "2": {Name: "first bar", Status: "", Type: 1, Done: 5, Total: 10},
				id + "3": {Name: "finished bar", Status: "", Type: 1, Done: 1, Total: 1},
				id + "4": {Name: "oversize not display bar", Status: "", Type: 1, Done: 3, Total: 10},
			}
		})

		time.AfterFunc(time.Duration(tt.seconds)*tt.args.d, func() {
			FinishProgress()
		})
		StartProgress(tt.args.d, tt.args.maxBarCount)
		monkey.UnpatchAll()
		assert.Equal(t, tt.activeBarCount, pbGroup.GetActiveCount(), tt.name)
		assert.Equal(t, tt.barCount, len(pbGroup.bars), tt.name)
		for _, b := range pbGroup.bars {
			if b.bar != nil {
				b.bar.Abort(true)
			}
		}
	}
}
