package logmesssages

import "fmt"

const (
	InfoLog    = "INFO"
	ErrorLog   = "ERROR"
	WarningLog = "WARNING"
	DebugLog   = "DEBUG"
	TraceLog   = "TRACE"

	InfoStartingServiceDiscovery = "Starting service discovery tool with server IP: %s"

	ErrorLoadingEnvFile = "Error loading .env file"

	InfoMemcachedKeyNotFound = "Memcached key not found: %s"

	ErrorMemcachedFailedToRetrieveServerAddress = "Failed to retrieve Memcached server address %s"
	ErrorMemcachedConnectionFailed              = "Failed to connect to Memcached server %s"
	ErrorMemcachedPingFailed                    = "Memcached ping failed"
)

func MessageR(template string, args ...interface{}) string {
	if len(args) == 0 {
		return template
	}
	return fmt.Sprintf(template, args...)
}
