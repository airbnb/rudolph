package syncstate

import (
	"testing"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/types"
	"github.com/stretchr/testify/assert"
)

var frozenTime, _ = clock.ParseRFC3339("2000-01-01T00:00:00Z")
var timeProvider = clock.FrozenTimeProvider{
	Current: frozenTime,
}

func Test_CreateSyncState(t *testing.T) {
	type test struct {
		machineID              string
		timeProvider           clock.TimeProvider
		expectedBatchSize      int
		expectedCleanSync      bool
		expectedLastCleanSync  string
		expectedPreflightAt    string
		expectedFeedSyncCursor string
		expectedExpiresAfter   int64
		expectedDataType       types.DataType
	}

	cases := []test{
		{
			machineID:              "AAAAAAAA-A00A-1234-1234-5864377B4831", // We've pre-populated the mock db with some values
			timeProvider:           timeProvider,
			expectedBatchSize:      50,
			expectedCleanSync:      false,
			expectedLastCleanSync:  "",
			expectedPreflightAt:    clock.RFC3339(frozenTime),
			expectedFeedSyncCursor: "",
			expectedExpiresAfter:   GetSyncStateExpiresAfter(timeProvider),
			expectedDataType:       GetDataType(),
		},
	}

	for _, test := range cases {
		result := CreateNewSyncState(timeProvider, test.machineID, test.expectedCleanSync, test.expectedLastCleanSync, test.expectedBatchSize, test.expectedFeedSyncCursor)

		if test.expectedBatchSize != 0 {
			assert.NotEmpty(t, result)
			assert.Equal(t, test.expectedBatchSize, result.BatchSize)
		}

		if test.expectedCleanSync != false {
			assert.NotEmpty(t, result)
			assert.Equal(t, test.expectedCleanSync, result.CleanSync)
			assert.Equal(t, test.expectedLastCleanSync, result.LastCleanSync)
		}

		if test.expectedPreflightAt != "" {
			assert.NotEmpty(t, result)
			assert.Equal(t, test.expectedPreflightAt, result.PreflightAt)
		}

		if test.expectedFeedSyncCursor != "" {
			assert.NotEmpty(t, result)
			assert.Equal(t, test.expectedFeedSyncCursor, result.FeedSyncCursor)
		}

		if test.expectedDataType != "" {
			assert.NotEmpty(t, result)
			assert.Equal(t, test.expectedDataType, result.DataType)
		}

		if test.expectedExpiresAfter != 0 {
			assert.NotEmpty(t, result)
			assert.Equal(t, test.expectedExpiresAfter, result.ExpiresAfter)
		}

	}
}
