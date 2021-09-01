package syncstate

import (
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
)

func CreateNewSyncState(timeProvider clock.TimeProvider, machineID string, requestCleanSync bool, lastCleanSync string, batchSize int, feedSyncCursor string) SyncStateRow {
	return SyncStateRow{
		PrimaryKey: dynamodb.PrimaryKey{
			PartitionKey: syncStatePK(machineID),
			SortKey:      syncStateSK,
		},
		SyncState: SyncState{
			MachineID:      machineID,
			CleanSync:      requestCleanSync,
			BatchSize:      batchSize,
			LastCleanSync:  lastCleanSync,
			PreflightAt:    clock.RFC3339(timeProvider.Now()),
			FeedSyncCursor: feedSyncCursor,
			ExpiresAfter:   GetSyncStateExpiresAfter(timeProvider),
			DataType:       GetDataType(),
		},
	}
}
