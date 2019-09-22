package storage

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yunify/qsctl/v2/constants"
)

// Available preset object for mock storage.
const (
	Mock0BObject = "0b"
	MockMBObject = "mb"
	MockGBObject = "gb"
	MockTBObject = "tb"

	MockZoneAlpha = "mock-alpha"
	MockZoneBeta  = "mock-beta"
)

// MockObjectStorage will implement ObjectStorage interface.
type MockObjectStorage struct {
	Meta          map[string]*ObjectMeta
	Multipart     map[string]*multipart
	Buckets       map[string]*bucketMeta
	currentBucket *bucketMeta
}

type multipart struct {
	Length int64
	Parts  []int
}

type bucketMeta struct {
	Created  time.Time
	Location string
	Name     string
	URL      string
	OwnerID  string
}

// NewMockObjectStorage will create a new mock object storage.
func NewMockObjectStorage() *MockObjectStorage {
	s := &MockObjectStorage{
		Meta:      make(map[string]*ObjectMeta),
		Multipart: make(map[string]*multipart),
		Buckets:   make(map[string]*bucketMeta),
	}

	// Adding persist keys.
	s.Meta[Mock0BObject] = &ObjectMeta{
		Key:           Mock0BObject,
		ContentLength: 0,
		LastModified:  time.Unix(612889200, 0),
	}
	s.Meta[MockMBObject] = &ObjectMeta{
		Key:           MockMBObject,
		ContentLength: 1024 * 1024,
		LastModified:  time.Unix(612889200, 0),
	}
	s.Meta[MockGBObject] = &ObjectMeta{
		Key:           MockGBObject,
		ContentLength: 1024 * 1024 * 1024,
		LastModified:  time.Unix(612889200, 0),
	}
	s.Meta[MockTBObject] = &ObjectMeta{
		Key:           MockTBObject,
		ContentLength: 1024 * 1024 * 1024 * 1024,
		LastModified:  time.Unix(612889200, 0),
	}

	// Adding test Buckets.
	s.Buckets[MockZoneAlpha] = &bucketMeta{
		Name:     MockZoneAlpha,
		Created:  time.Unix(612889200, 0),
		Location: MockZoneAlpha,
		OwnerID:  MockZoneAlpha + "user",
	}
	s.Buckets[MockZoneBeta] = &bucketMeta{
		Name:     MockZoneBeta,
		Created:  time.Unix(612889200, 0),
		Location: MockZoneBeta,
		OwnerID:  MockZoneBeta + "user",
	}
	return s
}

// SetupBucket implements ObjectStorage.SetupBucket
func (m *MockObjectStorage) SetupBucket(bucketName, zone string) error {
	if zone != "" {
		m.currentBucket = &bucketMeta{
			Name:     bucketName,
			Created:  time.Unix(612889200, 0),
			Location: zone,
			OwnerID:  zone + "user",
		}
		return nil
	}
	m.currentBucket = m.Buckets[bucketName]
	return nil
}

// PutBucket implements ObjectStorage.PutBucket
func (m *MockObjectStorage) PutBucket() error {
	if _, ok := m.Buckets[m.currentBucket.Name]; ok {
		return constants.ErrorBucketAlreadyExists
	}
	m.Buckets[m.currentBucket.Name] = m.currentBucket
	return nil
}

// DeleteBucket implements ObjectStorage.DeleteBucket
func (m *MockObjectStorage) DeleteBucket() error {
	delete(m.Buckets, m.currentBucket.Name)
	return nil
}

// HeadObject implements ObjectStorage.HeadObject
func (m *MockObjectStorage) HeadObject(objectKey string) (om *ObjectMeta, err error) {
	if om, ok := m.Meta[objectKey]; ok {
		return om, nil
	}
	return nil, constants.ErrorQsPathNotFound
}

// GetObject implements ObjectStorage.GetObject
func (m *MockObjectStorage) GetObject(objectKey string) (r io.Reader, err error) {
	om, ok := m.Meta[objectKey]
	if !ok {
		return nil, constants.ErrorQsPathNotFound
	}

	r = io.LimitReader(rand.Reader, om.ContentLength)
	return
}

// InitiateMultipartUpload implements ObjectStorage.InitiateMultipartUpload
func (m *MockObjectStorage) InitiateMultipartUpload(objectKey string) (uploadID string, err error) {
	id := uuid.New().String()
	m.Multipart[objectKey] = &multipart{
		Parts: make([]int, 0),
	}
	return id, nil
}

// UploadMultipart implements ObjectStorage.UploadMultipart
func (m *MockObjectStorage) UploadMultipart(objectKey, uploadID string, size int64, partNumber int, md5sum []byte, r io.Reader) (err error) {
	_, ok := m.Multipart[objectKey]
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

	m.Multipart[objectKey].Length += size
	m.Multipart[objectKey].Parts = append(m.Multipart[objectKey].Parts, partNumber)
	return nil
}

