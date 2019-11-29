package i18n

import (
	"io"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var p *message.Printer

// Setup will set the global printer.
func Setup(lang language.Tag) {
	Init(lang)
	p = message.NewPrinter(lang)
}

// Init will init i18n support via input language.
func Init(lang language.Tag) {
	switch lang {
	case language.AmericanEnglish, language.English:
		initEnUS(lang)
	case language.SimplifiedChinese, language.Chinese:
		initZhCN(lang)
	default:
		initEnUS(lang)
	}
}

// Fprintf is like fmt.Fprintf, but using language-specific formatting.
func Fprintf(w io.Writer, key message.Reference, a ...interface{}) (n int, err error) {
	return p.Fprintf(w, key, a...)
}

// Printf is like fmt.Printf, but using language-specific formatting.
func Printf(format string, a ...interface{}) {
	_, _ = p.Printf(format, a...)
}

// Sprintf formats according to a format specifier and returns the resulting string.
func Sprintf(format string, a ...interface{}) string {
	return p.Sprintf(format, a...)
}

// Sprint is like fmt.Sprint, but using language-specific formatting.
func Sprint(a ...interface{}) string {
	return p.Sprint(a...)
}

func init() {
	p = message.NewPrinter(language.English)
}
