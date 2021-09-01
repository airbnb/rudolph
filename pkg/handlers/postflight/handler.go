package postflight

import (
	"log"
	"net/http"
	"os"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/response"
	"github.com/aws/aws-lambda-go/events"
)

// PostflightHandler is the entry point for the /postflight API call
type PostPostflightHandler struct {
	booted           bool
	ruleDestroyer    staleRuleDestroyer
	syncStateUpdater syncStateUpdater
}

//
func (h *PostPostflightHandler) Boot() (err error) {
	if h.booted {
		return
	}

	dynamodbTableName := os.Getenv("DYNAMODB_NAME")
	awsRegion := os.Getenv("REGION")

	client := dynamodb.GetClient(dynamodbTableName, awsRegion)

	h.ruleDestroyer = concreteRuleDestroyer{
		queryer: client,
		deleter: client,
	}
	h.syncStateUpdater = concreteSyncStateUpdater{
		timeProvider: clock.ConcreteTimeProvider{},
		updater:      client,
	}

	h.booted = true
	return
}

//
func (h *PostPostflightHandler) Handles(request events.APIGatewayProxyRequest) bool {
	return request.Resource == "/postflight/{machine_id}" && request.HTTPMethod == "POST"
}

//
func (h *PostPostflightHandler) Handle(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	machineID, ok := request.PathParameters["machine_id"]
	if !ok {
		// Unreachable code; API Gateway will never allow {machine_id} to be blank
		log.Printf("ASSERTION FAILED: Received blank {machine_id}")
		return response.APIResponse(http.StatusBadRequest, nil)
	}

	// FIXME (derek.wang) This postflight is not properly setting the FeedSyncCursor; it's just setting it to
	// like "10 minutes ago and leaving it like that. Pretty dumb, derek!"
	err := h.syncStateUpdater.updatePostflightDate(machineID)
	if err != nil {
		log.Printf("Failed to set final PostflightAt")
		return response.APIResponse(http.StatusInternalServerError, err)
	}

	// Delete rules marked for deletion
	err = h.ruleDestroyer.destroyMachineRulesMarkedForDeletion(machineID)
	if err != nil {
		return response.APIResponse(http.StatusInternalServerError, err)
	}

	// Optional: Archive the sync state
	// _ = archiveSyncState(h.dynamodbClient, machineID)

	return response.APIResponse(http.StatusOK, map[string]string{"status": "ok"})
}
