package types

import (
	"github.com/Xuanwo/navvy"
)

// TodoFunc will be used be TODO to add new todo
type TodoFunc func(Todoist) navvy.Task

// Todo is the struct which inplements Todoist.
type Todo struct {
	v []TodoFunc
}

// Todoist will holding a todo list for tasks.
type Todoist interface {
	AddTODOs(...TodoFunc)
	NextTODO() TodoFunc
}

// NextTODO will return next task to do.
func (o *Todo) NextTODO() TodoFunc {
	if len(o.v) == 0 {
		return nil
	}
	fn := o.v[0]
	o.v[0] = nil
	o.v = o.v[1:]
	return fn
}

// AddTODOs will add new todos.
func (o *Todo) AddTODOs(v ...TodoFunc) {
	o.v = v
}
