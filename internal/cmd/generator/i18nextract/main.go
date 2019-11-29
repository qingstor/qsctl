package main

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
)

func main() {
	data := make(map[string]string)

	fi, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range fi {
		if v.IsDir() {
			continue
		}

		content, err := ioutil.ReadFile(v.Name())
		if err != nil {
			log.Fatal(err)
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, v.Name(), string(content), 0)
		if err != nil {
			log.Fatal(err)
		}

		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			fn, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			pack, ok := fn.X.(*ast.Ident)
			if !ok {
				return true
			}
			if pack.Name != "i18n" {
				return true
			}
			if len(call.Args) == 0 {
				return true
			}
			str, ok := call.Args[0].(*ast.BasicLit)
			if !ok {
				return true
			}

			// Keep this for later debug usage.
			// log.Printf("%v", str.Value)
			data[str.Value] = str.Value
			return true
		})
	}

	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("../../translations/en_US/data.json", content, 0664)
	if err != nil {
		log.Fatal(err)
	}
}
