package response

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// APIResponse constructs an AWS Lambda API response object from the given data
// This object can be returned from the Lambda's main handler and will be picked up by
// AWS Lambda and rendered into an HTTP Response
func APIResponse(status int, rawBody interface{}) (*events.APIGatewayProxyResponse, error) {

	resp := &events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		StatusCode: status,
	}

	body, err := json.Marshal(rawBody)
	if err != nil {
		resp.StatusCode = 500
		resp.Body = "value could not be serialized to JSON"

		return resp, err
	}

	resp.Body = string(body)

	return resp, nil
}
