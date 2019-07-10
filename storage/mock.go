package storage

import "io"

// MockObjectStorage will implement ObjectStorage interface.
type MockObjectStorage struct {
}

// SetupBucket implements ObjectStorage.SetupBucket
func (m *MockObjectStorage) SetupBucket(bucketName, zone string) error {
	panic("implement me")
}

// PutBucket implements ObjectStorage.PutBucket
func (m *MockObjectStorage) PutBucket(bucketName, zone string) error {
	panic("implement me")
}

// HeadObject implements ObjectStorage.HeadObject
func (m *MockObjectStorage) HeadObject(objectKey string) (om *ObjectMeta, err error) {
	panic("implement me")
}

// GetObject implements ObjectStorage.GetObject
func (m *MockObjectStorage) GetObject(objectKey string) (r io.Reader, err error) {
	panic("implement me")
}

// InitiateMultipartUpload implements ObjectStorage.InitiateMultipartUpload
func (m *MockObjectStorage) InitiateMultipartUpload(objectKey string) (uploadID string, err error) {
	panic("implement me")
}

// UploadMultipart implements ObjectStorage.UploadMultipart
func (m *MockObjectStorage) UploadMultipart(objectKey, uploadID string, size int64, partNumber int, md5sum [16]byte, r io.Reader) (err error) {
	panic("implement me")
}

// CompleteMultipartUpload implements ObjectStorage.CompleteMultipartUpload
func (m *MockObjectStorage) CompleteMultipartUpload(objectKey, uploadID string, totalNumber int) (err error) {
	panic("implement me")
}
