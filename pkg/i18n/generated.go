package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// initEnUS will init en_US support.
func initEnUS(tag language.Tag) {
	_ = message.SetString(tag, "Key <%s> copied.\n", "Key <%s> copied.\n")
}

// initZhCN will init zh_CN support.
func initZhCN(tag language.Tag) {
	_ = message.SetString(tag, "Key <%s> copied.\n", "文件 <%s> 已复制\n")
}
