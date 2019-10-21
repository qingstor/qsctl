package common

import (
	"errors"

	"github.com/Xuanwo/storage/pkg/iterator"
	"github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *ObjectListTask) run() {
	log.Debugf("Task <%s> for key <%s> started", "ObjectListTask", t.GetPrefix())
	pairs := make([]*types.Pair, 0)

	if !t.GetRecursive() {
		pairs = append(pairs, types.WithDelimiter("/"))
	}

	it := t.GetDestinationStorage().ListDir(t.GetPrefix(), pairs...)

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

	log.Debugf("Task <%s> for key <%s> finished", "ObjectListTask", t.GetPrefix())
}
