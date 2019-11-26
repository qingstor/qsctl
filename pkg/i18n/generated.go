package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func Init(lang string) {
	switch lang {
	case "en_US":
		languageTag := language.MustParse("en_US")
		_ = message.SetString(languageTag, `Copy file: qsctl cp /path/to/file qs://prefix/a`, `Copy file: qsctl cp /path/to/file qs://prefix/a`)
		_ = message.SetString(languageTag, `Key <%s> copied.\n`, `Key <%s> copied.\n`)
	case "zh_CN":
		languageTag := language.MustParse("zh_CN")
		_ = message.SetString(languageTag, `Key <%s> copied.
`, `文件 <%s> 已复制
`)
	default:
		languageTag := language.MustParse("en_US")
		_ = message.SetString(languageTag, `Copy file: qsctl cp /path/to/file qs://prefix/a`, `Copy file: qsctl cp /path/to/file qs://prefix/a`)
		_ = message.SetString(languageTag, `Key <%s> copied.\n`, `Key <%s> copied.\n`)
	}
}
