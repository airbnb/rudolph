package preflight

import (
	"fmt"
	"testing"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/sensordata"
	"github.com/airbnb/rudolph/pkg/model/syncstate"
	rudolphtypes "github.com/airbnb/rudolph/pkg/types"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

type mockPutter func(item interface{}) (*awsdynamodb.PutItemOutput, error)

func (m mockPutter) PutItem(item interface{}) (*awsdynamodb.PutItemOutput, error) { return m(item) }

var _ dynamodb.PutItemAPI = mockPutter(nil)

func Test_concreteSensorDataSaver_OK(t *testing.T) {
	type test struct {
		comment              string
		hostname             string
		machineID            string
		serialNumber         string
		osBuild              string
		osVersion            string
		santaVersion         string
		binaryRuleCount      int
		certRuleCount        int
		transitiveRuleCount  int
		compilerRuleCount    int
		ruleCount            int
		primaryUser          string
		expectedTime         string
		expectedExpiresAfter int64
		expectedDataType     rudolphtypes.DataType
	}

	timeProvider := clock.ConcreteTimeProvider{}
	machineID := "AAAAAAAA-A00A-1234-1234-5864377B4831"
	pk, sk := sensordata.MachineIDSensorDataPKSK(machineID)

	expected := test{
		comment:              fmt.Sprintf("%s %s", "Testing", machineID),
		hostname:             "macbook.pro.localhost",
		machineID:            machineID,
		osBuild:              "20A21",
		osVersion:            "12.34",
		santaVersion:         "2021.1",
		serialNumber:         "123456789ABC",
		binaryRuleCount:      4,
		certRuleCount:        3,
		transitiveRuleCount:  2,
		compilerRuleCount:    1,
		ruleCount:            10,
		primaryUser:          "john_doe",
		expectedTime:         clock.RFC3339(timeProvider.Now()),
		expectedExpiresAfter: clock.Unixtimestamp(timeProvider.Now().UTC().AddDate(0, 0, 90)),
		expectedDataType:     rudolphtypes.DataTypeSensorData,
	}

	saver := concreteSensorDataSaver{
		putter: mockPutter(
			func(item interface{}) (*awsdynamodb.PutItemOutput, error) {
				sensorData := item.(sensordata.SensorData)
				assert.Equal(t, expected.machineID, sensorData.MachineID)
				assert.Equal(t, expected.osBuild, sensorData.OSBuild)
				assert.Equal(t, expected.osVersion, sensorData.OSVersion)
				assert.Equal(t, expected.binaryRuleCount, sensorData.BinaryRuleCount)
				assert.Equal(t, expected.certRuleCount, sensorData.CertificateRuleCount)
				assert.Equal(t, expected.compilerRuleCount, sensorData.CompilerRuleCount)
				assert.Equal(t, expected.transitiveRuleCount, sensorData.TransitiveRuleCount)
				assert.Equal(t, expected.ruleCount, sensorData.RuleCount)
				assert.Equal(t, expected.primaryUser, sensorData.PrimaryUser)
				assert.Equal(t, expected.serialNumber, sensorData.SerialNum)
				assert.Equal(t, expected.expectedDataType, sensorData.DataType)
				// assert.Equal(t, expected.santaVersion, sensorData.SantaVersion)
				assert.Equal(t, pk, sensorData.PartitionKey)
				assert.Equal(t, sk, sensorData.SortKey)
				// assert.Equal(t, "osversion", sensorData.Hostname) // Missing??
				// assert.Equal(t, types.Lockdown, sensorData.ClientMode) // missing??

				return &awsdynamodb.PutItemOutput{}, nil
			},
		),
	}
	request := &PreflightRequest{
		OSBuild:              expected.osBuild,
		OSVersion:            expected.osVersion,
		Hostname:             expected.hostname,
		SantaVersion:         expected.santaVersion,
		ClientMode:           rudolphtypes.Lockdown,
		BinaryRuleCount:      expected.binaryRuleCount,
		CertificateRuleCount: expected.certRuleCount,
		TransitiveRuleCount:  expected.transitiveRuleCount,
		CompilerRuleCount:    expected.compilerRuleCount,
		PrimaryUser:          expected.primaryUser,
		SerialNumber:         expected.serialNumber,
	}
	err := saver.saveSensorDataFromRequest(timeProvider, machineID, request)

	assert.Empty(t, err)
}

func Test_machineIDToInt(t *testing.T) {
	type test struct {
		comment            string
		machineID          string
		expectMachineIDInt int
		expectToMatch      bool
	}
	// Verified randomized output //
	// MachineID: 18bd9617-dfe7-4239-9c1f-dd986e0e7647 | Result: 0
	// MachineID: cf5181ed-2941-4a11-bfc6-6442b3684a09 | Result: 1
	// MachineID: 89e6a5ae-d04d-4662-a285-e0386a2a613c | Result: 3
	// MachineID: 297d2635-8296-4a79-be0b-c577ab6da363 | Result: 10

	tests := []test{
		{
			comment:            "Machine ID Matches 6",
			machineID:          "52c9e6f1-046d-46db-9bc0-0fe142920093",
			expectMachineIDInt: 6,
			expectToMatch:      true,
		},
		{
			comment:            "Machine ID Matches 0",
			machineID:          "18bd9617-dfe7-4239-9c1f-dd986e0e7647",
			expectMachineIDInt: 0,
			expectToMatch:      true,
		},
		{
			comment:            "Machine ID Matches 1",
			machineID:          "cf5181ed-2941-4a11-bfc6-6442b3684a09",
			expectMachineIDInt: 1,
			expectToMatch:      true,
		},
		{
			comment:            "Machine ID Matches 10",
			machineID:          "297d2635-8296-4a79-be0b-c577ab6da363",
			expectMachineIDInt: 10,
			expectToMatch:      true,
		},
		{
			comment:            "Machine ID Does Not Match 3",
			machineID:          "89e6a5ae-d04d-4662-a285-e0386a2a613c",
			expectMachineIDInt: 1,
			expectToMatch:      false,
		},
	}

	for _, tc := range tests {
		result, _ := machineIDToInt(tc.machineID)

		if !tc.expectToMatch {
			assert.NotEqual(t, result, tc.expectMachineIDInt, tc.comment)
		} else {
			assert.Equal(t, tc.expectMachineIDInt, result, tc.comment)
		}
	}
}

func Test_performCleanSync(t *testing.T) {
	var daysElapseUntilCleanSync int = 7

	type test struct {
		comment                  string
		machineID                string
		daysElapseUntilCleanSync int
		daysSinceLastSync        int
		expectToCleanSync        bool
	}
	// Verified randomized output //
	// MachineID: 18bd9617-dfe7-4239-9c1f-dd986e0e7647 | Result: 0
	// MachineID: 89e6a5ae-d04d-4662-a285-e0386a2a613c | Result: 3
	// MachineID: 297d2635-8296-4a79-be0b-c577ab6da363 | Result: 10

	tests := []test{
		{
			comment:                  "No Clean Sync - Last clean sync occurred 0 days ago - machine ID converted int is 0",
			machineID:                "18bd9617-dfe7-4239-9c1f-dd986e0e7647",
			daysElapseUntilCleanSync: daysElapseUntilCleanSync,
			daysSinceLastSync:        0,
			expectToCleanSync:        false,
		},
		{
			comment:                  "No Clean Sync - Last clean sync occurred 6 days ago - machine ID converted int is 0",
			machineID:                "18bd9617-dfe7-4239-9c1f-dd986e0e7647",
			daysElapseUntilCleanSync: daysElapseUntilCleanSync,
			daysSinceLastSync:        6,
			expectToCleanSync:        false,
		},
		{
			comment:                  "No Clean Sync - Last clean sync occurred 6 days ago - machine ID converted int is 3",
			machineID:                "89e6a5ae-d04d-4662-a285-e0386a2a613c",
			daysElapseUntilCleanSync: daysElapseUntilCleanSync,
			daysSinceLastSync:        6,
			expectToCleanSync:        false,
		},
		{
			comment:                  "No Clean Sync - Last clean sync occurred 1 day ago - machine ID converted int is 10",
			machineID:                "297d2635-8296-4a79-be0b-c577ab6da363",
			daysElapseUntilCleanSync: daysElapseUntilCleanSync,
			daysSinceLastSync:        1,
			expectToCleanSync:        false,
		},
		{
			comment:                  "No Clean Sync - Last clean sync occurred 7 days ago - machine ID converted int is 10",
			machineID:                "297d2635-8296-4a79-be0b-c577ab6da363",
			daysElapseUntilCleanSync: daysElapseUntilCleanSync,
			daysSinceLastSync:        7,
			expectToCleanSync:        false,
		},
		{
			comment:                  "Expect Clean Sync - Last clean sync occurred 10 days ago - machine ID converted int is 3",
			machineID:                "89e6a5ae-d04d-4662-a285-e0386a2a613c",
			daysElapseUntilCleanSync: daysElapseUntilCleanSync,
			daysSinceLastSync:        10,
			expectToCleanSync:        true,
		},
	}

	for _, tc := range tests {
		result, err := shouldPerformCleanSync(tc.machineID, tc.daysSinceLastSync, tc.daysElapseUntilCleanSync)
		assert.Equal(t, nil, err)
		assert.Equal(t, tc.expectToCleanSync, result, tc.comment)
	}
}

func Test_timeSinceLastSync_NothingInDB(t *testing.T) {
	type test struct {
		comment     string
		currentTime string
		syncState   *syncstate.SyncStateRow
		expectation int
	}

	tests := []test{
		{currentTime: "2020-01-01T00:00:00Z", syncState: nil, expectation: 99999999},
		{currentTime: "2020-01-01T00:00:00Z", syncState: &syncstate.SyncStateRow{}, expectation: 99999999},
		{
			comment:     "Same time",
			currentTime: "2020-01-01T00:00:00Z",
			syncState: &syncstate.SyncStateRow{
				SyncState: syncstate.SyncState{
					LastCleanSync: "2020-01-01T00:00:00Z",
				},
			},
			expectation: 0,
		},
		{
			comment:     "1 hour diff",
			currentTime: "2020-01-01T01:00:00Z",
			syncState: &syncstate.SyncStateRow{
				SyncState: syncstate.SyncState{
					LastCleanSync: "2020-01-01T00:00:00Z",
				},
			},
			expectation: 0,
		},
		{
			comment:     "24 hour diff",
			currentTime: "2020-01-02T00:00:00Z",
			syncState: &syncstate.SyncStateRow{
				SyncState: syncstate.SyncState{
					LastCleanSync: "2020-01-01T00:00:00Z",
				},
			},
			expectation: 1,
		},
		{
			comment:     "36 hour diff",
			currentTime: "2020-01-02T12:00:00Z",
			syncState: &syncstate.SyncStateRow{
				SyncState: syncstate.SyncState{
					LastCleanSync: "2020-01-01T00:00:00Z",
				},
			},
			expectation: 1,
		},
		{
			comment:     "72 hour diff",
			currentTime: "2020-01-04T00:00:00Z",
			syncState: &syncstate.SyncStateRow{
				SyncState: syncstate.SyncState{
					LastCleanSync: "2020-01-01T00:00:00Z",
				},
			},
			expectation: 3,
		},
		{
			comment:     "1 month",
			currentTime: "2020-02-01T00:00:00Z",
			syncState: &syncstate.SyncStateRow{
				SyncState: syncstate.SyncState{
					LastCleanSync: "2020-01-01T00:00:00Z",
				},
			},
			expectation: 31,
		},
	}

	for _, tc := range tests {
		currT, _ := clock.ParseRFC3339(tc.currentTime)
		provider := clock.FrozenTimeProvider{Current: currT}
		actual := daysSinceLastSync(provider, tc.syncState)

		assert.Equal(t, tc.expectation, actual, tc.comment)
	}
}
