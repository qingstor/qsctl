package utils

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	log "github.com/sirupsen/logrus"
)

// StrDoubleChecker double check a string and return satisfied or not
type StrDoubleChecker interface {
	DoubleCheckString() (bool, error)
}

// Survey implements StrDoubleChecker
type Survey struct {
	msg    string
	expect string
}

// DoubleCheckString is the public func for string double check
func DoubleCheckString(expect, msg string) (bool, error) {
	sv := newSurvey(withMsg(msg), withExpect(expect))
	return sv.DoubleCheckString()
}

type surveyOptFn func(*Survey)

// newSurvey return a Survey struct
func newSurvey(opts ...surveyOptFn) Survey {
	sv := Survey{}
	for _, fn := range opts {
		fn(&sv)
	}
	return sv
}

func withMsg(msg string) surveyOptFn {
	return func(s *Survey) {
		s.msg = msg
	}
}

func withExpect(expect string) surveyOptFn {
	return func(s *Survey) {
		s.expect = expect
	}
}

// DoubleCheckString implements StrDoubleChecker.DoubleCheckString()
func (s Survey) DoubleCheckString() (bool, error) {
	name := ""
	prompt := &survey.Input{
		Message: s.msg,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		if err == terminal.InterruptErr {
			log.Debug("interrupted")
			os.Exit(0)
		}
		return false, err
	}

	return name == s.expect, nil
}
