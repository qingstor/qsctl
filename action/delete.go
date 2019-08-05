package action

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
)

// Delete will delete a remote object.
func Delete(ctx context.Context) (err error) {
	// Get params from context
	zone := contexts.FromContext(ctx, constants.ZoneFlag).(string)
	remote := contexts.FromContext(ctx, "remote").(string)

	bucketName, objectKey, err := ParseQsPath(remote)
	if err != nil {
		return
	}

	if objectKey == "" {
		return constants.ErrorQsPathObjectKeyRequired
	}

	err = contexts.Storage.SetupBucket(bucketName, zone)
	if err != nil {
		return
	}
	// Head to check whether object not found or forbidden
	if _, err = contexts.Storage.HeadObject(objectKey); err != nil {
		switch err {
		case constants.ErrorQsPathNotFound:
			log.Errorf("object key <%s> not found", objectKey)
		case constants.ErrorQsPathAccessForbidden:
			log.Errorf("not enough permission for object <%s>", objectKey)
		}
		return
	}
	if err = contexts.Storage.DeleteObject(objectKey); err != nil {
		return
	}
	log.Infof("Object <%s> removed.", objectKey)
	return
}
