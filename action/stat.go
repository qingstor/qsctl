package action

import (
	"github.com/c2h5oh/datasize"

	"github.com/yunify/qsctl/constants"
	"github.com/yunify/qsctl/contexts"
	"github.com/yunify/qsctl/utils"
)

// Stat will handle all stat actions.
func Stat(remote string) (err error) {
	bucketName, objectKey, err := ParseQsPath(remote)
	if err != nil {
		return err
	}
	if objectKey == "" {
		return constants.ErrorQsPathObjectKeyRequired
	}
	err = contexts.Storage.SetupBucket(bucketName, "")
	if err != nil {
		return
	}

	return StatRemoteObject(objectKey)
}

// StatRemoteObject will stat a remote object.
func StatRemoteObject(objectKey string) (err error) {
	om, err := contexts.Storage.HeadObject(objectKey)
	if err != nil {
		return
	}

	content := []string{
		"Key: " + om.Key,
		"Size: " + datasize.ByteSize(om.ContentLength).String(),
		"Type: " + om.ContentType,
		"Modify: " + om.LastModified.String(),
		"StorageClass: " + om.StorageClass,
	}

	if om.ETag != "" {
		content = append(content, "MD5: "+om.ETag)
	}

	println(utils.AlignPrintWithColon(content...))
	return
}
