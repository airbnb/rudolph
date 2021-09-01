package syncstate

import (
	"fmt"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/types"
)

const (
	machineInfoPKPrefix             = "Machine#"
	syncStateSK                     = "SyncState"
	syncStateExpiresAfterInDays int = 90
)

// SyncStateRow persists the sync state to the database
type SyncStateRow struct {
	dynamodb.PrimaryKey
	SyncState
}

// SyncState encapulsates data within the sync state row that's not the Primary key
type SyncState struct {
	MachineID              string         `dynamodbav:"MachineID"`
	BatchSize              int            `dynamodbav:"BatchSize"`
	CleanSync              bool           `dynamodbav:"CleanSync"`
	LastCleanSync          string         `dynamodbav:"LastCleanSync"`
	FeedSyncCursor         string         `dynamodbav:"FeedSyncCursor"`
	PreflightAt            string         `dynamodbav:"PreflightAt"`
	RuledownloadStartedAt  string         `dynamodbav:"RuledownloadStartedAt"`
	RuledownloadFinishedAt string         `dynamodbav:"RuledownloadFinishedAt"`
	PostflightAt           string         `dynamodbav:"PostflightAt"`
	ExpiresAfter           int64          `dynamodbav:"ExpiresAfter,omitempty"`
	DataType               types.DataType `dynamodbav:"DataType"`
}

// Update fragments
type UpdatePostflightItem struct {
	PostflightAt   string `dynamodbav:"PostflightAt"`
	FeedSyncCursor string `dynamodbav:"FeedSyncCursor"`
	ExpiresAfter   int64  `dynamodbav:"ExpiresAfter,omitempty"`
}
type updateRuledownloadAtItem struct {
	RuledownloadStartedAt string `dynamodbav:"RuledownloadStartedAt"`
}
type updateRuledownloadFinishedAtItem struct {
	RuledownloadFinishedAt string `dynamodbav:"RuledownloadFinishedAt"`
}

func syncStatePK(machineID string) string {
	return fmt.Sprintf("%s%s", machineInfoPKPrefix, machineID)
}

func GetSyncStateExpiresAfter(timeProvider clock.TimeProvider) int64 {
	return clock.Unixtimestamp(timeProvider.Now().UTC().AddDate(0, 0, syncStateExpiresAfterInDays))
}

func GetDataType() types.DataType {
	return types.DataTypeSyncState
}
