package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qingstor/qsctl/v2/constants"
)

func Test_multipartFlags_parse(t *testing.T) {
	type fields struct {
		partThresholdStr string
		partThreshold    int64
		partSizeStr      string
		partSize         int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "all valid",
			fields: fields{
				partThresholdStr: "1M",
				partThreshold:    1024 * 1024,
				partSizeStr:      "1K",
				partSize:         1024,
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
				partThresholdStr: "",
				partThreshold:    constants.MaximumAutoMultipartSize,
				partSizeStr:      "",
				partSize:         0,
			},
			wantErr: false,
		},
		{
			name: "invalid part size",
			fields: fields{
				partThresholdStr: "",
				partThreshold:    0,
				partSizeStr:      "1x",
				partSize:         0,
			},
			wantErr: true,
		},
		{
			name: "invalid part threshold",
			fields: fields{
				partThresholdStr: "1x",
				partThreshold:    0,
				partSizeStr:      "1K",
				partSize:         1024,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &multipartFlags{
				partThresholdStr: tt.fields.partThresholdStr,
				partSizeStr:      tt.fields.partSizeStr,
			}
			err := f.parse()
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.fields.partSize, f.partSize)
			assert.Equal(t, tt.fields.partThreshold, f.partThreshold)
		})
	}
}

func Test_inExcludeFlags_parse(t *testing.T) {
	validRegxStr, invalidRegxStr := "abc", "(abc"
	validRegx := regexp.MustCompile(validRegxStr)

	type fields struct {
		includeRegxStr string
		includeRegx    *regexp.Regexp
		excludeRegxStr string
		excludeRegx    *regexp.Regexp
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "all valid",
			fields: fields{
				includeRegxStr: validRegxStr,
				includeRegx:    validRegx,
				excludeRegxStr: validRegxStr,
				excludeRegx:    validRegx,
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
				includeRegxStr: "",
				includeRegx:    nil,
				excludeRegxStr: "",
				excludeRegx:    nil,
			},
			wantErr: false,
		},
		{
			name: "invalid exclude",
			fields: fields{
				includeRegxStr: validRegxStr,
				includeRegx:    nil,
				excludeRegxStr: invalidRegxStr,
				excludeRegx:    nil,
			},
			wantErr: true,
		},
		{
			name: "invalid include",
			fields: fields{
				includeRegxStr: invalidRegxStr,
				includeRegx:    nil,
				excludeRegxStr: validRegxStr,
				excludeRegx:    validRegx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &inExcludeFlags{
				excludeRegxStr: tt.fields.excludeRegxStr,
				includeRegxStr: tt.fields.includeRegxStr,
			}
			err := f.parse()
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.fields.excludeRegx, f.excludeRegx)
			assert.Equal(t, tt.fields.includeRegx, f.includeRegx)
		})
	}
}
