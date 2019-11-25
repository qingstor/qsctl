package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/endpoint"
	"github.com/Xuanwo/storage/services/posixfs"
	"github.com/Xuanwo/storage/services/qingstor"
	typ "github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/pairs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/types"
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
func ParseLocalPath(p string) (pathType typ.ObjectType, err error) {
	// Use - means we will read from stdin.
	if p == "-" {
		return typ.ObjectTypeStream, nil
	}

	fi, err := os.Stat(p)
	if os.IsNotExist(err) {
		// if not exist, we use path's suffix to determine object type
		if strings.HasSuffix(p, string(os.PathSeparator)) {
			return typ.ObjectTypeDir, nil
		}
		return typ.ObjectTypeFile, nil
	}
	if err != nil {
		return typ.ObjectTypeInvalid, fmt.Errorf("parse path failed: {%w}", types.NewErrUnhandled(err))
	}
	if fi.IsDir() {
		return typ.ObjectTypeDir, nil
	}
	return typ.ObjectTypeFile, nil
}

// ParseQsPath will parse a key into different key type.
func ParseQsPath(p string) (keyType typ.ObjectType, bucketName, objectKey string, err error) {
	// qs-path includes three part: "qs://" prefix, bucket name and object key.
	// "qs://" prefix could be emit.
	p = strings.TrimPrefix(p, "qs://")

	s := strings.SplitN(p, "/", 2)

	// Only have bucket name or object key is "/"
	// For example: "qs://testbucket/"

	if len(s) == 1 || s[1] == "" {
		return typ.ObjectTypeDir, s[0], "", nil
	}

	if strings.HasSuffix(p, "/") {
		return typ.ObjectTypeDir, s[0], s[1], nil
	}
	return typ.ObjectTypeFile, s[0], s[1], nil
}

// ParseStorageInput will parse storage input and return a initiated storager.
func ParseStorageInput(input string, storageType typ.StoragerType) (path string, objectType typ.ObjectType, store storage.Storager, err error) {
	switch storageType {
	case posixfs.StoragerType:
		objectType, err = ParseLocalPath(input)
		if err != nil {
			return
		}
		path = input
		store = posixfs.NewClient()
		return
	case qingstor.StoragerType:
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
		panic(fmt.Errorf("no supported storager type %s", storageType))
	}
}

// ParseServiceInput will parse service input.
func ParseServiceInput(serviceType typ.ServicerType) (service storage.Servicer, err error) {
	switch serviceType {
	case qingstor.ServicerType:
		service, err = NewQingStorService()
		if err != nil {
			return
		}
		return
	default:
		panic(fmt.Errorf("no supported servicer type %s", serviceType))
	}
}

// ParseAtServiceInput will parse single args and setup service.
func ParseAtServiceInput(t interface {
	types.ServiceSetter
}) (err error) {
	service, err := ParseServiceInput(qingstor.ServicerType)
	if err != nil {
		return
	}
	setupService(t, service)
	return
}

// ParseAtStorageInput will parse single args and setup path, type, storager.
func ParseAtStorageInput(t interface {
	types.PathSetter
	types.StorageSetter
	types.TypeSetter
}, input string) (err error) {
	flow := ParseFlow(input, "")
	if flow != constants.FlowAtRemote {
		panic("invalid flow")
	}

	dstPath, dstType, dstStore, err := ParseStorageInput(input, qingstor.StoragerType)
	if err != nil {
		return
	}
	setupStorage(t, dstPath, dstType, dstStore)
	return
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
		srcType, dstType   typ.ObjectType
		srcStore, dstStore storage.Storager
	)

	switch flow {
	case constants.FlowToRemote:
		srcPath, srcType, srcStore, err = ParseStorageInput(src, posixfs.StoragerType)
		if err != nil {
			return
		}
		dstPath, dstType, dstStore, err = ParseStorageInput(dst, qingstor.StoragerType)
		if err != nil {
			return
		}
		dstPath = "/" + dstPath // Add / on qingstor path for base.
	case constants.FlowToLocal:
		srcPath, srcType, srcStore, err = ParseStorageInput(src, qingstor.StoragerType)
		if err != nil {
			return
		}
		srcPath = "/" + srcPath // Add / on qingstor path for base.
		dstPath, dstType, dstStore, err = ParseStorageInput(dst, posixfs.StoragerType)
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
}, path string, objectType typ.ObjectType, store storage.Storager) {
	t.SetSourcePath(path)
	t.SetSourceType(objectType)
	t.SetSourceStorage(store)
}

