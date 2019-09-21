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
	Name           string
	Type           string   `json:"type"`
	Path           string   `json:"path"`
	Depend         string   `json:"depend"`
	Description    string   `json:"description"`
	InheritedValue []string `json:"inherited_value"`
	RuntimeValue   []string `json:"runtime_value"`
}

var funcs = template.FuncMap{
	"getTypeName": func(s string) string {
		if strings.HasSuffix(s, "etter") {
			return s[:len(s)-6]
		}
		return s
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
		taskNames = append(taskNames, k)
	}
	sort.Strings(taskNames)

	// Set task name and categorized via path.
	pages := make(map[string][]*task)
	for _, v := range taskNames {
		tasks[v].Name = v
		pages[tasks[v].Path] = append(pages[tasks[v].Path], tasks[v])
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

var pageTmpl = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
package {{ .Package }}

import (
	"github.com/Xuanwo/navvy"

	"github.com/yunify/qsctl/v2/task/types"
	"github.com/yunify/qsctl/v2/task/utils"
)

var _ navvy.Pool
var _ types.Pool
var _ = utils.SubmitNextTask
`))

var requiredTaskTmpl = template.Must(template.New("").Funcs(funcs).Parse(`
// {{ .Name }}TaskRequirement is the requirement for execute {{ .Name }}Task.
type {{ .Name }}TaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter

{{- range $k, $v := .RuntimeValue }}
	types.{{$v}}
{{- end }}
}

// {{ .Name }}Task will {{ .Description }}.
type {{ .Name }}Task struct {
	{{ .Name }}TaskRequirement
}

// mock{{ .Name }}Task is the mock task for {{ .Name }}Task.
type mock{{ .Name }}Task struct {
	types.Todo
	types.Pool
{{- range $k, $v := .RuntimeValue }}
	types.{{getTypeName $v}}
{{- end }}
}

func (t *mock{{ .Name }}Task) Run() {
	panic("mock{{ .Name }}Task should not be run.")
}

// New{{ .Name }}Task will create a new {{ .Name }}Task.
func New{{ .Name }}Task(task types.Todoist) navvy.Task {
	return &{{ .Name }}Task{task.({{ .Name }}TaskRequirement)}
}
`))

var taskTmpl = map[string]*template.Template{
	"required":  requiredTaskTmpl,
	"dependent": dependentTaskTmpl,
}

var dependentTaskTmpl = template.Must(template.New("").Parse(`
// {{ .Name }}Task will {{ .Description }}.
type {{ .Name }}Task struct {
	// Inherited value
{{- if .Depend }}
	types.Pool
{{- end }}
{{- range $k, $v := .InheritedValue }}
	types.{{$v}}
{{- end }}

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
func init{{ .Name }}Task(task types.Todoist) (t *{{ .Name }}Task, o *{{ .Depend }}Task) {
	o = task.(*{{ .Depend }}Task)

	t = &{{ .Name }}Task{}
	t.SetPool(o.GetPool())
{{- range $k, $v := .InheritedValue }}
	t.Set{{$v}}(o.Get{{$v}}())
{{- end }}
	return
}
{{- else }}
// Wait will wait until {{ .Name }}Task has been finished
func (t *{{ .Name }}Task) Wait() {
	t.GetPool().Wait()
}
{{- end }}
`))
