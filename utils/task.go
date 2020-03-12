package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/endpoint"
	"github.com/Xuanwo/storage/services/fs"
	"github.com/Xuanwo/storage/services/qingstor"
	typ "github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/pairs"
	"github.com/qingstor/noah/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/qingstor/qsctl/v2/constants"
)

// ErrStoragerTypeInvalid returned when storager type invalid
var ErrStoragerTypeInvalid = errors.New("storager type no valid")

// ErrInvalidFlow returned when parsed flow not valid
var ErrInvalidFlow = errors.New("invalid flow")

// StoragerType is the alias for the type in storager
type StoragerType = string

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
func ParseStorageInput(input string, storageType StoragerType) (
	workDir, path string, objectType typ.ObjectType, store storage.Storager, err error) {
	switch storageType {
	case fs.Type:
		objectType, err = ParseLocalPath(input)
		if err != nil {
			return
		}
		workDir, path, err = ParseWorkDir(input, string(os.PathSeparator))
		if err != nil {
			return
		}
		log.Debugf("%s work dir: %s", fs.Type, workDir)
		_, store, err = fs.New(pairs.WithWorkDir(workDir))
		if err != nil {
			return
		}
	case qingstor.Type:
		var bucketName, objectKey string
		var srv storage.Servicer

		objectType, bucketName, objectKey, err = ParseQsPath(input)
		if err != nil {
			return
		}
		// always treat qs path as abs path, so add "/" before
		workDir, path, err = ParseWorkDir("/"+objectKey, "/")
		if err != nil {
			return
		}
		log.Debugf("%s work dir: %s", qingstor.Type, workDir)
		srv, err = NewQingStorService()
		if err != nil {
			return
		}
		store, err = srv.Get(bucketName, pairs.WithWorkDir(workDir))
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("%w %s", ErrStoragerTypeInvalid, storageType)
	}
	return
}

// ParseServiceInput will parse service input.
func ParseServiceInput(serviceType StoragerType) (service storage.Servicer, err error) {
	switch serviceType {
	case qingstor.Type:
		service, err = NewQingStorService()
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("%w %s", ErrStoragerTypeInvalid, serviceType)
	}
	return
}

// ParseAtServiceInput will parse single args and setup service.
func ParseAtServiceInput(t interface {
	types.ServiceSetter
}) (err error) {
	service, err := ParseServiceInput(qingstor.Type)
	if err != nil {
		return
	}
	setupService(t, service)
	return
}

// ParseAtStorageInput will parse single args and setup path, type, storager.
func ParseAtStorageInput(t interface {
	types.WorkDirSetter
	types.PathSetter
	types.StorageSetter
	types.TypeSetter
}, input string) (err error) {
	flow := ParseFlow(input, "")
	if flow != constants.FlowAtRemote {
		err = ErrInvalidFlow
		return
	}

	dstWorkDir, dstPath, dstType, dstStore, err := ParseStorageInput(input, qingstor.Type)
	if err != nil {
		return
	}
	setupStorage(t, dstWorkDir, dstPath, dstType, dstStore)
	return
}

// ParseBetweenStorageInput will parse two args into flow, path and key.
func ParseBetweenStorageInput(t interface {
	types.SourceWorkDirSetter
	types.SourcePathSetter
	types.SourceStorageSetter
	types.SourceTypeSetter
	types.DestinationWorkDirSetter
	types.DestinationPathSetter
	types.DestinationStorageSetter
	types.DestinationTypeSetter
}, src, dst string) (err error) {
	flow := ParseFlow(src, dst)
	var (
		srcWorkDir, dstWorkDir string
		srcPath, dstPath       string
		srcType, dstType       typ.ObjectType
		srcStore, dstStore     storage.Storager
	)

	switch flow {
	case constants.FlowToRemote:
		srcWorkDir, srcPath, srcType, srcStore, err = ParseStorageInput(src, fs.Type)
		if err != nil {
			return
		}
		dstWorkDir, dstPath, dstType, dstStore, err = ParseStorageInput(dst, qingstor.Type)
		if err != nil {
			return
		}
	case constants.FlowToLocal:
		srcWorkDir, srcPath, srcType, srcStore, err = ParseStorageInput(src, qingstor.Type)
		if err != nil {
			return
		}
		dstWorkDir, dstPath, dstType, dstStore, err = ParseStorageInput(dst, fs.Type)
		if err != nil {
			return
		}
	default:
		err = ErrInvalidFlow
		return
	}

	// if dstPath is blank while srcPath not,
	// it means copy file/dir to dst with the same name,
	// so set dst path to the src path
	if dstPath == "" && srcPath != "" {
		dstPath = srcPath
	}
	setupSourceStorage(t, srcWorkDir, srcPath, srcType, srcStore)
	setupDestinationStorage(t, dstWorkDir, dstPath, dstType, dstStore)
	return
}

