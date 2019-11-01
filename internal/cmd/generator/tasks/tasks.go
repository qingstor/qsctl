// The following directive is necessary to make the package coherent:
// This program generates types, It can be invoked by running
// go generate
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"
)

type task struct {
	Name           string   `json:"-"`
	Description    string   `json:"description"`
	InheritedValue []string `json:"inherited_value,omitempty"`
	MutableValue   []string `json:"mutable_value,omitempty"`
	RuntimeValue   []string `json:"runtime_value,omitempty"`

	SatisfiedParametricRequirement []string `json:"-"`
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
	"endwith": func(x, y string) bool {
		return strings.HasSuffix(x, y)
	},
	"merge": func(x, y []string) []string {
		a := make(map[string]struct{})
		for _, v := range x {
			a[v] = struct{}{}
		}
		for _, v := range y {
			a[v] = struct{}{}
		}
		o := make([]string, 0)
		for x := range a {
			o = append(o, x)
		}

		sort.Strings(o)
		return o
	},
}

func hasBetweenSubString(s []string, x string) bool {
	str := strings.Join(s, " ")
	return strings.Contains(str, "Source"+x) && strings.Contains(str, "Destination"+x)
}

func getRequirements(x *task) (getter, setter map[string]struct{}) {
	getter, setter = make(map[string]struct{}), make(map[string]struct{})

	for _, v := range x.InheritedValue {
		getter[v] = struct{}{}
	}
	for _, v := range x.MutableValue {
		setter[v] = struct{}{}
	}

	return
}

func getAbility(x *task) (getter, setter map[string]struct{}) {
	getter, setter = make(map[string]struct{}), make(map[string]struct{})

	for _, v := range x.InheritedValue {
		getter[v] = struct{}{}
	}
	for _, v := range x.MutableValue {
		setter[v] = struct{}{}
	}
	for _, v := range x.RuntimeValue {
		getter[v] = struct{}{}
		setter[v] = struct{}{}
	}

	return
}

// Check whether x contains y.
func isContain(x, y map[string]struct{}) bool {
	for i := range y {
		if _, ok := x[i]; !ok {
			return false
		}
	}

	return true
}

// Check whether y satisfy the x's requirement.
func satisfyRequirement(x *task, y *task) {
	abilityGetter, abilitySetter := getAbility(x)
	requirementGetter, requirementSetter := getRequirements(y)

	if !isContain(abilityGetter, requirementGetter) || !isContain(abilitySetter, requirementSetter) {
		return
	}
	x.SatisfiedParametricRequirement = append(x.SatisfiedParametricRequirement, y.Name)
}

func executeTemplate(tmpl *template.Template, w io.Writer, v interface{}) {
	err := tmpl.Execute(w, v)
	if err != nil {
		log.Fatal(err)
	}
}

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

	// Set task names.
	for k, v := range tasks {
		v.Name = k
	}

	// Add shims for storage.
	shimTasks := make(map[string]*task)
	for k, v := range tasks {
		if strings.HasSuffix(k, "Shim") {
			continue
		}

		t := *v
		t.Name = k + "Shim"
		t.Description = fmt.Sprintf("Storage shim task for %s", v.Name)
		t.InheritedValue = append(t.InheritedValue, t.RuntimeValue...)
		t.MutableValue = append(t.MutableValue, t.RuntimeValue...)

		t.RuntimeValue = nil
		for _, s := range []string{"Storage", "Path", "Type"} {
			if hasBetweenSubString(t.InheritedValue, s) {
				t.RuntimeValue = append(t.RuntimeValue, s)
			}
		}
		if len(t.RuntimeValue) == 0 {
			continue
		}
		shimTasks[t.Name] = &t
	}
	for k, v := range shimTasks {
		v := v
		tasks[k] = v
		taskNames = append(taskNames, k)
	}
	sort.Strings(taskNames)

	parametricTasks := make([]*task, 0)
	for _, v := range tasks {
		v := v

		if !strings.HasSuffix(v.Name, "Parametric") {
			continue
		}
		parametricTasks = append(parametricTasks, v)
	}

	requirementFile, err := os.Create("../pkg/types/requirements.go")
	if err != nil {
		log.Fatal(err)
	}
	defer requirementFile.Close()

	executeTemplate(requirementPageTmpl, requirementFile, nil)

	mockFile, err := os.Create("../pkg/types/mocks.go")
	if err != nil {
		log.Fatal(err)
	}
	defer mockFile.Close()

	executeTemplate(mockPageTmpl, mockFile, nil)

	taskFile, err := os.Create("generated.go")
	if err != nil {
		log.Fatal(err)
	}
	defer taskFile.Close()

	executeTemplate(pageTmpl, taskFile, nil)

	testFile, err := os.Create("generated_test.go")
	if err != nil {
		log.Fatal(err)
	}
	defer testFile.Close()

	executeTemplate(testPageTmpl, testFile, nil)

	for _, taskName := range taskNames {
		v := tasks[taskName]

		for _, rv := range parametricTasks {
			satisfyRequirement(v, rv)
		}

		if strings.HasSuffix(taskName, "Shim") {
			executeTemplate(requirementTmpl, requirementFile, v)
			executeTemplate(shimTaskTmpl, taskFile, v)
		} else if strings.HasSuffix(taskName, "Parametric") {
			executeTemplate(requirementTmpl, requirementFile, v)
			executeTemplate(parametricTmpl, requirementFile, v)
		} else {
			executeTemplate(requirementTmpl, requirementFile, v)
			executeTemplate(mockTmpl, mockFile, v)
			executeTemplate(taskTmpl, taskFile, v)
			executeTemplate(taskTestTmpl, testFile, v)
		}
	}
}

