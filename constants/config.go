package constants

// Available config.
const (
	// Optional config.
	ConfigHost     = "host"
	ConfigPort     = "port"
	ConfigProtocol = "protocol"
	ConfigLogLevel = "log_level"

	// Required config.
	ConfigAccessKeyID     = "access_key_id"
	ConfigSecretAccessKey = "secret_access_key"

	// Runtime config.
	ConfigZone = "zone"
)

// Default config values.
const (
	DefaultHost     = "qingstor.com"
	DefaultPort     = "443"
	DefaultProtocol = "https"
	DefaultLogLevel = "info"
)

const (
	// EnvPrefix indicates the prefix of config env
	EnvPrefix = "qsctl"
)
