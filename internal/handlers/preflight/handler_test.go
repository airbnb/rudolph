package preflight

import (
	"encoding/json"
	"testing"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/machineconfiguration"
	"github.com/airbnb/rudolph/pkg/model/sensordata"
	"github.com/airbnb/rudolph/pkg/model/syncstate"
	"github.com/airbnb/rudolph/pkg/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPreflightHandler_InvalidMethod(t *testing.T) {
	// The API endpoint only accepts HTTP POST calls
	// https://github.com/aws/aws-lambda-go/blob/master/events/apigw.go#L5-L19
	var request = events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
	}

	h := &PostPreflightHandler{}
	assert.False(t, h.Handles(request))
}

func TestPreflightHandler_IncorrectType(t *testing.T) {
	// If the request contains a mediatype that's not json reject it
	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/preflight/{machine_id}",
		PathParameters: map[string]string{"machine_id": "AAAAAAAA-A00A-1234-1234-5864377B4831"},
		Headers:        map[string]string{"Content-Type": "application/xml"},
	}

	h := &PostPreflightHandler{}
	resp, _ := h.Handle(request)

	assert.Equal(t, 415, resp.StatusCode)
	assert.Equal(t, `{"error":"Invalid mediatype"}`, resp.Body)
}

func TestPreflightHandler_InvalidPathParameter(t *testing.T) {
	// If the request contains a non-valid path parameter
	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/preflight/{machine_id}",
		PathParameters: map[string]string{"machine_id": "AAAAAAAA-A00A-1234-1234-5864377B48311"},
		Headers:        map[string]string{"Content-Type": "application/json"},
	}

	h := &PostPreflightHandler{}
	resp, _ := h.Handle(request)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, `{"error":"Invalid path parameter"}`, resp.Body)
}

func TestPreflightHandler_BlankPathParameter(t *testing.T) {
	// If the request contains a blank path parameter
	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/preflight/{machine_id}",
		PathParameters: map[string]string{"machine_id": ""},
		Headers:        map[string]string{"Content-Type": "application/json"},
	}

	h := &PostPreflightHandler{}
	resp, _ := h.Handle(request)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, `{"error":"No path parameter"}`, resp.Body)
}

func TestPreflightHandler_EmptyBody(t *testing.T) {
	// If the request contains a mediatype that's not json reject it
	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/preflight/{machine_id}",
		PathParameters: map[string]string{"machine_id": "AAAAAAAA-A00A-1234-1234-5864377B4831"},
		Headers:        map[string]string{"Content-Type": "application/json"},
		Body:           ``,
	}

	h := &PostPreflightHandler{}
	resp, _ := h.Handle(request)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, `{"error":"Invalid request body"}`, resp.Body)
}

func TestPreflightHandler_InvalidBody(t *testing.T) {
	// If the request contains a mediatype that's not json reject it
	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/preflight/{machine_id}",
		PathParameters: map[string]string{"machine_id": "AAAAAAAA-A00A-1234-1234-5864377B4831"},
		Headers:        map[string]string{"Content-Type": "application/json"},
		Body:           `{`,
	}

	h := &PostPreflightHandler{}
	resp, _ := h.Handle(request)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, `{"error":"Invalid request body"}`, resp.Body)
}

type MockDynamodb struct {
	dynamodb.DynamoDBClient
	machineconfiguration.ConcreteConfigurationFetcher
	mock.Mock
}

func (m *MockDynamodb) GetItem(key dynamodb.PrimaryKey, consistentRead bool) (*awsdynamodb.GetItemOutput, error) {
	args := m.Called(key, consistentRead)
	return args.Get(0).(*awsdynamodb.GetItemOutput), args.Error(1)
}

func (m *MockDynamodb) PutItem(item interface{}) (*awsdynamodb.PutItemOutput, error) {
	args := m.Called(item)
	return args.Get(0).(*awsdynamodb.PutItemOutput), args.Error(1)
}

func (m *MockDynamodb) UpdateItem(key dynamodb.PrimaryKey, item interface{}) (*awsdynamodb.UpdateItemOutput, error) {
	args := m.Called(key, item)
	return args.Get(0).(*awsdynamodb.UpdateItemOutput), args.Error(1)
}

