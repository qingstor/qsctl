package utils

import (
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/types"
)

// bucketNameRegexp is the bucket name regexp, which indicates:
// 1. length: 6-63;
// 2. contains lowercase letters, digits and strikethrough;
// 3. starts and ends with letter or digit.
var bucketNameRegexp = regexp.MustCompile(`^[a-z\d][a-z-\d]{4,61}[a-z\d]$`)

// ParseFlow will parse the data flow
func ParseFlow(src, dst string) (flow constants.FlowType) {
	if dst == "" {
		if strings.HasPrefix(src, "qs://") {
			return constants.FlowAtRemote
		}
		return constants.FlowAtLocal
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

// ParsePath will parse a path into different path type.
func ParsePath(p string) (pathType constants.PathType, err error) {
	// Use - means we will read from stdin.
	if p == "-" {
		return constants.PathTypeStream, nil
	}

	fi, err := os.Stat(p)
	if os.IsNotExist(err) {
		log.Infof("File <%s> is not exist, please check your input", p)
		return constants.PathTypeInvalid, constants.ErrorFileNotExist
	}
	if err != nil {
		log.Errorf("Stat file failed [%s]", err)
		return
	}
	if fi.IsDir() {
		return constants.PathTypeLocalDir, nil
	}
	return constants.PathTypeFile, nil
}

// ParseKey will parse a key into different key type.
func ParseKey(p string) (keyType constants.KeyType, bucketName, objectKey string, err error) {
	// qs-path includes three part: "qs://" prefix, bucket name and object key.
	// "qs://" prefix could be emit.
	if strings.HasPrefix(p, "qs://") {
		p = p[5:]
	}
	s := strings.SplitN(p, "/", 2)

	if !IsValidBucketName(s[0]) {
		return constants.KeyTypeInvalid, "", "", constants.ErrorQsPathInvalid
	}

	// Only have bucket name or object key is "/"
	// For example: "qs://testbucket/"

	if len(s) == 1 || s[1] == "" {
		return constants.KeyTypeBucket, s[0], "", nil
	}

	if strings.HasSuffix(p, "/") {
		return constants.KeyTypePseudoDir, s[0], s[1], nil
	}
	return constants.KeyTypeObject, s[0], s[1], nil
}

// IsValidBucketName will check whether given string is a valid bucket name.
func IsValidBucketName(s string) bool {
	return bucketNameRegexp.MatchString(s)
}

// ParseInput will parse two args into flow, path and key.
func ParseInput(t interface {
	types.FlowTypeSetter
	types.PathSetter
	types.StreamSetter
	types.PathTypeSetter
	types.KeySetter
	types.KeyTypeSetter
	types.BucketNameSetter
}, src, dst string) (err error) {
	flow := ParseFlow(src, dst)
	t.SetFlowType(flow)

	var path string
	var pathType constants.PathType
	switch flow {
	case constants.FlowToRemote:
		pathType, err = ParsePath(src)
		if err != nil {
			return err
		}
		t.SetPathType(pathType)
		path = src

		keyType, bucketName, objectKey, err := ParseKey(dst)
		if err != nil {
			return err
		}
		t.SetKeyType(keyType)
		t.SetKey(objectKey)
		t.SetBucketName(bucketName)
	case constants.FlowToLocal, constants.FlowAtRemote:
		pathType, err = ParsePath(dst)
		if err != nil {
			return err
		}
		t.SetPathType(pathType)
		path = dst

		keyType, bucketName, objectKey, err := ParseKey(src)
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
