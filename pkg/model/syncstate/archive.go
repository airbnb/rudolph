package syncstate

import (
	"fmt"
	"time"

	"github.com/airbnb/rudolph/pkg/dynamodb"
)

// Archive will take the current item and duplicate it with a new sort key. This creates a unique item
// every time it is called, so can be used to save older copies of this row as it is mutated over time
// in order to preserve history.
func Archive(client dynamodb.PutItemAPI, syncState SyncStateRow) error {
	clone := syncState
	clone.SortKey = fmt.Sprintf("%s@%s", syncState.SortKey, time.Now().UTC().Format(time.RFC3339))

	_, err := client.PutItem(clone)

	return err
}
