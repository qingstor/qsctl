package action

import (
	"fmt"

	"github.com/c2h5oh/datasize"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/utils"
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

	fmt.Println(utils.AlignPrintWithColon(content...))
	return
}
