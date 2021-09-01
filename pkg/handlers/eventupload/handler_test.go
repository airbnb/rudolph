package eventupload

import (
	"errors"
	"testing"

	"github.com/airbnb/rudolph/pkg/kinesis"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestEventuploadHandler_InvalidMethod(t *testing.T) {
	var request = events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
	}

	h := &PostEventuploadHandler{}
	assert.False(t, h.Handles(request))

	resp, _ := h.Handle(request)
	assert.Equal(t, 405, resp.StatusCode)
}

func TestEventuploadHandler_IncorrectType(t *testing.T) {
	// If the request contains a mediatype that's not json reject it
	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/eventupload/{machine_id}",
		PathParameters: map[string]string{"machine_id": "AAAAAAAA-A00A-1234-1234-5864377B4831"},
		Headers:        map[string]string{"Content-Type": "application/xml"},
	}

	h := &PostEventuploadHandler{}
	resp, _ := h.Handle(request)

	assert.Equal(t, 415, resp.StatusCode)
	assert.Equal(t, `{"error":"Invalid mediatype"}`, resp.Body)
}

func TestEventuploadHandler_InvalidPathParameter(t *testing.T) {
	// If the request contains a non-valid path parameter
	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/eventupload/{machine_id}",
		PathParameters: map[string]string{"machine_id": "AAAAAAAA-A00A-1234-1234-5864377B48311"},
		Headers:        map[string]string{"Content-Type": "application/json"},
	}

	h := &PostEventuploadHandler{}
	resp, _ := h.Handle(request)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, `{"error":"Invalid path parameter"}`, resp.Body)
}

func TestEventuploadHandler_BlankPathParameter(t *testing.T) {
	// If the request contains a blank path parameter
	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/eventupload/{machine_id}",
		PathParameters: map[string]string{"machine_id": ""},
		Headers:        map[string]string{"Content-Type": "application/json"},
	}

	h := &PostEventuploadHandler{}
	resp, _ := h.Handle(request)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, `{"error":"No path parameter"}`, resp.Body)
}

func TestEventuploadHandler_EmptyBody(t *testing.T) {
	// If the request contains a mediatype that's not json reject it
	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/eventupload/{machine_id}",
		PathParameters: map[string]string{"machine_id": "AAAAAAAA-A00A-1234-1234-5864377B4831"},
		Headers:        map[string]string{"Content-Type": "application/json"},
		Body:           ``,
	}

	h := &PostEventuploadHandler{}
	resp, _ := h.Handle(request)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, `{"error":"Invalid request body"}`, resp.Body)
}

func TestEventuploadHandler_InvalidBody(t *testing.T) {
	// If the request contains a mediatype that's not json reject it
	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/eventupload/{machine_id}",
		PathParameters: map[string]string{"machine_id": "AAAAAAAA-A00A-1234-1234-5864377B4831"},
		Headers:        map[string]string{"Content-Type": "application/json"},
		Body:           `{`,
	}

	h := &PostEventuploadHandler{}
	resp, _ := h.Handle(request)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, `{"error":"Invalid request body"}`, resp.Body)
}

type testKinesisClient func(machineID string, events kinesis.KinesisEvents) error

func (c testKinesisClient) Send(machineID string, events kinesis.KinesisEvents) (err error) {
	return c(machineID, events)
}

func TestEventuploadHandler_Kinesis_InternalServerError(t *testing.T) {
	t.Run("Internal kinesis error", func(t *testing.T) {
		h := &PostEventuploadHandler{
			enableKinesis: true,
			kinesisClient: testKinesisClient(
				func(machineID string, events kinesis.KinesisEvents) error {
					return errors.New("A A A A A A A")
				},
			),
		}

		var request = events.APIGatewayProxyRequest{
			HTTPMethod:     "POST",
			Resource:       "/eventupload/{machine_id}",
			PathParameters: map[string]string{"machine_id": "AAAAAAAA-A00A-1234-1234-5864377B4831"},
			Headers:        map[string]string{"Content-Type": "application/json"},
			Body: `{"events": [{
	"parent_name": "launchd",
	"ppid": 3472,
	"pid": 24832,
	"file_path": "/Applications/My Application.app/Contents/Library/LoginItems/LauncherApplication.app/Contents/MacOS",
	"quarantine_timestamp": 0,
	"logged_in_users": [
		"john_doe"
	],
	"current_sessions": [
		"john_doe@console",
		"john_doe@ttys000",
		"john_doe@ttys001"
	],
	"executing_user": "john_doe",
	"execution_time": 1619729340.537646,
	"file_sha256": "35de834c7f280df703f57ff75b3486b9a04d73c0df96f9f6968db15fa86b8962",
	"file_name": "LauncherApplication",
	"decision": "ALLOW_UNKNOWN",
	"machine_id": "AAAAAAAA-BBBB-CCCC-DDDD-12345689012"
}]}`,
		}

		resp, _ := h.Handle(request)

		assert.Equal(t, 500, resp.StatusCode)
		assert.Equal(t, `{"error":"Internal server error"}`, resp.Body)
	})
}