// CompleteMultipartUpload implements ObjectStorage.CompleteMultipartUpload
func (m *MockObjectStorage) CompleteMultipartUpload(objectKey, uploadID string, totalNumber int) (err error) {
	mp, ok := m.Multipart[objectKey]
	if !ok {
		return constants.ErrorQsPathNotFound
	}

	if len(mp.Parts) != totalNumber {
		return fmt.Errorf("parts length is not match, expected %d, actual %d", totalNumber, len(mp.Parts))
	}

	m.Meta[objectKey] = &ObjectMeta{
		Key:           objectKey,
		ContentLength: mp.Length,
		LastModified:  time.Now(),
	}
	delete(m.Multipart, objectKey)
	return nil
}

// DeleteObject implements ObjectStorage.DeleteObject
func (m *MockObjectStorage) DeleteObject(objectKey string) (err error) {
	_, ok := m.Meta[objectKey]
	if !ok {
		return constants.ErrorQsPathNotFound
	}

	delete(m.Meta, objectKey)
	return nil
}

// ListBuckets implements ObjectStorage.ListBuckets
func (m *MockObjectStorage) ListBuckets(zone string) (buckets []string, err error) {
	buckets = make([]string, 0, len(m.Buckets))
	for _, b := range m.Buckets {
		if zone == "" || b.Location == zone {
			buckets = append(buckets, b.Name)
		}
	}
	return
}

// ListObjects implements ObjectStorage.ListObjects
func (m *MockObjectStorage) ListObjects(prefix, delimiter string, marker *string) (om []*ObjectMeta, err error) {
	om = make([]*ObjectMeta, 0)
	// delimiter blank means no directory concept
	if delimiter == "" {
		for k, obj := range m.Meta {
			if strings.HasPrefix(k, prefix) {
				om = append(om, obj)
			}
		}
		return
	}
	for _, obj := range m.Meta {
		// Determine whether obj is the sub-object of prefix
		// obj.Key must start with prefix
		if strings.HasPrefix(obj.Key, prefix) &&
			// obj.Key contains as same amount as prefix (corresponding to obj is not a directory) OR
			(strings.Count(obj.Key, "/") == strings.Count(prefix, "/") ||
				// obj.Key contains one more '/' than prefix AND end with '/' (corresponding to obj is a directory)
				(strings.Count(obj.Key, "/") == strings.Count(prefix, "/")+1 &&
					strings.HasSuffix(obj.Key, "/"))) {
			om = append(om, obj)
		}
	}
	return
}

// ResetMockObjects reset mock objects with specific prefix for test
func (m *MockObjectStorage) ResetMockObjects(prefix string, num int) {
	dirKey := prefix + "/"
	// obj/
	m.Meta[dirKey] = &ObjectMeta{
		Key:         dirKey,
		ContentType: constants.DirectoryContentType,
	}
	// obj/obj/
	m.Meta[dirKey+dirKey] = &ObjectMeta{
		Key:         dirKey + dirKey,
		ContentType: constants.DirectoryContentType,
	}
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("%s_%d", prefix, i)
		// obj_0 ... obj_19
		m.Meta[key] = &ObjectMeta{
			Key:           key,
			ContentLength: int64(i * 1024),
		}
		// obj/obj_0/ ... obj/obj_19
		secondLvlDir := fmt.Sprintf("%s/%s_%d/", prefix, prefix, i)
		m.Meta[secondLvlDir] = &ObjectMeta{
			Key:           secondLvlDir,
			ContentLength: int64(0),
			ContentType:   constants.DirectoryContentType,
		}
		// obj/obj_0/obj ... obj/obj_19/obj
		secondLvlKey := fmt.Sprintf("%s%s", secondLvlDir, prefix)
		m.Meta[secondLvlKey] = &ObjectMeta{
			Key:           secondLvlKey,
			ContentLength: int64(i * 1024),
		}
		// obj/obj/obj_0/ ... obj/obj/obj_19/
		thirdLvlDir := fmt.Sprintf("%s/%s/%s_%d/", prefix, prefix, prefix, i)
		m.Meta[thirdLvlDir] = &ObjectMeta{
			Key:           thirdLvlDir,
			ContentLength: int64(0),
			ContentType:   constants.DirectoryContentType,
		}
		// obj/obj/obj_0/obj ... obj/obj/obj_19/obj
		thirdLvlKey := fmt.Sprintf("%s%s", thirdLvlDir, prefix)
		m.Meta[thirdLvlKey] = &ObjectMeta{
			Key:           thirdLvlKey,
			ContentLength: int64(i * 1024),
		}
	}
}

// GetBucketACL implements ObjectStorage.GetBucketACL
func (m *MockObjectStorage) GetBucketACL() (ar *ACLResp, err error) {
	return &ACLResp{
		OwnerID: m.currentBucket.OwnerID,
	}, nil
}

// PutObject will put a object.
func (m *MockObjectStorage) PutObject(objectKey string, md5sum []byte, r io.Reader) (err error) {
	h := md5.New()

	n, err := io.Copy(h, r)
	if err != nil {
		panic(err)
	}
	realMD5 := h.Sum(nil)
	if bytes.Compare(realMD5, md5sum) != 0 {
		return fmt.Errorf("content md5 is not match, expected %s, actual %s", md5sum, realMD5)
	}

	m.Meta[objectKey] = &ObjectMeta{
		Key:           objectKey,
		ContentLength: n,
	}
	return nil
}
