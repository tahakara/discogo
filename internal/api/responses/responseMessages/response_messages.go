package responsemessages

import "fmt"

const (
	ErrorResponseType   = "error"
	VersionResponseType = "version"

	FailedToGetTemplate = "Failed to get %s"

	ServiceHealthCheckFailed  = "Service health check failed"
	ServiceHealthCheckSuccess = "Service health check successful"
)

func MessageR(template string, args ...interface{}) string {
	return fmt.Sprintf(template, args...)
}
