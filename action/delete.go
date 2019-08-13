package action

import (
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
)

// DeleteHandler is all params for Delete func
type DeleteHandler struct {
	// Remote is the remote qs path
	Remote string `json:"remote"`
	// Zone specifies the zone for delete action
	Zone string `json:"zone"`
}

// WithRemote sets the Remote field with given remote path
func (dh *DeleteHandler) WithRemote(path string) *DeleteHandler {
	dh.Remote = path
	return dh
}

// WithZone sets the Zone field with given zone
func (dh *DeleteHandler) WithZone(z string) *DeleteHandler {
	dh.Zone = z
	return dh
}

// Delete will delete a remote object.
func (dh *DeleteHandler) Delete() (err error) {
	bucketName, objectKey, err := ParseQsPath(dh.Remote)
	if err != nil {
		return
	}

	if objectKey == "" {
		return constants.ErrorQsPathObjectKeyRequired
	}

	err = contexts.Storage.SetupBucket(bucketName, dh.Zone)
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
