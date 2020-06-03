package i18n

import (
	"io"

	"github.com/Xuanwo/go-locale"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var p *message.Printer

func newMatcher(t []language.Tag) *matcher {
	tags := &matcher{make(map[language.Tag]int)}
	for i, tag := range t {
		ct, err := language.All.Canonicalize(tag)
		if err != nil {
			ct = tag
		}
		tags.index[ct] = i
	}
	return tags
}

type matcher struct {
	index map[language.Tag]int
}

func (m matcher) Match(want ...language.Tag) (language.Tag, int, language.Confidence) {
	for _, t := range want {
		ct, err := language.All.Canonicalize(t)
		if err != nil {
			ct = t
		}
		conf := language.Exact
		for {
			if index, ok := m.index[ct]; ok {
				return ct, index, conf
			}
			if ct == language.Und {
				break
			}
			ct = ct.Parent()
			conf = language.High
		}
	}
	return language.Und, 0, language.No
}

var supported = newMatcher([]language.Tag{
	language.AmericanEnglish,
	language.English,
	language.SimplifiedChinese,
	language.Chinese,
})

// Init will init i18n support via input language.
func Init(lang language.Tag) {
	tag, _, _ := supported.Match(lang)
	switch tag {
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
	tag, err := locale.Detect()
	if err != nil {
		log.Errorf("locale detect failed: [%v], use en_us instead", err)
		tag = language.AmericanEnglish
	}
	Init(tag)
	p = message.NewPrinter(tag)
}
