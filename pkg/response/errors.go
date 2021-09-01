package response

var ErrInvalidPathParameterResponse = ErrorResponse{Error: "Invalid path parameter"}
var ErrBlankPathParameterResponse = ErrorResponse{Error: "No path parameter"}
var ErrInvalidContentTypeResponse = ErrorResponse{Error: "Invalid request content-type"}
var ErrInvalidMediaTypeResponse = ErrorResponse{Error: "Invalid mediatype"}
var ErrInvalidBodyResponse = ErrorResponse{Error: "Invalid request body"}
var ErrInvalidBodyNoSerialResponse = ErrorResponse{Error: "No serial number provided"}
var ErrInternalServerErrorResponse = ErrorResponse{Error: "Internal server error"}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}
