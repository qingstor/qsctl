package task

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/yunify/qsctl/v2/constants"
)

func (t *ObjectPresignTask) run() {
	if _, err := t.GetStorage().HeadObject(t.GetKey()); err != nil {
		panic(err)
	}
	// get acl response
	ar, err := t.GetStorage().GetBucketACL()
	if err != nil {
		panic(err)
	}
	// check whether the bucket is public or not
	var isPublic bool
	for _, acl := range ar.ACLs {
		if acl.GranteeName == constants.PublicBucketACL {
			isPublic = true
			break
		}
	}
	if isPublic {
		t.AddTODOs(NewObjectPresignPublicTask)
		return
	}
	t.AddTODOs(NewObjectPresignPrivateTask)
}

func (t *ObjectPresignPublicTask) run() {
	// compose the public bucket url
	url := fmt.Sprintf("%s://%s.%s.%s:%d/%s",
		viper.GetString(constants.ConfigProtocol),
		t.GetBucketName(),
		t.GetStorage().GetBucketZone(),
		viper.GetString(constants.ConfigHost),
		viper.GetInt(constants.ConfigPort),
		t.GetKey(),
	)
	t.SetURL(url)
	log.Debugf("Task <%s> for key <%s> finished, get signed URL <%s>",
		"ObjectPresignPublicTask", t.GetKey(), url)
}

func (t *ObjectPresignPrivateTask) run() {
	url, err := t.GetStorage().PresignObject(t.GetKey(), t.GetExpire())
	if err != nil {
		panic(err)
	}
	t.SetURL(url)
	log.Debugf("Task <%s> for key <%s> finished, get signed URL <%s>",
		"ObjectPresignPrivateTask", t.GetKey(), url)
}
