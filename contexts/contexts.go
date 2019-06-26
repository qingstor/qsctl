package contexts

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/yunify/qingstor-sdk-go/v3/config"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/yunify/qsctl/constants"
)

var (
	// Service is the global service.
	Service *service.Service
	// Bucket is the bucket for bucket operation.
	Bucket *service.Bucket
)

// SetupServices will setup services.
func SetupServices() (err error) {
	if Service != nil {
		return
	}

	cfg, err := config.New(
		viper.GetString(constants.ConfigAccessKeyID),
		viper.GetString(constants.ConfigSecretAccessKey),
	)
	if err != nil {
		return
	}
	cfg.Host = viper.GetString(constants.ConfigHost)
	cfg.Port = viper.GetInt(constants.ConfigPort)
	cfg.Protocol = viper.GetString(constants.ConfigProtocol)
	cfg.ConnectionRetries = viper.GetInt(constants.ConfigConnectionRetries)
	cfg.LogLevel = viper.GetString(constants.ConfigLogLevel)

	Service, err = service.Init(cfg)
	if err != nil {
		return
	}
	return
}

// SetupBuckets will create a new bucket instance.
func SetupBuckets(name, zone string) (bucket *service.Bucket, err error) {
	err = SetupServices()
	if err != nil {
		return
	}

	if Bucket != nil {
		return Bucket, nil
	}

	if zone != "" {
		Bucket, _ = Service.Bucket(name, zone)
		return Bucket, nil
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
		return
	}

	// Example URL: https://bucket.zone.qingstor.com
	zone = strings.Split(r.Header.Get("Location"), ".")[1]
	Bucket, _ = Service.Bucket(name, zone)
	return Bucket, nil
}
