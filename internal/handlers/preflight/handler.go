package preflight

import (
	"net/http"
	"os"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/machineconfiguration"
	apiRequest "github.com/airbnb/rudolph/pkg/request"
	"github.com/airbnb/rudolph/pkg/response"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

type PostPreflightHandler struct {
	booted                     bool
	daysElapseUntilCleanSync   int
	sensorDataSaver            sensorDataSaver
	machineConfigurationGetter machineConfigurationGetter
	syncStateManager           syncStateManager
	timeProvider               clock.TimeProvider
}

//
func (h *PostPreflightHandler) Boot() (err error) {
	if h.booted {
		return
	}

	dynamodbTableName := os.Getenv("DYNAMODB_NAME")
	awsRegion := os.Getenv("REGION")
	client := dynamodb.GetClient(dynamodbTableName, awsRegion)
	h.timeProvider = clock.ConcreteTimeProvider{}

	h.sensorDataSaver = concreteSensorDataSaver{
		putter: client,
	}
	h.machineConfigurationGetter = concreteMachineConfigurationGetter{
		fetcher: machineconfiguration.GetMachineConfigurationService(client, h.timeProvider),
	}
	h.syncStateManager = concreteSyncStateManager{
		getter: client,
		putter: client,
	}

	h.daysElapseUntilCleanSync = 7
	h.booted = true
	return
}

//
func (h *PostPreflightHandler) Handles(request events.APIGatewayProxyRequest) bool {
	return request.Resource == "/preflight/{machine_id}" && request.HTTPMethod == "POST"
}

//
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
	err := h.sensorDataSaver.saveSensorDataFromRequest(h.timeProvider, machineID, preflightRequest)
	if err != nil {
		return response.APIResponse(http.StatusInternalServerError, response.ErrInternalServerErrorResponse)
	}

	// Get the config rules using the machineID from the preflight request
	// if no results, GetDesiredConfig will return the global_config
	machineConfiguration, err := h.machineConfigurationGetter.getDesiredConfig(machineID)
	if err != nil {
		return response.APIResponse(http.StatusInternalServerError, response.ErrInternalServerErrorResponse)
	}

	// Get the previous sync state
	prevSyncState, err := h.syncStateManager.getSyncState(machineID)
	if err != nil {
		return response.APIResponse(http.StatusInternalServerError, response.ErrInternalServerErrorResponse)
	}

	// Here we figure out where to set the "feedSync" cursor in the newly created sync state
	var feedSyncCursor string
	var doCleanSync bool = false

	if !preflightRequest.RequestCleanSync {
		// Check when the last preflight request took place
		if prevSyncState != nil && prevSyncState.FeedSyncCursor != "" {
			// Inherit the feed feed sync cursor from the previous sync state to kind of "pick up where it left off"
			feedSyncCursor = prevSyncState.FeedSyncCursor
		} else {
			// If there is no previous sync state, or if no cursor exists, then we assume the client either
			// has never sync'd before or something went horribly wrong. Always force it to clean sync and just set
			// the feed sync cursor to "now"
			feedSyncCursor = clock.RFC3339(h.timeProvider.Now())
			doCleanSync = true
		}
	}

	// Re-queue for a clean sync if its been awhile >= daysElapseUntilCleanSync setting
	// If prevSyncState has not happened yet or lastCleanSync is nil, perform a clean sync
	if preflightRequest.RequestCleanSync {
		doCleanSync = true
	} else {
		// If a clean sync was not explicitly requested, we instead figure out when the last one was
		// and then force a clean sync if it hasn't happened in the last 7 days
		daysSinceLastSync := daysSinceLastSync(h.timeProvider, prevSyncState)

		// To reduce stampeding, we introduce a bit of dithering by using the machineID as the seed to randomize the number of days required to elapse before performing a clean sync.
		// Given the same MachineID, performCleanSync will always provide the same chaos int and evenly space all clients out to require clean sync
		// 7 days + 1d10 (Based on the MachineID input) * 10 minutes
		shouldPerformCleanSync, err := shouldPerformCleanSync(machineID, daysSinceLastSync, h.daysElapseUntilCleanSync)
		if err != nil {
			return response.APIResponse(http.StatusInternalServerError, err)
		}

		// Determine if a CleanSync is required
		if shouldPerformCleanSync {
			doCleanSync = true
		}
	}

	// Set up a syncState object which will track the progress of the currently requested sync
	// Here we use dynamodb:PutItem to restart the whole process, wipe out any previous sync
	newLastCleanSync := ""
	if doCleanSync {
		newLastCleanSync = clock.RFC3339(h.timeProvider.Now())
	} else if prevSyncState != nil {
		newLastCleanSync = prevSyncState.LastCleanSync
	}
	err = h.syncStateManager.saveNewSyncState(
		h.timeProvider,
		machineID,
		// If the CleanSync is going to be forced, log this request in the SensorState, so we can indicate
		// in the /ruledownload step which strategy to use
		doCleanSync,
		newLastCleanSync,
		machineConfiguration.BatchSize,
		feedSyncCursor,
	)
	if err != nil {
		err = errors.Wrapf(err, "Encountered error trying to save new sync state")
		return response.APIResponse(http.StatusInternalServerError, err)
	}

	// Construct the response
	preflightResponse := ConstructPreflightResponse(machineConfiguration, doCleanSync)
	if err != nil {
		return response.APIResponse(http.StatusInternalServerError, err)
	}

	return response.APIResponse(http.StatusOK, preflightResponse)
}
