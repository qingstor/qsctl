package i18n

import (
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func Test(t *testing.T) {
	message.SetString(language.Chinese, `Key <%s> copied.\n`, `文件 <%s> 已复制\n`)
	p := message.NewPrinter(language.SimplifiedChinese)
	t.Log(p.Sprintf(`Key <%s> copied.\n`, "test"))
}
