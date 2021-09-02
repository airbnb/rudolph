package eventupload

import (
	"log"
	"net/http"
	"os"

	"github.com/airbnb/rudolph/pkg/firehose"
	"github.com/airbnb/rudolph/pkg/kinesis"
	"github.com/airbnb/rudolph/pkg/lambda"

	"github.com/airbnb/rudolph/pkg/response"
	"github.com/aws/aws-lambda-go/events"
)

type PostEventuploadHandler struct {
	booted         bool
	firehoseClient firehose.FirehoseClient
	enableFirehose bool
	kinesisClient  kinesis.KinesisClient
	enableKinesis  bool

	lambdaClient lambda.LambdaClient
	enableLambda bool
}

//
func (h *PostEventuploadHandler) Boot() (err error) {
	if h.booted {
		return
	}

	handler := os.Getenv("HANDLER")
	region := os.Getenv("REGION")

	switch handler {
	case "FIREHOSE":
		firehoseName := os.Getenv("FIREHOSE_NAME")
		h.firehoseClient = firehose.GetClient(firehoseName, region)
		h.enableFirehose = true
	case "KINESIS":
		kinesisName := os.Getenv("KINESIS_NAME")
		h.kinesisClient = kinesis.GetClient(kinesisName, region)
		h.enableKinesis = true
	case "NONE":
		fallthrough
	default:
		// No handling of eventupload.
	}

	lambdaName := os.Getenv("LAMBDA_NAME")
	if lambdaName != "" {
		h.enableLambda = true
		h.lambdaClient = lambda.GetClient(lambdaName, "$LATEST", region)
	}

	h.booted = true
	return
}

//
func (h *PostEventuploadHandler) Handles(request events.APIGatewayProxyRequest) bool {
	return request.Resource == "/eventupload/{machine_id}" && request.HTTPMethod == "POST"
}

//
func (h *PostEventuploadHandler) Handle(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("EventUploadHandler request:\n%+v\n", request)

	machineID, eventsRequest, errorResponse, err := parseRequest(request)
	if errorResponse != nil || err != nil {
		return errorResponse, err
	}

	if !h.enableFirehose && !h.enableKinesis && !h.enableLambda {
		// Shortcircuit if no handlers are enabled
		log.Printf("No eventupload handlers are enabled")
		return response.APIResponse(http.StatusOK, map[string]string{"status": "ok"})
	}

	if h.enableFirehose {
		err = sendToFirehose(h.firehoseClient, machineID, eventsRequest.Events)
	}

	if err == nil && h.enableKinesis {
		err = sendToKinesis(h.kinesisClient, machineID, eventsRequest.Events)
	}

	if err == nil && h.enableLambda {
		err = sendToLambda(h.lambdaClient, machineID, eventsRequest.Events)
	}

	if err != nil {
		return response.APIResponse(http.StatusInternalServerError, response.ErrInternalServerErrorResponse)
	}

	return response.APIResponse(http.StatusOK, map[string]string{"status": "ok"})
}
