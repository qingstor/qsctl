package storage

import (
	"io"
	"time"
)

// ObjectMeta is the metadata for an object.
type ObjectMeta struct {
	Key string

	ContentLength int64
	ContentType   string
	ETag          string
	LastModified  time.Time
	StorageClass  string
}

// ObjectStorage is the interface to communicate with object storage service.
type ObjectStorage interface {
	SetupBucket(name, zone string) (err error)
	PutBucket() error

	DeleteObject(objectKey string) (err error)
	HeadObject(objectKey string) (om *ObjectMeta, err error)
	GetObject(objectKey string) (r io.Reader, err error)

	InitiateMultipartUpload(objectKey string) (uploadID string, err error)
	UploadMultipart(objectKey, uploadID string, size int64, partNumber int, md5sum [16]byte, r io.Reader) (err error)
	CompleteMultipartUpload(objectKey, uploadID string, totalNumber int) (err error)
}