func setupDestinationStorage(t interface {
	types.DestinationPathSetter
	types.DestinationStorageSetter
	types.DestinationTypeSetter
}, path string, objectType typ.ObjectType, store storage.Storager) {
	t.SetDestinationPath(path)
	t.SetDestinationType(objectType)
	t.SetDestinationStorage(store)
}

func setupStorage(t interface {
	types.PathSetter
	types.StorageSetter
	types.TypeSetter
}, path string, objectType typ.ObjectType, store storage.Storager) {
	t.SetPath(path)
	t.SetType(objectType)
	t.SetStorage(store)
}

func setupService(t interface {
	types.ServiceSetter
}, store storage.Servicer) {
	t.SetService(store)
}

// NewQingStorService will create a new qingstor service.
func NewQingStorService() (*qingstor.Service, error) {
	srv := qingstor.New()
	err := srv.Init(
		pairs.WithEndpoint(endpoint.NewStaticFromParsedURL(
			viper.GetString(constants.ConfigProtocol),
			viper.GetString(constants.ConfigHost),
			viper.GetInt(constants.ConfigPort),
		)),
		pairs.WithCredential(credential.NewStatic(
			viper.GetString(constants.ConfigAccessKeyID),
			viper.GetString(constants.ConfigSecretAccessKey),
		)),
	)
	return srv, err
}

// ChooseDestinationStorage will choose the destination storage to fill.
func ChooseDestinationStorage(x interface {
	types.PathSetter
	types.StorageSetter
}, y interface {
	types.DestinationPathGetter
	types.DestinationStorageGetter
}) {
	x.SetPath(y.GetDestinationPath())
	x.SetStorage(y.GetDestinationStorage())
}

// ChooseSourceStorage will choose the source storage to fill.
func ChooseSourceStorage(x interface {
	types.PathSetter
	types.StorageSetter
}, y interface {
	types.SourcePathGetter
	types.SourceStorageGetter
}) {
	x.SetPath(y.GetSourcePath())
	x.SetStorage(y.GetSourceStorage())
}

// ChooseDestinationStorageAsSegmenter will choose the destination storage as a segmenter.
func ChooseDestinationStorageAsSegmenter(x interface {
	types.PathSetter
	types.SegmenterSetter
}, y interface {
	types.DestinationPathGetter
	types.DestinationStorageGetter
}) (err error) {
	x.SetPath(y.GetDestinationPath())

	segmenter, ok := y.GetDestinationStorage().(storage.Segmenter)
	if !ok {
		return types.NewErrStorageInsufficientAbility(nil)
	}
	x.SetSegmenter(segmenter)
	return
}

// ChooseDestinationSegmenter will choose the destination storage as a segmenter.
func ChooseDestinationSegmenter(x interface {
	types.DestinationPathSetter
	types.DestinationSegmenterSetter
}, y interface {
	types.DestinationPathGetter
	types.DestinationStorageGetter
}) (err error) {
	x.SetDestinationPath(y.GetDestinationPath())

	segmenter, ok := y.GetDestinationStorage().(storage.Segmenter)
	if !ok {
		return types.NewErrStorageInsufficientAbility(nil)
	}
	x.SetDestinationSegmenter(segmenter)
	return
}
