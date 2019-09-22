package utils

import (
	"os"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/task/types"
	"github.com/yunify/qsctl/v2/task/utils"
)

// SetupStorage will setup storage for task.
func SetupStorage(t interface {
	types.StorageSetter
	types.BucketNameGetter
}) (err error) {
	stor, err := storage.NewQingStorObjectStorage()
	if err != nil {
		return err
	}
	t.SetStorage(stor)

	return stor.SetupBucket(t.GetBucketName(), "")
}

// ParseTowArgs will parse two args into flow, path and key.
func ParseTowArgs(t interface {
	types.FlowTypeSetter
	types.PathSetter
	types.StreamSetter
	types.PathTypeSetter
	types.KeySetter
	types.KeyTypeSetter
	types.BucketNameSetter
}, args []string) (err error) {
	src, dst := args[0], args[1]
	flow := utils.ParseFlow(src, dst)
	t.SetFlowType(flow)

	var path string
	var pathType constants.PathType
	switch flow {
	case constants.FlowToRemote:
		pathType, err = utils.ParsePath(src)
		if err != nil {
			return err
		}
		t.SetPathType(pathType)
		path = src

		keyType, bucketName, objectKey, err := utils.ParseKey(dst)
		if err != nil {
			return err
		}
		t.SetKeyType(keyType)
		t.SetKey(objectKey)
		t.SetBucketName(bucketName)
	case constants.FlowToLocal, constants.FlowAtRemote:
		pathType, err = utils.ParsePath(dst)
		if err != nil {
			return err
		}
		t.SetPathType(pathType)
		path = dst

		keyType, bucketName, objectKey, err := utils.ParseKey(src)
		if err != nil {
			return err
		}
		t.SetKeyType(keyType)
		t.SetKey(objectKey)
		t.SetBucketName(bucketName)
	default:
		panic("this case should never be switched")
	}

	t.SetPath(path)

	switch pathType {
	case constants.PathTypeFile:
		t.SetPath(path)
	case constants.PathTypeStream:
		// TODO: we could support other stream type, for example, read from a socket.
		t.SetStream(os.Stdin)
	case constants.PathTypeLocalDir:
	default:
		panic("invalid path type")
	}

	return
}
