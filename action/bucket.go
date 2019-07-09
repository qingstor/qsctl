package action

import (
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/contexts"
)

// MakeBucket will make a bucket with specific name.
func MakeBucket(bucketPath string) (err error) {
	// Get bucket name from path
	bucketName, _, err := ParseQsPath(bucketPath)
	if err != nil {
		log.Errorf("Parse qs path <%s> failed [%v]",
			bucketPath, err)
		return
	}
	// Init bucket
	bucket, err := contexts.Service.Bucket(bucketName, contexts.Zone)
	if err != nil {
		log.Errorf("Initial bucket <%s> in zone <%s> failed [%v]",
			bucketName, contexts.Zone, err)
		return
	}
	// Request and create bucket
	if _, err = bucket.Put(); err != nil {
		log.Errorf("Make bucket <%s> in zone <%s> failed [%v]",
			bucketName, contexts.Zone, err)
		return
	}
	log.Infof("Bucket <%s> created.", bucketName)
	return nil
}
