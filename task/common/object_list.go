package common

import (
	"errors"

	"github.com/Xuanwo/storage/pkg/iterator"
	"github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *ObjectListTask) run() {
	pairs := make([]*types.Pair, 0)

	if !t.GetRecursive() {
		pairs = append(pairs, types.WithDelimiter("/"))
	}

	it := t.GetDestinationStorage().ListDir(t.GetKey(), pairs...)

	// Always close the object channel.
	defer close(t.GetObjectChannel())

	for {
		o, err := it.Next()
		if err != nil && errors.Is(err, iterator.ErrDone) {
			break
		}
		if err != nil {
			t.TriggerFault(fault.NewUnhandled(err))
		}
		t.GetObjectChannel() <- o
	}

	log.Debugf("Task <%s> for key <%s> finished", "ObjectListTask", t.GetKey())
}
