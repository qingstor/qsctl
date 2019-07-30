package action

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/c2h5oh/datasize"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/utils"
)

// ownerID record the current bucket's owner ID
var ownerID string

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

	// Setting delimiter to "/" will emulate visiting as directory structure (not recursively for next level)
	delimiter := "/"

	// construct the object tree
	root, err := listObjects(objectKey, delimiter)
	if err != nil {
		return err
	}

	// if long format (-l), set bucket owner for printing
	if contexts.LongFormat {
		if err = getBucketOwner(); err != nil {
			return err
		}
	}
	// print first level children keys
	if err = printChildrenKeys(root); err != nil {
		return err
	}

	// if recursive (-R), print next level keys recursively
	if contexts.Recursive {
		for _, om := range root.Children {
			if om.IsDir() {
				if err := printChildrenKeysRecursively(om); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// listObjects list objects with specific prefix and delimiter from a bucket,
// return the root of object tree.
func listObjects(prefix, delimiter string) (root *storage.ObjectMeta, err error) {
	oms, err := contexts.Storage.ListObjects(prefix, delimiter, nil)
	if err != nil {
		return
	}
	root = &storage.ObjectMeta{
		Key: prefix,
	}
	// if prefix end with "/", handle it as a directory
	if strings.HasSuffix(prefix, "/") {
		root.ContentType = constants.DirectoryContentType
	}

	// append children to root
	for _, om := range oms {
		// if om is a dir and same with the prefix, not add as a child
		if om.IsDir() && om.Equal(prefix) {
			continue
		}
		root.Children = append(root.Children, om)
	}

	if !contexts.Recursive && !contexts.LongFormat {
		return
	}

	var once bool
	// if long-format (-l) and not recursive (-R), list only one more level for counting contentNum
	if contexts.LongFormat && !contexts.Recursive {
		once = true
	}
	// recursively list keys appended from each dir
	for _, om := range root.Children {
		// cuz all children om is not same with the prefix,
		// so we need only determine whether om is a dir
		if om.IsDir() {
			if err = recursiveListObjects(om, once); err != nil {
				return nil, err
			}
		}
	}
	return root, nil
}

// recursiveListObjects list objects recursively for each dir,
// if once is true, only recurse once, for contentNum count.
func recursiveListObjects(root *storage.ObjectMeta, once bool) error {
	oms, err := contexts.Storage.ListObjects(root.Key, "/", nil)
	if err != nil {
		return err
	}
	// for every om, if is dir and equal with root, not add as children
	for _, om := range oms {
		if om.IsDir() && om.Equal(root.Key) {
			continue
		}
		root.Children = append(root.Children, om)
	}
	// if recurse once set true, not list objects for next level
	if once {
		return nil
	}
	for _, om := range root.Children {
		// list children for dir next level
		if om.IsDir() {
			if err := recursiveListObjects(om, once); err != nil {
				return err
			}
		}
	}
	return nil
}

// getBucketOwner will assign the owner id of current bucket
func getBucketOwner() error {
	ar, err := contexts.Storage.GetBucketACL()
	if err != nil {
		return err
	}
	ownerID = ar.OwnerID
	return nil
}

// printChildrenKeys will handle the main logic of printing the children info of root
func printChildrenKeys(root *storage.ObjectMeta) (err error) {
	// if no children, return
	if root.Children == nil {
		return
	}
	sortOms(root.Children)

	// if not long-format (-l), only print key
	if !contexts.LongFormat {
		for _, om := range root.Children {
			// if root is dir, trim prefix, as well as suffix "/"
			if root.IsDir() {
				fmt.Printf("%s\n", strings.TrimSuffix(strings.TrimPrefix(om.Key, root.Key), "/"))
			} else {
				fmt.Printf("%s\n", om.Key)
			}
		}
		return nil
	}
	// if long-format (-l), print key's detail info
	res := make([][]string, 0)
	var total int64
	for _, om := range root.Children {
		total += om.ContentLength
		key := om.Key
		// if root is dir, trim prefix, as well as suffix "/"
		if root.IsDir() {
			key = strings.TrimSuffix(strings.TrimPrefix(om.Key, root.Key), "/")
		}
		// format this line
		line, err := omInfoSlice(om)
		if err != nil {
			return err
		}
		// before append this line into res, append key to the end of line
		res = append(res, append(line, key))
	}

	// print total
	if contexts.HumanReadable {
		totalSize, err := utils.UnixReadableSize(datasize.ByteSize(total).HR())
		if err != nil {
			return err
		}
		fmt.Println("total", totalSize)
	} else {
		fmt.Println("total", strconv.FormatInt(total, 10))
	}
	// align the result and print
	for _, line := range utils.AlignLinux(res...) {
		fmt.Println(strings.Join(line, " "))
	}
	return nil
}

// printChildrenKeysRecursively will recursively print keys
func printChildrenKeysRecursively(root *storage.ObjectMeta) (err error) {
	dirKey := root.Key
	fmt.Println()
	fmt.Printf("%s:\n", dirKey)

	if err = printChildrenKeys(root); err != nil {
		return err
	}

	for _, om := range root.Children {
		if om.IsDir() && !om.Equal(dirKey) {
			if err = printChildrenKeysRecursively(om); err != nil {
				return err
			}
		}
	}
	return nil
}

// sortOms sort the oms slice by contexts.Reverse
// if true, desc; if false, asc (default)
func sortOms(oms []*storage.ObjectMeta) {
	sort.Slice(oms, func(i, j int) bool {
		if contexts.Reverse {
			return oms[i].Key > oms[j].Key
		}
		return oms[i].Key < oms[j].Key
	})
}

// omInfoSlice returns the om detail info slice
func omInfoSlice(om *storage.ObjectMeta) (line []string, err error) {
	// if om is a dir, set size to 0 and last modified blank
	if om.IsDir() {
		contentNum := 0
		if om.Children != nil {
			contentNum = len(om.Children)
		}
		return []string{constants.ACLDirectory, strconv.Itoa(contentNum), ownerID, ownerID, "0", ""}, nil
	}
	size := ""
	if contexts.HumanReadable {
		// if human readable flag true, print size as human readable format
		size, err = utils.UnixReadableSize(datasize.ByteSize(om.ContentLength).HR())
		if err != nil {
			return nil, err
		}
	} else {
		// otherwise print size by bytes
		size = strconv.FormatInt(om.ContentLength, 10)
	}
	// if om is a obj, set content num to 1
	return []string{constants.ACLObject, "1", ownerID, ownerID, size,
		om.FormatLastModified(constants.LsDefaultFormat)}, nil
}
