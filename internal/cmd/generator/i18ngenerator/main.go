package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
)

var funcs = template.FuncMap{
	"funcName": func(lang string) string {
		lang = strings.ReplaceAll(lang, "_", "")
		lang = strings.ToUpper(lang[:1]) + lang[1:]
		return lang
	},
}

func main() {
	const translationPath = "../../translations"

	fi, err := ioutil.ReadDir(translationPath)
	if err != nil {
		log.Fatal(err)
	}

	goFile, err := os.Create("../../pkg/i18n/generated.go")
	if err != nil {
		log.Fatal(err)
	}
	data := make(map[string]*map[string]string)

	for _, v := range fi {
		if !v.IsDir() {
			continue
		}

		dataFiles, err := ioutil.ReadDir(path.Join(translationPath, v.Name()))
		if err != nil {
			log.Fatal(err)
		}
		data[v.Name()] = new(map[string]string)

		for _, file := range dataFiles {
			content, err := ioutil.ReadFile(path.Join(translationPath, v.Name(), file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(content, data[v.Name()])
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	err = i18nTmpl.Execute(goFile, struct {
		Data      map[string]*map[string]string
		BackQuote string
	}{
		data,
		"`",
	})
	if err != nil {
		log.Fatal(err)
	}
}

var i18nTmpl = template.Must(template.New("i18n").Funcs(funcs).Parse(`package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Init will init i18n support via input language.
func Init(lang string) {
	switch lang {
{{- range $k, $v := .Data }}
	case "{{$k}}":
		init{{ funcName $k }}()
{{- end }}
	default:
		initEnUS()
	}
}

{{- range $k, $v := .Data }}
// init{{ funcName $k }} will init {{ $k }} support.
func init{{ funcName $k }}() {
	languageTag := language.MustParse("{{ $k }}")
	{{- range $k, $v := $v }}
	_ = message.SetString(languageTag, {{$.BackQuote}}{{$k}}{{$.BackQuote}}, {{$.BackQuote}}{{$v}}{{$.BackQuote}})
{{- end }}
}
{{- end }}
`))
