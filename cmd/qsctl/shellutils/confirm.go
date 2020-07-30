package shellutils

import "regexp"

var yesRx = regexp.MustCompile("^(?i:y(?:es)?)$")

// CheckYes indicate whether input is yes regexp
func CheckYes(input string) bool {
	return yesRx.MatchString(input)
}
