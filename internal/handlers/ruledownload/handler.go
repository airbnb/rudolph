package ruledownload

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/response"
	"github.com/aws/aws-lambda-go/events"
)

// RuleDownloadHandler handles requests to the /ruledownload and /ruledownload/{machine_id} API endpoints
//
//	During every sync, Santa sensors make successive POST requests to the /ruledownload endpoint to paginate
//	through all rules.
//	When given a blank postbody (e.g. {}), it indicates the very first request in a sequence. If the
//	API returns a "cursor" in the response body, this cursor will be sent back verbatim in a subsequent postbody.
//	When a response does not return a "cursor" in the body, it signals that there are no more items to page
//	through, and the sensor will stop sending requests.
type PostRuledownloadHandler struct {
	booted        bool
	cursorService ruledownloadCursorService
	ghandler      globalRuleDownloader
	fhandler      feedRuleDownloader
	mhandler      machineRuleDownloder
}

func (h *PostRuledownloadHandler) Boot() (err error) {
	if h.booted {
		return
	}

	dynamodbTableName := os.Getenv("DYNAMODB_NAME")
	awsRegion := os.Getenv("REGION")

	client := dynamodb.GetClient(dynamodbTableName, awsRegion)

	h.cursorService = concreteRuledownloadCursorService{
		timer:   clock.ConcreteTimeProvider{},
		updater: client,
		getter:  client,
	}
	h.ghandler = concreteGlobalRuleDownloader{
		queryer: client,
	}
	h.fhandler = concreteFeedRuleDownloader{
		queryer: client,
	}
	h.mhandler = concreteMachineRuleDownloader{
		queryer: client,
		updater: client,
		timer:   clock.ConcreteTimeProvider{},
	}

	h.booted = true
	return
}

func (h *PostRuledownloadHandler) Handles(request events.APIGatewayProxyRequest) bool {
	return request.Resource == "/ruledownload/{machine_id}" && request.HTTPMethod == "POST"
}

func (h *PostRuledownloadHandler) Handle(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	machineID, ok := request.PathParameters["machine_id"]
	if !ok {
		// Unreachable code; API Gateway will never allow {machine_id} to be blank
		log.Printf("ASSERTION FAILED: Received blank {machine_id}")
		return response.APIResponse(http.StatusBadRequest, nil)
	}

	// Parse the request
	var ruledownloadRequest *RuledownloadRequest
	err := json.Unmarshal([]byte(request.Body), &ruledownloadRequest)
	if err != nil {
		log.Printf("  Failed to unmarshall ruledownload request Error %s", err.Error())
		return response.APIResponse(http.StatusBadRequest, response.ErrInvalidBodyResponse)
	}

	if ruledownloadRequest.RawCursor != "" {
		ruledownloadRequest.Cursor = &ruledownloadCursor{}
		err = json.Unmarshal([]byte(ruledownloadRequest.RawCursor), ruledownloadRequest.Cursor)
		if err != nil {
			log.Printf("  Failed to unmarshall cursor Error %s", err.Error())
			return response.APIResponse(http.StatusBadRequest, response.ErrInvalidBodyResponse)
		}
	}

	return h.handleRuleDownload(machineID, ruledownloadRequest)
}

// The "meat and potatoes" of the ruledownload flow.
func (h *PostRuledownloadHandler) handleRuleDownload(machineID string, ruledownloadRequest *RuledownloadRequest) (*events.APIGatewayProxyResponse, error) {
	cursor, err := h.cursorService.ConstructCursor(*ruledownloadRequest, machineID)
	if err != nil {
		return response.APIResponse(http.StatusInternalServerError, response.ErrInternalServerErrorResponse)
	}

	// Whenever the client requests a clean sync, we need to echo it back to the client in the response;
	// this gives the client permission to wipe out its local rules database and restart from scratch
	switch cursor.Strategy {
	case ruledownloadStrategyClean:
		return h.ghandler.handle(machineID, cursor)

	case ruledownloadStrategyIncremental:
		return h.fhandler.handle(machineID, cursor)

	case ruledownloadStrategyMachine:
		return h.mhandler.handle(machineID, ruledownloadRequest)
	}

	// How did you get here??
	log.Printf("Unreachable code reached in handleRuleDownload()!")
	return response.APIResponse(http.StatusInternalServerError, response.ErrInternalServerErrorResponse)
}
