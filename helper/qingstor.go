package helper

import (
	"encoding/hex"
	"io"

	"github.com/pengsrc/go-shared/convert"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/yunify/qsctl/constants"
	"github.com/yunify/qsctl/contexts"
)

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

// InitiateMultipartUpload will initiate a multipart upload.
func InitiateMultipartUpload(objectKey string) (uploadID string, err error) {
	resp, err := contexts.Bucket.InitiateMultipartUpload(objectKey, nil)
	if err != nil {
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
		return err
	}
	return nil
}

// calculatePartSize will calculate the object's part size.
func calculatePartSize(size int64) (partSize int64, err error) {
	partSize = constants.DefaultPartSize

	if size > constants.MaximumObjectSize {
		err = constants.ErrorFileTooLarge
		return
	}

	for size/partSize >= int64(constants.MaximumMultipartNumber) {
		if partSize < constants.MaximumAutoMultipartSize {
			partSize = partSize << 1
			continue
		}
		// Try to adjust partSize if it is too small and account for
		// integer division truncation.
		partSize = size/int64(constants.MaximumMultipartNumber) + 1
		break
	}
	return
}
