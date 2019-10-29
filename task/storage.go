package task

import (
	"errors"
	"testing"

	"github.com/Xuanwo/navvy"
	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/pkg/types"
)

// NewMakeBucketTask will create a make bucket task.
func NewMakeBucketTask(fn func(t *MakeBucketTask)) *MakeBucketTask {
	t := &MakeBucketTask{}

	pool := navvy.NewPool(10)
	t.SetPool(pool)

	fn(t)
	t.AddTODOs(NewBucketCreateTask)
	return t
}

func (t *BucketCreateTask) run() {
	_, err := t.GetDestinationService().Create(t.GetBucketName(), typ.WithLocation(t.GetZone()))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	log.Debugf("Task <%s> for Bucket <%s> finished.", "BucketCreateTask", t.GetBucketName())
}

func (t *BucketDeleteTask) run() {
	log.Debugf("Task <%s> for Bucket <%s> started.", "BucketDeleteTask", t.GetBucketName())
	err := t.GetDestinationService().Delete(t.GetBucketName())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	log.Debugf("Task <%s> for Bucket <%s> finished.", "BucketDeleteTask", t.GetBucketName())
}

func (t *BucketListTask) run() {
	resp, err := t.GetDestinationService().List(typ.WithLocation(t.GetZone()))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	buckets := make([]string, 0, len(resp))
	for _, v := range resp {
		b, err := v.Metadata()
		if err != nil {
			t.TriggerFault(fault.NewUnhandled(err))
			return
		}
		if name, ok := b.GetName(); ok {
			buckets = append(buckets, name)
		}
	}
	t.SetBucketList(buckets)
	log.Debugf("Task <%s> in zone <%s> finished.", "BucketListTask", t.GetZone())
}

func (t *RemoveBucketForceTask) new() {
	oc := make(chan *typ.Object)
	t.SetObjectChannel(oc)
	// done to notify get object from channel has done
	done := false
	t.SetDone(&done)
	// set recursive for list async task to list recursively
	t.SetRecursive(true)

	t.SetScheduler(types.NewScheduler(NewObjectDeleteScheduledTask))

	t.GetScheduler().Sync(t, NewObjectListAsyncTask)
	t.GetScheduler().Sync(t, NewObjectDeleteIterateTask)
	t.GetScheduler().Sync(t, NewAbortMultipartTask)
	t.GetScheduler().Sync(t, NewBucketDeleteTask)
}

// NewRemoveBucketTask will create a remove bucket task
func NewRemoveBucketTask(fn func(*RemoveBucketTask)) *RemoveBucketTask {
	t := &RemoveBucketTask{}

	pool := navvy.NewPool(10)
	t.SetPool(pool)

	fn(t)

	if t.ValidateFault() {
		return t
	}

	// check force flag, if true, do rm -r, then delete bucket
	if t.GetForce() {
		t.AddTODOs(NewRemoveBucketForceTask)
		return t
	}
	t.AddTODOs(NewBucketDeleteTask)
	return t
}
func TestNewRemoveBucketTask(t *testing.T) {
	removeBucketErr := errors.New("remove bucket error")
	cases := []struct {
		input            string
		force            bool
		expectedTodoFunc types.TaskFunc
		expectErr        error
	}{
		{"qs://test-bucket/obj", false, NewBucketDeleteTask, nil},
		{"qs://test-bucket/obj", true, NewRemoveBucketForceTask, nil},
		{"error", false, nil, removeBucketErr},
	}

	for _, v := range cases {
		pt := NewRemoveBucketTask(func(task *RemoveBucketTask) {
			task.SetForce(v.force)
			if v.expectErr != nil {
				task.TriggerFault(v.expectErr)
			}
		})

		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))

		if v.expectErr != nil {
			assert.Equal(t, true, pt.ValidateFault())
			assert.Equal(t, true, errors.Is(pt.GetFault(), v.expectErr))
		} else {
			assert.Equal(t, false, pt.ValidateFault())
		}
	}
}
