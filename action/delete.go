package action

import (
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
)

// DeleteHandler is all params for Delete func
type DeleteHandler struct {
	// Recursive means whether recursively delete objects
	Recursive bool `json:"recursive"`
	// Remote is the remote qs path
	Remote string `json:"remote"`
	// Zone specifies the zone for delete action
	Zone string `json:"zone"`
}

// WithRecursive sets the Recursive field with given bool value
func (dh *DeleteHandler) WithRecursive(r bool) *DeleteHandler {
	dh.Recursive = r
	return dh
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
	om, err := contexts.Storage.HeadObject(objectKey)
	if err != nil {
		switch err {
		case constants.ErrorQsPathNotFound:
			log.Errorf("object key <%s> not found", objectKey)
		case constants.ErrorQsPathAccessForbidden:
			log.Errorf("not enough permission for object <%s>", objectKey)
		}
		return
	}

	// If om is not a directory, delete the single object.
	if !om.IsDir() {
		if err = contexts.Storage.DeleteObject(objectKey); err != nil {
			return
		}
		log.Infof("Object <%s> removed.", objectKey)
		return
	}

	// If om is a directory, and recursive flag not set, return error.
	if !dh.Recursive {
		log.Errorf("directory should be deleted with -r")
		return constants.ErrorRecursiveRequired
	}

	err = contexts.Storage.DeleteMultipleObjects(objectKey)
	if err != nil {
		return err
	}
	log.Infof("Directory <%s> removed.", objectKey)
	return
}
