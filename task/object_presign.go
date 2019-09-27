package task

import (
	log "github.com/sirupsen/logrus"
)

func (t *ObjectPresignTask) run() {
	if _, err := t.GetStorage().HeadObject(t.GetKey()); err != nil {
		panic(err)
	}

	url, err := t.GetStorage().PresignObject(t.GetKey(), t.GetExpire())
	if err != nil {
		panic(err)
	}
	t.SetURL(url)
	log.Debugf("Task <%s> for key <%s> finished, get signed URL <%s>",
		"ObjectPresignTask", t.GetKey(), url)
}
