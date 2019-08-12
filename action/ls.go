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

// ListHandler is all params for List func
type ListHandler struct {
	BaseHandler
	// Remote is the remote qs path
	Remote string `json:"remote"`
	// Prefix is the prefix to list
	Prefix string `json:"prefix"`
	// Delimiter puts all keys that share a common prefix into a list
	Delimiter string `json:"delimiter"`
	// Root is the root node of the file tree structure
	Root *storage.ObjectMeta `json:"root"`
}

// WithHumanReadable rewrite the WithHumanReadable method
func (lh *ListHandler) WithHumanReadable(h bool) *ListHandler {
	lh.HumanReadable = h
	return lh
}

// WithLongFormat rewrite the WithLongFormat method
func (lh *ListHandler) WithLongFormat(l bool) *ListHandler {
	lh.LongFormat = l
	return lh
}

// WithRecursive rewrite the WithRecursive method
func (lh *ListHandler) WithRecursive(r bool) *ListHandler {
	lh.Recursive = r
	return lh
}

// WithReverse rewrite the WithReverse method
func (lh *ListHandler) WithReverse(r bool) *ListHandler {
	lh.Reverse = r
	return lh
}

// WithZone rewrite the WithZone method
func (lh *ListHandler) WithZone(z string) *ListHandler {
	lh.Zone = z
	return lh
}

// WithRemote sets the Remote field with given remote
func (lh *ListHandler) WithRemote(remote string) *ListHandler {
	lh.Remote = remote
	return lh
}

// WithPrefix sets the Prefix field with given prefix
func (lh *ListHandler) WithPrefix(prefix string) *ListHandler {
	lh.Prefix = prefix
	return lh
}

// WithDelimiter sets the Delimiter field with given delimiter
func (lh *ListHandler) WithDelimiter(delimiter string) *ListHandler {
	lh.Delimiter = delimiter
	return lh
}

// WithRoot sets the Root field with given om
func (lh *ListHandler) WithRoot(om *storage.ObjectMeta) *ListHandler {
	lh.Root = om
	return lh
}

// ListObjects will handle all ls actions.
func (lh *ListHandler) ListObjects() (err error) {
	bucketName, objectKey, err := ParseQsPath(lh.Remote)
	if err != nil {
		return err
	}

	err = contexts.Storage.SetupBucket(bucketName, lh.Zone)
	if err != nil {
		return
	}
	// Setting delimiter to "/" will emulate visiting as directory structure (not recursively for next level)
	// construct the object tree
	root, err := lh.WithPrefix(objectKey).WithDelimiter("/").listObjects()
	if err != nil {
		return err
	}

	// if long format (-l), set bucket owner for printing
	if lh.LongFormat {
		if err = getBucketOwner(); err != nil {
			return err
		}
	}
	// print first level children keys
	if err = lh.WithRoot(root).printChildrenKeys(); err != nil {
		return err
	}

	// if recursive (-R), print next level keys recursively
	if lh.Recursive {
		for _, om := range root.Children {
			if om.IsDir() {
				if err := lh.WithRoot(om).printChildrenKeysRecursively(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// listObjects list objects with specific prefix and delimiter from a bucket,
// return the root of object tree.
func (lh *ListHandler) listObjects() (root *storage.ObjectMeta, err error) {
	oms, err := contexts.Storage.ListObjects(lh.Prefix, lh.Delimiter, nil)
	if err != nil {
		return
	}
	root = &storage.ObjectMeta{
		Key: lh.Prefix,
	}
	// if prefix end with "/", handle it as a directory
	if strings.HasSuffix(lh.Prefix, "/") {
		root.ContentType = constants.DirectoryContentType
	}

	// append children to root
	for _, om := range oms {
		// if om is a dir and same with the prefix, not add as a child
		// this is because qs will return the same with prefix as the object key,
		// which should not be considered as the expected child.
		if om.IsDir() && om.Equal(lh.Prefix) {
			continue
		}
		root.Children = append(root.Children, om)
	}

	// if not recursive (-R) and not long-format (-l), stop here and return.
	if !lh.Recursive && !lh.LongFormat {
		return
	}

	var once bool
	// if long-format (-l) and not recursive (-R), list only one more level for counting contentNum
	if lh.LongFormat && !lh.Recursive {
		once = true
	}
	// recursively list keys appended from each dir
	for _, om := range root.Children {
		// cuz all children oms are not same with the prefix,
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
func (lh *ListHandler) printChildrenKeys() (err error) {
	// Get root from handler
	root := lh.Root
	// if no children, return
	if root.Children == nil {
		return
	}
	lh.sortChildren()

	// if not long-format (-l), only print key
	if !lh.LongFormat {
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
		line, err := lh.WithRoot(om).omInfoSlice()
		if err != nil {
			return err
		}
		// before append this line into res, append key to the end of line
		res = append(res, append(line, key))
	}

	// print total
	if lh.HumanReadable {
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
func (lh *ListHandler) printChildrenKeysRecursively() (err error) {
	// Get params from handler
	root := lh.Root

	dirKey := root.Key
	fmt.Println()
	fmt.Printf("%s:\n", dirKey)
	if err = lh.printChildrenKeys(); err != nil {
		return err
	}

	for _, om := range root.Children {
		if om.IsDir() && !om.Equal(dirKey) {
			if err = lh.WithRoot(om).printChildrenKeysRecursively(); err != nil {
				return err
			}
		}
	}
	return nil
}

// sortChildren sort the oms slice by reverse flag
// if true, desc; if false, asc (default)
func (lh *ListHandler) sortChildren() {
	// Get oms from handler
	oms := lh.Root.Children
	sort.Slice(oms, func(i, j int) bool {
		if lh.Reverse {
			return oms[i].Key > oms[j].Key
		}
		return oms[i].Key < oms[j].Key
	})
}

// omInfoSlice returns the om detail info slice
func (lh *ListHandler) omInfoSlice() (line []string, err error) {
	// Get root from handler
	om := lh.Root
	// if om is a dir, set size to 0 and last modified blank
	if om.IsDir() {
		contentNum := 0
		if om.Children != nil {
			contentNum = len(om.Children)
		}
		return []string{constants.ACLDirectory, strconv.Itoa(contentNum), ownerID, ownerID, "0", ""}, nil
	}
	size := ""
	if lh.HumanReadable {
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
