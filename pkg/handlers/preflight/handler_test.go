package preflight

import (
	"testing"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/model/machineconfiguration"
	"github.com/airbnb/rudolph/pkg/model/syncstate"
	"github.com/airbnb/rudolph/pkg/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
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

type mockSensorDataSaver func(time clock.TimeProvider, machineID string, request *PreflightRequest) error

func (m mockSensorDataSaver) saveSensorDataFromRequest(time clock.TimeProvider, machineID string, request *PreflightRequest) error {
	return m(time, machineID, request)
}

type mockMachineConfigurationGetter func(machineID string) (config machineconfiguration.MachineConfiguration, err error)

func (m mockMachineConfigurationGetter) getDesiredConfig(machineID string) (config machineconfiguration.MachineConfiguration, err error) {
	return m(machineID)
}

type mockSyncStateManager struct {
	a func(machineID string) (syncState *syncstate.SyncStateRow, err error)
	b func(time clock.TimeProvider, machineID string, requestCleanSync bool, lastCleanSync string, batchSize int, feedSyncCursor string) error
}

func (m mockSyncStateManager) getSyncState(machineID string) (syncState *syncstate.SyncStateRow, err error) {
	return m.a(machineID)
}
func (m mockSyncStateManager) saveNewSyncState(time clock.TimeProvider, machineID string, requestCleanSync bool, lastCleanSync string, batchSize int, feedSyncCursor string) error {
	return m.b(time, machineID, requestCleanSync, lastCleanSync, batchSize, feedSyncCursor)
}

// Coerce the mock types
var _ sensorDataSaver = mockSensorDataSaver(nil)
var _ machineConfigurationGetter = mockMachineConfigurationGetter(nil)
var _ syncStateManager = mockSyncStateManager{}

// OK
// Tests the basic positive flow.
func TestHandler_OK(t *testing.T) {
	now, _ := clock.ParseRFC3339("2000-01-01T00:00:00Z")
	inputMachineID := "AAAAAAAA-A00A-1234-1234-5864377B4831"
	var request = events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Resource:       "/eventupload/{machine_id}",
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

	h := &PostPreflightHandler{
		daysElapseUntilCleanSync: 6,
		timeProvider: clock.FrozenTimeProvider{
			Current: now,
		},
		sensorDataSaver: mockSensorDataSaver(
			func(time clock.TimeProvider, machineID string, request *PreflightRequest) error {
				// Ensures that the request body is parsed and saved properly.
				assert.Equal(t, inputMachineID, machineID)
				assert.Equal(t, "nobody", request.PrimaryUser)
				assert.Equal(t, types.Monitor, request.ClientMode)
				assert.Equal(t, "C02123456789", request.SerialNumber)
				assert.Equal(t, "20D5029f", request.OSBuild)
				assert.Equal(t, "11.2", request.OSVersion)
				assert.Equal(t, 3, request.BinaryRuleCount)
				assert.Equal(t, 2, request.CertificateRuleCount)
				assert.Equal(t, time.Now(), now)

				return nil
			},
		),
		machineConfigurationGetter: mockMachineConfigurationGetter(
			func(machineID string) (config machineconfiguration.MachineConfiguration, err error) {
				assert.Equal(t, inputMachineID, machineID)
				return machineconfiguration.MachineConfiguration{
					ClientMode:       types.Lockdown,
					BatchSize:        37,
					UploadLogsURL:    "/aaa",
					EnableBundles:    true,
					AllowedPathRegex: "",
					CleanSync:        false,
				}, nil
			},
		),
		syncStateManager: mockSyncStateManager{
			a: func(machineID string) (syncState *syncstate.SyncStateRow, err error) {
				assert.Equal(t, inputMachineID, machineID)
				return nil, nil
			},
			b: func(time clock.TimeProvider, machineID string, requestCleanSync bool, lastCleanSync string, batchSize int, feedSyncCursor string) error {
				assert.Equal(t, inputMachineID, machineID)
				assert.Equal(t, 37, batchSize)
				assert.True(t, requestCleanSync)
				assert.Equal(t, "2000-01-01T00:00:00Z", lastCleanSync)
				assert.Equal(t, "2000-01-01T00:00:00Z", feedSyncCursor)
				assert.Equal(t, now, time.Now())

				return nil
			},
		},
	}

	resp, err := h.Handle(request)

	assert.Empty(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Ensure that the response matches the configuration returned
	assert.Equal(t, `{"client_mode":"LOCKDOWN","blocked_path_regex":"","allowed_path_regex":"","batch_size":37,"enable_bundles":true,"enable_transitive_rules":false,"clean_sync":true,"upload_logs_url":"/aaa"}`, resp.Body)
}

// func TestPreflightHandler_OK(t *testing.T) {
// 	var putSensorData *dynamodb.PutItemInput
// 	var putItemCall *dynamodb.PutItemInput
// 	var updateItemCall *dynamodb.UpdateItemInput
// 	var getItemCall *dynamodb.GetItemInput
// 	var getItemCallFirstRun *dynamodb.GetItemInput
// 	var getItemCallPostInitialRun *dynamodb.GetItemInput
// 	var getItemCallPastSyncExpirationRun *dynamodb.GetItemInput

// 	// Track the current time when the test is run
// 	currentTime := time.Now().UTC().Format(time.RFC3339)

// 	// Track a time in the past that creates a mock environment where a sync operation occurred two days ago and a clean sync one day ago
// 	lastSyncTimePostInitial := time.Now().UTC().AddDate(0, 0, -2).Format(time.RFC3339)
// 	lastCleanSyncTimePostInitial := time.Now().UTC().AddDate(0, 0, -1).Format(time.RFC3339)

// 	// Track a time in the past that creates a mock environment where a sync operation occurred 31 days ago and a clean sync 30 days ago
// 	lastSyncTimePastSyncExpiration := time.Now().UTC().AddDate(0, 0, -31).Format(time.RFC3339)
// 	lastCleanSyncTimePastSyncExpiration := time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)

// 	machineConfigInitial := map[string]*dynamodb.AttributeValue{
// 		"PK":                    {S: aws.String("MachineConfig#AAAAAAAA-A00A-1234-1234-5864377B0001")},
// 		"SK":                    {S: aws.String("Current")},
// 		"ClientMode":            {N: aws.String("1")},
// 		"BlockedPathRegex":      {S: aws.String("")},
// 		"AllowedPathRegex":      {S: aws.String("")},
// 		"BatchSize":             {N: aws.String("01")},
// 		"EnableBundles":         {BOOL: aws.Bool(false)},
// 		"EnableTransitiveRules": {BOOL: aws.Bool(false)},
// 	}

// 	// Sync State for testing
// 	syncStatePostInitial := map[string]*dynamodb.AttributeValue{
// 		"PK":             {S: aws.String("MachineInfo#AAAAAAAA-A00A-1234-1234-5864377B0002")},
// 		"SK":             {S: aws.String("SyncState")},
// 		"MachineID":      {S: aws.String("AAAAAAAA-A00A-1234-1234-5864377B0002")},
// 		"BatchSize":      {N: aws.String("2")},
// 		"LastCleanSync":  {S: aws.String(lastCleanSyncTimePostInitial)},
// 		"FeedSyncCursor": {S: aws.String(lastSyncTimePostInitial)},
// 		"PreflightAt":    {S: aws.String(currentTime)},
// 	}

// 	machineConfigPostInitial := map[string]*dynamodb.AttributeValue{
// 		"PK":                    {S: aws.String("MachineConfig#AAAAAAAA-A00A-1234-1234-5864377B0002")},
// 		"SK":                    {S: aws.String("Current")},
// 		"ClientMode":            {N: aws.String("1")},
// 		"BlockedPathRegex":      {S: aws.String("")},
// 		"AllowedPathRegex":      {S: aws.String("")},
// 		"BatchSize":             {N: aws.String("02")},
// 		"EnableBundles":         {BOOL: aws.Bool(false)},
// 		"EnableTransitiveRules": {BOOL: aws.Bool(false)},
// 	}

// 	// Sync State for testing
// 	syncStatePastSyncExpiration := map[string]*dynamodb.AttributeValue{
// 		"PK":             {S: aws.String("MachineInfo#AAAAAAAA-A00A-1234-1234-5864377B0030")},
// 		"SK":             {S: aws.String("SyncState")},
// 		"MachineID":      {S: aws.String("AAAAAAAA-A00A-1234-1234-5864377B0030")},
// 		"BatchSize":      {N: aws.String("30")},
// 		"LastCleanSync":  {S: aws.String(lastCleanSyncTimePastSyncExpiration)},
// 		"FeedSyncCursor": {S: aws.String(lastSyncTimePastSyncExpiration)},
// 		"PreflightAt":    {S: aws.String(currentTime)},
// 	}

// 	machineConfigPastSyncExpiration := map[string]*dynamodb.AttributeValue{
// 		"PK":                    {S: aws.String("MachineConfig#AAAAAAAA-A00A-1234-1234-5864377B0030")},
// 		"SK":                    {S: aws.String("Current")},
// 		"ClientMode":            {N: aws.String("1")},
// 		"BlockedPathRegex":      {S: aws.String("")},
// 		"AllowedPathRegex":      {S: aws.String("")},
// 		"BatchSize":             {N: aws.String("30")},
// 		"EnableBundles":         {BOOL: aws.Bool(false)},
// 		"EnableTransitiveRules": {BOOL: aws.Bool(false)},
// 	}

// 	store.ActivateMock(store.MockDBCallbackClient{
// 		PutItemFunction: func(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
// 			switch *in.Item["SK"].S {
// 			case "Current": // SensorData
// 				putSensorData = in
// 			case "SyncState": // SyncState
// 				putItemCall = in
// 			}

// 			return &dynamodb.PutItemOutput{}, nil
// 		},
// 		UpdateItemFunction: func(in *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
// 			updateItemCall = in
// 			return &dynamodb.UpdateItemOutput{}, nil
// 		},
// 		GetItemFunction: func(in *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
// 			if *in.Key["PK"].S == "MachineInfo#AAAAAAAA-A00A-1234-1234-5864377B0002" {
// 				getItemCall = in
// 				return &dynamodb.GetItemOutput{
// 					Item: syncStatePostInitial,
// 				}, nil
// 			} else if *in.Key["PK"].S == "MachineInfo#AAAAAAAA-A00A-1234-1234-5864377B0030" {
// 				getItemCall = in
// 				return &dynamodb.GetItemOutput{
// 					Item: syncStatePastSyncExpiration,
// 				}, nil
// 			} else if *in.Key["PK"].S == "MachineConfig#AAAAAAAA-A00A-1234-1234-5864377B0001" {
// 				getItemCallFirstRun = in
// 				return &dynamodb.GetItemOutput{
// 					Item: machineConfigInitial,
// 				}, nil
// 			} else if *in.Key["PK"].S == "MachineConfig#AAAAAAAA-A00A-1234-1234-5864377B0002" {
// 				getItemCallPostInitialRun = in
// 				return &dynamodb.GetItemOutput{
// 					Item: machineConfigPostInitial,
// 				}, nil
// 			} else if *in.Key["PK"].S == "MachineConfig#AAAAAAAA-A00A-1234-1234-5864377B0030" {
// 				getItemCallPastSyncExpirationRun = in
// 				return &dynamodb.GetItemOutput{
// 					Item: machineConfigPastSyncExpiration,
// 				}, nil
// 			} else if strings.HasPrefix(*in.Key["PK"].S, "MachineConfig#") {
// 				getItemCall = in
// 				return &dynamodb.GetItemOutput{Item: map[string]*dynamodb.AttributeValue{
// 					"PK":                    {S: aws.String("MachineConfig#AAAAAAAA-A00A-1234-1234-5864377B4831")},
// 					"SK":                    {S: aws.String("Current")},
// 					"ClientMode":            {N: aws.String("1")},
// 					"BlockedPathRegex":      {S: aws.String("")},
// 					"AllowedPathRegex":      {S: aws.String("")},
// 					"BatchSize":             {N: aws.String("50")},
// 					"EnableBundles":         {BOOL: aws.Bool(false)},
// 					"EnableTransitiveRules": {BOOL: aws.Bool(false)},
// 				}}, nil
// 			} else {
// 				return &dynamodb.GetItemOutput{Item: map[string]*dynamodb.AttributeValue{}}, nil
// 			}
// 		},
// 	})

// 	t.Run("Not-Clean Sync - first run", func(t *testing.T) {
// 		var request = events.APIGatewayProxyRequest{
// 			HTTPMethod:     "POST",
// 			Resource:       "/preflight/{machine_id}",
// 			PathParameters: map[string]string{"machine_id": "AAAAAAAA-A00A-1234-1234-5864377B0001"},
// 			Headers:        map[string]string{"Content-Type": "application/json"},
// 			Body: `{
// 	"os_build":"20D5029f",
// 	"santa_version":"2021.1",
// 	"hostname":"my-awesome-macbook-pro.attlocal.net",
// 	"transitive_rule_count":0,
// 	"os_version":"11.2",
// 	"certificate_rule_count":0,
// 	"client_mode":"MONITOR",
// 	"serial_num":"C02123456789",
// 	"binary_rule_count":0,
// 	"primary_user":"nobody",
// 	"compiler_rule_count":0
// }`,
// 		}

// 		h := &PostPreflightHandler{}
// 		resp, _ := h.Handle(request)

// 		// Assert that we save the sensordata
// 		assert.Contains(t, putSensorData.Item, "PK")
// 		assert.Contains(t, putSensorData.Item, "SK")
// 		assert.Contains(t, putSensorData.Item, "OSVersion")
// 		assert.Equal(t, "11.2", *putSensorData.Item["OSVersion"].S)
// 		assert.Contains(t, putSensorData.Item, "SerialNum")
// 		assert.Equal(t, "C02123456789", *putSensorData.Item["SerialNum"].S)
// 		assert.Contains(t, putSensorData.Item, "OSBuild")
// 		assert.Equal(t, "20D5029f", *putSensorData.Item["OSBuild"].S)
// 		assert.Contains(t, putSensorData.Item, "RequestCleanSync")
// 		assert.Equal(t, false, *putSensorData.Item["RequestCleanSync"].BOOL)
// 		assert.Contains(t, putSensorData.Item, "PrimaryUser")
// 		assert.Equal(t, "nobody", *putSensorData.Item["PrimaryUser"].S)
// 		assert.Contains(t, putSensorData.Item, "MachineID")
// 		assert.Equal(t, "AAAAAAAA-A00A-1234-1234-5864377B0001", *putSensorData.Item["MachineID"].S)

// 		// Assert that we create a new sync state
// 		assert.Contains(t, putItemCall.Item, "PK")
// 		assert.Contains(t, putItemCall.Item, "SK")
// 		assert.Contains(t, putItemCall.Item, "FeedSyncCursor")
// 		assert.Contains(t, putItemCall.Item, "PreflightAt")
// 		assert.Contains(t, putItemCall.Item, "LastCleanSync")
// 		assert.Contains(t, putItemCall.Item, "MachineID")
// 		assert.Contains(t, putItemCall.Item, "BatchSize")
// 		// Since LastCleanSync was never run before, this will contain the same time of the first preflight
// 		assert.Contains(t, putItemCall.Item, "CleanSync")
// 		assert.Equal(t, "MachineInfo#AAAAAAAA-A00A-1234-1234-5864377B0001", *putItemCall.Item["PK"].S)
// 		assert.Equal(t, "SyncState", *putItemCall.Item["SK"].S)
// 		assert.Equal(t, "AAAAAAAA-A00A-1234-1234-5864377B0001", *putItemCall.Item["MachineID"].S)
// 		assert.Equal(t, "1", *putItemCall.Item["BatchSize"].N)

// 		// Assert nothing gets updated
// 		assert.Empty(t, updateItemCall)

// 		// Assert that when we query machineConfiguration consistently
// 		assert.Equal(t, "MachineConfig#AAAAAAAA-A00A-1234-1234-5864377B0001", *getItemCallFirstRun.Key["PK"].S)
// 		assert.Equal(t, "Current", *getItemCallFirstRun.Key["SK"].S)
// 		assert.False(t, *getItemCallFirstRun.ConsistentRead)

// 		assert.Equal(t, 200, resp.StatusCode)
// 		assert.Equal(t, `{"client_mode":"MONITOR","blocked_path_regex":"","allowed_path_regex":"","batch_size":1,"enable_bundles":false,"enable_transitive_rules":false,"clean_sync":true}`, resp.Body)
// 	})

// 	t.Run("Not-Clean Sync - Post-Initial-Run", func(t *testing.T) {
// 		var request = events.APIGatewayProxyRequest{
// 			HTTPMethod:     "POST",
// 			Resource:       "/preflight/{machine_id}",
// 			PathParameters: map[string]string{"machine_id": "AAAAAAAA-A00A-1234-1234-5864377B0002"},
// 			Headers:        map[string]string{"Content-Type": "application/json"},
// 			Body: `{
// 	"os_build":"20D5029f",
// 	"santa_version":"2021.1",
// 	"hostname":"my-awesome-macbook-pro.attlocal.net",
// 	"transitive_rule_count":0,
// 	"os_version":"11.2",
// 	"certificate_rule_count":0,
// 	"client_mode":"MONITOR",
// 	"serial_num":"C02123456789",
// 	"binary_rule_count":0,
// 	"primary_user":"nobody",
// 	"compiler_rule_count":0
// }`,
// 		}
// 		h := &PostPreflightHandler{}
// 		resp, _ := h.Handle(request)

// 		// Assert that we save the sensordata
// 		assert.Contains(t, putSensorData.Item, "PK")
// 		assert.Contains(t, putSensorData.Item, "SK")
// 		assert.Contains(t, putSensorData.Item, "OSVersion")
// 		assert.Equal(t, "11.2", *putSensorData.Item["OSVersion"].S)
// 		assert.Contains(t, putSensorData.Item, "SerialNum")
// 		assert.Equal(t, "C02123456789", *putSensorData.Item["SerialNum"].S)
// 		assert.Contains(t, putSensorData.Item, "OSBuild")
// 		assert.Equal(t, "20D5029f", *putSensorData.Item["OSBuild"].S)
// 		assert.Contains(t, putSensorData.Item, "RequestCleanSync")
// 		assert.Equal(t, false, *putSensorData.Item["RequestCleanSync"].BOOL)
// 		assert.Contains(t, putSensorData.Item, "PrimaryUser")
// 		assert.Equal(t, "nobody", *putSensorData.Item["PrimaryUser"].S)
// 		assert.Contains(t, putSensorData.Item, "MachineID")
// 		assert.Equal(t, "AAAAAAAA-A00A-1234-1234-5864377B0002", *putSensorData.Item["MachineID"].S)

// 		// Assert that we create a new sync state
// 		assert.Contains(t, putItemCall.Item, "PK")
// 		assert.Contains(t, putItemCall.Item, "SK")
// 		assert.Contains(t, putItemCall.Item, "FeedSyncCursor")
// 		assert.Contains(t, putItemCall.Item, "PreflightAt")
// 		assert.Contains(t, putItemCall.Item, "LastCleanSync")

// 		assert.Equal(t, lastCleanSyncTimePostInitial, *putItemCall.Item["LastCleanSync"].S)
// 		assert.Contains(t, putItemCall.Item, "MachineID")
// 		assert.Contains(t, putItemCall.Item, "BatchSize")
// 		assert.NotContains(t, putItemCall.Item, "CleanSync")
// 		assert.Equal(t, "MachineInfo#AAAAAAAA-A00A-1234-1234-5864377B0002", *putItemCall.Item["PK"].S)
// 		assert.Equal(t, "SyncState", *putItemCall.Item["SK"].S)
// 		assert.Equal(t, "AAAAAAAA-A00A-1234-1234-5864377B0002", *putItemCall.Item["MachineID"].S)
// 		assert.Equal(t, "2", *putItemCall.Item["BatchSize"].N)

// 		// Assert that when we query sync state it works
// 		assert.Equal(t, "MachineInfo#AAAAAAAA-A00A-1234-1234-5864377B0002", *getItemCall.Key["PK"].S)

// 		// Assert nothing gets updated
// 		assert.Empty(t, updateItemCall)

// 		// Assert that when we query machineConfiguration consistently
// 		assert.Equal(t, "MachineConfig#AAAAAAAA-A00A-1234-1234-5864377B0002", *getItemCallPostInitialRun.Key["PK"].S)
// 		assert.Equal(t, "Current", *getItemCallPostInitialRun.Key["SK"].S)
// 		assert.False(t, *getItemCallPostInitialRun.ConsistentRead)

// 		assert.Equal(t, 200, resp.StatusCode)
// 		assert.Equal(t, `{"client_mode":"MONITOR","blocked_path_regex":"","allowed_path_regex":"","batch_size":2,"enable_bundles":false,"enable_transitive_rules":false}`, resp.Body)
// 	})

// 	t.Run("Not-Clean Sync - 30+ days - clean sync should be forced", func(t *testing.T) {
// 		var request = events.APIGatewayProxyRequest{
// 			HTTPMethod:     "POST",
// 			Resource:       "/preflight/{machine_id}",
// 			PathParameters: map[string]string{"machine_id": "AAAAAAAA-A00A-1234-1234-5864377B0030"},
// 			Headers:        map[string]string{"Content-Type": "application/json"},
// 			Body: `{
// 	"os_build":"20D5029f",
// 	"santa_version":"2021.1",
// 	"hostname":"my-awesome-macbook-pro.attlocal.net",
// 	"transitive_rule_count":0,
// 	"os_version":"11.2",
// 	"certificate_rule_count":0,
// 	"client_mode":"MONITOR",
// 	"serial_num":"C02123456789",
// 	"binary_rule_count":0,
// 	"primary_user":"nobody",
// 	"compiler_rule_count":0
// }`,
// 		}

// 		h := &PostPreflightHandler{}
// 		resp, _ := h.Handle(request)

// 		// Assert that we save the sensordata
// 		assert.Contains(t, putSensorData.Item, "PK")
// 		assert.Contains(t, putSensorData.Item, "SK")
// 		assert.Contains(t, putSensorData.Item, "OSVersion")
// 		assert.Equal(t, "11.2", *putSensorData.Item["OSVersion"].S)
// 		assert.Contains(t, putSensorData.Item, "SerialNum")
// 		assert.Equal(t, "C02123456789", *putSensorData.Item["SerialNum"].S)
// 		assert.Contains(t, putSensorData.Item, "OSBuild")
// 		assert.Equal(t, "20D5029f", *putSensorData.Item["OSBuild"].S)
// 		assert.Contains(t, putSensorData.Item, "RequestCleanSync")
// 		assert.Equal(t, false, *putSensorData.Item["RequestCleanSync"].BOOL)
// 		assert.Contains(t, putSensorData.Item, "PrimaryUser")
// 		assert.Equal(t, "nobody", *putSensorData.Item["PrimaryUser"].S)
// 		assert.Contains(t, putSensorData.Item, "MachineID")
// 		assert.Equal(t, "AAAAAAAA-A00A-1234-1234-5864377B0030", *putSensorData.Item["MachineID"].S)

// 		// Assert that we create a new sync state
// 		assert.Contains(t, putItemCall.Item, "PK")
// 		assert.Contains(t, putItemCall.Item, "SK")
// 		assert.Contains(t, putItemCall.Item, "FeedSyncCursor")
// 		assert.Contains(t, putItemCall.Item, "PreflightAt")
// 		assert.Contains(t, putItemCall.Item, "LastCleanSync")
// 		assert.Equal(t, lastCleanSyncTimePastSyncExpiration, *putItemCall.Item["LastCleanSync"].S)
// 		assert.Contains(t, putItemCall.Item, "MachineID")
// 		assert.Contains(t, putItemCall.Item, "BatchSize")
// 		// Since 30 days have elapsed, CleanSync will be requested
// 		assert.Contains(t, putItemCall.Item, "CleanSync")
// 		assert.Equal(t, "MachineInfo#AAAAAAAA-A00A-1234-1234-5864377B0030", *putItemCall.Item["PK"].S)
// 		assert.Equal(t, "SyncState", *putItemCall.Item["SK"].S)
// 		assert.Equal(t, "AAAAAAAA-A00A-1234-1234-5864377B0030", *putItemCall.Item["MachineID"].S)
// 		assert.Equal(t, "30", *putItemCall.Item["BatchSize"].N)

// 		// Assert nothing gets updated
// 		assert.Empty(t, updateItemCall)

// 		// Assert that when we query machineConfiguration consistently
// 		assert.Equal(t, "MachineConfig#AAAAAAAA-A00A-1234-1234-5864377B0030", *getItemCallPastSyncExpirationRun.Key["PK"].S)
// 		assert.Equal(t, "Current", *getItemCallPastSyncExpirationRun.Key["SK"].S)
// 		assert.False(t, *getItemCallPastSyncExpirationRun.ConsistentRead)

// 		assert.Equal(t, 200, resp.StatusCode)
// 		assert.Equal(t, `{"client_mode":"MONITOR","blocked_path_regex":"","allowed_path_regex":"","batch_size":30,"enable_bundles":false,"enable_transitive_rules":false,"clean_sync":true}`, resp.Body)
// 	})

// 	t.Run("Clean Sync", func(t *testing.T) {
// 		var request = events.APIGatewayProxyRequest{
// 			HTTPMethod:     "POST",
// 			Resource:       "/preflight/{machine_id}",
// 			PathParameters: map[string]string{"machine_id": "AAAAAAAA-A00A-1234-1234-5864377B4831"},
// 			Headers:        map[string]string{"Content-Type": "application/json"},
// 			Body: `{
// 	"os_build":"20D5029f",
// 	"santa_version":"2021.1",
// 	"hostname":"my-awesome-macbook-pro.attlocal.net",
// 	"transitive_rule_count":0,
// 	"os_version":"11.2",
// 	"certificate_rule_count":0,
// 	"client_mode":"MONITOR",
// 	"serial_num":"C02123456789",
// 	"binary_rule_count":0,
// 	"primary_user":"",
// 	"compiler_rule_count":0,
// 	"request_clean_sync":true
// }`,
// 		}

// 		h := &PostPreflightHandler{}
// 		resp, _ := h.Handle(request)

// 		// Assert that we create a new sync state with cleansync = true
// 		assert.Contains(t, putItemCall.Item, "CleanSync")
// 		assert.Equal(t, true, *putItemCall.Item["CleanSync"].BOOL)

// 		// Ensure "clean_sync":true is present in the response
// 		assert.Equal(t, 200, resp.StatusCode)
// 		assert.Equal(t, `{"client_mode":"MONITOR","blocked_path_regex":"","allowed_path_regex":"","batch_size":50,"enable_bundles":false,"enable_transitive_rules":false,"clean_sync":true}`, resp.Body)
// 	})
// }
