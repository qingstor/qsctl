package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"

	fs "github.com/aos-dev/go-service-fs"
	qingstor "github.com/aos-dev/go-service-qingstor"
	"github.com/aos-dev/go-storage/v2"
	"github.com/aos-dev/go-storage/v2/pkg/credential"
	"github.com/aos-dev/go-storage/v2/pkg/endpoint"
	typ "github.com/aos-dev/go-storage/v2/types"
	"github.com/aos-dev/go-storage/v2/types/pairs"
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
		workDir, path, err = ParseFsWorkDir(input)
		if err != nil {
			return
		}
		log.Debugf("%s work dir: %s", fs.Type, workDir)
		store, err = fs.NewStorager(pairs.WithWorkDir(workDir))
		if err != nil {
			return
		}
	case qingstor.Type:
		var bucketName, objectKey string

		objectType, bucketName, objectKey, err = ParseQsPath(input)
		if err != nil {
			return
		}
		workDir, path = ParseQsWorkDir(objectKey)
		log.Debugf("%s work dir: %s", qingstor.Type, workDir)
		store, err = NewQingStorStorage(pairs.WithName(bucketName), pairs.WithWorkDir(workDir))
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
	types.PathSetter
	types.StorageSetter
	types.TypeSetter
}, input string) (dstWorkDir string, err error) {
	flow := ParseFlow(input, "")
	if flow != constants.FlowAtRemote {
		err = ErrInvalidFlow
		return
	}

	var (
		dstPath  string
		dstType  typ.ObjectType
		dstStore storage.Storager
	)
	dstWorkDir, dstPath, dstType, dstStore, err = ParseStorageInput(input, qingstor.Type)
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
}, src, dst string) (srcWorkDir, dstWorkDir string, err error) {
	flow := ParseFlow(src, dst)
	var (
		srcPath, dstPath   string
		srcType, dstType   typ.ObjectType
		srcStore, dstStore storage.Storager
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
func NewQingStorService() (storage.Servicer, error) {
	return qingstor.NewServicer(getQsServicePairs()...)
}

// NewQingStorStorage will create a new qingstor storage.
func NewQingStorStorage(pairs ...*typ.Pair) (storage.Storager, error) {
	srvPairs := getQsServicePairs()
	srvPairs = append(srvPairs, pairs...)
	return qingstor.NewStorager(srvPairs...)
}

func getQsServicePairs() []*typ.Pair {
	// init pairs with cap for less memory allocate
	ps := make([]*typ.Pair, 0, 5)
	ps = append(ps, pairs.WithCredential(credential.MustNewHmac(
		viper.GetString(constants.ConfigAccessKeyID),
		viper.GetString(constants.ConfigSecretAccessKey),
	)))

	ps = append(ps, qingstor.WithEnableVirtualStyle(viper.GetBool(constants.ConfigEnableVirtualStyle)))
	ps = append(ps, qingstor.WithDisableURICleaning(viper.GetBool(constants.ConfigDisableURICleaning)))

	if zone := viper.GetString(constants.ConfigZone); zone != "" {
		ps = append(ps, pairs.WithLocation(zone))
	}

	// add endpoint by different protocol https/http
	switch protocol := viper.GetString(constants.ConfigProtocol); protocol {
	case endpoint.ProtocolHTTPS:
		ps = append(ps, pairs.WithEndpoint(
			endpoint.NewHTTPS(
				viper.GetString(constants.ConfigHost),
				viper.GetInt(constants.ConfigPort),
			),
		))
	default: // endpoint.ProtocolHTTP:
		ps = append(ps, pairs.WithEndpoint(
			endpoint.NewHTTP(
				viper.GetString(constants.ConfigHost),
				viper.GetInt(constants.ConfigPort),
			),
		))
	}
	return ps
}

// IsQsPath check whether a path is qingstor path
func IsQsPath(s string) bool {
	return strings.HasPrefix(s, "qs://")
}