// OK
// Tests the basic positive flow.
func TestHandler_OK(t *testing.T) {
	now, _ := clock.ParseRFC3339("2000-01-01T00:00:00Z")
	inputMachineID := "AAAAAAAA-A00A-1234-1234-5864377B4831"
	timeProvider := clock.FrozenTimeProvider{
		Current: now,
	}

	preflightRequestBody := &PreflightRequest{
		OSBuild:              "23E214",
		OSVersion:            "14.4",
		SantaVersion:         "2024.2",
		Hostname:             "my-awesome-macbook-pro",
		ClientMode:           types.Lockdown,
		SerialNumber:         "C02123456789",
		PrimaryUser:          "nobody",
		CertificateRuleCount: 2364,
		TeamIDRuleCount:      0,
		SigningIDRuleCount:   12,
		CDHashRuleCount:      34,
		BinaryRuleCount:      43676,
		RequestCleanSync:     false,
	}

	b, err := json.Marshal(preflightRequestBody)
	if err != nil {
		t.Fatal(err)
	}

	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/preflight/{machine_id}",
		PathParameters: map[string]string{"machine_id": inputMachineID},
		Headers:        map[string]string{"Content-Type": "application/json"},
		Body:           string(b),
	}

	mockedConfigurationFetcher := &MockDynamodb{}

	config := machineconfiguration.MachineConfiguration{
		ClientMode:       types.Lockdown,
		BatchSize:        37,
		UploadLogsURL:    "/aaa",
		EnableBundles:    true,
		AllowedPathRegex: "(^/Applications)",
		CleanSync:        false,
	}

	returnedConfig, err := attributevalue.MarshalMap(config)
	if err != nil {
		t.Fatal(err)
	}
	mockedConfigurationFetcher.On("GetItem", mock.Anything, mock.Anything).Return(&awsdynamodb.GetItemOutput{
		Item: returnedConfig,
	}, nil)

	mockedStateTracking := &MockDynamodb{}
	mockedStateTracking.On("GetItem", mock.Anything, mock.Anything).Return(&awsdynamodb.GetItemOutput{
		Item: nil,
	}, nil)

	mockedStateTracking.On("PutItem", mock.MatchedBy(func(syncState syncstate.SyncStateRow) bool {
		return syncState.MachineID == inputMachineID && syncState.BatchSize == 37 && syncState.LastCleanSync == "2000-01-01T00:00:00Z" && syncState.FeedSyncCursor == "2000-01-01T00:00:00Z"
	})).Return(&awsdynamodb.PutItemOutput{}, nil)

	mockedStateTracking.On("PutItem", mock.MatchedBy(func(sensorData sensordata.SensorData) bool {
		if sensorData.OSBuild != preflightRequestBody.OSBuild {
			return false
		}
		if sensorData.OSVersion != preflightRequestBody.OSVersion {
			return false
		}
		if sensorData.SantaVersion != preflightRequestBody.SantaVersion {
			return false
		}
		if sensorData.ClientMode != preflightRequestBody.ClientMode {
			return false
		}
		if sensorData.SerialNum != preflightRequestBody.SerialNumber {
			return false
		}
		if sensorData.PrimaryUser != preflightRequestBody.PrimaryUser {
			return false
		}
		if sensorData.CertificateRuleCount != preflightRequestBody.CertificateRuleCount {
			return false
		}
		if sensorData.TeamIDRuleCount != preflightRequestBody.TeamIDRuleCount {
			return false
		}
		if sensorData.SigningIDRuleCount != preflightRequestBody.SigningIDRuleCount {
			return false
		}
		if sensorData.CDHashRuleCount != preflightRequestBody.CDHashRuleCount {
			return false
		}
		if sensorData.BinaryRuleCount != preflightRequestBody.BinaryRuleCount {
			return false
		}
		if sensorData.RequestCleanSync != preflightRequestBody.RequestCleanSync {
			return false
		}
		return true
	})).Return(&awsdynamodb.PutItemOutput{}, nil)

	h := &PostPreflightHandler{
		timeProvider:                timeProvider,
		machineConfigurationService: machineconfiguration.GetMachineConfigurationService(mockedConfigurationFetcher, timeProvider),
		stateTrackingService:        getStateTrackingService(mockedStateTracking, timeProvider),
		cleanSyncService:            getCleanSyncService(timeProvider),
	}

	resp, err := h.Handle(request)

	assert.Empty(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Ensure that the response matches the configuration returned
	assert.Equal(t, `{"client_mode":"LOCKDOWN","blocked_path_regex":"","allowed_path_regex":"(^/Applications)","batch_size":37,"enable_bundles":true,"enable_transitive_rules":false,"upload_logs_url":"/aaa","sync_type":"clean"}`, resp.Body)
}

