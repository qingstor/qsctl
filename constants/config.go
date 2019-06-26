package constants

// Available config.
const (
	// Optional config.
	ConfigHost              = "host"
	ConfigPort              = "port"
	ConfigProtocol          = "protocol"
	ConfigConnectionRetries = "connection_retries"
	ConfigLogLevel          = "level"

	// Required config.
	ConfigAccessKeyID     = "access_key_id"
	ConfigSecretAccessKey = "secret_access_key"

	// Runtime config.
	ConfigZone       = "zone"
	ConfigBucketName = "bucket_name"
)
