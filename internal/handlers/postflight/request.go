package postflight

import (
	"encoding/json"
	"log"
	"net/http"

	apirequest "github.com/airbnb/rudolph/pkg/request"
	"github.com/airbnb/rudolph/pkg/response"
	"github.com/aws/aws-lambda-go/events"
)

// Parses the HTTP Request into the appropriate request type, or returns a HTTP Response if something is wrong
func parseRequest(request events.APIGatewayProxyRequest) (machineID string, parsedRequest *PostflightRequest, errorResponse *events.APIGatewayProxyResponse, err error) {
	if request.Headers["content-type"] != "application/json" && request.Headers["Content-Type"] != "application/json" {
		errorResponse, err = response.APIResponse(http.StatusUnsupportedMediaType, response.ErrInvalidMediaTypeResponse)
		return
	}

	machineID, errorResponse, err = apirequest.GetMachineID(request)
	if errorResponse != nil || err != nil {
		// Unreachable code; API Gateway should never encounter an error attempting to GetMachineID
		return
	}

	if len(request.Body) > 0 {
		// Parse the request
		err = json.Unmarshal([]byte(request.Body), &parsedRequest)
		if err != nil {
			log.Printf("%s\n%s", err.Error(), "request body unmarshal was not successful")
			errorResponse, err = response.APIResponse(http.StatusBadRequest, response.ErrInvalidBodyResponse)
			return
		}
	}

	return
}

type PostflightRequest struct {
	RulesReceived  int `json:"rules_received"`
	RulesProcessed int `json:"rules_processed"`
}
