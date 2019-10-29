package task

import (
	"errors"

	"github.com/Xuanwo/navvy"
	"github.com/Xuanwo/storage/pkg/iterator"
	"github.com/yunify/qsctl/v2/pkg/fault"

	typ "github.com/Xuanwo/storage/types"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/types"
)

var listTaskConstructor = map[constants.ListType]types.TaskFunc{
	constants.ListTypeBucket: NewBucketListTask,
	constants.ListTypeKey:    NewObjectListTask,
}

// NewListTask will create a list task.
func NewListTask(fn func(*ListTask)) *ListTask {
	t := &ListTask{}

	pool := navvy.NewPool(10)
	t.SetPool(pool)

	fn(t)

	oc := make(chan *typ.Object)
	t.SetObjectChannel(oc)

	todo := listTaskConstructor[t.GetListType()]
	if todo == nil {
		panic("invalid todo func")
	}
	t.AddTODOs(todo)
	return t
}

func (t *ObjectListTask) run() {
	log.Debugf("Task <%s> for key <%s> started", "ObjectListTask", t.GetDestinationPath())
	pairs := make([]*typ.Pair, 0)

	if !t.GetRecursive() {
		pairs = append(pairs, typ.WithDelimiter("/"))
	}

	it := t.GetDestinationStorage().ListDir(t.GetDestinationPath(), pairs...)

	// Always close the object channel.
	defer close(t.GetObjectChannel())

	for {
		o, err := it.Next()
		if err != nil && errors.Is(err, iterator.ErrDone) {
			break
		}
		if err != nil {
			t.TriggerFault(fault.NewUnhandled(err))
			return
		}
		t.GetObjectChannel() <- o
	}

	log.Debugf("Task <%s> for key <%s> finished", "ObjectListTask", t.GetDestinationPath())
}

func (t *ObjectListAsyncTask) run() {
	log.Debugf("Task <%s> for key <%s> started", "ObjectListAsyncTask", t.GetDestinationPath())

	go func() {
		pairs := make([]*typ.Pair, 0)

		if !t.GetRecursive() {
			pairs = append(pairs, typ.WithDelimiter("/"))
		}

		it := t.GetDestinationStorage().ListDir(t.GetDestinationPath(), pairs...)
		defer close(t.GetObjectChannel())
		for {
			o, err := it.Next()
			if err != nil && errors.Is(err, iterator.ErrDone) {
				break
			}
			if err != nil {
				t.TriggerFault(fault.NewUnhandled(err))
				return
			}
			t.GetObjectChannel() <- o
		}
	}()

	log.Debugf("Task <%s> for key <%s> finished", "ObjectListAsyncTask", t.GetDestinationPath())
}
