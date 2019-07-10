package storage

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/yunify/qsctl/v2/constants"
)

// MockObjectStorage will implement ObjectStorage interface.
type MockObjectStorage struct {
	meta      map[string]*ObjectMeta
	multipart map[string]*multipart
}

type multipart struct {
	Length int64
	Parts  []int
}

// NewMockObjectStorage will create a new mock object storage.
func NewMockObjectStorage() *MockObjectStorage {
	return &MockObjectStorage{
		meta:      make(map[string]*ObjectMeta),
		multipart: make(map[string]*multipart),
	}
}

// SetupBucket implements ObjectStorage.SetupBucket
func (m *MockObjectStorage) SetupBucket(bucketName, zone string) error {
	return nil
}

// PutBucket implements ObjectStorage.PutBucket
func (m *MockObjectStorage) PutBucket() error {
	return nil
}

// HeadObject implements ObjectStorage.HeadObject
func (m *MockObjectStorage) HeadObject(objectKey string) (om *ObjectMeta, err error) {
	if om, ok := m.meta[objectKey]; ok {
		return om, nil
	}
	return nil, constants.ErrorQsPathNotFound
}

// GetObject implements ObjectStorage.GetObject
func (m *MockObjectStorage) GetObject(objectKey string) (r io.Reader, err error) {
	om, ok := m.meta[objectKey]
	if !ok {
		return nil, constants.ErrorQsPathNotFound
	}

	r = io.LimitReader(rand.Reader, om.ContentLength)
	return
}

// InitiateMultipartUpload implements ObjectStorage.InitiateMultipartUpload
func (m *MockObjectStorage) InitiateMultipartUpload(objectKey string) (uploadID string, err error) {
	id := uuid.New().String()
	m.multipart[objectKey] = &multipart{
		Parts: make([]int, 0),
	}
	return id, nil
}

// UploadMultipart implements ObjectStorage.UploadMultipart
func (m *MockObjectStorage) UploadMultipart(objectKey, uploadID string, size int64, partNumber int, md5sum []byte, r io.Reader) (err error) {
	_, ok := m.multipart[objectKey]
	if !ok {
		return constants.ErrorQsPathNotFound
	}

	h := md5.New()

	n, err := io.Copy(h, r)
	if err != nil {
		panic(err)
	}
	if n != size {
		return fmt.Errorf("content length is not match, expected %d, actual %d", size, n)
	}
	realMD5 := h.Sum(nil)
	if bytes.Compare(realMD5, md5sum) != 0 {
		return fmt.Errorf("content md5 is not match, expected %s, actual %s", md5sum, realMD5)
	}

	m.multipart[objectKey].Length += size
	m.multipart[objectKey].Parts = append(m.multipart[objectKey].Parts, partNumber)
	return nil
}

// CompleteMultipartUpload implements ObjectStorage.CompleteMultipartUpload
func (m *MockObjectStorage) CompleteMultipartUpload(objectKey, uploadID string, totalNumber int) (err error) {
	mp, ok := m.multipart[objectKey]
	if !ok {
		return constants.ErrorQsPathNotFound
	}

	if len(mp.Parts) != totalNumber {
		return fmt.Errorf("parts length is not match, expected %d, actual %d", totalNumber, len(mp.Parts))
	}

	m.meta[objectKey] = &ObjectMeta{
		Key:           objectKey,
		ContentLength: mp.Length,
		LastModified:  time.Now(),
	}
	delete(m.multipart, objectKey)
	return nil
}

// DeleteObject implements ObjectStorage.DeleteObject
func (m *MockObjectStorage) DeleteObject(objectKey string) (err error) {
	_, ok := m.meta[objectKey]
	if !ok {
		return constants.ErrorQsPathNotFound
	}

	delete(m.meta, objectKey)
	return nil
}