func TestEventuploadHandler_Kinesis_OK(t *testing.T) {
	t.Run("With signing certificate", func(t *testing.T) {
		h := &PostEventuploadHandler{
			enableKinesis: true,
			kinesisClient: testKinesisClient(
				func(machineID string, events kinesis.KinesisEvents) error {
					assert.Equal(t, 1, len(events.Items))

					firstEvent := events.Items[0].(ForwardedEventUploadEvent) // Coerce type

					assert.Equal(t, "ALLOW_UNKNOWN", firstEvent.Decision)
					assert.Equal(t, "35de834c7f280df703f57ff75b3486b9a04d73c0df96f9f6968db15fa86b8962", firstEvent.FileSHA256)
					assert.Equal(t, "AAAAAAAA-A00A-1234-1234-5864377B4831", machineID)
					assert.Equal(t, "AAAAAAAA-A00A-1234-1234-5864377B4831", firstEvent.MachineID)

					assert.Equal(t, 3, len(firstEvent.SigningChain))
					assert.Equal(t, "0000000b28b738354c43a11486651ca33266e2b7454477d6b351df09c2e97faf", firstEvent.SigningChain[0].SHA256)

					return nil
				},
			),
		}

		var request = events.APIGatewayProxyRequest{
			HTTPMethod:     "POST",
			Resource:       "/eventupload/{machine_id}",
			PathParameters: map[string]string{"machine_id": "AAAAAAAA-A00A-1234-1234-5864377B4831"},
			Headers:        map[string]string{"Content-Type": "application/json"},
			Body: `{"events": [{
	"parent_name": "launchd",
	"ppid": 3472,
	"pid": 24832,
	"file_path": "/Applications/My Application.app/Contents/Library/LoginItems/LauncherApplication.app/Contents/MacOS",
	"quarantine_timestamp": 0,
	"logged_in_users": [
		"john_doe"
	],
	"current_sessions": [
		"john_doe@console",
		"john_doe@ttys000",
		"john_doe@ttys001"
	],
	"executing_user": "derjohn_doeek_wang",
	"execution_time": 1619729340.537646,
	"file_sha256": "35de834c7f280df703f57ff75b3486b9a04d73c0df96f9f6968db15fa86b8962",
	"file_name": "LauncherApplication",
	"decision": "ALLOW_UNKNOWN",
	"machine_id": "AAAAAAAA-BBBB-CCCC-DDDD-12345689012",
	"signing_chain": [
		{
			"cn":"Developer ID Application: My Application, Inc. (FNN8Z5JMFP)",
			"org":"My Application, Inc.",
			"ou":"FNN8Z5JMFP",
			"sha256":"0000000b28b738354c43a11486651ca33266e2b7454477d6b351df09c2e97faf",
			"valid_from":1492010408,
			"valid_until":1649863208
		},
		{
			"cn": "Developer ID Certification Authority",
			"org":"Apple Inc.",
			"ou":"Apple Certification Authority",
			"sha256":"7afc9d01a62f03a2de9637936d4afe68090d2de18d03f29c88cfb0b1ba63587f",
			"valid_from":1328134335,
			"valid_until":1801519935
		},
		{
			"cn":"Apple Root CA",
			"org":"Apple Inc.",
			"ou":"Apple Certification Authority",
			"sha256":"b0b1730ecbc7ff4505142c49f1295e6eda6bcaed7e2c68c5be91b5a11001f024",
			"valid_from":1146001236,
			"valid_until":2054670036
		}
	]
}]}`,
		}

		resp, _ := h.Handle(request)

		// Assert that we save the sensordata

		// Assert that we create a new sync state

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, `{"status":"ok"}`, resp.Body)
	})
}
