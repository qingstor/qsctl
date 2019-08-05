package action

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/c2h5oh/datasize"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/utils"
)

// Stat will handle all stat actions.
func Stat(ctx context.Context) (err error) {
	// Get params from context
	zone := contexts.FromContext(ctx, constants.ZoneFlag).(string)
	remote := contexts.FromContext(ctx, "remote").(string)

	bucketName, objectKey, err := ParseQsPath(remote)
	if err != nil {
		return err
	}
	if objectKey == "" {
		return constants.ErrorQsPathObjectKeyRequired
	}
	err = contexts.Storage.SetupBucket(bucketName, zone)
	if err != nil {
		return
	}
	ctx = contexts.SetContext(ctx, "objectKey", objectKey)
	return StatRemoteObject(ctx)
}

// StatRemoteObject will stat a remote object.
func StatRemoteObject(ctx context.Context) (err error) {
	// Get params from context
	format := contexts.FromContext(ctx, constants.FormatFlag).(string)
	objectKey := contexts.FromContext(ctx, "objectKey").(string)

	om, err := contexts.Storage.HeadObject(objectKey)
	if err != nil {
		return
	}

	if format != "" {
		fmt.Println(statFormat(format, om))
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
