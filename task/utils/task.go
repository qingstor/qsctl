package utils

import (
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/constants"
)

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
	pattern := "^(?:qs://)?([a-z\\d][a-z-\\d]{4,61}[a-z\\d])?(.*)?$"

	x := regexp.MustCompile(pattern).FindStringSubmatch(p)
	if len(x) != 3 || x[1] == "" {
		return constants.KeyTypeInvalid, "", "", constants.ErrorQsPathInvalid
	}

	bucketName, objectKey = x[1], x[2]

	// TODO: add bucket name and object key check here.

	if objectKey == "" || objectKey == "/" {
		return constants.KeyTypeBucket, bucketName, objectKey, nil
	}

	if strings.HasSuffix(p, "/") {
		return constants.KeyTypePseudoDir, bucketName, objectKey, nil
	}
	return constants.KeyTypeObject, bucketName, objectKey, nil
}
