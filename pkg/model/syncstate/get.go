package syncstate

import (
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/pkg/errors"
)

// Returns the machine's current sync state
func GetByMachineID(client dynamodb.GetItemAPI, machineID string) (syncState *SyncStateRow, err error) {
	output, err := client.GetItem(
		dynamodb.PrimaryKey{
			PartitionKey: syncStatePK(machineID),
			SortKey:      syncStateSK,
		},
		// The sync state must be retrieved consistently because it is set by a /preflight request and then immediately
		// retrieved again in the /ruledownload step. Unfortunately it cannot be passed from the /preflight step directly
		// to the ruledownload via some sort of "initial cursor" (At least not as far as I can tell).
		true,
	)
	if err != nil {
		return
	}

	if len(output.Item) == 0 {
		return
	}

	err = attributevalue.UnmarshalMap(output.Item, &syncState)

	if err != nil {
		err = errors.Wrap(err, "succeeded GetItem but failed to unmarshalMap into SyncState")
		return
	}

	return
}
