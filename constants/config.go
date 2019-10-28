package constants

// Available config.
const (
	// Optional config.
	ConfigHost              = "host"
	ConfigPort              = "port"
	ConfigProtocol          = "protocol"
	ConfigConnectionRetries = "connection_retries"
	ConfigLogLevel          = "log_level"

	// Required config.
	ConfigAccessKeyID     = "access_key_id"
	ConfigSecretAccessKey = "secret_access_key"

	// Runtime config.
	// ConfigZone = "zone"
)

// Default config values.
const (
	DefaultHost              = "qingstor.com"
	DefaultPort              = "443"
	DefaultProtocol          = "https"
	DefaultConnectionRetries = 3
	DefaultLogLevel          = "info"
)
