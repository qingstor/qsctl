package task

import (
	log "github.com/sirupsen/logrus"
)

func (t *ObjectStatTask) run() {
	om, err := t.GetStorage().HeadObject(t.GetKey())
	if err != nil {
		panic(err)
	}
	oriOm := t.GetObjectMeta()
	// replace the original om
	*oriOm = *om
	log.Debugf("Task <%s> for Key <%s> finished.", "StatObjectTask", t.GetKey())
}
