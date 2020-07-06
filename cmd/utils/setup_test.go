package utils

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestCredentialConfig_String(t *testing.T) {
	type fields struct {
		AccessKeyID     string
		SecretAccessKey string
		Host            string
		Port            string
		Protocol        string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "normal",
			fields: fields{
				AccessKeyID:     "access-key",
				SecretAccessKey: "secret-key",
				Host:            "qingstor.com",
				Port:            "443",
				Protocol:        "https",
			},
			want: "ak: access-key, sk: secret-key, url: https://qingstor.com:443",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CredentialConfig{
				AccessKeyID:     tt.fields.AccessKeyID,
				SecretAccessKey: tt.fields.SecretAccessKey,
				Host:            tt.fields.Host,
				Port:            tt.fields.Port,
				Protocol:        tt.fields.Protocol,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newCredentialConfig(t *testing.T) {
	option := struct {
		ak       string
		sk       string
		host     string
		port     string
		protocol string

		inputers   []StringInputer
		checker    InputChecker
		writerFunc func() error
	}{
		ak:       uuid.New().String(),
		sk:       uuid.New().String(),
		host:     uuid.New().String(),
		port:     uuid.New().String(),
		protocol: uuid.New().String(),

		inputers: []StringInputer{},
		checker: InputChecker{
			Confirm: mockBoolInputer{},
			Test:    mockBoolInputer{},
			TestFunc: func() error {
				return nil
			},
		},
		writerFunc: func() error {
			return nil
		},
	}

	type args struct {
		opts []ConfigOption
	}
	tests := []struct {
		name  string
		args  args
		check func(CredentialConfig) bool
	}{
		{
			name: "blank opt",
			args: args{},
			check: func(c CredentialConfig) bool {
				defaultConfig := DefaultCredentialConfig()
				return c.String() == defaultConfig.String()
			},
		},
		{
			name: "with ak",
			args: args{opts: []ConfigOption{withAccessKeyID(option.ak)}},
			check: func(c CredentialConfig) bool {
				return option.ak == c.AccessKeyID
			},
		},
		{
			name: "all",
			args: args{opts: []ConfigOption{
				withAccessKeyID(option.ak),
				withSecretAccessKey(option.sk),
				withHost(option.host),
				withPort(option.port),
				withProtocol(option.protocol),
				withInputers(option.inputers...),
				withConfirmChecker(option.checker.Confirm),
				withTestChecker(option.checker.Test),
				withTestFunc(option.checker.TestFunc),
				withWriteFunc(option.writerFunc),
			}},
			check: func(c CredentialConfig) bool {
				return option.ak == c.AccessKeyID &&
					option.sk == c.SecretAccessKey &&
					option.host == c.Host &&
					option.port == c.Port &&
					option.protocol == c.Protocol &&
					len(option.inputers) == len(c.inputers) &&
					reflect.DeepEqual(option.checker.Test, c.checker.Test) &&
					reflect.DeepEqual(option.checker.Confirm, c.checker.Confirm) &&
					option.checker.TestFunc() == c.checker.TestFunc() &&
					option.writerFunc() == c.writerFunc()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCredentialConfig(tt.args.opts...); !tt.check(got) {
				t.Errorf("newCredentialConfig() = %v, check failed", got)
			}
		})
	}
}

type mockStringInputer struct {
	name  string
	key   *string
	value string
	err   error
}

func (m mockStringInputer) SetStringP(s *string) error {
	if m.err != nil {
		return m.err
	}
	*s = m.value
	return nil
}

func (m mockStringInputer) GetStringP() *string {
	return m.key
}

func (m mockStringInputer) GetName() string {
	return m.name
}

type mockBoolInputer struct {
	name  string
	value bool
	err   error
}

func (m mockBoolInputer) GetBool() (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.value, nil
}

func (m mockBoolInputer) GetName() string {
	return m.name
}

func Test_setupConfig(t *testing.T) {
	tmpErr := errors.New("test error")
	ak, sk := uuid.New().String(), uuid.New().String()
	config := DefaultCredentialConfig()

	type args struct {
		inputers   []StringInputer
		checker    InputChecker
		writerFunc func() error
	}
	tests := []struct {
		name      string
		args      args
		checkFunc func(c CredentialConfig) bool
		err       error
	}{
		{
			name: "all without error",
			args: args{
				inputers: []StringInputer{
					mockStringInputer{
						name:  "ak",
						key:   &config.AccessKeyID,
						value: ak,
					},
					mockStringInputer{
						name:  "sk",
						key:   &config.SecretAccessKey,
						value: sk,
					},
				},
				checker: InputChecker{
					Confirm: mockBoolInputer{
						name:  "confirm",
						value: true,
					},
					Test: mockBoolInputer{
						name:  "test",
						value: true,
					},
					TestFunc: func() error {
						return nil
					},
				},
				writerFunc: func() error {
					return nil
				},
			},
			checkFunc: func(c CredentialConfig) bool {
				return c.AccessKeyID == ak && c.SecretAccessKey == sk
			},
		},
		{
			name: "inputer failed",
			args: args{
				inputers: []StringInputer{
					mockStringInputer{
						name:  "ak",
						key:   &config.AccessKeyID,
						value: ak,
						err:   tmpErr,
					},
				},
			},
			err: tmpErr,
		},
		{
			name: "test check not set, and confirm not set, no error",
			args: args{
				inputers: []StringInputer{
					mockStringInputer{
						name:  "ak",
						key:   &config.AccessKeyID,
						value: ak,
					},
					mockStringInputer{
						name:  "sk",
						key:   &config.SecretAccessKey,
						value: sk,
					},
				},
				writerFunc: func() error {
					return nil
				},
			},
			checkFunc: func(c CredentialConfig) bool {
				return c.AccessKeyID == ak && c.SecretAccessKey == sk
			},
		},
		{
			name: "test check not set, and confirm with error",
			args: args{
				checker: InputChecker{
					Confirm: mockBoolInputer{
						name:  "confirm",
						value: true,
						err:   tmpErr,
					},
				},
			},
			err: tmpErr,
		},
		{
			name: "test check not set, and confirm, write with error",
			args: args{
				checker: InputChecker{
					Confirm: mockBoolInputer{
						name:  "confirm",
						value: true,
					},
				},
				writerFunc: func() error {
					return tmpErr
				},
			},
			err: tmpErr,
		},
		{
			name: "test check set, get test with error",
			args: args{
				checker: InputChecker{
					Test: mockBoolInputer{
						name:  "test",
						value: false,
						err:   tmpErr,
					},
				},
			},
			err: tmpErr,
		},
		{
			name: "test check set, get test but test func not set",
			args: args{
				checker: InputChecker{
					Test: mockBoolInputer{
						name:  "test",
						value: true,
					},
				},
				writerFunc: func() error {
					return nil
				},
			},
			checkFunc: func(c CredentialConfig) bool {
				return c.AccessKeyID == "" && c.SecretAccessKey == ""
			},
		},
		{
			name: "test check set, get test, test func get error",
			args: args{
				checker: InputChecker{
					Test: mockBoolInputer{
						name:  "test",
						value: true,
					},
					TestFunc: func() error {
						return tmpErr
					},
				},
				writerFunc: func() error {
					return nil
				},
			},
			checkFunc: func(c CredentialConfig) bool {
				return c.AccessKeyID == "" && c.SecretAccessKey == ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config = DefaultCredentialConfig()
			config.inputers = tt.args.inputers
			config.checker = tt.args.checker
			config.writerFunc = tt.args.writerFunc
			got := config.setupConfig()
			if tt.err == nil {
				if got != nil {
					t.Errorf("setupConfig want not error, but got %v", got)
					return
				}
				if !tt.checkFunc(config) {
					t.Errorf("config check failed, got %v", config)
				}
				return
			}
			if got == nil {
				t.Errorf("setupConfig want error, but got nil")
				return
			}
			if !errors.Is(got, tt.err) {
				t.Errorf("setupConfig want err %v, but got %v", tt.err, got)
				return
			}
		})
	}
}
