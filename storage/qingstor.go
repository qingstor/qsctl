package storage

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pengsrc/go-shared/convert"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yunify/qingstor-sdk-go/v3/config"
	"github.com/yunify/qingstor-sdk-go/v3/request/errors"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/yunify/qsctl/v2/constants"
)

// QingStorObjectStorage will implement ObjectStorage interface.
type QingStorObjectStorage struct {
	service *service.Service
	bucket  *service.Bucket
}

// NewQingStorObjectStorage will create a new qingstor object storage.
func NewQingStorObjectStorage() (q *QingStorObjectStorage, err error) {
	cfg, err := config.New(
		viper.GetString(constants.ConfigAccessKeyID),
		viper.GetString(constants.ConfigSecretAccessKey),
	)
	if err != nil {
		log.Errorf("Init config failed [%v]", err)
		return
	}

	cfg.Host = viper.GetString(constants.ConfigHost)
	cfg.Port = viper.GetInt(constants.ConfigPort)
	cfg.Protocol = viper.GetString(constants.ConfigProtocol)
	cfg.ConnectionRetries = viper.GetInt(constants.ConfigConnectionRetries)
	cfg.LogLevel = viper.GetString(constants.ConfigLogLevel)

	q = &QingStorObjectStorage{}

	q.service, err = service.Init(cfg)
	if err != nil {
		log.Errorf("Init service failed [%v]", err)
		return
	}
	return
}

// SetupBucket implements ObjectStorage.SetupBucket
func (q *QingStorObjectStorage) SetupBucket(name, zone string) (err error) {
	if zone == "" {
		zone = viper.GetString(constants.ConfigZone)
	}

	if zone != "" {
		q.bucket, err = q.service.Bucket(name, zone)
		if err != nil {
			log.Errorf("Init bucket <%s> in zone <%s> failed [%v]", name, zone, err)
			return constants.ErrorExternalServiceError
		}
		return
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url := fmt.Sprintf("%s://%s.%s:%d",
		viper.GetString(constants.ConfigProtocol),
		name,
		viper.GetString(constants.ConfigHost),
		viper.GetInt(constants.ConfigPort))

	r, err := client.Head(url)
	if err != nil {
		log.Errorf("Head location failed [%v]", err)
		return constants.ErrorExternalServiceError
	}
	if r.StatusCode != http.StatusTemporaryRedirect {
		log.Infof("Detect bucket location failed, please check your input")
		return constants.ErrorQsPathNotFound
	}

	// Example URL: https://bucket.zone.qingstor.com
	zone = strings.Split(r.Header.Get("Location"), ".")[1]
	q.bucket, err = q.service.Bucket(name, zone)
	if err != nil {
		log.Errorf("Init bucket <%s> in zone <%s> failed [%v]", name, zone, err)
		return constants.ErrorExternalServiceError
	}
	return
}

// HeadObject will head object.
func (q *QingStorObjectStorage) HeadObject(objectKey string) (om *ObjectMeta, err error) {
	resp, err := q.bucket.HeadObject(objectKey, nil)
	if err != nil {
		if e, ok := err.(*errors.QingStorError); ok {
			if e.StatusCode == http.StatusNotFound {
				return nil, constants.ErrorQsPathNotFound
			} else if e.StatusCode == http.StatusForbidden {
				return nil, constants.ErrorQsPathAccessForbidden
			}
		}
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
func (q *QingStorObjectStorage) InitiateMultipartUpload(objectKey string) (uploadID string, err error) {
	resp, err := q.bucket.InitiateMultipartUpload(objectKey, nil)
	if err != nil {
		log.Errorf("Object <%s> InitiateMultipartUpload failed [%v]", objectKey, err)
		return
	}

	uploadID = *resp.UploadID
	return
}

// UploadMultipart will upload a multipart.
func (q *QingStorObjectStorage) UploadMultipart(
	objectKey, uploadID string, size int64, partNumber int, md5sum []byte, r io.Reader,
) (err error) {
	_, err = q.bucket.UploadMultipart(objectKey, &service.UploadMultipartInput{
		Body:          r,
		ContentLength: convert.Int64(size),
		UploadID:      convert.String(uploadID),
		PartNumber:    convert.Int(partNumber),
		ContentMD5:    convert.String(hex.EncodeToString(md5sum[:])),
	})
	if err != nil {
		log.Errorf("Object <%s> part <%d> UploadMultipart failed [%v]", objectKey, partNumber, err)
		return
	}
	return
}

// CompleteMultipartUpload will complete a multipart upload.
func (q *QingStorObjectStorage) CompleteMultipartUpload(objectKey, uploadID string, totalNumber int) (err error) {
	parts := make([]*service.ObjectPartType, totalNumber)
	for i := 0; i < totalNumber; i++ {
		parts[i] = &service.ObjectPartType{
			PartNumber: convert.Int(i),
		}
	}

	_, err = q.bucket.CompleteMultipartUpload(
		objectKey, &service.CompleteMultipartUploadInput{
			UploadID:    convert.String(uploadID),
			ObjectParts: parts,
		})
	if err != nil {
		log.Errorf("Object <%s> CompleteMultipartUpload failed [%v]", objectKey, err)
		return err
	}
	return nil
}

// GetObject will get an object.
func (q *QingStorObjectStorage) GetObject(objectKey string) (r io.Reader, err error) {
	resp, err := q.bucket.GetObject(objectKey, nil)
	if err != nil {
		log.Errorf("Object <%s> GetObject failed [%v]", objectKey, err)
		return nil, err
	}
	return resp.Body, nil
}

// PutBucket will make a bucket with specific name.
func (q *QingStorObjectStorage) PutBucket() error {
	// Request and create bucket
	_, err := q.bucket.Put()
	if err != nil {
		log.Errorf("Make bucket <%s> in zone <%s> failed [%v]",
			*q.bucket.Properties.BucketName, *q.bucket.Properties.Zone, err)
		return err
	}
	return nil
}

// DeleteObject will delete an object with specific key.
func (q *QingStorObjectStorage) DeleteObject(objectKey string) (err error) {
	if _, err = q.bucket.DeleteObject(objectKey); err != nil {
		log.Errorf("Delete object <%s> failed [%v]", objectKey, err)
		return
	}
	return nil
}
