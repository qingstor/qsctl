package common

import (
	log "github.com/sirupsen/logrus"
)

func (t *ObjectPresignTask) run() {
	// if _, err := t.GetDestinationStorage().HeadObject(t.GetKey()); err != nil {
	// 	t.TriggerFault(fault.NewUnhandled(err))
	// 	return
	// }
	//
	// url, err := t.GetDestinationStorage().PresignObject(t.GetKey(), t.GetExpire())
	// if err != nil {
	// 	t.TriggerFault(fault.NewUnhandled(err))
	// 	return
	// }
	// t.SetURL(url)
	log.Debugf("Task <%s> for key <%s> finished, get signed URL <%s>",
		"ObjectPresignTask", t.GetKey(), "")
}
