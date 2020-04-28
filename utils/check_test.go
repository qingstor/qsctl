package utils

import (
	"testing"

	"bou.ke/monkey"
	"github.com/AlecAivazis/survey/v2"
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
			defer monkey.UnpatchAll()
			monkey.Patch(survey.AskOne, func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
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
			s := Survey{
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
