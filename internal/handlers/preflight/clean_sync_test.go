package preflight

import (
	"testing"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/model/syncstate"
	"github.com/stretchr/testify/assert"
)

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

func Test_determinePeriodicRefreshCleanSync(t *testing.T) {
	// Uses the daysToElapseUntilRefreshCleanSync = 7
	type test struct {
		comment           string
		machineID         string
		daysSinceLastSync int
		expectToCleanSync bool
	}
	// Verified randomized output //
	// MachineID: 18bd9617-dfe7-4239-9c1f-dd986e0e7647 | Result: 0
	// MachineID: 89e6a5ae-d04d-4662-a285-e0386a2a613c | Result: 3
	// MachineID: 297d2635-8296-4a79-be0b-c577ab6da363 | Result: 10

	tests := []test{
		{
			comment:           "No Clean Sync - Last clean sync occurred 0 days ago - machine ID converted int is 0",
			machineID:         "18bd9617-dfe7-4239-9c1f-dd986e0e7647",
			daysSinceLastSync: 0,
			expectToCleanSync: false,
		},
		{
			comment:           "No Clean Sync - Last clean sync occurred 6 days ago - machine ID converted int is 0",
			machineID:         "18bd9617-dfe7-4239-9c1f-dd986e0e7647",
			daysSinceLastSync: 6,
			expectToCleanSync: false,
		},
		{
			comment:           "No Clean Sync - Last clean sync occurred 6 days ago - machine ID converted int is 3",
			machineID:         "89e6a5ae-d04d-4662-a285-e0386a2a613c",
			daysSinceLastSync: 6,
			expectToCleanSync: false,
		},
		{
			comment:           "No Clean Sync - Last clean sync occurred 1 day ago - machine ID converted int is 10",
			machineID:         "297d2635-8296-4a79-be0b-c577ab6da363",
			daysSinceLastSync: 1,
			expectToCleanSync: false,
		},
		{
			comment:           "No Clean Sync - Last clean sync occurred 7 days ago - machine ID converted int is 10",
			machineID:         "297d2635-8296-4a79-be0b-c577ab6da363",
			daysSinceLastSync: 7,
			expectToCleanSync: false,
		},
		{
			comment:           "Expect Clean Sync - Last clean sync occurred 10 days ago - machine ID converted int is 3",
			machineID:         "89e6a5ae-d04d-4662-a285-e0386a2a613c",
			daysSinceLastSync: 10,
			expectToCleanSync: true,
		},
	}

	for _, tc := range tests {
		expectPerformCleanSync, err := determinePeriodicRefreshCleanSync(tc.machineID, tc.daysSinceLastSync)
		assert.Equal(t, nil, err)
		assert.Equal(t, tc.expectToCleanSync, expectPerformCleanSync, tc.comment)
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
		actual := daysSinceLastCleanSync(provider, tc.syncState)

		assert.Equal(t, tc.expectation, actual, tc.comment)
	}
}
