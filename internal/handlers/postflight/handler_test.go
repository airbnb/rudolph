package postflight

import (
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler_InvalidMethod(t *testing.T) {
	var request = events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
	}

	h := &PostPostflightHandler{}
	assert.False(t, h.Handles(request))
}

type mockRuleDestroyer func(machineID string) error

func (m mockRuleDestroyer) destroyMachineRulesMarkedForDeletion(machineID string) error {
	return m(machineID)
}

type mockSyncStateUpdater func(machineID string) error

func (m mockSyncStateUpdater) updatePostflightDate(machineID string) error {
	return m(machineID)
}

// Coerce the mock types
var _ staleRuleDestroyer = mockRuleDestroyer(nil)
var _ syncStateUpdater = mockSyncStateUpdater(nil)

func TestHandler_OK(t *testing.T) {
	inputMachineID := "AAAAAAAA-A00A-1234-1234-5864377B4831"
	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/eventupload/{machine_id}",
		PathParameters: map[string]string{"machine_id": inputMachineID},
		Headers:        map[string]string{"Content-Type": "application/json"},
	}

	h := &PostPostflightHandler{
		ruleDestroyer: mockRuleDestroyer(
			func(machineID string) error {
				assert.Equal(t, inputMachineID, machineID)
				return nil
			},
		),
		syncStateUpdater: mockSyncStateUpdater(
			func(machineID string) error {
				assert.Equal(t, inputMachineID, machineID)
				return nil
			},
		),
	}

	resp, err := h.Handle(request)

	assert.Empty(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, `{"status":"ok"}`, resp.Body)
}

func TestHandler_Whoops(t *testing.T) {
	inputMachineID := "AAAAAAAA-A00A-1234-1234-5864377B4831"
	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/eventupload/{machine_id}",
		PathParameters: map[string]string{"machine_id": inputMachineID},
		Headers:        map[string]string{"Content-Type": "application/json"},
	}

	h := &PostPostflightHandler{
		ruleDestroyer: mockRuleDestroyer(
			func(machineID string) error {
				return errors.New("Yep an error.")
			},
		),
		syncStateUpdater: mockSyncStateUpdater(
			func(machineID string) error {
				return nil
			},
		),
	}

	resp, err := h.Handle(request)

	assert.Empty(t, err)
	assert.Equal(t, 500, resp.StatusCode)
	assert.Equal(t, `{}`, resp.Body)
}