func TestHandler_OK_Refresh_CleanSync(t *testing.T) {
	now, _ := clock.ParseRFC3339("2001-01-01T00:00:00Z")
	inputMachineID := "AAAAAAAA-A00A-1234-1234-5864377B4831"
	timeProvider := clock.FrozenTimeProvider{
		Current: now,
	}

	preflightRequestBody := &PreflightRequest{
		OSBuild:              "23E214",
		OSVersion:            "14.4",
		SantaVersion:         "2024.2",
		Hostname:             "my-awesome-macbook-pro",
		ClientMode:           types.Lockdown,
		SerialNumber:         "C02123456789",
		PrimaryUser:          "nobody",
		CertificateRuleCount: 2364,
		TeamIDRuleCount:      0,
		SigningIDRuleCount:   12,
		CDHashRuleCount:      34,
		BinaryRuleCount:      43676,
		RequestCleanSync:     false,
	}

	b, err := json.Marshal(preflightRequestBody)
	if err != nil {
		t.Fatal(err)
	}

	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/preflight/{machine_id}",
		PathParameters: map[string]string{"machine_id": inputMachineID},
		Headers:        map[string]string{"Content-Type": "application/json"},
		Body:           string(b),
	}
	mockedConfigurationFetcher := &MockDynamodb{}

	config := machineconfiguration.MachineConfiguration{
		ClientMode:       types.Lockdown,
		BatchSize:        37,
		UploadLogsURL:    "/aaa",
		EnableBundles:    true,
		AllowedPathRegex: "(^/Applications)",
		CleanSync:        false,
	}

	returnedConfig, err := attributevalue.MarshalMap(config)
	if err != nil {
		t.Fatal(err)
	}
	mockedConfigurationFetcher.On("GetItem", mock.Anything, mock.Anything).Return(&awsdynamodb.GetItemOutput{
		Item: returnedConfig,
	}, nil)

	mockedStateTracking := &MockDynamodb{}
	syncState := syncstate.SyncStateRow{
		PrimaryKey: dynamodb.PrimaryKey{PartitionKey: "Machine#AAAAAAAA-A00A-1234-1234-5864377B4831", SortKey: "SyncState"},
		SyncState: syncstate.SyncState{
			MachineID:      inputMachineID,
			BatchSize:      37,
			LastCleanSync:  "2000-12-01T00:00:00Z",
			FeedSyncCursor: "2000-12-15T00:00:00Z",
			DataType:       types.DataTypeSyncState,
		},
	}
	returnedSyncState, err := attributevalue.MarshalMap(syncState)
	if err != nil {
		t.Fatal(err)
	}
	mockedStateTracking.On("GetItem", mock.Anything, mock.Anything).Return(&awsdynamodb.GetItemOutput{
		Item: returnedSyncState,
	}, nil)

	// mockedStateTracking.On("PutItem", mock.MatchedBy(func(item interface{}) bool {
	mockedStateTracking.On("PutItem", mock.MatchedBy(func(syncState syncstate.SyncStateRow) bool {
		return syncState.PrimaryKey.PartitionKey == "Machine#AAAAAAAA-A00A-1234-1234-5864377B4831" && syncState.MachineID == inputMachineID && syncState.BatchSize == 37 && syncState.LastCleanSync == "2001-01-01T00:00:00Z" && syncState.FeedSyncCursor == "2000-12-15T00:00:00Z" && syncState.CleanSync == true
	})).Return(&awsdynamodb.PutItemOutput{}, nil)

	mockedStateTracking.On("PutItem", mock.MatchedBy(func(sensorData sensordata.SensorData) bool {
		if sensorData.OSBuild != preflightRequestBody.OSBuild {
			return false
		}
		if sensorData.OSVersion != preflightRequestBody.OSVersion {
			return false
		}
		if sensorData.SantaVersion != preflightRequestBody.SantaVersion {
			return false
		}
		if sensorData.ClientMode != preflightRequestBody.ClientMode {
			return false
		}
		if sensorData.SerialNum != preflightRequestBody.SerialNumber {
			return false
		}
		if sensorData.PrimaryUser != preflightRequestBody.PrimaryUser {
			return false
		}
		if sensorData.CertificateRuleCount != preflightRequestBody.CertificateRuleCount {
			return false
		}
		if sensorData.TeamIDRuleCount != preflightRequestBody.TeamIDRuleCount {
			return false
		}
		if sensorData.SigningIDRuleCount != preflightRequestBody.SigningIDRuleCount {
			return false
		}
		if sensorData.CDHashRuleCount != preflightRequestBody.CDHashRuleCount {
			return false
		}
		if sensorData.BinaryRuleCount != preflightRequestBody.BinaryRuleCount {
			return false
		}
		if sensorData.RequestCleanSync != preflightRequestBody.RequestCleanSync {
			return false
		}
		return true
	})).Return(&awsdynamodb.PutItemOutput{}, nil)

	h := &PostPreflightHandler{
		timeProvider:                timeProvider,
		machineConfigurationService: machineconfiguration.GetMachineConfigurationService(mockedConfigurationFetcher, timeProvider),
		stateTrackingService:        getStateTrackingService(mockedStateTracking, timeProvider),
		cleanSyncService:            getCleanSyncService(timeProvider),
	}

	resp, err := h.Handle(request)

	assert.Empty(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Ensure that the response matches the configuration returned
	assert.Equal(t, `{"client_mode":"LOCKDOWN","blocked_path_regex":"","allowed_path_regex":"(^/Applications)","batch_size":37,"enable_bundles":true,"enable_transitive_rules":false,"upload_logs_url":"/aaa","sync_type":"clean"}`, resp.Body)
}

