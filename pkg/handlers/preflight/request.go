package preflight

import (
	"encoding/json"
	"net/http"

	"github.com/airbnb/rudolph/pkg/response"
	"github.com/airbnb/rudolph/pkg/types"
	"github.com/aws/aws-lambda-go/events"
)

// Parses the HTTP Request into the appropriate request type, or returns a HTTP Response if something is wrong
func parseRequest(request events.APIGatewayProxyRequest) (parsedRequest *PreflightRequest, errorResponse *events.APIGatewayProxyResponse, err error) {
	if request.Headers["content-type"] != "application/json" && request.Headers["Content-Type"] != "application/json" {
		errorResponse, err = response.APIResponse(http.StatusUnsupportedMediaType, response.ErrInvalidMediaTypeResponse)
		return
	}

	if len(request.Body) <= 0 {
		errorResponse, err = response.APIResponse(http.StatusBadRequest, response.ErrInvalidBodyResponse)
		return
	}

	err = json.Unmarshal([]byte(request.Body), &parsedRequest)
	if err != nil {
		errorResponse, err = response.APIResponse(http.StatusBadRequest, response.ErrInvalidBodyResponse)
		return
	}

	return
}

// PreflightRequest represents sync payload sent to a sync server from a Santa client.
type PreflightRequest struct {
	OSBuild              string           `json:"os_build"`
	SantaVersion         string           `json:"santa_version"`
	Hostname             string           `json:"hostname"`
	OSVersion            string           `json:"os_version"`
	CertificateRuleCount int              `json:"certificate_rule_count"`
	BinaryRuleCount      int              `json:"binary_rule_count"`
	ClientMode           types.ClientMode `json:"client_mode"`
	SerialNumber         string           `json:"serial_num"`
	PrimaryUser          string           `json:"primary_user"`
	CompilerRuleCount    int              `json:"compiler_rule_count"`
	TransitiveRuleCount  int              `json:"transitive_rule_count"`
	RequestCleanSync     bool             `json:"request_clean_sync"`
}
