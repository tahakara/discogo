package responses

type ResponseBase struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type DataResponse struct {
	ResponseBase
	Data interface{} `json:"data,omitempty"`
}

func NewDataResponse(data interface{}, message string) DataResponse {
	if message == "" {
		message = "Ok"
	}
	return DataResponse{
		ResponseBase: ResponseBase{
			Success: true,
			Message: message,
		},
		Data: data,
	}
}

type ErrorResponse struct {
	ResponseBase
	Error string `json:"error"`
}

// ErrorResponse için constructor fonksiyonu
func NewErrorResponse(message, errorDetail string) ErrorResponse {
	if errorDetail == "" {
		errorDetail = message
	}
	return ErrorResponse{
		ResponseBase: ResponseBase{
			Success: false,
			Message: message,
		},
		Error: errorDetail,
	}
}

type SuccessResponse struct {
	ResponseBase
}

func NewSuccessResponse(message string) SuccessResponse {
	if message == "" {
		message = "Ok"
	}
	return SuccessResponse{
		ResponseBase: ResponseBase{
			Success: true,
			Message: message,
		},
	}
}

type VersionResponse struct {
	SuccessResponse
	AppName     string `json:"appName"`
	Version     string `json:"version"`
	VersionName string `json:"versionName"`
}

func NewVersionResponse(appName, version, versionName string) VersionResponse {
	return VersionResponse{
		SuccessResponse: NewSuccessResponse(""),
		AppName:         appName,
		Version:         version,
		VersionName:     versionName,
	}
}