var requirementPageTmpl = template.Must(template.New("requirementPage").Funcs(funcs).Parse(`// Code generated by go generate; DO NOT EDIT.
package types

import (
	"github.com/Xuanwo/navvy"
)
`))

var requirementTmpl = template.Must(template.New("requirement").Funcs(funcs).Parse(`
// {{ .Name }}Requirement is the requirement for {{ .Name }}Task.
type {{ .Name }}Requirement interface {
	navvy.Task

	// Predefined inherited value
	PoolGetter
	FaultGetter

	// Inherited value
{{- range $k, $v := .InheritedValue }}
	{{$v}}Getter
{{- end }}

	// Mutable value
{{- range $k, $v := .MutableValue }}
	{{$v}}Setter
{{- end }}
}
`))

var parametricTmpl = template.Must(template.New("requirement").Funcs(funcs).Parse(`
type {{ .Name }}Func func(navvy.Task){{ .Name }}Requirement
`))

var mockPageTmpl = template.Must(template.New("mockPage").Funcs(funcs).Parse(`// Code generated by go generate; DO NOT EDIT.
package types
`))

var mockTmpl = template.Must(template.New("mock").Funcs(funcs).Parse(`
// <ock{{ .Name }}Task is the mock task for {{ .Name }}Task.
type Mock{{ .Name }}Task struct {
	Pool
	Fault
	ID

	// Inherited and mutable values.
{{- range $k, $v := merge .InheritedValue .MutableValue}}
	{{$v}}
{{- end }}
}

func (t *Mock{{ .Name }}Task) Run() {
	panic("mock{{ .Name }}Task should not be run.")
}
`))

var pageTmpl = template.Must(template.New("page").Parse(`// Code generated by go generate; DO NOT EDIT.
package task

import (
	"fmt"

	"github.com/Xuanwo/navvy"
	"github.com/google/uuid"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/pkg/schedule"
)

var _ navvy.Pool
var _ types.Pool
var _ = uuid.New()
`))

var taskTmpl = template.Must(template.New("task").Funcs(funcs).Parse(`
// {{ .Name }}Task will {{ .Description }}.
type {{ .Name }}Task struct {
	types.{{ .Name }}Requirement

	// Predefined runtime value
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
	t.GetScheduler().Wait()
}

func (t *{{ .Name }}Task) TriggerFault(err error) {
	t.GetFault().Append(fmt.Errorf("Task {{ .Name }} failed: {%w}", err))
}

// New{{ .Name }} will create a {{ .Name }}Task struct and fetch inherited data from parent task.
func New{{ .Name }}(task navvy.Task) *{{ .Name }}Task {
	t := &{{ .Name }}Task{
		{{ .Name }}Requirement: task.(types.{{ .Name }}Requirement),
	}
	t.SetID(uuid.New().String())
	t.SetScheduler(schedule.NewScheduler(t.GetPool()))

	t.new()
	return t
}

// New{{ .Name }}Task will create a {{ .Name }}Task which meets navvy.Task.
func New{{ .Name }}Task(task navvy.Task) navvy.Task {
	return New{{ .Name }}(task)
}

{{- range $_, $v := .SatisfiedParametricRequirement }}
// New{{ $.Name }}{{ $v }}Task will create a {{ $.Name }}Task which meets types.{{ $v }}Requirement.
func New{{ $.Name }}{{ $v }}Task(task navvy.Task) types.{{ $v }}Requirement {
	return New{{ $.Name }}(task)
}
{{- end }}
`))

var shimTaskTmpl = template.Must(template.New("shimTask").Funcs(funcs).Parse(`
{{ $name := .Name | lowerFirst }}
// {{ $name }}Task will {{ .Description }}.
type {{ $name }}Task struct {
	types.{{ .Name }}Requirement

	// Runtime value
{{- range $k, $v := .RuntimeValue }}
	types.{{$v}}
{{- end }}
}

// Run implement navvy.Task
func (t *{{ $name }}Task) Run() {}

// New{{ .Name }} will create a {{ $name }}Task struct and fetch inherited data from parent task.
func New{{ .Name }}(task navvy.Task) *{{ $name }}Task {
	t := &{{ $name }}Task{
		{{ .Name }}Requirement: task.(types.{{ .Name }}Requirement),
	}
	return t
}
`))

var testPageTmpl = template.Must(template.New("testPage").Parse(`// Code generated by go generate; DO NOT EDIT.
package task

import (
	"errors"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/pkg/fault"
)

var _ navvy.Pool
var _ types.Pool
`))

var taskTestTmpl = template.Must(template.New("taskTest").Funcs(funcs).Parse(`
func Test{{ .Name }}Task_TriggerFault(t *testing.T) {
	m := &types.Mock{{ .Name }}Task{}
	m.SetFault(fault.New())
	task := &{{ .Name }}Task{ {{ .Name }}Requirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMock{{ .Name }}Task_Run(t *testing.T) {
	task := &types.Mock{{ .Name }}Task{}
	assert.Panics(t, func() {
		task.Run()
	})
}
`))
