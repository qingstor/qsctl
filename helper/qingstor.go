package helper

import (
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
