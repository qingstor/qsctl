package action

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/c2h5oh/datasize"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/utils"
)

// ListObjects will handle all ls actions.
func ListObjects(remote string) (err error) {
	bucketName, objectKey, err := ParseQsPath(remote)
	if err != nil {
		return err
	}

	err = contexts.Storage.SetupBucket(bucketName, contexts.Zone)
	if err != nil {
		return
	}

	// Setting delimiter to "/" will emulate visiting as directory structure (not recursively)
	delimiter := "/"
	if contexts.Recursive {
		delimiter = ""
	}

	if contexts.LongFormat {
		return listObjectsLong(objectKey, delimiter)
	}
	return listObjects(objectKey, delimiter)
}

// listObjects list objects with specific prefix and delimiter from a bucket.
func listObjects(prefix, delimiter string) (err error) {
	oms, err := contexts.Storage.ListObjects(prefix, delimiter, nil)
	if err != nil {
		return
	}
	for _, om := range oms {
		fmt.Println(om.Key)
	}
	return nil
}

// listObjects list objects in long format with specific prefix and delimiter from a bucket.
func listObjectsLong(prefix, delimiter string) (err error) {
	oms, err := contexts.Storage.ListObjects(prefix, delimiter, nil)
	if err != nil {
		return
	}
	curUser, err := getBucketOwner()
	if err != nil {
		return err
	}
	res := make([][]string, 0, len(oms))
	var (
		acl         string
		size        string
		lasModified string
		contentNum  int
	)
	for _, om := range oms {
		if om.IsDir() {
			acl = constants.ACLDirectory
			lasModified = "" // directory will not show last modified time
			contentNum, err = getDirectoryObjCount(om.Key)
			if err != nil {
				return err
			}
		} else {
			acl = constants.ACLObject
			lasModified = om.FormatLastModified(constants.LsDefaultFormat) // format time
			contentNum = 1                                                 // object will show content num only 1
		}

		if contexts.HumanReadable {
			// if human readable flag true, print size as human readable format
			size, err = utils.UnixReadableSize(datasize.ByteSize(om.ContentLength).HR())
			if err != nil {
				return err
			}
		} else {
			// otherwise print size by bytes
			size = strconv.FormatInt(om.ContentLength, 10)
		}
		res = append(res, []string{acl, strconv.Itoa(contentNum), curUser, curUser, size, lasModified, om.Key})
	}

	// align the result and print
	res = utils.AlignLinux(res...)
	fmt.Println("total", len(res))
	for _, line := range res {
		fmt.Println(strings.Join(line, " "))
	}
	return
}

// getBucketOwner will return the owner id of current bucket
func getBucketOwner() (name string, err error) {
	ar, err := contexts.Storage.GetBucketACL()
	if err != nil {
		return "", err
	}
	return ar.OwnerID, nil
}

// getDirectoryObjCount will return objects count from a specific directory (prefix actually)
func getDirectoryObjCount(key string) (count int, err error) {
	oms, err := contexts.Storage.ListObjects(key, "/", nil)
	if err != nil {
		return 0, err
	}
	return len(oms), nil
}
