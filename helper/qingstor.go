package helper

import (
	"encoding/hex"
	"io"
	"time"

	"github.com/pengsrc/go-shared/convert"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/yunify/qsctl/constants"
	"github.com/yunify/qsctl/contexts"
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

// ListObjects will list object by prefix.
func ListObjects(prefix string) (
	keys chan service.KeyType, err error,
) {
	keys = make(chan service.KeyType)
	defer close(keys)

	marker := ""
	for {
		resp, err := contexts.Bucket.ListObjects(&service.ListObjectsInput{
			Limit:  service.Int(200),
			Marker: service.String(marker),
			Prefix: service.String(prefix),
		})
		if err != nil {
			log.Errorf("prefix %s ListObjects failed [%v]", prefix, err)
			return nil, err
		}

		for _, v := range resp.Keys {
			if *v.MimeType == constants.DirectoryContentType {
				continue
			}

			keys <- *v
		}

		marker = *resp.NextMarker

		if marker == "" {
			break
		}
	}

	return
}

// HeadObject will head object.
func HeadObject(objectKey string) (om *ObjectMeta, err error) {
	resp, err := contexts.Bucket.HeadObject(objectKey, nil)
	// TODO: handle not found here.
	if err != nil {
		return
	}

	om = &ObjectMeta{
		Key:           objectKey,
		ContentLength: convert.Int64Value(resp.ContentLength),
		ContentType:   convert.StringValue(resp.ContentType),
		ETag:          convert.StringValue(resp.ETag),
		LastModified:  convert.TimeValue(resp.LastModified),
		StorageClass:  convert.StringValue(resp.XQSStorageClass),
	}
	return
}

// InitiateMultipartUpload will initiate a multipart upload.
func InitiateMultipartUpload(objectKey string) (uploadID string, err error) {
	resp, err := contexts.Bucket.InitiateMultipartUpload(objectKey, nil)
	if err != nil {
		log.Errorf("object %s InitiateMultipartUpload failed [%v]", objectKey, err)
		return
	}

	uploadID = *resp.UploadID
	return
}

// UploadMultipart will upload a multipart.
func UploadMultipart(
	objectKey, uploadID string, size int64, partNumber int, md5sum [16]byte, r io.Reader,
) (err error) {
	_, err = contexts.Bucket.UploadMultipart(objectKey, &service.UploadMultipartInput{
		Body:          r,
		ContentLength: convert.Int64(size),
		UploadID:      convert.String(uploadID),
		PartNumber:    convert.Int(partNumber),
		ContentMD5:    convert.String(hex.EncodeToString(md5sum[:])),
	})
	if err != nil {
		log.Errorf("object %s part %d UploadMultipart failed [%v]", objectKey, partNumber, err)
		return
	}
	return
}

// CompleteMultipartUpload will complete a multipart upload.
func CompleteMultipartUpload(objectKey, uploadID string, totalNumber int) (err error) {
	parts := make([]*service.ObjectPartType, totalNumber)
	for i := 0; i < totalNumber; i++ {
		parts[i] = &service.ObjectPartType{
			PartNumber: convert.Int(i),
		}
	}

	_, err = contexts.Bucket.CompleteMultipartUpload(
		objectKey, &service.CompleteMultipartUploadInput{
			UploadID:    convert.String(uploadID),
			ObjectParts: parts,
		})
	if err != nil {
		log.Errorf("object %s CompleteMultipartUpload failed [%v]", objectKey, err)
		return err
	}
	return nil
}

// GetObject will get an object.
func GetObject(objectKey string) (r io.Reader, err error) {
	resp, err := contexts.Bucket.GetObject(objectKey, nil)
	if err != nil {
		log.Errorf("object %s GetObject failed [%v]", objectKey, err)
		return nil, err
	}
	return resp.Body, nil
}
