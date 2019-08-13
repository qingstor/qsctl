package storage

import (
	"io"
	"time"

	"github.com/yunify/qsctl/v2/constants"
)

// ObjectMeta is the metadata for an object.
type ObjectMeta struct {
	Key string

	ContentLength int64
	ContentType   string
	ETag          string
	LastModified  time.Time
	StorageClass  string

	Children []*ObjectMeta
}

// ACLResp is the response struct for acl request.
type ACLResp struct {
	OwnerID   string
	OwnerName string
	ACLs      []*ACLMeta
}

// ACLMeta is the metadata for acl info.
type ACLMeta struct {
	GranteeType string
	GranteeID   string
	GranteeName string
	Permission  string
}

// ObjectStorage is the interface to communicate with object storage service.
type ObjectStorage interface {
	SetupBucket(name, zone string) (err error)
	ListBuckets(zone string) (buckets []string, err error)
	PutBucket() error
	DeleteBucket() error
	GetBucketACL() (ar *ACLResp, err error)

	DeleteObject(objectKey string) (err error)
	HeadObject(objectKey string) (om *ObjectMeta, err error)
	GetObject(objectKey string) (r io.Reader, err error)
	ListObjects(prefix, delimiter string, marker *string) (om []*ObjectMeta, err error)
	DeleteMultipleObjects(prefix string) (err error)

	InitiateMultipartUpload(objectKey string) (uploadID string, err error)
	UploadMultipart(objectKey, uploadID string, size int64, partNumber int, md5sum []byte, r io.Reader) (err error)
	CompleteMultipartUpload(objectKey, uploadID string, totalNumber int) (err error)
}

// FormatLastModified transfer last modified from time to formatted string
func (om *ObjectMeta) FormatLastModified(format string) string {
	zero := time.Time{}
	if om.LastModified == zero {
		return ""
	}
	return om.LastModified.Format(format)
}

// IsDir will return whether the obj is a directory
func (om ObjectMeta) IsDir() bool {
	return om.ContentType == constants.DirectoryContentType
}

// Equal will return whether the obj's is equal with the given key
func (om ObjectMeta) Equal(key string) bool {
	return om.Key == key
}

// ChildrenCount will return all children objs count with the root of om
func (om *ObjectMeta) ChildrenCount() int {
	total := len(om.Children)
	for _, obj := range om.Children {
		total += obj.ChildrenCount()
	}
	return total
}