func TestHandler_OK_No_Refresh_CleanSync(t *testing.T) {
	now, _ := clock.ParseRFC3339("2001-01-01T00:00:00Z")
	inputMachineID := "AAAAAAAA-A00A-1234-1234-5864377B4831"
	timeProvider := clock.FrozenTimeProvider{
		Current: now,
	}
	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/preflight/{machine_id}",
		PathParameters: map[string]string{"machine_id": inputMachineID},
		Headers:        map[string]string{"Content-Type": "application/json"},
		Body: `{
	"os_build":"20D5029f",
	"santa_version":"2021.1",
	"hostname":"my-awesome-macbook-pro.attlocal.net",
	"transitive_rule_count":0,
	"os_version":"11.2",
	"certificate_rule_count":2,
	"client_mode":"MONITOR",
	"serial_num":"C02123456789",
	"binary_rule_count":3,
	"primary_user":"nobody",
	"compiler_rule_count":0
}`,
	}
	mockedConfigurationFetcher := &MockDynamodb{}

	config := machineconfiguration.MachineConfiguration{
		ClientMode:       types.Lockdown,
		BatchSize:        37,
		UploadLogsURL:    "/aaa",
		EnableBundles:    true,
		AllowedPathRegex: "",
		CleanSync:        false,
	}

	returnedConfig, err := attributevalue.MarshalMap(config)
	if err != nil {
		t.Fatal(err)
	}
	mockedConfigurationFetcher.On("GetItem", mock.Anything, mock.Anything).Return(&awsdynamodb.GetItemOutput{
		Item: returnedConfig,
	}, nil)

	mockedStateTracking := &MockDynamodb{}
	syncState := syncstate.SyncStateRow{
		PrimaryKey: dynamodb.PrimaryKey{PartitionKey: "Machine#AAAAAAAA-A00A-1234-1234-5864377B4831", SortKey: "SyncState"},
		SyncState: syncstate.SyncState{
			MachineID:      inputMachineID,
			BatchSize:      37,
			LastCleanSync:  "2000-12-31T00:00:00Z",
			FeedSyncCursor: "2000-12-31T00:00:00Z",
			DataType:       types.DataTypeSyncState,
		},
	}
	returnedSyncState, err := attributevalue.MarshalMap(syncState)
	if err != nil {
		t.Fatal(err)
	}
	mockedStateTracking.On("GetItem", mock.Anything, mock.Anything).Return(&awsdynamodb.GetItemOutput{
		Item: returnedSyncState,
	}, nil)

	// mockedStateTracking.On("PutItem", mock.MatchedBy(func(item interface{}) bool {
	mockedStateTracking.On("PutItem", mock.MatchedBy(func(syncState syncstate.SyncStateRow) bool {
		return syncState.PrimaryKey.PartitionKey == "Machine#AAAAAAAA-A00A-1234-1234-5864377B4831" && syncState.MachineID == inputMachineID && syncState.BatchSize == 37 && syncState.LastCleanSync == "2000-12-31T00:00:00Z" && syncState.FeedSyncCursor == "2000-12-31T00:00:00Z" && syncState.CleanSync == false
	})).Return(&awsdynamodb.PutItemOutput{}, nil)

	mockedStateTracking.On("PutItem", mock.MatchedBy(func(sensorData sensordata.SensorData) bool {
		return sensorData.OSBuild == "20D5029f" && sensorData.SerialNum == "C02123456789" && sensorData.MachineID == inputMachineID && sensorData.PrimaryUser == "nobody" && sensorData.BinaryRuleCount == 3 && sensorData.CompilerRuleCount == 0
	})).Return(&awsdynamodb.PutItemOutput{}, nil)

	h := &PostPreflightHandler{
		timeProvider:                timeProvider,
		machineConfigurationService: machineconfiguration.GetMachineConfigurationService(mockedConfigurationFetcher, timeProvider),
		stateTrackingService:        getStateTrackingService(mockedStateTracking, timeProvider),
		cleanSyncService:            getCleanSyncService(timeProvider),
	}

	resp, err := h.Handle(request)

	assert.Empty(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Ensure that the response matches the configuration returned
	assert.Equal(t, `{"client_mode":"LOCKDOWN","blocked_path_regex":"","allowed_path_regex":"","batch_size":37,"enable_bundles":true,"enable_transitive_rules":false,"upload_logs_url":"/aaa","sync_type":"normal"}`, resp.Body)
}
