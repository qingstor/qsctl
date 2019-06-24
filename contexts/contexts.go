package contexts

import (
	"fmt"
	"github.com/yunify/qingstor-sdk-go/service"
	"net/http"
	"strings"
)

var (
	Service *service.Service
)

func SetupContexts() error {
	return nil
}

func GetBucket(name string) (bucket *service.Bucket, err error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url := fmt.Sprintf("%s://%s.%s:%d", c.Protocol, c.BucketName, c.Host, c.Port)

	r, err := client.Head(url)
	if err != nil {
		return
	}

	// Example URL: https://bucket.zone.qingstor.com
	zone := strings.Split(r.Header.Get("Location"), ".")[1]
	bucket, _ = Service.Bucket(name, zone)
	return
}
