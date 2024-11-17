package ruledownload

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

// Test Mocks
type mockCursorService func(req RuledownloadRequest, machineID string) (cursor ruledownloadCursor, err error)

func (m mockCursorService) ConstructCursor(req RuledownloadRequest, machineID string) (cursor ruledownloadCursor, err error) {
	return m(req, machineID)
}

type mockGlobalRuleDownloader func(machineID string, cursor ruledownloadCursor) (*events.APIGatewayProxyResponse, error)

func (m mockGlobalRuleDownloader) handle(machineID string, cursor ruledownloadCursor) (*events.APIGatewayProxyResponse, error) {
	return m(machineID, cursor)
}

type mockFeedRuleDownloader func(machineID string, cursor ruledownloadCursor) (*events.APIGatewayProxyResponse, error)

func (m mockFeedRuleDownloader) handle(machineID string, cursor ruledownloadCursor) (*events.APIGatewayProxyResponse, error) {
	return m(machineID, cursor)
}

type mockMachineRuleDownloader func(machineID string, ruledownloadRequest *RuledownloadRequest) (*events.APIGatewayProxyResponse, error)

func (m mockMachineRuleDownloader) handle(machineID string, ruledownloadRequest *RuledownloadRequest) (*events.APIGatewayProxyResponse, error) {
	return m(machineID, ruledownloadRequest)
}

// Coerce our mocks to conform to the interfaces that they are intended to implement
var _ ruledownloadCursorService = mockCursorService(nil)
var _ globalRuleDownloader = mockGlobalRuleDownloader(nil)
var _ feedRuleDownloader = mockFeedRuleDownloader(nil)
var _ machineRuleDownloder = mockMachineRuleDownloader(nil)

// Actual Tests
func Test_PostRuledownloadHandler_SendToCorrectHandler(t *testing.T) {
	type test struct {
		cursor         ruledownloadCursor
		ghandlerCalled bool
		fhandlerCalled bool
		mhandlerCalled bool
	}

	cases := []test{
		{
			cursor: ruledownloadCursor{
				Strategy:   ruledownloadStrategyClean,
				BatchSize:  3,
				PageNumber: 1,
			},
			ghandlerCalled: true,
		},
		{
			cursor: ruledownloadCursor{
				Strategy:   ruledownloadStrategyIncremental,
				BatchSize:  3,
				PageNumber: 1,
			},
			fhandlerCalled: true,
		},
		{
			cursor: ruledownloadCursor{
				Strategy:   ruledownloadStrategyMachine,
				BatchSize:  3,
				PageNumber: 1,
			},
			mhandlerCalled: true,
		},
	}

	for _, test := range cases {
		handler := PostRuledownloadHandler{
			cursorService: mockCursorService(
				func(req RuledownloadRequest, machineID string) (ruledownloadCursor, error) {
					return test.cursor, nil
				},
			),
			ghandler: mockGlobalRuleDownloader(
				func(machineID string, cursor ruledownloadCursor) (*events.APIGatewayProxyResponse, error) {
					assert.True(t, test.ghandlerCalled)
					assert.Equal(t, 3, cursor.BatchSize)
					return &events.APIGatewayProxyResponse{
						StatusCode: http.StatusOK,
						Body:       "blah",
					}, nil
				},
			),
			fhandler: mockFeedRuleDownloader(
				func(machineID string, cursor ruledownloadCursor) (*events.APIGatewayProxyResponse, error) {
					assert.True(t, test.fhandlerCalled)
					assert.Equal(t, 3, cursor.BatchSize)
					return &events.APIGatewayProxyResponse{
						StatusCode: http.StatusOK,
						Body:       "blah",
					}, nil
				},
			),
			mhandler: mockMachineRuleDownloader(
				func(machineID string, ruledownloadRequest *RuledownloadRequest) (*events.APIGatewayProxyResponse, error) {
					assert.True(t, test.mhandlerCalled)
					return &events.APIGatewayProxyResponse{
						StatusCode: http.StatusOK,
						Body:       "blah",
					}, nil
				},
			),
		}

		machineID := "AAAA-BBBB-CCCC-DDDD"
		request := events.APIGatewayProxyRequest{
			HTTPMethod: "POST",
			Resource:   "/ruledownload/{machine_id}",
			Headers:    map[string]string{"Content-Type": "application/json"},
			PathParameters: map[string]string{
				"machine_id": machineID,
			},
			Body: "{}",
		}
		resp, err := handler.Handle(request)

		assert.Empty(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, `blah`, resp.Body)
	}
}

func Test_PostRuledownloadHandler_CursorServiceReceivesCorrectRequest(t *testing.T) {
	machineID := "AAAA-BBBB-CCCC-DDDD"
	var called bool
	handler := PostRuledownloadHandler{
		cursorService: mockCursorService(
			func(req RuledownloadRequest, mID string) (ruledownloadCursor, error) {
				called = true
				assert.Equal(t, machineID, mID)

				assert.Equal(t, ruledownloadStrategyIncremental, req.Cursor.Strategy)
				assert.Equal(t, 7, req.Cursor.BatchSize)
				assert.Equal(t, 2, req.Cursor.PageNumber)
				assert.Equal(t, "AAAA", req.Cursor.PartitionKey)
				assert.Equal(t, "eeeeee", req.Cursor.SortKey)

				// Return anything; doesn't matter.
				return ruledownloadCursor{Strategy: ruledownloadStrategyIncremental}, nil
			},
		),
		ghandler: mockGlobalRuleDownloader(
			func(machineID string, cursor ruledownloadCursor) (*events.APIGatewayProxyResponse, error) {
				return nil, nil
			},
		),
		fhandler: mockFeedRuleDownloader(
			func(machineID string, cursor ruledownloadCursor) (*events.APIGatewayProxyResponse, error) {
				return nil, nil
			},
		),
		mhandler: mockMachineRuleDownloader(
			func(machineID string, ruledownloadRequest *RuledownloadRequest) (*events.APIGatewayProxyResponse, error) {
				return nil, nil
			},
		),
	}

	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Resource:   "/ruledownload/{machine_id}",
		Headers:    map[string]string{"Content-Type": "application/json"},
		PathParameters: map[string]string{
			"machine_id": machineID,
		},
		Body: `{"cursor": "{\"strategy\":2, \"batch_size\":7,\"page\":2,\"pk\":\"AAAA\",\"sk\":\"eeeeee\"}"}`,
	}

	_, err := handler.Handle(request)

	assert.True(t, called)
	assert.Empty(t, err)
}
