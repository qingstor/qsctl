package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// initEnUS will init en_US support.
func initEnUS() {
	languageTag := language.MustParse("en_US")
	_ = message.SetString(languageTag, `Copy file: qsctl cp /path/to/file qs://prefix/a`, `Copy file: qsctl cp /path/to/file qs://prefix/a`)
	_ = message.SetString(languageTag, `Key <%s> copied.\n`, `Key <%s> copied.\n`)
}

// initZhCN will init zh_CN support.
func initZhCN() {
	languageTag := language.MustParse("zh_CN")
	_ = message.SetString(languageTag, `Key <%s> copied.
`, `文件 <%s> 已复制
`)
}
