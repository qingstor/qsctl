package helper

import (
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/yunify/qsctl/constants"
	"github.com/yunify/qsctl/contexts"
)

// ListObjects will list object by prefix.
func ListObjects(prefix string) (
	commonPrefix chan string, keys chan service.KeyType, err error,
) {
	bucket, _ := contexts.Service.Bucket("test", "test")

	commonPrefix = make(chan string)
	keys = make(chan service.KeyType)
	defer close(keys)
	defer close(commonPrefix)

	marker := ""
	for {
		resp, err := bucket.ListObjects(&service.ListObjectsInput{
			Limit:  service.Int(200),
			Marker: service.String(marker),
			Prefix: service.String(prefix),
		})
		if err != nil {
			return nil, nil, err
		}

		for _, v := range resp.CommonPrefixes {
			commonPrefix <- *v
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
