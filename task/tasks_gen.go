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
	"path"
	"sort"
	"strings"
	"text/template"
)

type task struct {
	Name           string   `json:"-"`
	Type           string   `json:"type"`
	Path           string   `json:"path"`
	Depend         string   `json:"depend,omitempty"`
	Description    string   `json:"description"`
	InheritedValue []string `json:"inherited_value,omitempty"`
	RuntimeValue   []string `json:"runtime_value,omitempty"`
}

var funcs = template.FuncMap{
	"lowerFirst": func(s string) string {
		if len(s) == 0 {
			return ""
		}
		if s[0] < 'A' || s[0] > 'Z' {
			return s
		}
		return string(s[0]+'a'-'A') + s[1:]
	},
}

//go:generate go run tasks_gen.go
func main() {
	data, err := ioutil.ReadFile("tasks.json")
	if err != nil {
		log.Fatal(err)
	}

	tasks := make(map[string]*task)
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		log.Fatal(err)
	}

	// Do sort to all tasks via name.
	taskNames := make([]string, 0)
	for k := range tasks {
		sort.Strings(tasks[k].InheritedValue)
		sort.Strings(tasks[k].RuntimeValue)

		taskNames = append(taskNames, k)
	}
	sort.Strings(taskNames)

	// Set task name and categorized via path.
	pages := make(map[string][]*task)
	for _, v := range taskNames {
		tasks[v].Name = v
		pages[tasks[v].Path] = append(pages[tasks[v].Path], tasks[v])
	}

	// Format input tasks.json
	data, err = json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("tasks.json", data, 0664)
	if err != nil {
		log.Fatal(err)
	}

	for pathName, page := range pages {
		err := os.MkdirAll(path.Dir(pathName), 0664)
		if err != nil {
			log.Fatal(err)
		}
		f, err := os.Create(pathName + "_generated.go")
		if err != nil {
			log.Fatal(err)
		}

		packageName := "task"
		if strings.Contains(pathName, "/") {
			packageName = path.Dir(pathName)
		}

		// Write page temple firstly.
		err = pageTmpl.Execute(f, struct {
			Package string
		}{
			packageName,
		})
		if err != nil {
			log.Fatal(err)
		}

		// Write task.
		for _, task := range page {
			err = requirementTmpl.Execute(f, task)
			if err != nil {
				log.Fatal(err)
			}
			err = mockTmpl.Execute(f, task)
			if err != nil {
				log.Fatal(err)
			}
			err = taskTmpl[task.Type].Execute(f, task)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}

var taskTmpl = map[string]*template.Template{
	"required":  requiredTaskTmpl,
	"dependent": dependentTaskTmpl,
}

var pageTmpl = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
package {{ .Package }}

import (
	"github.com/Xuanwo/navvy"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/task/utils"
)

var _ navvy.Pool
var _ types.Pool
var _ = utils.SubmitNextTask
`))

var requirementTmpl = template.Must(template.New("").Funcs(funcs).Parse(`
// {{ .Name | lowerFirst }}TaskRequirement is the requirement for execute {{ .Name }}Task.
type {{ .Name | lowerFirst }}TaskRequirement interface {
	navvy.Task
{{- if eq .Type "required" }}
	types.Todoist
	types.PoolGetter
{{ else }}
{{- if .Depend }}
	types.PoolGetter
{{- end }}
{{- end }}

	// Inherited value
{{- range $k, $v := .InheritedValue }}
	types.{{$v}}Getter
{{- end }}

{{- if eq .Type "required" }}
	// Runtime value
{{- range $k, $v := .RuntimeValue }}
	types.{{$v}}Setter
{{- end }}
{{- end }}
}
`))

var mockTmpl = template.Must(template.New("").Funcs(funcs).Parse(`
// mock{{ .Name }}Task is the mock task for {{ .Name }}Task.
type mock{{ .Name }}Task struct {
	types.Todo
	types.Pool

	// Inherited value
{{- range $k, $v := .InheritedValue }}
	types.{{$v}}
{{- end }}

{{- if eq .Type "required" }}
	// Runtime value
{{- range $k, $v := .RuntimeValue }}
	types.{{$v}}
{{- end }}
{{- end }}
}

func (t *mock{{ .Name }}Task) Run() {
	panic("mock{{ .Name }}Task should not be run.")
}
`))

var requiredTaskTmpl = template.Must(template.New("").Funcs(funcs).Parse(`
// {{ .Name }}Task will {{ .Description }}.
type {{ .Name }}Task struct {
	{{ .Name | lowerFirst }}TaskRequirement
}

// Run implement navvy.Task.
func (t *{{ .Name }}Task) Run() {
	t.run()
	utils.SubmitNextTask(t.{{ .Name | lowerFirst }}TaskRequirement)
}

// New{{ .Name }}Task will create a new {{ .Name }}Task.
func New{{ .Name }}Task(task types.Todoist) navvy.Task {
	return &{{ .Name }}Task{task.({{ .Name | lowerFirst }}TaskRequirement)}
}
`))

var dependentTaskTmpl = template.Must(template.New("").Funcs(funcs).Parse(`
// {{ .Name }}Task will {{ .Description }}.
type {{ .Name }}Task struct {
	{{ .Name | lowerFirst }}TaskRequirement

	// Runtime value
	types.Todo
{{- range $k, $v := .RuntimeValue }}
	types.{{$v}}
{{- end }}
}

// Run implement navvy.Task
func (t *{{ .Name }}Task) Run() {
	utils.SubmitNextTask(t)
}

{{- if .Depend }}
// init{{ .Name }}Task will create a {{ .Name }}Task and fetch inherited data from {{ .Depend }}Task.
func New{{ .Name }}Task(task types.Todoist) navvy.Task {
	t := &{{ .Name }}Task{
		{{ .Name | lowerFirst }}TaskRequirement: task.({{ .Name | lowerFirst }}TaskRequirement),
	}
	t.new()
	return t
}
{{- else }}
// Wait will wait until {{ .Name }}Task has been finished
func (t *{{ .Name }}Task) Wait() {
	t.GetPool().Wait()
}
{{- end }}
`))
