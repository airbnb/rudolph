package request

import (
	"log"
	"net/http"

	"github.com/airbnb/rudolph/pkg/response"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

// IsValidUUID returns if a machineID is a properly formatted UUID string bool
func IsValidUUID(machineID string) bool {
	_, err := uuid.Parse(machineID)
	return err == nil
}

func GetMachineID(req events.APIGatewayProxyRequest) (machineID string, errorResponse *events.APIGatewayProxyResponse, err error) {
	machineID, found := req.PathParameters["machine_id"]

	if !found || len(machineID) == 0 {
		// Unreachable code; API Gateway will never allow {machine_id} to be blank
		log.Printf("ASSERTION FAILED: Received blank {machine_id}")
		errorResponse, err = response.APIResponse(http.StatusBadRequest, response.ErrBlankPathParameterResponse)
		return
	}

	// Check if UUID is valid
	if !(IsValidUUID(machineID)) {
		// Refuse any non-valid UUID as a path parameter
		log.Printf("ASSERTION FAILED: Received blank {machine_id}")
		errorResponse, err = response.APIResponse(http.StatusBadRequest, response.ErrInvalidPathParameterResponse)
		return
	}

	return
}
