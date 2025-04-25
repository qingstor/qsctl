package utils

import (
	"regexp"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/agiledragon/gomonkey/v2"
)

func TestSurvey_DoubleCheckString(t *testing.T) {
	type fields struct {
		msg    string
		expect string
	}
	tests := []struct {
		name    string
		fields  fields
		input   string
		want    bool
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				msg:    "test msg",
				expect: "abc",
			},
			input:   "abc",
			want:    true,
			wantErr: false,
		},
		{
			name: "not match",
			fields: fields{
				msg:    "test msg",
				expect: "abcd",
			},
			input:   "abc",
			want:    false,
			wantErr: false,
		},
		{
			name: "ask error",
			fields: fields{
				msg:    "test msg",
				expect: "abc",
			},
			input:   "abc",
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()
			patches.ApplyFunc(survey.AskOne, func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
				if tt.wantErr {
					return errTmp
				}
				res, ok := response.(*string)
				if !ok {
					t.Fatal("param invalid")
				}
				*res = tt.input
				return nil
			})
			s := InputCheck{
				msg:    tt.fields.msg,
				expect: tt.fields.expect,
			}

			got, err := s.DoubleCheckString()
			if (err != nil) != tt.wantErr {
				t.Errorf("DoubleCheckString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DoubleCheckString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfirmCheck_CheckConfirm(t *testing.T) {
	yesRx := regexp.MustCompile("^(?i:y(?:es)?)$")
	type fields struct {
		msg string
	}
	tests := []struct {
		name    string
		fields  fields
		input   string
		want    bool
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				msg: "test msg",
			},
			input:   "y",
			want:    true,
			wantErr: false,
		},
		{
			name: "not match n",
			fields: fields{
				msg: "test msg",
			},
			input:   "n",
			want:    false,
			wantErr: false,
		},
		{
			name: "not match other",
			fields: fields{
				msg: "test msg",
			},
			input:   "a",
			want:    false,
			wantErr: false,
		},
		{
			name: "ask error",
			fields: fields{
				msg: "test msg",
			},
			input:   "abc",
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()
			patches.ApplyFunc(survey.AskOne, func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
				if tt.wantErr {
					return errTmp
				}
				res, ok := response.(*bool)
				if !ok {
					t.Fatal("param invalid")
				}
				*res = yesRx.Match([]byte(tt.input))
				return nil
			})
			c := ConfirmCheck{
				msg: tt.fields.msg,
			}
			got, err := c.CheckConfirm()
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckConfirm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckConfirm() got = %v, want %v", got, tt.want)
			}
		})
	}
}
