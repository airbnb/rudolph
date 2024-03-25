package preflight

import (
	"fmt"
	"net/http"
	"os"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/machineconfiguration"
	apiRequest "github.com/airbnb/rudolph/pkg/request"
	"github.com/airbnb/rudolph/pkg/response"
	"github.com/aws/aws-lambda-go/events"
)

type PostPreflightHandler struct {
	booted                      bool
	rudolphDynamoDBClient       dynamodb.DynamoDBClient
	machineConfigurationService machineconfiguration.MachineConfigurationService
	stateTrackingService        stateTrackingService
	cleanSyncService            cleanSyncService
	timeProvider                clock.TimeProvider
}

func (h *PostPreflightHandler) Boot() (err error) {
	if h.booted {
		return
	}

	dynamodbTableName := os.Getenv("DYNAMODB_NAME")
	awsRegion := os.Getenv("REGION")
	h.rudolphDynamoDBClient = dynamodb.GetClient(dynamodbTableName, awsRegion)
	h.timeProvider = clock.ConcreteTimeProvider{}

	h.stateTrackingService = getStateTrackingService(h.rudolphDynamoDBClient, h.timeProvider)

	h.cleanSyncService = getCleanSyncService(h.timeProvider)

	h.machineConfigurationService = machineconfiguration.GetMachineConfigurationService(h.rudolphDynamoDBClient, h.timeProvider)

	h.booted = true
	return
}

func (h *PostPreflightHandler) Handles(request events.APIGatewayProxyRequest) bool {
	return request.Resource == "/preflight/{machine_id}" && request.HTTPMethod == "POST"
}

func (h *PostPreflightHandler) Handle(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Get machineID from req.PathParameter of "machine_id" and it must be a valid UUID
	machineID, errResponse, err := apiRequest.GetMachineID(request)
	if errResponse != nil {
		// API Gateway will never allow {machine_id} to be blank,
		// this will fail if the path parameter of "machine_id" does not match a valid UUID
		return errResponse, nil
	}
	if err != nil {
		// Unreachable code; API Gateway should never encounter an error attempting to GetMachineID
		return errResponse, err
	}

	preflightRequest, errorResponse, err := parseRequest(request)
	if errorResponse != nil {
		return errorResponse, nil
	}
	if err != nil {
		return response.APIResponse(http.StatusBadRequest, response.ErrInvalidBodyResponse)
	}
	return h.handlePreflight(machineID, preflightRequest)
}

// The main controller function for API calls to /preflight
func (h *PostPreflightHandler) handlePreflight(machineID string, preflightRequest *PreflightRequest) (*events.APIGatewayProxyResponse, error) {
	// Here we figure out where to set the "feedSync" cursor in the newly created sync state
	var feedSyncCursor string
	var performCleanSync bool = false

	// Save the state of the preflight request - this is the sensor state
	err := h.stateTrackingService.saveSensorDataFromPreflightRequest(machineID, preflightRequest)
	if err != nil {
		return response.APIResponse(http.StatusInternalServerError, response.ErrInternalServerErrorResponse)
	}

	// Retrieve the intended configuration for the machine
	machineConfiguration, err := h.machineConfigurationService.GetIntendedConfig(machineID)
	if err != nil {
		return response.APIResponse(http.StatusInternalServerError, response.ErrInternalServerErrorResponse)
	}

	// Get the previous sync state
	prevSyncState, err := h.stateTrackingService.getSyncState(machineID)
	if err != nil {
		return response.APIResponse(http.StatusInternalServerError, response.ErrInternalServerErrorResponse)
	}

	// Determine if a Clean sync should be performed based on the preflight request
	// if the machine needs a periodic refresh
	// or if the machine is new
	switch preflightRequest.RequestCleanSync {
	case true:
		performCleanSync = true
	case false:
		// Retrieve the current feed sync cursor
		feedSyncCursor, performCleanSync = h.stateTrackingService.getFeedSyncStateCursor(prevSyncState)
		// If a clean sync should be forced, break out and do it now
		if performCleanSync {
			break
		}
		// Determine if a refresh clean sync should be performed
		performCleanSync, err = h.cleanSyncService.determineCleanSync(
			machineID,
			preflightRequest,
			prevSyncState,
		)
		if err != nil {
			return response.APIResponse(http.StatusInternalServerError, err)
		}
	default:
		performCleanSync = true
	}

	// Set up a syncState object which will track the progress of the currently requested sync
	// Here we use dynamodb:PutItem to restart the whole process, wipe out any previous sync
	var lastCleanSyncTime string
	switch performCleanSync {
	case true:
		lastCleanSyncTime = clock.RFC3339(h.timeProvider.Now())
	case false:
		lastCleanSyncTime = prevSyncState.LastCleanSync
	}

	err = h.stateTrackingService.saveSyncState(
		machineID,
		// If the CleanSync is going to be forced, log this request in the SensorState, so we can indicate
		// in the /ruledownload step which strategy to use
		performCleanSync,
		lastCleanSyncTime,
		machineConfiguration.BatchSize,
		feedSyncCursor,
	)

	if err != nil {
		err = fmt.Errorf("encountered error trying to save new sync state: %w", err)
		return response.APIResponse(http.StatusInternalServerError, err)
	}

	// Construct the response
	preflightResponse := ConstructPreflightResponse(machineConfiguration, performCleanSync)

	return response.APIResponse(http.StatusOK, preflightResponse)
}
