// The following directive is necessary to make the package coherent:

// +build ignore

// This program generates types, It can be invoked by running
// go generate
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

//go:generate go run types_gen.go
func main() {
	const filePath = "types.go"

	data, err := ioutil.ReadFile("types.json")
	if err != nil {
		log.Fatal(err)
	}
	types := make(map[string]string)
	err = json.Unmarshal(data, &types)
	if err != nil {
		log.Fatal(err)
	}

	// Format input tasks.json
	data, err = json.MarshalIndent(types, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("types.json", data, 0664)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = tmpl.Execute(f, struct {
		Data map[string]string
	}{
		types,
	})
	if err != nil {
		log.Fatal(err)
	}
}

var tmpl = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
package types

import (
	"bytes"
	"sync"
	"io"

	"github.com/Xuanwo/navvy"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/types/storage"
)

{{- range $k, $v := .Data }}

type {{$k}}Getter interface {
	Get{{$k}}() {{$v}}
}

type {{$k}}Setter interface {
	Set{{$k}}({{$v}})
}

type {{$k}} struct {
	valid bool
	v {{$v}}
}

func (o *{{$k}}) Get{{$k}}() {{$v}} {
	if !o.valid {
		panic("{{$k}} value is not valid")
	}
	return o.v
}

func (o *{{$k}}) Set{{$k}}(v {{$v}}) {
	o.v = v
	o.valid = true
}
{{- end }}

`))
