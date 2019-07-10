package action

import (
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/constants"
	"github.com/yunify/qsctl/contexts"
	"github.com/yunify/qsctl/helper"
)

// DeleteObject will delete a remote object.
func DeleteObject(remote string) (err error) {
	bucketName, objectKey, err := ParseQsPath(remote)
	if err != nil {
		return
	}

	if objectKey == "" {
		return constants.ErrorQsPathObjectKeyRequired
	}

	_, err = contexts.SetupBuckets(bucketName, "")
	if err != nil {
		return
	}
	// Head to check whether object not found or forbidden
	if _, err = helper.HeadObject(objectKey); err != nil {
		switch err {
		case constants.ErrorQsPathNotFound:
			log.Errorf("object key <%s> not found", objectKey)
		case constants.ErrorQsPathAccessForbidden:
			log.Errorf("not enough permission for object <%s>", objectKey)
		}
		return
	}
	if err = helper.DeleteObject(objectKey); err != nil {
		return
	}
	log.Infof("Object <%s> removed.", objectKey)
	return
}
