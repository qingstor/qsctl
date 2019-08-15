package task

import (
	"github.com/Xuanwo/navvy"
)

// Todoist will holding a todo list for tasks.
type Todoist interface {
	AddTODOs(...func(Todoist) navvy.Task)
	NextTODO() func(Todoist) navvy.Task
}

// Todo is the struct which inplements Todoist.
type Todo struct {
	v []func(Todoist) navvy.Task
}

// NextTODO will return next task to do.
func (o *Todo) NextTODO() func(Todoist) navvy.Task {
	if len(o.v) == 0 {
		return nil
	}
	fn := o.v[0]
	o.v[0] = nil
	o.v = o.v[1:]
	return fn
}

// AddTODOs will add new todos.
func (o *Todo) AddTODOs(v ...func(Todoist) navvy.Task) {
	o.v = v
}
