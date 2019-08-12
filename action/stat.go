package action

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/c2h5oh/datasize"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/utils"
)

// StatHandler is all params for Stat func
type StatHandler struct {
	BaseHandler
	// Remote is the remote qs path
	Remote string `json:"remote"`
	// ObjectKey is the remote object key
	ObjectKey string `json:"object_key"`
}

// WithZone rewrite the WithZone method
func (sh *StatHandler) WithZone(z string) *StatHandler {
	sh.Zone = z
	return sh
}

// WithFormat rewrite the WithFormat method
func (sh *StatHandler) WithFormat(f string) *StatHandler {
	sh.Format = f
	return sh
}

// WithRemote sets the Remote field with given remote path
func (sh *StatHandler) WithRemote(path string) *StatHandler {
	sh.Remote = path
	return sh
}

// WithObjectKey sets the ObjectKey field with given key
func (sh *StatHandler) WithObjectKey(key string) *StatHandler {
	sh.ObjectKey = key
	return sh
}

// Stat will handle all stat actions.
func (sh *StatHandler) Stat() (err error) {
	bucketName, objectKey, err := ParseQsPath(sh.Remote)
	if err != nil {
		return err
	}
	if objectKey == "" {
		return constants.ErrorQsPathObjectKeyRequired
	}
	err = contexts.Storage.SetupBucket(bucketName, sh.Zone)
	if err != nil {
		return
	}
	return sh.WithObjectKey(objectKey).StatRemoteObject()
}

// StatRemoteObject will stat a remote object.
func (sh *StatHandler) StatRemoteObject() (err error) {
	om, err := contexts.Storage.HeadObject(sh.ObjectKey)
	if err != nil {
		return
	}

	if sh.Format != "" {
		fmt.Println(statFormat(sh.Format, om))
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

func statFormat(input string, om *storage.ObjectMeta) string {
	input = strings.ReplaceAll(input, "%F", om.ContentType)
	input = strings.ReplaceAll(input, "%h", om.ETag)
	input = strings.ReplaceAll(input, "%n", om.Key)
	input = strings.ReplaceAll(input, "%s", strconv.FormatInt(om.ContentLength, 10))
	input = strings.ReplaceAll(input, "%y", om.LastModified.String())
	input = strings.ReplaceAll(input, "%Y", strconv.FormatInt(om.LastModified.Unix(), 10))

	return input
}
