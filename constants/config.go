package constants

// Available config.
const (
	// Optional config.
	ConfigHost               = "host"
	ConfigPort               = "port"
	ConfigProtocol           = "protocol"
	ConfigLogLevel           = "log_level"
	ConfigEnableVirtualStyle = "enable_virtual_style"
	ConfigDisableURICleaning = "disable_uri_cleaning"

	// Required config.
	ConfigAccessKeyID     = "access_key_id"
	ConfigSecretAccessKey = "secret_access_key"

	// Runtime config.
	ConfigZone = "zone"
)

// Default config values.
const (
	DefaultHost               = "qingstor.com"
	DefaultPort               = "443"
	DefaultProtocol           = "https"
	DefaultLogLevel           = "info"
	DefaultEnableVirtualStyle = false
)
