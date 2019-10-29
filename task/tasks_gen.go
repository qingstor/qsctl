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
	"sort"
	"text/template"
)

type task struct {
	Name           string   `json:"-"`
	Description    string   `json:"description"`
	InheritedValue []string `json:"inherited_value,omitempty"`
	MutableValue   []string `json:"mutable_value,omitempty"`
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

	// Format input tasks.json
	data, err = json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("tasks.json", data, 0664)
	if err != nil {
		log.Fatal(err)
	}

	taskFile, err := os.Create("generated.go")
	if err != nil {
		log.Fatal(err)
	}
	defer taskFile.Close()

	testFile, err := os.Create("generated_test.go")
	if err != nil {
		log.Fatal(err)
	}
	defer testFile.Close()

	// Write page temple firstly.
	err = pageTmpl.Execute(taskFile, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Write page temple firstly.
	err = testPageTmpl.Execute(testFile, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, taskName := range taskNames {
		v := tasks[taskName]
		v.Name = taskName

		// Write task.
		err = requirementTmpl.Execute(taskFile, v)
		if err != nil {
			log.Fatal(err)
		}
		err = mockTmpl.Execute(taskFile, v)
		if err != nil {
			log.Fatal(err)
		}
		err = taskTmpl.Execute(taskFile, v)
		if err != nil {
			log.Fatal(err)
		}

		// Write test.
		err = taskTestTmpl.Execute(testFile, v)
		if err != nil {
			log.Fatal(err)
		}
	}
}

var pageTmpl = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
package task

import (
	"fmt"

	"github.com/Xuanwo/navvy"
	"github.com/google/uuid"

	"github.com/yunify/qsctl/v2/pkg/types"
)

var _ navvy.Pool
var _ types.Pool
var _ = uuid.New()
`))

var testPageTmpl = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
package task

import (
	"errors"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
)

var _ navvy.Pool
var _ types.Pool
`))

var requirementTmpl = template.Must(template.New("").Funcs(funcs).Parse(`
// {{ .Name | lowerFirst }}TaskRequirement is the requirement for execute {{ .Name }}Task.
type {{ .Name | lowerFirst }}TaskRequirement interface {
	navvy.Task

	// Predefined inherited value
	types.PoolGetter

	// Inherited value
{{- range $k, $v := .InheritedValue }}
	types.{{$v}}Getter
{{- end }}

	// Mutable value
{{- range $k, $v := .MutableValue }}
	types.{{$v}}Setter
{{- end }}
}
`))

var mockTmpl = template.Must(template.New("").Funcs(funcs).Parse(`
// mock{{ .Name }}Task is the mock task for {{ .Name }}Task.
type mock{{ .Name }}Task struct {
	types.Pool
	types.Fault
	types.ID

	// Inherited value
{{- range $k, $v := .InheritedValue }}
	types.{{$v}}
{{- end }}

	// Mutable value
{{- range $k, $v := .MutableValue }}
	types.{{$v}}
{{- end }}
}

func (t *mock{{ .Name }}Task) Run() {
	panic("mock{{ .Name }}Task should not be run.")
}
`))

var taskTmpl = template.Must(template.New("").Funcs(funcs).Parse(`
// {{ .Name }}Task will {{ .Description }}.
type {{ .Name }}Task struct {
	{{ .Name | lowerFirst }}TaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Scheduler

	// Runtime value
{{- range $k, $v := .RuntimeValue }}
	types.{{$v}}
{{- end }}
}

// Run implement navvy.Task
func (t *{{ .Name }}Task) Run() {
	t.run()
}

func (t *{{ .Name }}Task) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task {{ .Name }} failed: {%w}", err))
}

// New{{ .Name }}Task will create a {{ .Name }}Task and fetch inherited data from parent task.
func New{{ .Name }}Task(task navvy.Task) navvy.Task {
	t := &{{ .Name }}Task{
		{{ .Name | lowerFirst }}TaskRequirement: task.({{ .Name | lowerFirst }}TaskRequirement),
	}
	t.SetID(uuid.New().String())
	t.new()
	return t
}
`))

var taskTestTmpl = template.Must(template.New("").Funcs(funcs).Parse(`
func TestNew{{ .Name }}Task(t *testing.T) {
	m := &mock{{ .Name }}Task{}
	task := New{{ .Name }}Task(m)
	assert.NotNil(t, task)
}

func Test{{ .Name }}Task_Run(t *testing.T) {
	cases := []struct {
		name     string
		hasFault bool
		hasCall  bool
		gotCall  bool
	}{
		{
			"has fault",
			true,
			false,
			false,
		},
		{
			"no fault",
			false,
			true,
			false,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			pool := navvy.NewPool(10)

			m := &mock{{ .Name }}Task{}
			m.SetPool(pool)
			task := &{{ .Name }}Task{ {{ .Name | lowerFirst }}TaskRequirement: m}

			err := errors.New("test error")
			if v.hasFault {
				task.SetFault(err)
			}
			task.GetScheduler.Sync(task, 
				func(todoist types.TaskFunc) navvy.Task {
				x := utils.NewCallbackTask(func() {
					v.gotCall = true
				})
				return x
			})

			task.Run()
			pool.Wait()

			assert.Equal(t, v.hasCall, v.gotCall)
		})
	}
}

func Test{{ .Name }}Task_TriggerFault(t *testing.T) {
	m := &mock{{ .Name }}Task{}
	task := &{{ .Name }}Task{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.{{ .Name | lowerFirst }}TaskRequirement.ValidateFault())
}

func TestMock{{ .Name }}Task_Run(t *testing.T) {
	task := &mock{{ .Name }}Task{}
	assert.Panics(t, func() {
		task.Run()
	})
}
`))