func setupSourceStorage(t interface {
	types.SourceWorkDirSetter
	types.SourcePathSetter
	types.SourceStorageSetter
	types.SourceTypeSetter
}, workDir, path string, objectType typ.ObjectType, store storage.Storager) {
	t.SetSourceWorkDir(workDir)
	t.SetSourcePath(path)
	t.SetSourceType(objectType)
	t.SetSourceStorage(store)
}

func setupDestinationStorage(t interface {
	types.DestinationWorkDirSetter
	types.DestinationPathSetter
	types.DestinationStorageSetter
	types.DestinationTypeSetter
}, workDir, path string, objectType typ.ObjectType, store storage.Storager) {
	t.SetDestinationWorkDir(workDir)
	t.SetDestinationPath(path)
	t.SetDestinationType(objectType)
	t.SetDestinationStorage(store)
}

func setupStorage(t interface {
	types.WorkDirSetter
	types.PathSetter
	types.StorageSetter
	types.TypeSetter
}, workDir, path string, objectType typ.ObjectType, store storage.Storager) {
	t.SetWorkDir(workDir)
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
func NewQingStorService() (storage.Servicer, error) {
	var ep endpoint.Static
	switch protocol := viper.GetString(constants.ConfigProtocol); protocol {
	case endpoint.ProtocolHTTPS:
		ep = endpoint.NewHTTPS(
			viper.GetString(constants.ConfigHost),
			viper.GetInt(constants.ConfigPort))
	default: // endpoint.ProtocolHTTP:
		ep = endpoint.NewHTTP(
			viper.GetString(constants.ConfigHost),
			viper.GetInt(constants.ConfigPort))
	}
	srv, _, err := qingstor.New(
		pairs.WithEndpoint(ep),
		pairs.WithCredential(credential.MustNewHmac(
			viper.GetString(constants.ConfigAccessKeyID),
			viper.GetString(constants.ConfigSecretAccessKey),
		)),
	)
	return srv, err
}

// GetAbsPath combines the work dir and path together
func GetAbsPath(t interface {
	types.WorkDirGetter
	types.PathGetter
}) string {
	return filepath.Join(t.GetWorkDir(), t.GetPath())
}

// GetAbsSourcePath combines the source work dir and source path together
func GetAbsSourcePath(t interface {
	types.SourceWorkDirGetter
	types.SourcePathGetter
}) string {
	return filepath.Join(t.GetSourceWorkDir(), t.GetSourcePath())
}

// GetAbsDestinationPath combines the destination work dir and destination path together
func GetAbsDestinationPath(t interface {
	types.DestinationWorkDirGetter
	types.DestinationPathGetter
}) string {
	return filepath.Join(t.GetDestinationWorkDir(), t.GetDestinationPath())
}

// GetAbsBetweenPath combine abs source path and abs destination path separately
func GetAbsBetweenPath(t interface {
	types.DestinationWorkDirGetter
	types.DestinationPathGetter
	types.SourceWorkDirGetter
	types.SourcePathGetter
}) (string, string) {
	return GetAbsSourcePath(t), GetAbsDestinationPath(t)
}
