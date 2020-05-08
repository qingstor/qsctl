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

// ConfirmChecker try to confirm return confirmed or not
type ConfirmChecker interface {
	CheckConfirm() (bool, error)
}

// DoubleCheckString is the public func for string double check
func DoubleCheckString(expect, msg string) (bool, error) {
	sc := newInputCheck(withInputMsg(msg), withInputExpect(expect))
	return sc.DoubleCheckString()
}

// CheckConfirm is the public func for confirm check
func CheckConfirm(msg string) (bool, error) {
	cc := newConfirmCheck(withConfirmMsg(msg))
	return cc.CheckConfirm()
}

// InputCheck implements StrDoubleChecker
type InputCheck struct {
	msg    string
	expect string
}

type inputCheckOptFn func(*InputCheck)

// newInputCheck return a InputCheck struct
func newInputCheck(opts ...inputCheckOptFn) InputCheck {
	sc := InputCheck{}
	for _, fn := range opts {
		fn(&sc)
	}
	return sc
}

func withInputMsg(msg string) inputCheckOptFn {
	return func(s *InputCheck) {
		s.msg = msg
	}
}

func withInputExpect(expect string) inputCheckOptFn {
	return func(s *InputCheck) {
		s.expect = expect
	}
}

// DoubleCheckString implements StrDoubleChecker.DoubleCheckString()
func (s InputCheck) DoubleCheckString() (bool, error) {
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

// ConfirmCheck implements ConfirmChecker
type ConfirmCheck struct {
	msg string
}

type confirmCheckOptFn func(*ConfirmCheck)

// newConfirmCheck return a ConfirmCheck struct
func newConfirmCheck(opts ...confirmCheckOptFn) ConfirmCheck {
	cc := ConfirmCheck{}
	for _, fn := range opts {
		fn(&cc)
	}
	return cc
}

func withConfirmMsg(msg string) confirmCheckOptFn {
	return func(s *ConfirmCheck) {
		s.msg = msg
	}
}

// CheckConfirm implements ConfirmChecker.CheckConfirm
func (c ConfirmCheck) CheckConfirm() (bool, error) {
	var ok bool
	prompt := &survey.Confirm{
		Message: c.msg,
	}
	err := survey.AskOne(prompt, &ok)
	if err != nil {
		if err == terminal.InterruptErr {
			log.Debug("interrupted")
			os.Exit(0)
		}
		return false, err
	}

	return ok, nil
}
