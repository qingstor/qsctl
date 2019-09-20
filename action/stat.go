package action

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/c2h5oh/datasize"
	utils2 "github.com/yunify/qsctl/v2/task/utils"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/utils"
)

// StatHandler is all params for Stat func
type StatHandler struct {
	// Format is the user-specified output format
	Format string `json:"format"`
	// Key is the remote object key
	ObjectKey string `json:"object_key"`
	// Remote is the remote qs path
	Remote string `json:"remote"`
	// Zone specifies the zone for stat action
	Zone string `json:"zone"`
}

// WithFormat sets the Format field with given format string
func (sh *StatHandler) WithFormat(f string) *StatHandler {
	sh.Format = f
	return sh
}

// WithObjectKey sets the Key field with given key
func (sh *StatHandler) WithObjectKey(key string) *StatHandler {
	sh.ObjectKey = key
	return sh
}

// WithRemote sets the Remote field with given remote path
func (sh *StatHandler) WithRemote(path string) *StatHandler {
	sh.Remote = path
	return sh
}

// WithZone sets the Zone field with given zone
func (sh *StatHandler) WithZone(z string) *StatHandler {
	sh.Zone = z
	return sh
}

// Stat will handle all stat actions.
func (sh *StatHandler) Stat() (err error) {
	bucketName, objectKey, err := utils2.ParseQsPath(sh.Remote)
	if err != nil {
		return err
	}
	if objectKey == "" {
		return constants.ErrorQsPathObjectKeyRequired
	}
	err = stor.SetupBucket(bucketName, sh.Zone)
	if err != nil {
		return
	}
	return sh.WithObjectKey(objectKey).StatRemoteObject()
}

// StatRemoteObject will stat a remote object.
func (sh *StatHandler) StatRemoteObject() (err error) {
	om, err := stor.HeadObject(sh.ObjectKey)
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
