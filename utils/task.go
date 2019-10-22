package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/services/posixfs"
	"github.com/Xuanwo/storage/services/qingstor"
	stypes "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/pkg/types"
)

// Current supported path type
const (
	StorageTypePOSIXFs  = "posixfs"
	StorageTypeQingStor = "qingstor"
)

// ParseFlow will parse the data flow
func ParseFlow(src, dst string) (flow constants.FlowType) {
	if dst == "" {
		return constants.FlowAtRemote
	}

	// If src and dst both local file or both remote object, the path is invalid.
	if strings.HasPrefix(src, "qs://") == strings.HasPrefix(dst, "qs://") {
		log.Errorf("Action between <%s> and <%s> is invalid", src, dst)
		return constants.FlowInvalid
	}

	if strings.HasPrefix(src, "qs://") {
		return constants.FlowToLocal
	}
	return constants.FlowToRemote
}

// ParseLocalPath will parse a path into different path type.
func ParseLocalPath(p string) (pathType stypes.ObjectType, err error) {
	// Use - means we will read from stdin.
	if p == "-" {
		return stypes.ObjectTypeStream, nil
	}

	fi, err := os.Stat(p)
	if os.IsNotExist(err) {
		return stypes.ObjectTypeInvalid, fmt.Errorf("parse path failed: {%w}", fault.NewLocalFileNotExist(err, p))
	}
	if err != nil {
		return stypes.ObjectTypeInvalid, fmt.Errorf("parse path failed: {%w}", fault.NewUnhandled(err))
	}
	if fi.IsDir() {
		return stypes.ObjectTypeDir, nil
	}
	return stypes.ObjectTypeFile, nil
}

// ParseQsPath will parse a key into different key type.
func ParseQsPath(p string) (keyType stypes.ObjectType, bucketName, objectKey string, err error) {
	// qs-path includes three part: "qs://" prefix, bucket name and object key.
	// "qs://" prefix could be emit.
	if strings.HasPrefix(p, "qs://") {
		p = p[5:]
	}
	s := strings.SplitN(p, "/", 2)

	// Only have bucket name or object key is "/"
	// For example: "qs://testbucket/"

	if len(s) == 1 || s[1] == "" {
		return stypes.ObjectTypeDir, s[0], "", nil
	}

	if strings.HasSuffix(p, "/") {
		return stypes.ObjectTypeDir, s[0], s[1], nil
	}
	return stypes.ObjectTypeFile, s[0], s[1], nil
}

// ParseStorageInput will parse storage input and return a initiated storager.
func ParseStorageInput(input, storageType string) (path string, objectType stypes.ObjectType, store storage.Storager, err error) {
	switch storageType {
	case StorageTypePOSIXFs:
		objectType, err = ParseLocalPath(input)
		if err != nil {
			return
		}
		path = input
		store = posixfs.NewClient()
		return
	case StorageTypeQingStor:
		var bucketName, objectKey string
		var srv *qingstor.Service

		objectType, bucketName, objectKey, err = ParseQsPath(input)
		if err != nil {
			return
		}
		srv, err = NewQingStorService()
		if err != nil {
			return
		}
		store, err = srv.Get(bucketName)
		if err != nil {
			return
		}
		path = objectKey
		return
	default:
		panic("error")
	}
}

// ParseBetweenStorageInput will parse two args into flow, path and key.
func ParseBetweenStorageInput(t interface {
	types.SourcePathSetter
	types.SourceStorageSetter
	types.SourceTypeSetter
	types.DestinationPathSetter
	types.DestinationStorageSetter
	types.DestinationTypeSetter
}, src, dst string) (err error) {
	flow := ParseFlow(src, dst)
	var (
		srcPath, dstPath   string
		srcType, dstType   stypes.ObjectType
		srcStore, dstStore storage.Storager
	)

	switch flow {
	case constants.FlowToRemote:
		srcPath, srcType, srcStore, err = ParseStorageInput(src, StorageTypePOSIXFs)
		if err != nil {
			return
		}
		dstPath, dstType, dstStore, err = ParseStorageInput(dst, StorageTypeQingStor)
		if err != nil {
			return
		}
	case constants.FlowToLocal:
		srcPath, srcType, srcStore, err = ParseStorageInput(src, StorageTypeQingStor)
		if err != nil {
			return
		}
		dstPath, dstType, dstStore, err = ParseStorageInput(dst, StorageTypePOSIXFs)
		if err != nil {
			return
		}
	default:
		panic("invalid flow")
	}

	setupSourceStorage(t, srcPath, srcType, srcStore)
	setupDestinationStorage(t, dstPath, dstType, dstStore)
	return
}

func setupSourceStorage(t interface {
	types.SourcePathSetter
	types.SourceStorageSetter
	types.SourceTypeSetter
}, path string, objectType stypes.ObjectType, store storage.Storager) {
	t.SetSourcePath(path)
	t.SetSourceType(objectType)
	t.SetSourceStorage(store)
}

func setupDestinationStorage(t interface {
	types.DestinationPathSetter
	types.DestinationStorageSetter
	types.DestinationTypeSetter
}, path string, objectType stypes.ObjectType, store storage.Storager) {
	t.SetDestinationPath(path)
	t.SetDestinationType(objectType)
	t.SetDestinationStorage(store)
}

// NewQingStorService will create a new qingstor service.
func NewQingStorService() (*qingstor.Service, error) {
	srv := qingstor.New()
	err := srv.Init(
		stypes.WithAccessKey(viper.GetString(constants.ConfigAccessKeyID)),
		stypes.WithSecretKey(viper.GetString(constants.ConfigSecretAccessKey)),
		stypes.WithHost(viper.GetString(constants.ConfigHost)),
		stypes.WithPort(viper.GetInt(constants.ConfigPort)),
		stypes.WithProtocol(viper.GetString(constants.ConfigProtocol)),
	)
	return srv, err
}
