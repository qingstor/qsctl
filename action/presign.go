package action

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/task"
)

// PresignHandler is all params for Presign func
type PresignHandler struct {
	// Expire is the seconds for presign url expires
	Expire int `json:"expire"`
	// Remote is the remote qs path
	Remote string `json:"remote"`
	// Zone specifies the zone for presign action
	Zone string `json:"zone"`
}

// WithExpire sets the Expire field with given seconds
func (ph *PresignHandler) WithExpire(s int) *PresignHandler {
	ph.Expire = s
	return ph
}

// WithRemote sets the Remote field with given remote path
func (ph *PresignHandler) WithRemote(path string) *PresignHandler {
	ph.Remote = path
	return ph
}

// WithZone sets the Zone field with given zone
func (ph *PresignHandler) WithZone(z string) *PresignHandler {
	ph.Zone = z
	return ph
}

// Presign will handle the presign action.
func (ph *PresignHandler) Presign() (err error) {
	bucketName, objectKey, err := task.ParseQsPath(ph.Remote)
	if err != nil {
		return err
	}
	if objectKey == "" {
		return constants.ErrorQsPathObjectKeyRequired
	}
	err = contexts.Storage.SetupBucket(bucketName, ph.Zone)
	if err != nil {
		return
	}

	_, err = contexts.Storage.HeadObject(objectKey)
	if err != nil {
		return err
	}

	var isPublic bool
	// check whether the bucket is public
	ar, err := contexts.Storage.GetBucketACL()
	if err != nil {
		return err
	}
	for _, acl := range ar.ACLs {
		if acl.GranteeName == constants.PublicBucketACL {
			isPublic = true
			break
		}
	}

	// if public bucket, splice URL with bucket name, zone, host, port and its name,
	if isPublic {
		publicURL := signPublicBucket(bucketName, objectKey)
		fmt.Println(publicURL)
		return
	}

	// if the bucket is non-public, generate the link with signature,
	// expire seconds and other formatted parameters
	url, err := signPrivateBucket(objectKey, ph.Expire)
	if err != nil {
		return err
	}
	fmt.Println(url)
	return nil
}

// signPublicBucket sign a object belongs to a public bucket
func signPublicBucket(bucket, objectKey string) string {
	// splice URL with bucket name, zone, host, port and its name
	return fmt.Sprintf("%s://%s.%s.%s:%d/%s",
		viper.GetString(constants.ConfigProtocol),
		bucket,
		contexts.Storage.GetBucketZone(),
		viper.GetString(constants.ConfigHost),
		viper.GetInt(constants.ConfigPort),
		objectKey,
	)
}

// signPrivateBucket sign a object belons to a private bucket
func signPrivateBucket(objectKey string, expire int) (url string, err error) {
	return contexts.Storage.PresignObject(objectKey, expire)
}
